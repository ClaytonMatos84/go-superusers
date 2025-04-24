package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

type project struct {
	Name      string `json:"nome"`
	Concluded bool   `json:"concluido"`
}

type team struct {
	Name     string    `json:"nome"`
	Leader   bool      `json:"lider"`
	Projects []project `json:"projetos"`
}

type logData struct {
	Date   string `json:"data"`
	Action string `json:"acao"`
}

type User struct {
	ID      string    `json:"id"`
	Name    string    `json:"nome"`
	Age     int       `json:"idade"`
	Score   int       `json:"score"`
	Active  bool      `json:"ativo"`
	Country string    `json:"pais"`
	Team    team      `json:"equipe"`
	Logs    []logData `json:"logs"`
}

var users []User

func UploadUsers(w http.ResponseWriter, r *http.Request) {
	var memStatus runtime.MemStats
	start_time := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println("Error getting file from form: ", err)
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if file == nil {
		http.Error(w, "File is empty", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	log.Println("Received users: ", len(users))
	runtime.ReadMemStats(&memStatus)
	duration := time.Since(start_time)
	info := fmt.Sprintf("Elapsed time = %s. Total memory(KB) consumed = %v", duration, memStatus.Sys/1024)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(info))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}
