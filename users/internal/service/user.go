package service

import (
	"context"
	"net/url"

	"github.com/PratikKumar125/go-microservices/pkg/logging"
	"github.com/PratikKumar125/go-microservices/users/internal/config"
	"github.com/PratikKumar125/go-microservices/users/internal/models"
	"github.com/PratikKumar125/go-microservices/users/internal/repositories"
)

type UserService struct {
	appConfig      *config.AppConfig
	logger         *logging.Logger
	userRepository *repositories.UserRepository
}

func NewUserService(cfg *config.AppConfig, logger *logging.Logger, repo *repositories.UserRepository) *UserService {
	return &UserService{
		appConfig:      cfg,
		logger:         logger,
		userRepository: repo,
	}
}

func (service UserService) GetUsers(ctx context.Context, req url.Values) ([]models.User, error) {
	filters := &models.UserSearchFilters{
		Email: req.Get("email"),
		Name:  req.Get("name"),
	}
	users, err := service.userRepository.Search(ctx, filters)
	if err != nil {
		service.logger.Error("Error getting users", "error", err)
		return nil, err
	}

	return users, err
}
