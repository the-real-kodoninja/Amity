package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Client
}

func (a *App) Initialize() error {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Update with your MongoDB URI
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	a.DB = client
	a.Router = mux.NewRouter()
	a.initializeRoutes()

	return nil
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/health", a.healthCheck).Methods("GET")
}

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Amity Backend is running!")
}

func (a *App) Run(addr string) {
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func main() {
	app := &App{}
	err := app.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(":8080")
}