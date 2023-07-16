package lib

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoggingMiddleware struct {
	database  Database
	infoLevel InfoLevel
	// Prints logs to stdout.
	debug bool
}

type LoggingMiddlewareBuilder struct {
	database  Database
	infoLevel InfoLevel
	// Prints logs to stdout.
	debug bool
}

func (b LoggingMiddlewareBuilder) Build() *LoggingMiddleware {
	return &LoggingMiddleware{b.database, b.infoLevel, b.debug}
}

func (b LoggingMiddlewareBuilder) WithInfoLevel(infoLevel InfoLevel) *LoggingMiddlewareBuilder {
	b.infoLevel = infoLevel
	return &b
}

func (b LoggingMiddlewareBuilder) WithDebug(debug bool) *LoggingMiddlewareBuilder {
	b.debug = debug
	return &b
}

func NewLoggingMiddlewareBuilder(database Database) *LoggingMiddlewareBuilder {
	return &LoggingMiddlewareBuilder{database, Minimal, false}
}

// LogRoute middleware logs the request to the database.
func (q LoggingMiddleware) LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stamp := time.Now().UTC()

		info := "Request to " + r.URL.Path + " from " + r.RemoteAddr
		switch q.infoLevel {
		case Verbose:
			info += "\nBody:\n"
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				info += string(bodyBytes)
			} else {
				info += "Error reading body: " + err.Error()
			}
		case Normal:
			info += "\nHeaders:\n"
			for k, v := range r.Header {
				info += k + ": " + v[0] + "\n"
			}
		default: // Minimal, don't do anything
		}

		if q.debug {
			fmt.Println(stamp.String() + ": \n" + info)
		}

		// Truncate info to 255 characters.
		if len(info) > 255 {
			info = info[:255]
		}

		q.database.Write(r.Context(), LogSchema{
			ID:          uuid.NewString(),
			TimeWritten: stamp,
			Category:    "api",
			Info:        info,
		})

		next.ServeHTTP(w, r)
	})
}
