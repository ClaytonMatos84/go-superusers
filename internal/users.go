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

type responseBody struct {
	Timestamp     string `json:"timestamp"`
	ExecutionTime int64  `json:"execution_time_ms"`
	Message       string `json:"message"`
}

type responseUploadUsers struct {
	responseBody
	Count int `json:"user_count"`
}

type responseUsers struct {
	responseBody
	Count int    `json:"user_count"`
	Data  []User `json:"data"`
}

var users []User

func UploadUsers(w http.ResponseWriter, r *http.Request) {
	var memStatus runtime.MemStats
	start_time := time.Now()

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
	milliseconds := duration.Milliseconds()
	info := fmt.Sprintf("Elapsed time = %vms. Total memory(KB) consumed = %v", milliseconds, memStatus.Sys/1024)

	response := responseUploadUsers{
		responseBody: responseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Count: len(users),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var memStatus runtime.MemStats
	start_time := time.Now()

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	runtime.ReadMemStats(&memStatus)
	duration := time.Since(start_time)
	milliseconds := duration.Milliseconds()
	info := fmt.Sprintf("Elapsed time = %vms. Total memory(KB) consumed = %v", milliseconds, memStatus.Sys/1024)

	response := responseUsers{
		responseBody: responseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Count: len(users),
		Data:  users,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}

func GetSuperUsers(w http.ResponseWriter, r *http.Request) {
	var memStatus runtime.MemStats
	start_time := time.Now()

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	superUsers := make([]User, 0)
	for _, user := range users {
		if user.Score >= 900 && user.Active {
			superUsers = append(superUsers, user)
		}
	}

	runtime.ReadMemStats(&memStatus)
	duration := time.Since(start_time)
	milliseconds := duration.Milliseconds()
	info := fmt.Sprintf("Elapsed time = %vms. Total memory(KB) consumed = %v", milliseconds, memStatus.Sys/1024)
	response := responseUsers{
		responseBody: responseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Count: len(superUsers),
		Data:  superUsers,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode super users", http.StatusInternalServerError)
		return
	}
}
