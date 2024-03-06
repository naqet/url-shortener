package handlers

import (
	"net/http"

	"github.com/naqet/url-shortener/services"
	"github.com/naqet/url-shortener/views/auth"
	"github.com/naqet/url-shortener/views/dashboard"
	"github.com/naqet/url-shortener/views/home"
)

type PagesHandler struct {
	authService *services.AuthService
	urlService *services.UrlService
}

func NewPagesHandler(authService *services.AuthService, urlService *services.UrlService) *PagesHandler {
    return &PagesHandler{authService, urlService}
}

func (han *PagesHandler) Home(w http.ResponseWriter, r *http.Request) {
    isLogged := han.authService.IsUserLogged(r);
    home.Index(isLogged).Render(r.Context(), w);
}

func (han *PagesHandler) SignUp(w http.ResponseWriter, r *http.Request) {
    auth.SignUp().Render(r.Context(), w);
}

func (han *PagesHandler) Login(w http.ResponseWriter, r *http.Request) {
    auth.Login().Render(r.Context(), w);
}

func (han *PagesHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
    isLogged := han.authService.IsUserLogged(r);
    dashboard.Index(isLogged).Render(r.Context(), w);
}
