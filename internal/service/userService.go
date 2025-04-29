package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ClaytonMatos84/go-superusers/internal/model"
	"github.com/ClaytonMatos84/go-superusers/internal/model/dto"
	"github.com/ClaytonMatos84/go-superusers/pkg"
)

var users []model.User

func UploadLogs(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := pkg.InitControlRequest()

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
	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseUploadUsers{
		ResponseBody: dto.ResponseBody{
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

func GetLogs(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := pkg.InitControlRequest()

	emptyList := model.ValidateUserList(users, &w)
	if emptyList {
		return
	}

	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseUsers{
		ResponseBody: dto.ResponseBody{
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
	memStatus, start_time := pkg.InitControlRequest()

	emptyList := model.ValidateUserList(users, &w)
	if emptyList {
		return
	}

	superUsers := model.FindSuperUsers(users)
	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseUsers{
		ResponseBody: dto.ResponseBody{
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

func GetTopCountries(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := pkg.InitControlRequest()

	emptyList := model.ValidateUserList(users, &w)
	if emptyList {
		return
	}

	superUsers := model.FindSuperUsers(users)
	countries := make(map[string]int)
	for _, superUser := range superUsers {
		countries[superUser.Country]++
	}

	countriesList := make([]dto.CountCountry, 0)
	for country, count := range countries {
		countriesList = append(countriesList, dto.CountCountry{
			Country: country,
			Count:   count,
		})
	}

	query := r.URL.Query()
	querySize := query.Get("size")
	countriesList = orderCountriesList(countriesList, querySize)

	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseTopCountries{
		ResponseBody: dto.ResponseBody{
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

func orderCountriesList(countriesList []dto.CountCountry, querySize string) []dto.CountCountry {
	responseSize := 5
	if querySize != "" {
		responseSize, _ = strconv.Atoi(querySize)
	}

	sort.Slice(countriesList, func(i, j int) bool {
		return countriesList[i].Count > countriesList[j].Count
	})

	if len(countriesList) > int(responseSize) {
		countriesList = countriesList[:responseSize]
	}

	return countriesList
}

func GetTeamInsights(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := pkg.InitControlRequest()

	emptyUserList := model.ValidateUserList(users, &w)
	if emptyUserList {
		return
	}

	teamMap := make(map[string]dto.ResponseTeams)
	for _, user := range users {
		if user.Team.Name == "" {
			continue
		}

		teamName := user.Team.Name
		if _, exists := teamMap[teamName]; !exists {
			teamMap[teamName] = dto.ResponseTeams{
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

	teamsList := make([]dto.ResponseTeams, 0)
	for _, team := range teamMap {
		activeCount := 0
		for _, user := range users {
			if user.Team.Name == team.Team && user.Active {
				activeCount++
			}
		}
		if team.Count > 0 {
			team.ActivePercent = pkg.RoundFloat(float64(activeCount)/float64(team.Count)*100, 2)
		}
		teamsList = append(teamsList, team)
	}

	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseTeam{
		ResponseBody: dto.ResponseBody{
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

func GetLoginsPerDay(w http.ResponseWriter, r *http.Request) {
	memStatus, start_time := pkg.InitControlRequest()

	emptyUserList := model.ValidateUserList(users, &w)
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

	loginsList := make([]dto.CountLogins, 0)
	for date, count := range loginsPerDay {
		loginsList = append(loginsList, dto.CountLogins{
			Date:  date,
			Count: count,
		})
	}

	query := r.URL.Query()
	queryMin := query.Get("min")
	if queryMin != "" {
		filteredLoginsList := make([]dto.CountLogins, 0)
		for _, login := range loginsList {
			if minValue, err := strconv.ParseInt(queryMin, 10, 64); err == nil && int64(login.Count) > minValue {
				filteredLoginsList = append(filteredLoginsList, login)
			}
		}
		loginsList = filteredLoginsList
	}

	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseLogins{
		ResponseBody: dto.ResponseBody{
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
