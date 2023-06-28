package lib

import (
	"net/http"
	"time"
)

type QLogger struct {
	database Database
}

func NewQLogger(database Database) *QLogger {
	return &QLogger{database}
}

// LogRoute middleware logs the request to the database.
func (q QLogger) LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q.database.Write(r.Context(), LogSchema{
			TimeWritten: time.Now().UTC(),
			Message:     "Request to " + r.URL.Path + " from " + r.RemoteAddr,
			Severity:    0,
			Category:    "api",
		})

		next.ServeHTTP(w, r)
	})
}
