package repository

import (
	repository "agrigation_api/internal/database/postgres"
	"agrigation_api/pkg/models"
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSubscription(context.Context, models.CreateOrUpdateRequest) (*models.Subscription, error)
	UpdateSubscription(context.Context, models.CreateOrUpdateRequest) (*models.Subscription, error)
	GetSubscription(context.Context, uuid.UUID, string) (*models.Subscription, error)
	DeleteSubscription(context.Context, uuid.UUID, string) error
	ListUserSubscriptions(context.Context, uuid.UUID) ([]models.Subscription, error)
	CalculateTotal(context.Context, models.CalculateTotalRequest) (int, error)
	CloseConnection()
}

func InitRepository() (Repository, error) {
	return repository.InitPostgres()
}
