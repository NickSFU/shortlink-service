package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/NickSFU/shortlink-service/internal/entity/click"
	"github.com/NickSFU/shortlink-service/internal/entity/shortlink"
	"github.com/NickSFU/shortlink-service/internal/entity/user"
	"github.com/NickSFU/shortlink-service/internal/middleware"
)

func NewRouter(
	db *pgxpool.Pool,
	cache *redis.Client,
) http.Handler {

	r := chi.NewRouter()

	// shortlink
	shortRepo := shortlink.NewRepository(db)
	shortService := shortlink.NewService(
		shortRepo,
		cache,
	)

	// click
	clickRepo := click.NewRepository(db)
	clickService := click.NewService(clickRepo)

	// shortlink handler
	shortHandler := shortlink.NewHandler(
		shortService,
		clickService,
	)

	// user
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// ping
	r.Get("/ping", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// shortlink
	r.With(middleware.Auth).
		Post("/shorten", shortHandler.CreateShortLink)

	r.With(middleware.Auth).
		Get("/my-links", shortHandler.GetMyLinks)

	r.With(middleware.Auth).
		Delete("/links/{code}", shortHandler.DeleteLink)

	r.With(middleware.Auth).
		Patch("/links/{code}", shortHandler.UpdateLink)

	// analytics
	r.With(middleware.Auth).
		Get("/stats/{code}", shortHandler.GetStats)
	// redirect
	r.Get("/{code}", shortHandler.Redirect)

	// auth
	r.Post("/auth/register", userHandler.Register)
	r.Post("/auth/login", userHandler.Login)

	return r
}
