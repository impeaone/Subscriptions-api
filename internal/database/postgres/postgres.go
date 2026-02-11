package postgres

import (
	"agrigation_api/pkg/models"
	"agrigation_api/pkg/tools"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func InitPostgres() (*Repository, error) {
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
	pgHost := tools.GetEnv("PG_HOST", "192.168.3.92")
	pgPort := tools.GetEnvAsInt("PG_PORT", 5432)
	pgDatabase := tools.GetEnv("PG_DATABASE", "storage")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", pgUser, pgPassword, pgHost, pgPort, pgDatabase)

	pool, errPGX := pgxpool.New(context.Background(), connStr)
	if errPGX != nil {
		return nil, errPGX
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

// CreateSubscription - создать подписку
func (r *Repository) CreateSubscription(ctx context.Context, req models.CreateOrUpdateRequest) (*models.Subscription, error) {
	query := `
    INSERT INTO subscriptions 
    (user_id, service_name, price, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING user_id, service_name, price, start_date, end_date, created_at, updated_at`

	var sub models.Subscription
	var endDate sql.NullString

	var endDateVal interface{}
	if req.EndDate == "" {
		endDateVal = nil
	} else {
		endDateVal = req.EndDate
	}
	var end *time.Time = nil
	if endDateVal != nil {
		*end, _ = tools.ParseMonthYear(req.EndDate)
	}

	start, errSt := tools.ParseMonthYear(req.StartDate)
	if errSt != nil {
		return nil, errSt
	}

	err := r.pool.QueryRow(ctx, query,
		req.UserID,
		req.ServiceName,
		req.Price,
		start,
		end,
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
		*sub.EndDate, _ = time.Parse("01-2006", endDate.String)
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

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	/*
		if endDate.Valid {
			sub.EndDate = &endDate.String
		}*/
	if endDate.Valid {
		*sub.EndDate, _ = time.Parse("01-2006", endDate.String)
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
		/*
			if endDate.Valid {
				sub.EndDate = &endDate.String
			}*/
		if endDate.Valid {
			*sub.EndDate, _ = time.Parse("01-2006", endDate.String)
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *Repository) CalculateTotal(ctx context.Context, req models.CalculateTotalRequest) (int, error) {
	if req.StartMonth.After(req.EndMonth) {
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

	args = append(args, req.EndMonth, req.StartMonth)

	var total int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		fmt.Println(req)
		return 0, fmt.Errorf("calculate total: %w", err)
	}

	return total, nil
}

func (r *Repository) CloseConnection() {
	r.pool.Close()
}
