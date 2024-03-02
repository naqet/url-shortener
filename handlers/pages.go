package handlers

import (
	"net/http"

	"github.com/naqet/url-shortener/services"
	"github.com/naqet/url-shortener/views/auth"
	"github.com/naqet/url-shortener/views/home"
)

type PagesHandler struct {
	service *services.AuthService
}

func NewPagesHandler(service *services.AuthService) *PagesHandler {
    return &PagesHandler{service}
}

func (han *PagesHandler) Home(w http.ResponseWriter, r *http.Request) {
    home.Index().Render(r.Context(), w);
}

func (han *PagesHandler) SignUp(w http.ResponseWriter, r *http.Request) {
    auth.SignUp().Render(r.Context(), w);
}

func (han *PagesHandler) Login(w http.ResponseWriter, r *http.Request) {
    auth.Login().Render(r.Context(), w);
}
