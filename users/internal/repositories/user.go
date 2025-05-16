package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/PratikKumar125/go-microservices/users/internal/config"
	"github.com/PratikKumar125/go-microservices/users/internal/db"
	"github.com/PratikKumar125/go-microservices/users/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryContract interface {
	All(ctx context.Context) ([]models.User, error)
}

type UserRepository struct {
	db        *db.Database
	appConfig *config.AppConfig
}

func NewUserRepository(config *config.AppConfig, db *db.Database) *UserRepository {
	return &UserRepository{
		appConfig: config,
		db:        db,
	}
}

func (repo *UserRepository) All(ctx context.Context) ([]models.User, error) {
	query := `SELECT * FROM USERS LIMIT 15`
	rows, err := repo.db.Pool().Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
}

func (repo *UserRepository) Search(ctx context.Context, filters *models.UserSearchFilters) ([]models.User, error) {
	query := `SELECT id,name,email FROM USERS`

	var conditions []string

	if filters.Email != "" {
		current := fmt.Sprintf("%s='%s'", "LOWER(email)", strings.ToLower(filters.Email))
		conditions = append(conditions, current)
	}
	if filters.Name != "" {
		current := fmt.Sprintf("name ILIKE '%%%s%%'", filters.Name)
		conditions = append(conditions, current)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	} else {
		query += strings.Join(conditions, " AND ")
	}

	rows, err := repo.db.Pool().Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
}
