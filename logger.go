package lib

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoggingMiddleware struct {
	database Database
	print    bool
}

func NewLoggingMiddleware(database Database, print ...bool) *LoggingMiddleware {
	if len(print) == 0 {
		return &LoggingMiddleware{database, false}
	}
	return &LoggingMiddleware{database, print[0]}
}

// LogRoute middleware logs the request to the database.
func (q LoggingMiddleware) LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stamp := time.Now().UTC()

		if q.print {
			fmt.Println(stamp.String() + ": Request to " + r.URL.Path + " from " + r.RemoteAddr)
		}

		q.database.Write(r.Context(), LogSchema{
			ID:          uuid.NewString(),
			TimeWritten: stamp,
			Message:     "Request to " + r.URL.Path + " from " + r.RemoteAddr,
			Severity:    0,
			Category:    "api",
		})

		next.ServeHTTP(w, r)
	})
}
