package service

import (
	"agrigation_api/internal/database/repository"
	"agrigation_api/pkg/models"
	"context"
	"github.com/google/uuid"
)

type Subscriptions interface {
	CreateSubscription(context.Context, models.CreateOrUpdateRequest) (*models.Subscription, error)
	UpdateSubscription(context.Context, models.CreateOrUpdateRequest) (*models.Subscription, error)
	GetSubscription(context.Context, uuid.UUID, string) (*models.Subscription, error)
	DeleteSubscription(context.Context, uuid.UUID, string) error
	ListSubscriptions(context.Context, uuid.UUID) ([]models.Subscription, error)
	CalculateTotal(context.Context, models.CalculateTotalRequest) (int, error)
}

type SubscriptionService struct {
	rep repository.Repository
}

func NewSubscriptionService(rep repository.Repository) *SubscriptionService {
	return &SubscriptionService{rep}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req models.CreateOrUpdateRequest) (*models.Subscription, error) {
	return s.rep.CreateSubscription(ctx, req)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, req models.CreateOrUpdateRequest) (*models.Subscription, error) {
	return s.rep.UpdateSubscription(ctx, req)
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, req uuid.UUID, name string) (*models.Subscription, error) {
	return s.rep.GetSubscription(ctx, req, name)
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, req uuid.UUID, name string) error {
	return s.rep.DeleteSubscription(ctx, req, name)
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, req uuid.UUID) ([]models.Subscription, error) {
	return s.rep.ListUserSubscriptions(ctx, req)
}

func (s *SubscriptionService) CalculateTotal(ctx context.Context, req models.CalculateTotalRequest) (int, error) {
	return s.rep.CalculateTotal(ctx, req)
}
