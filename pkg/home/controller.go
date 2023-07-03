package home

import (
	"net/http"

	"github.com/akifkadioglu/askida-kod/models"
	"github.com/akifkadioglu/askida-kod/utils"
	"github.com/go-chi/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	model := models.JwtModel{
		ID:      "id",
		Email:   "email",
		Name:    "name",
		Picture: nil,
		Time:    "time",
	}

	tokenAsString, _ := utils.GenerateToken(model)

	render.JSON(w, r, map[string]string{
		"token": tokenAsString,
	})
}
