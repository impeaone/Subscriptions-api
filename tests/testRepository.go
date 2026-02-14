package tests

import (
	"agrigation_api/pkg/models"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TestRepositoryPool пул тестовых репозиториев
type TestRepositoryPool struct {
	mu           sync.RWMutex
	repositories map[string]*TestRepository
}

// TestRepository тестовый репозиторий
type TestRepository struct {
	mu            sync.RWMutex
	subscriptions map[string]map[string]*models.Subscription // userID -> subscriptionID -> subscription
	shouldFail    bool                                       // флаг для имитации ошибок
	failOnMethod  string                                     // на каком методе фейлить
	closeCalled   bool                                       // был ли вызван CloseConnection
}

// NewTestRepositoryPool создает новый пул тестовых репозиториев
func NewTestRepositoryPool() *TestRepositoryPool {
	return &TestRepositoryPool{
		repositories: make(map[string]*TestRepository),
	}
}

// GetRepository возвращает или создает репозиторий для пользователя
func (p *TestRepositoryPool) GetRepository(userID string) *TestRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if repo, exists := p.repositories[userID]; exists {
		return repo
	}

	repo, _ := NewTestRepository()
	p.repositories[userID] = repo
	return repo
}

// NewTestRepository создает новый тестовый репозиторий
func NewTestRepository() (*TestRepository, error) {
	return &TestRepository{
		subscriptions: make(map[string]map[string]*models.Subscription),
	}, nil
}

// CreateSubscription создает новую подписку
func (t *TestRepository) CreateSubscription(ctx context.Context, req models.CreateOrUpdateRequest) (*models.Subscription, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Проверка на имитацию ошибки
	if t.shouldFail && (t.failOnMethod == "" || t.failOnMethod == "CreateSubscription") {
		return nil, errors.New("simulated error in CreateSubscription")
	}

	// Парсим даты
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format, expected MM-YYYY")
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date format, expected MM-YYYY")
		}
		endDate = &parsed
	}

	now := time.Now()

	// Создаем новую подписку
	subscription := &models.Subscription{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Инициализируем мапу для пользователя если её нет
	userKey := req.UserID.String()
	if _, exists := t.subscriptions[userKey]; !exists {
		t.subscriptions[userKey] = make(map[string]*models.Subscription)
	}

	key := req.ServiceName
	t.subscriptions[userKey][key] = subscription

	return subscription, nil
}

// UpdateSubscription обновляет существующую подписку
func (t *TestRepository) UpdateSubscription(ctx context.Context, req models.CreateOrUpdateRequest) (*models.Subscription, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.shouldFail && (t.failOnMethod == "" || t.failOnMethod == "UpdateSubscription") {
		return nil, errors.New("simulated error in UpdateSubscription")
	}

	userKey := req.UserID.String()
	userSubs, exists := t.subscriptions[userKey]
	if !exists {
		return nil, errors.New("user has no subscriptions")
	}

	// Ищем подписку по сервису и дате старта
	key := req.ServiceName
	subscription, exists := userSubs[key]
	if !exists {
		return nil, errors.New("subscription not found")
	}

	// Обновляем поля (только те что пришли в запросе)
	if req.ServiceName != "" && req.ServiceName != subscription.ServiceName {
		// Если меняется сервис, нужно пересоздать ключ
		return nil, errors.New("service_name cannot be updated, delete and create new subscription")
	}

	if req.Price > 0 {
		subscription.Price = req.Price
	}

	if req.EndDate != "" {
		parsed, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date format, expected MM-YYYY")
		}
		subscription.EndDate = &parsed
	}

	subscription.UpdatedAt = time.Now()

	return subscription, nil
}

// GetSubscription получает подписку по сервису и дате старта
func (t *TestRepository) GetSubscription(ctx context.Context, userID uuid.UUID, serviceName string) (*models.Subscription, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.shouldFail && (t.failOnMethod == "" || t.failOnMethod == "GetSubscription") {
		return nil, errors.New("simulated error in GetSubscription")
	}

	userKey := userID.String()
	userSubs, exists := t.subscriptions[userKey]
	if !exists {
		return nil, errors.New("user has no subscriptions")
	}

	subscription, exists := userSubs[serviceName]
	if !exists {
		return nil, errors.New("subscription not found")
	}

	return subscription, nil
}

