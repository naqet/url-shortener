package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/naqet/url-shortener/services"
)

var (
	secretKey     = []byte("secretKey")
	refreshKey    = []byte("refreshKey")
	tokenExpire   = time.Hour
	refreshExpire = time.Hour * 24 * 7
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (handler *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if name == "" || password == "" || email == "" {
		http.Error(w, "Insufficient credentials", http.StatusBadRequest)
		return
	}

	err := handler.service.SignUp(name, password, email)

	if errors.Is(err, services.ErrorUserExists) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    w.Header().Add("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if password == "" || email == "" {
		http.Error(w, "Incorrect credentials", http.StatusBadRequest)
		return
	}

	id, err := handler.service.Login(email, password)

    if errors.Is(err, services.ErrorIncorrectEmailOrPass) {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
    } else if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = handler.handleTokens(w, id)

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    w.Header().Add("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")

		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := handler.service.ValidateToken(cookie.Value, secretKey)

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

        claims, ok := token.Claims.(jwt.MapClaims);

        if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized);
            return;
        }

        id, ok := claims["id"].(string);

        if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized);
            return;
        }

        ctx := context.WithValue(r.Context(), "id", id);

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil || cookie.Value == "" {
		http.Error(w, "Refresh token missing", http.StatusBadRequest)
		return
	}

	token, err := handler.service.ValidateToken(cookie.Value, refreshKey)

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		http.Error(w, "Invalid refresh token claims", http.StatusUnauthorized)
		return
	}

	id, ok := claims["id"].(string)

	if !ok {
		http.Error(w, "Invalid refresh token claims", http.StatusUnauthorized)
		return
	}

	err = handler.handleTokens(w, id)

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) handleTokens(w http.ResponseWriter, id string) error {
	newToken, err := handler.service.CreateToken(id, tokenExpire, secretKey)

	if err != nil {
		return err
	}

	newRefreshToken, err := handler.service.CreateToken(id, refreshExpire, refreshKey)

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newToken,
		Expires:  time.Now().Add(tokenExpire),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(refreshExpire),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}
