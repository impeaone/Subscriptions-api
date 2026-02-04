package postgres

import (
	"agrigation_api/pkg/tools"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func InitPostgres() (*Postgres, error) {
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

	if err := createTables(pool); err != nil {
		return nil, err
	}

	return &Postgres{
		pool: pool,
	}, nil

}
func createTables(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS users (
    		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    		email VARCHAR(255) UNIQUE NOT NULL,
    		name VARCHAR(100),
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS subscriptions (
    		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    		service_name VARCHAR(100) NOT NULL,
    		plan_name VARCHAR(50),
    		price_rub INTEGER NOT NULL CHECK (price_rub >= 0),
    		periodicity VARCHAR(20) DEFAULT 'monthly',
    		status VARCHAR(20) DEFAULT 'active',
    		start_date DATE DEFAULT CURRENT_DATE,
   			next_billing_date DATE,
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS payments (
    		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    		subscription_id UUID REFERENCES subscriptions(id),
    		user_id UUID REFERENCES users(id),
    		amount_rub INTEGER NOT NULL CHECK (amount_rub >= 0),
    		paid_at DATE DEFAULT CURRENT_DATE,
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
