package routers

import (
	"log/slog"
	"net/http"

	"github.com/ClaytonMatos84/go-superusers/internal/service"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/health", healthCheck).Methods("GET")

	mux.HandleFunc("/log-users", service.UploadLogs).Methods("POST")

	mux.HandleFunc("/log-users", service.GetLogs).Methods("GET")
	mux.HandleFunc("/superusers", service.GetSuperUsers).Methods("GET")
	mux.HandleFunc("/top-countries", service.GetTopCountries).Methods("GET")
	mux.HandleFunc("/team-insights", service.GetTeamInsights).Methods("GET")
	mux.HandleFunc("/active-users-per-day", service.GetLoginsPerDay).Methods("GET")

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info("Health check endpoint hit")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
