package services

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/naqet/url-shortener/db"
	"github.com/naqet/url-shortener/models"
)

var availableChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

type UrlService struct {
	database *db.DB
}

func NewUrlService(database *db.DB) (UrlService, error) {
	if database == nil {
		return UrlService{}, errors.New("Database has not been provided")
	}
	return UrlService{database}, nil
}

func (service *UrlService) CreateNewURL(originalUrl string) (string, error) {
    key := randomKey(20);
    err := service.database.SaveURL(key, originalUrl);

    if err != nil {
        return "", err;
    }
    return key, nil;
}

func (service *UrlService) GetURL(key string) (models.ShortUrl, error) {
    result, err := service.database.GetURL(key);
    if err != nil {
        return models.ShortUrl{}, err
    }

    return result, nil;
}

func randomKey(length int) string {
    sb := strings.Builder{}
    sb.Grow(length);

    for i := 0; i< length; i++ {
        sb.WriteRune(availableChars[rand.Intn(len(availableChars))])
    }

    return sb.String();
}
