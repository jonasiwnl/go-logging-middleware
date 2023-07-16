package lib

import (
	"context"
	"database/sql"
	"time"
)

type sqlWrapper struct {
	database *sql.DB
}

func (s sqlWrapper) Write(ctx context.Context, log LogSchema) error {
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(15*time.Second),
	)
	defer cancel()

	_, err := s.database.ExecContext(
		ctx,
		"INSERT INTO logs (id, time_written, category, info) VALUES (?, ?, ?, ?)",
		log.ID,
		log.TimeWritten,
		log.Info,
		log.Category,
	)
	return err
}

func NewSQLDatabase(database *sql.DB) Database {
	return sqlWrapper{database}
}
