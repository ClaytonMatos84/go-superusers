package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/ClaytonMatos84/go-superusers/internal/routers"
	"github.com/gorilla/handlers"
)

func main() {
	logger := slog.Default()
	mux := routers.GetRouter()

	logger.Info("Server started", slog.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(mux)); err != nil {
		slog.Error("Failed to start server", slog.String("address", ":8080"), slog.String("error", err.Error()))
		log.Fatalf("Failed to start server: %v", err)
	}
}
