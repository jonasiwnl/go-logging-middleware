package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	m "github.com/jonasiwnl/go-logging-middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("starting...")

	// Parse .env
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil || env["DSI"] == "" {
		log.Fatal(err)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.
		Client().
		ApplyURI(env["DSI"]).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	// Ping database to confirm connection.
	err = client.
		Database("LoggingMiddleware").
		RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).
		Err()

	if err != nil {
		log.Fatal(err)
	}

	logger := m.NewLoggingMiddlewareBuilder(
		m.NewMongoDatabase(client.Database("LoggingMiddleware").Collection("logs"))
	).WithInfoLevel(m.Verbose).Build()

	http.Handle("/", logger.LogRoute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
