package handler

import (
	"time"

	"github.com/PratikKumar125/go-microservices/users/internal/config"
	"github.com/PratikKumar125/go-microservices/users/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UserHandler struct {
	userService   *service.UserService
	appConfig *config.AppConfig
}

func NewUserHandler(cfg *config.AppConfig, userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService:   userService,
		appConfig: cfg,
	}
}

func (h *UserHandler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Timeout(time.Second * time.Duration(h.appConfig.ConfigService.Int64("app.timeout"))))
	r.Use(AuthorizeUserGuard)

	r.Mount("/debug", middleware.Profiler())
	r.Get("/users", h.getAllUsers)

	return r
}
