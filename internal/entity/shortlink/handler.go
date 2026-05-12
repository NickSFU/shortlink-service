package shortlink

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/NickSFU/shortlink-service/internal/entity/click"
	"github.com/NickSFU/shortlink-service/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service      *Service
	clickService *click.Service
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

type updateRequest struct {
	URL string `json:"url"`
}

type createResponse struct {
	Code string `json:"code"`
}

// POST /shorten
func (h *Handler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	userID := r.Context().
		Value(middleware.UserIDKey).(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	code, err := h.service.CreateShortLink(userID, req.URL)
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

func (h *Handler) GetMyLinks(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().
		Value(middleware.UserIDKey).(int)

	links, err := h.service.GetUserLinks(userID)
	if err != nil {
		http.Error(
			w,
			"internal error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(links)
}

func (h *Handler) DeleteLink(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := chi.URLParam(r, "code")

	link, err := h.service.GetLink(code)
	if err != nil {
		http.Error(
			w,
			"link not found",
			http.StatusNotFound,
		)
		return
	}

	userID := r.Context().
		Value(middleware.UserIDKey).(int)

	// ownership check
	if link.UserID != userID {
		http.Error(
			w,
			"forbidden",
			http.StatusForbidden,
		)
		return
	}

	err = h.service.DeleteLink(code)
	if err != nil {
		http.Error(
			w,
			"internal error",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateLink(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := chi.URLParam(r, "code")

	link, err := h.service.GetLink(code)
	if err != nil {
		http.Error(
			w,
			"link not found",
			http.StatusNotFound,
		)
		return
	}

	userID := r.Context().
		Value(middleware.UserIDKey).(int)

	// ownership check
	if link.UserID != userID {
		http.Error(
			w,
			"forbidden",
			http.StatusForbidden,
		)
		return
	}

	var req updateRequest

	err = json.NewDecoder(r.Body).
		Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid body",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.UpdateLink(
		code,
		req.URL,
	)
	if err != nil {
		http.Error(
			w,
			"internal error",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetStats(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := chi.URLParam(r, "code")

	link, err := h.service.GetLink(code)
	if err != nil {
		http.Error(
			w,
			"link not found",
			http.StatusNotFound,
		)
		return
	}

	userID := r.Context().
		Value(middleware.UserIDKey).(int)

	if link.UserID != userID {
		http.Error(
			w,
			"forbidden",
			http.StatusForbidden,
		)
		return
	}

	stats, err := h.clickService.GetStats(link.ID)
	if err != nil {
		http.Error(
			w,
			"internal error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(stats)
}