// DeleteSubscription удаляет подписку
func (t *TestRepository) DeleteSubscription(ctx context.Context, userID uuid.UUID, serviceName string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.shouldFail && (t.failOnMethod == "" || t.failOnMethod == "DeleteSubscription") {
		return errors.New("simulated error in DeleteSubscription")
	}

	userKey := userID.String()
	userSubs, exists := t.subscriptions[userKey]
	if !exists {
		return errors.New("user has no subscriptions")
	}

	key := serviceName
	if _, exists := userSubs[key]; !exists {
		return errors.New("subscription not found")
	}

	delete(userSubs, key)

	// Если у пользователя больше нет подписок, удаляем мапу
	if len(userSubs) == 0 {
		delete(t.subscriptions, userKey)
	}

	return nil
}

// ListUserSubscriptions возвращает все подписки пользователя
func (t *TestRepository) ListUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.shouldFail && (t.failOnMethod == "" || t.failOnMethod == "ListUserSubscriptions") {
		return nil, errors.New("simulated error in ListUserSubscriptions")
	}
	userKey := userID.String()
	userSubs, exists := t.subscriptions[userKey]
	if !exists {
		return []models.Subscription{}, nil // возвращаем пустой слайс, не ошибку
	}

	result := make([]models.Subscription, 0, len(userSubs))
	for _, sub := range userSubs {
		result = append(result, *sub)
	}

	return result, nil
}

// CalculateTotal вычисляет общую сумму за период
func (t *TestRepository) CalculateTotal(ctx context.Context, req models.CalculateTotalRequest) (int, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.shouldFail && (t.failOnMethod == "" || t.failOnMethod == "CalculateTotal") {
		return 0, errors.New("simulated error in CalculateTotal")
	}

	var total int
	var targetUserID = req.UserID
	var targetService = req.ServiceName
	var startMonth = req.StartMonth
	var endMonth = req.EndMonth

	// Перебираем всех пользователей или конкретного
	for userKey, userSubs := range t.subscriptions {
		if targetUserID != uuid.Nil {
			// Если указан конкретный пользователь, пропускаем остальных
			if userKey != targetUserID.String() {
				continue
			}
		}

		for _, sub := range userSubs {
			// Фильтр по сервису
			if targetService != "" && sub.ServiceName != targetService {
				continue
			}

			subStart := sub.StartDate
			subEnd := sub.EndDate

			// Нормализуем даты до первого числа месяца для сравнения
			startMonth = time.Date(startMonth.Year(), startMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
			endMonth = time.Date(endMonth.Year(), endMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
			subStart = time.Date(subStart.Year(), subStart.Month(), 1, 0, 0, 0, 0, time.UTC)

			var subEndNormalized *time.Time
			if subEnd != nil {
				endNorm := time.Date(subEnd.Year(), subEnd.Month(), 1, 0, 0, 0, 0, time.UTC)
				subEndNormalized = &endNorm
			}

			if !subStart.After(endMonth) {
				if subEndNormalized == nil || !subEndNormalized.Before(startMonth) {
					total += sub.Price
				}
			}
		}
	}

	return total, nil
}

// CloseConnection помечает соединение как закрытое
func (t *TestRepository) CloseConnection() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closeCalled = true
}

// Reset очищает все данные в репозитории
func (t *TestRepository) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.subscriptions = make(map[string]map[string]*models.Subscription)
	t.shouldFail = false
	t.failOnMethod = ""
	t.closeCalled = false
}

// GetSubscriptionCount возвращает количество подписок пользователя
func (t *TestRepository) GetSubscriptionCount(userID uuid.UUID) int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	userSubs, exists := t.subscriptions[userID.String()]
	if !exists {
		return 0
	}
	return len(userSubs)
}
