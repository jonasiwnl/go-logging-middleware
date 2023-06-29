# go-logging-middleware
i made this in 20 minutes on accident

# usage
```go
package main

import (
    "fmt"
    "net/http"

    glm "github.com/jonasiwnl/go-logging-middleware"
)

func main() {
    logger := glm.NewLoggingMiddleware(glm.NewMongoDatabase("YOUR MONGO COLLECTION"))

    http.Handle("/", logger.LogRoute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, world!")
    })))
    http.ListenAndServe(":8080", nil)
}
```

## TODO

- [x] IDs for logs into db
- [ ] Add more database drivers
