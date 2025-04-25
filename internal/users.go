package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"runtime"
	"sort"
	"strconv"
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

type countCountry struct {
	Country string `json:"country"`
	Count   int    `json:"total"`
}

// responses
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

type responseTopCountries struct {
	responseBody
	Countries []countCountry `json:"countries"`
}

type responseTeams struct {
	Team              string  `json:"team"`
	Count             int     `json:"total_members"`
	Leaders           int     `json:"leaders"`
	CompletedProjects int     `json:"completed_projects"`
	ActivePercent     float64 `json:"active_percentage"`
}
type responseTeam struct {
	responseBody
	Teams []responseTeams `json:"teams"`
}

type countLogins struct {
	Date  string `json:"date"`
	Count int    `json:"total"`
}

type responseLogins struct {
	responseBody
	Logins []countLogins `json:"logins"`
}

var users []User

func UploadUsers(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := initCheck()

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
	milliseconds, info := finishCheck(memStatus, start_time)
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

func initCheck() (runtime.MemStats, time.Time) {
	var memStatus runtime.MemStats
	start_time := time.Now()

	return memStatus, start_time
}

func finishCheck(memStatus runtime.MemStats, start_time time.Time) (int64, string) {
	runtime.ReadMemStats(&memStatus)
	duration := time.Since(start_time)

	milliseconds := duration.Milliseconds()
	info := fmt.Sprintf("Elapsed time = %vms. Total memory(KB) consumed = %v", milliseconds, memStatus.Sys/1024)

	return milliseconds, info
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := initCheck()

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	milliseconds, info := finishCheck(memStatus, start_time)
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
	memStatus, start_time := initCheck()

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	superUsers := findSuperUsers()
	milliseconds, info := finishCheck(memStatus, start_time)
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

func findSuperUsers() []User {
	superUsers := make([]User, 0)
	for _, user := range users {
		if user.Score >= 900 && user.Active {
			superUsers = append(superUsers, user)
		}
	}

	return superUsers
}

func GetTopCountries(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := initCheck()

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	superUsers := findSuperUsers()
	countries := make(map[string]int)
	for _, superUser := range superUsers {
		countries[superUser.Country]++
	}

	countriesList := make([]countCountry, 0)
	for country, count := range countries {
		countriesList = append(countriesList, countCountry{
			Country: country,
			Count:   count,
		})
	}

	sort.Slice(countriesList, func(i, j int) bool {
		return countriesList[i].Count > countriesList[j].Count
	})
	if len(countriesList) > 5 {
		countriesList = countriesList[:5]
	}

	milliseconds, info := finishCheck(memStatus, start_time)
	response := responseTopCountries{
		responseBody: responseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Countries: countriesList,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode top countries", http.StatusInternalServerError)
		return
	}
}

func GetTeamInsights(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := initCheck()

	emptyUserList := validateUserList(w)
	if emptyUserList {
		return
	}

	teamMap := make(map[string]responseTeams)
	for _, user := range users {
		if user.Team.Name == "" {
			continue
		}

		teamName := user.Team.Name
		if _, exists := teamMap[teamName]; !exists {
			teamMap[teamName] = responseTeams{
				Team:              teamName,
				Count:             0,
				Leaders:           0,
				CompletedProjects: 0,
				ActivePercent:     0,
			}
		}

		team := teamMap[teamName]
		team.Count++

		if user.Team.Leader {
			team.Leaders++
		}

		for _, project := range user.Team.Projects {
			if project.Concluded {
				team.CompletedProjects++
			}
		}

		teamMap[teamName] = team
	}

	teamsList := make([]responseTeams, 0)
	for _, team := range teamMap {
		activeCount := 0
		for _, user := range users {
			if user.Team.Name == team.Team && user.Active {
				activeCount++
			}
		}
		if team.Count > 0 {
			team.ActivePercent = roundFloat(float64(activeCount)/float64(team.Count)*100, 2)
		}
		teamsList = append(teamsList, team)
	}

	milliseconds, info := finishCheck(memStatus, start_time)
	response := responseTeam{
		responseBody: responseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Teams: teamsList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode team insights", http.StatusInternalServerError)
		return
	}
}

func validateUserList(w http.ResponseWriter) bool {
	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return true
	}
	return false
}

func roundFloat(num float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(num*p) / p
}

func GetLoginsPerDay(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := initCheck()

	emptyUserList := validateUserList(w)
	if emptyUserList {
		return
	}

	loginsPerDay := make(map[string]int)
	for _, user := range users {
		for _, log := range user.Logs {
			if log.Action == "login" {
				loginsPerDay[log.Date]++
			}
		}
	}

	loginsList := make([]countLogins, 0)
	for date, count := range loginsPerDay {
		loginsList = append(loginsList, countLogins{
			Date:  date,
			Count: count,
		})
	}

	query := r.URL.Query()
	queryMin := query.Get("min")
	if queryMin != "" {
		filteredLoginsList := make([]countLogins, 0)
		for _, login := range loginsList {
			if minValue, err := strconv.ParseInt(queryMin, 10, 64); err == nil && int64(login.Count) > minValue {
				filteredLoginsList = append(filteredLoginsList, login)
			}
		}
		loginsList = filteredLoginsList
	}

	milliseconds, info := finishCheck(memStatus, start_time)
	response := responseLogins{
		responseBody: responseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Logins: loginsList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode logins per day", http.StatusInternalServerError)
		return
	}
}
