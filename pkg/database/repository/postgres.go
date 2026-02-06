package repository

import (
	"agrigation_api/pkg/models"
	"agrigation_api/pkg/tools"
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func InitRepository() (*Repository, error) {
	pool, err := InitPGPool()
	if err != nil {
		return nil, err
	}
	return &Repository{
		pool: pool,
	}, nil

}

func InitPGPool() (*pgxpool.Pool, error) {
	pgUser := tools.GetEnv("PG_USER", "postgres")
	pgPassword := tools.GetEnv("PG_PASSWORD", "postgres")
	pgHost := tools.GetEnv("PG_HOST", "localhost")
	pgPort := tools.GetEnvAsInt("PG_PORT", 5432)
	pgDatabase := tools.GetEnv("PG_DATABASE", "aggregation")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", pgUser, pgPassword, pgHost, pgPort, pgDatabase)

	pool, errPGX := pgxpool.New(context.Background(), connStr)

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	if errPGX != nil {
		return nil, errPGX
	}
	return pool, nil
}

// UpsertSubscription - создать подписку
func (r *Repository) UpsertSubscription(ctx context.Context, req models.CreateOrUpdateRequest) (*models.Subscription, error) {
	query := `
    INSERT INTO subscriptions 
    (user_id, service_name, price, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (user_id, service_name) 
    DO UPDATE SET 
        price = EXCLUDED.price,
        start_date = EXCLUDED.start_date,
        end_date = EXCLUDED.end_date,
        updated_at = CURRENT_TIMESTAMP
    RETURNING user_id, service_name, price, start_date, end_date, created_at, updated_at`

	var sub models.Subscription
	var endDate sql.NullString

	var endDateVal interface{}
	if req.EndDate == "" {
		endDateVal = nil
	} else {
		endDateVal = req.EndDate
	}

	err := r.pool.QueryRow(ctx, query,
		req.UserID,
		req.ServiceName,
		req.Price,
		req.StartDate,
		endDateVal,
	).Scan(
		&sub.UserID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&endDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("upsert subscription: %w", err)
	}

	if endDate.Valid {
		sub.EndDate = &endDate.String
	}

	return &sub, nil
}

// GetSubscription - получение подписки у пользователя
func (r *Repository) GetSubscription(ctx context.Context, userID uuid.UUID, serviceName string) (*models.Subscription, error) {
	query := `
    SELECT user_id, service_name, price, start_date, end_date, created_at, updated_at
    FROM subscriptions 
    WHERE user_id = $1 AND service_name = $2`

	var sub models.Subscription
	var endDate sql.NullString

	err := r.pool.QueryRow(ctx, query, userID, serviceName).Scan(
		&sub.UserID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&endDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if endDate.Valid {
		sub.EndDate = &endDate.String
	}

	return &sub, nil
}

// DeleteSubscription - удаление подписки у пользователя
func (r *Repository) DeleteSubscription(ctx context.Context, userID uuid.UUID, serviceName string) error {
	query := `DELETE FROM subscriptions WHERE user_id = $1 AND service_name = $2`

	result, err := r.pool.Exec(ctx, query, userID, serviceName)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// ListUserSubscriptions - получение списка подписок у пользователя
func (r *Repository) ListUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	query := `
    SELECT user_id, service_name, price, start_date, end_date, created_at, updated_at
    FROM subscriptions 
    WHERE user_id = $1
    ORDER BY service_name`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list user subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		var endDate sql.NullString

		err := rows.Scan(
			&sub.UserID,
			&sub.ServiceName,
			&sub.Price,
			&sub.StartDate,
			&endDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}

		if endDate.Valid {
			sub.EndDate = &endDate.String
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *Repository) CalculateTotal(ctx context.Context, req models.CalculateTotalRequest) (int, error) {
	startTime, err := tools.ParseMonthYear(req.StartMonth)
	if err != nil {
		return 0, fmt.Errorf("invalid start_month: %w", err)
	}

	endTime, err := tools.ParseMonthYear(req.EndMonth)
	if err != nil {
		return 0, fmt.Errorf("invalid end_month: %w", err)
	}

	if startTime.After(endTime) {
		return 0, fmt.Errorf("start_month must be before end_month")
	}

	// Строим запрос
	query := `
    SELECT COALESCE(SUM(price), 0) 
    FROM subscriptions 
    WHERE 1=1`

	args := make([]interface{}, 0)
	argNum := 1

	// Фильтр по пользователю
	if req.UserID != uuid.Nil {
		query += " AND user_id = $" + strconv.Itoa(argNum)
		args = append(args, req.UserID)
		argNum++
	}

	// Фильтр по сервису
	if req.ServiceName != "" {
		query += " AND service_name = $" + strconv.Itoa(argNum)
		args = append(args, req.ServiceName)
		argNum++
	}

	// Фильтр по периоду
	// Подписка активна в период, если:
	// 1. start_date <= end_of_period (подписка началась до конца периода)
	// 2. end_date IS NULL OR end_date >= start_of_period (подписка активна в начале периода или бессрочная)
	query += " AND start_date <= $" + strconv.Itoa(argNum)
	argNum++
	query += " AND (end_date IS NULL OR end_date >= $" + strconv.Itoa(argNum) + ")"

	args = append(args, endTime.Format("01-2006"), startTime.Format("01-2006"))

	var total int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("calculate total: %w", err)
	}

	return total, nil
}
