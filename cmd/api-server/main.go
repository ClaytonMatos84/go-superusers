package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/ClaytonMatos84/go-superusers/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	logger := slog.Default()
	mux := mux.NewRouter()
	mux.HandleFunc("/health", healthCheck).Methods("GET")

	mux.HandleFunc("/log-users", service.UploadLogs).Methods("POST")

	mux.HandleFunc("/log-users", service.GetLogs).Methods("GET")
	mux.HandleFunc("/superusers", service.GetSuperUsers).Methods("GET")
	mux.HandleFunc("/top-countries", service.GetTopCountries).Methods("GET")
	mux.HandleFunc("/team-insights", service.GetTeamInsights).Methods("GET")
	mux.HandleFunc("/active-users-per-day", service.GetLoginsPerDay).Methods("GET")

	logger.Info("Server started", slog.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(mux)); err != nil {
		slog.Error("Failed to start server", slog.String("address", ":8080"), slog.String("error", err.Error()))
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info("Health check endpoint hit")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
