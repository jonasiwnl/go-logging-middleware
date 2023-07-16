package lib

import (
	"context"
	"time"
)

/*
 * Database interface for writing logs.
 */
type Database interface {
	Write(context.Context, LogSchema) error
}

/*
 * Schema for a single log entry.
 */
type LogSchema struct {
	ID          string    `bson:"_id"`
	TimeWritten time.Time `bson:"time_written"`
	Category    string    `bson:"category"`
	Info        string    `bson:"info"`
}

/*
 * Noisiness level for the logger.
 */
type InfoLevel int

const (
	// Just logs route hits.
	Minimal InfoLevel = iota

	// Logs body.
	Normal

	// Logs headers, body, and other information.
	Verbose
)
