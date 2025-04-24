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
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	mux.HandleFunc("/users", internal.UploadUsers).Methods("POST")
	mux.HandleFunc("/users", internal.GetUsers).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
