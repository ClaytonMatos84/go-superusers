package main

import (
	"log"
	"net/http"

	"github.com/ClaytonMatos84/go-superusers/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/health", HealthCheck).Methods("GET")
	mux.HandleFunc("/users", internal.UploadUsers).Methods("POST")
	mux.HandleFunc("/users", internal.GetUsers).Methods("GET")
	mux.HandleFunc("/superusers", internal.GetSuperUsers).Methods("GET")
	mux.HandleFunc("/top-countries", internal.GetTopCountries).Methods("GET")
	mux.HandleFunc("/team-insights", internal.GetTeamInsights).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
