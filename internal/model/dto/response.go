package dto

import (
	"github.com/ClaytonMatos84/go-superusers/internal/model"
	"github.com/ClaytonMatos84/go-superusers/pkg"
)

type ResponseBody struct {
	Timestamp     string `json:"timestamp"`
	ExecutionTime int64  `json:"execution_time_ms"`
	Message       string `json:"message"`
}

type ResponseUploadUsers struct {
	ResponseBody
	Count int `json:"user_count"`
}

type ResponseUsers struct {
	ResponseBody
	Count      int                `json:"user_count"`
	Pagination pkg.PaginationData `json:"pagination"`
	Data       []model.User       `json:"data"`
}

type CountCountry struct {
	Country string `json:"country"`
	Count   int    `json:"total"`
}

type ResponseTopCountries struct {
	ResponseBody
	Countries []CountCountry `json:"countries"`
}

type ResponseTeams struct {
	Team              string  `json:"team"`
	Count             int     `json:"total_members"`
	Leaders           int     `json:"leaders"`
	CompletedProjects int     `json:"completed_projects"`
	ActivePercent     float64 `json:"active_percentage"`
}
type ResponseTeam struct {
	ResponseBody
	Teams []ResponseTeams `json:"teams"`
}

type CountLogins struct {
	Date  pkg.CustomDate `json:"date"`
	Count int            `json:"total"`
}

type ResponseLogins struct {
	ResponseBody
	Logins []CountLogins `json:"logins"`
}
