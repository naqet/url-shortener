package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/naqet/url-shortener/services"
)

type UrlHandler struct {
	service *services.UrlService
}

func NewUrlHandler(service *services.UrlService) *UrlHandler {
	return &UrlHandler{service}
}

func (h *UrlHandler) Post(w http.ResponseWriter, r *http.Request) {
	originalUrl := r.FormValue("url")

	if originalUrl == "" {
		http.Error(w, "Missing URL value", http.StatusBadRequest)
		return
	}

	key, err := h.service.CreateNewURL(originalUrl)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"newUrl": "http://localhost:3000/" + key}

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Failed to convert response to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}

func (h *UrlHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	url, err := h.service.GetURL(key)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"id":          url.Id,
		"key":         "http://localhost:3000/" + url.Key,
		"originalUrl": url.OriginalUrl,
	}

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Failed to convert response to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}

func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	url, err := h.service.GetURL(key)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url.OriginalUrl, 302)
	return
}
