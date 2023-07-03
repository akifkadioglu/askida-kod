package test

import (
	"net/http"
	"testing"

	"github.com/akifkadioglu/askida-kod/pkg/auth"
	"github.com/akifkadioglu/askida-kod/route"
)

func TestAuth(t *testing.T) {
	setupTest()
	s := route.CreateNewServer()
	s.MountHandlers()

	t.Run("Register", func(t *testing.T) {
		body := setBodyForTest(auth.BodyRegister{
			Email:    "akifkadioglu@gmail.com",
			Name:     "Akif",
			Password: "abc",
		})
		req, _ := http.NewRequest("POST", "/api/public/register", body)
		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusOK, response.Code)
	})

	t.Run("Login", func(t *testing.T) {
		body := setBodyForTest(auth.BodyLogin{
			Email:    "akifkadioglu@gmail.com",
			Password: "abc",
		})
		req, _ := http.NewRequest("POST", "/api/public/login", body)
		response := executeRequest(req, s)
		checkResponseCode(t, http.StatusOK, response.Code)
	})
}
