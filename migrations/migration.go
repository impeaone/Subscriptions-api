package migrations

import (
	"agrigation_api/internal/database/postgres"
	"context"
)

// CheckAndCreateTables - создание таблиц если их нет
func CheckAndCreateTables() error {
	// Проверяем, существует ли таблица subscriptions
	db, err := postgres.InitPGPool()
	if err != nil {
		return err
	}
	defer db.Close()

	var tableExists bool
	errSearch := db.QueryRow(context.Background(), `
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'subscriptions'
        )
    `).Scan(&tableExists)

	if errSearch != nil {
		return err
	}

	if !tableExists {
		_, err := db.Exec(context.Background(), `
            CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
            CREATE TABLE subscriptions (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                user_id UUID NOT NULL,
                service_name VARCHAR(100) NOT NULL,
                price INTEGER NOT NULL CHECK (price > 0),
                start_date DATE NOT NULL,
                end_date DATE,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            );
            
            CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
            CREATE INDEX idx_subscriptions_service ON subscriptions(service_name);
        `)

		if err != nil {
			return err
		}

	}
	return nil
}
