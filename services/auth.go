package services

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/naqet/url-shortener/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUserExists           = errors.New("User with given email already exists")
	ErrorIncorrectEmailOrPass = errors.New("Incorrect email or password")
)

type AuthService struct {
	database *db.DB
}

func NewAuthService(database *db.DB) (*AuthService, error) {
	if database == nil {
		return &AuthService{}, errors.New("Database has not been provided")
	}
	return &AuthService{database}, nil
}

func (service *AuthService) SignUp(name, password, email string) error {
	_, err := service.database.GetUserByEmail(email)

	if err == nil {
		return ErrorUserExists
	} else if err != sql.ErrNoRows {
		return err
	}

	safePass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return err
	}

	return service.database.CreateUser(name, string(safePass), email)
}

func (service *AuthService) Login(email, password string) (string, error) {
	user, err := service.database.GetUserByEmail(email)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrorIncorrectEmailOrPass
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", ErrorIncorrectEmailOrPass
	}

	return user.Id, nil
}

func (service *AuthService) CreateToken(userId string, expireTime time.Duration, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(expireTime).Unix(),
	})

	return token.SignedString(secret)
}

func (service *AuthService) ValidateToken(tokenString string, key []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}

func (service *AuthService) IsUserLogged(r *http.Request) bool {
    _, ok := r.Context().Value("id").(string);
    return ok
}
