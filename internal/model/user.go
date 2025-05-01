package model

import (
	"net/http"

	"github.com/ClaytonMatos84/go-superusers/pkg"
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
	Date   pkg.CustomDate `json:"data"`
	Action Action         `json:"acao"`
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

func ValidateUserList(users []User, w *http.ResponseWriter) bool {
	if len(users) == 0 {
		http.Error(*w, "No users found", http.StatusNotFound)
		return true
	}
	return false
}

func FindSuperUsers(users []User) []User {
	superUsers := make([]User, 0)
	for _, user := range users {
		if user.Score >= 950 && user.Active {
			superUsers = append(superUsers, user)
		}
	}

	return superUsers
}
