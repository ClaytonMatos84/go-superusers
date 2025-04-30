package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ClaytonMatos84/go-superusers/internal/model"
	"github.com/ClaytonMatos84/go-superusers/internal/model/dto"
	"github.com/ClaytonMatos84/go-superusers/pkg"
)

var logger = slog.Default()
var users []model.User

func UploadLogs(w http.ResponseWriter, r *http.Request) {
	logger.Info("UploadLogs endpoint started")
	memStatus, start_time := pkg.InitControlRequest()

	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Error("Error getting file from form", slog.String("error", err.Error()))
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if file == nil {
		logger.Error("File is empty")
		http.Error(w, "File is empty", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		logger.Error("Error decoding JSON", slog.String("error", err.Error()))
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	logger.Info("File decoded successfully", slog.Int("user_count", len(users)))
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
		logger.Error("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
	logger.Info("UploadLogs endpoint finished successfully", slog.Int("user_count", len(users)))
}

func GetLogs(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetLogs endpoint started")
	memStatus, start_time := pkg.InitControlRequest()

	emptyList := model.ValidateUserList(users, &w)
	if emptyList {
		logger.Error("User logs list is empty")
		return
	}

	query := r.URL.Query()
	pagination, paginationError := pkg.Pagination(query.Get("page"), query.Get("items"), w, len(users))
	if paginationError {
		logger.Error("Error in pagination", slog.String("error", "Invalid pagination parameters"))
		return
	}

	userCollection := users[pagination.StartItems:pagination.EndItems]
	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseUsers{
		ResponseBody: dto.ResponseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Count:      len(userCollection),
		Pagination: pagination,
		Data:       userCollection,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
	logger.Info("GetLogs endpoint finished successfully", slog.Int("user_count", len(userCollection)))
}

func GetSuperUsers(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetSuperUsers endpoint started")
	memStatus, start_time := pkg.InitControlRequest()

	emptyList := model.ValidateUserList(users, &w)
	if emptyList {
		logger.Error("User logs list is empty")
		return
	}

	superUsers := model.FindSuperUsers(users)
	query := r.URL.Query()
	pagination, paginationError := pkg.Pagination(query.Get("page"), query.Get("items"), w, len(superUsers))
	if paginationError {
		logger.Error("Error in pagination", slog.String("error", "Invalid pagination parameters"))
		return
	}

	userCollection := superUsers[pagination.StartItems:pagination.EndItems]
	milliseconds, info := pkg.FinishControlCheck(memStatus, start_time)
	response := dto.ResponseUsers{
		ResponseBody: dto.ResponseBody{
			Timestamp:     time.Now().Format(time.RFC3339),
			ExecutionTime: milliseconds,
			Message:       info,
		},
		Pagination: pagination,
		Count:      len(userCollection),
		Data:       userCollection,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, "Failed to encode super users", http.StatusInternalServerError)
		return
	}
	logger.Info("GetSuperUsers endpoint finished successfully", slog.Int("user_count", len(superUsers)))
}

func GetTopCountries(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetTopCountries endpoint started")
	memStatus, start_time := pkg.InitControlRequest()

	emptyList := model.ValidateUserList(users, &w)
	if emptyList {
		logger.Error("User logs list is empty")
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
	var err bool
	countriesList, err = orderCountriesList(countriesList, query.Get("size"))
	if err {
		logger.Error("Error in ordering countries list", slog.String("error", "Invalid size parameter"))
		http.Error(w, "Invalid size parameter", http.StatusBadRequest)
		return
	}

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
		logger.Error("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, "Failed to encode top countries", http.StatusInternalServerError)
		return
	}
	logger.Info("GetTopCountries endpoint finished successfully", slog.Int("countries", len(countriesList)))
}

func orderCountriesList(countriesList []dto.CountCountry, querySize string) ([]dto.CountCountry, bool) {
	responseSize := 5
	var err error
	if querySize != "" {
		responseSize, err = strconv.Atoi(querySize)
		if err != nil || responseSize < 1 {
			return nil, true
		}
	}

	sort.Slice(countriesList, func(i, j int) bool {
		return countriesList[i].Count > countriesList[j].Count
	})

	if len(countriesList) > int(responseSize) {
		countriesList = countriesList[:responseSize]
	}

	return countriesList, false
}

func GetTeamInsights(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetTeamInsights endpoint started")
	memStatus, start_time := pkg.InitControlRequest()

	emptyUserList := model.ValidateUserList(users, &w)
	if emptyUserList {
		logger.Error("User logs list is empty")
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
		logger.Error("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, "Failed to encode team insights", http.StatusInternalServerError)
		return
	}
	logger.Info("GetTeamInsights endpoint finished successfully", slog.Int("teams", len(teamsList)))
}

func GetLoginsPerDay(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetLoginsPerDay endpoint started")
	memStatus, start_time := pkg.InitControlRequest()

	emptyUserList := model.ValidateUserList(users, &w)
	if emptyUserList {
		logger.Error("User logs list is empty")
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
		logger.Error("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, "Failed to encode logins per day", http.StatusInternalServerError)
		return
	}
	logger.Info("GetLoginsPerDay endpoint finished successfully", slog.Int("logins", len(loginsList)))
}
