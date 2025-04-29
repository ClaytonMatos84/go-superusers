package main

import (
	"log"
	"net/http"

	"github.com/ClaytonMatos84/go-superusers/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/health", healthCheck).Methods("GET")

	mux.HandleFunc("/log-users", service.UploadLogs).Methods("POST")

	mux.HandleFunc("/log-users", service.GetLogs).Methods("GET")
	mux.HandleFunc("/superusers", service.GetSuperUsers).Methods("GET")
	mux.HandleFunc("/top-countries", service.GetTopCountries).Methods("GET")
	mux.HandleFunc("/team-insights", service.GetTeamInsights).Methods("GET")
	mux.HandleFunc("/active-users-per-day", service.GetLoginsPerDay).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
