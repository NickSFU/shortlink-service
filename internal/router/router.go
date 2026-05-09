package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/NickSFU/shortlink-service/internal/entity/click"
	"github.com/NickSFU/shortlink-service/internal/entity/shortlink"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	// shortlink dependencies
	shortRepo := shortlink.NewRepository(db)
	shortService := shortlink.NewService(shortRepo)

	// click dependencies
	clickRepo := click.NewRepository(db)
	clickService := click.NewService(clickRepo)

	// handler
	handler := shortlink.NewHandler(
		shortService,
		clickService,
	)

	// health-check
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// shortlink routes
	r.Post("/shorten", handler.CreateShortLink)

	// statistics
	r.Get("/stats/{code}", handler.GetStats)

	// redirect
	r.Get("/{code}", handler.Redirect)

	return r
}
