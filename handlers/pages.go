package handlers

import (
	"net/http"

	"github.com/naqet/url-shortener/services"
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
