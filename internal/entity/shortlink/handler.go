package shortlink

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/NickSFU/shortlink-service/internal/entity/click"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service      *Service
	clickService *click.Service
}

type statsResponse struct {
	Clicks int `json:"clicks"`
}

func NewHandler(
	service *Service,
	clickService *click.Service,
) *Handler {
	return &Handler{
		service:      service,
		clickService: clickService,
	}
}

type createRequest struct {
	URL string `json:"url"`
}

type createResponse struct {
	Code string `json:"code"`
}

// POST /shorten
func (h *Handler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req createRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	code, err := h.service.CreateShortLink(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := createResponse{
		Code: code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(resp)
}

// GET /{code}
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := h.service.GetLink(code)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	clickData := &click.Click{
		ShortLinkID: link.ID,
		IP:          ip,
		UserAgent:   r.UserAgent(),
		Referer:     r.Referer(),
	}

	// статистика не должна ломать редирект
	go h.clickService.Create(clickData)

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := h.service.GetLink(code)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	count, err := h.clickService.CountByLink(link.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := statsResponse{
		Clicks: count,
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(resp)
}
