package handlers

import (
	"github.com/mr-time2028/WebChat/internal/helpers"
	"net/http"
	"testing"
)

func TestHandlerRepository_Register(t *testing.T) {
	var testCases = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
		expectedErr        bool
	}{
		{
			"valid data",
			`{
				"username": "david15",
				"password": "DavidPass1234",
				"confirm_password": "DavidPass1234"
			}`,
			http.StatusOK,
			false,
		},
		{
			"duplicate username",
			`{
				"username": "defaultUser",
				"password": "defaultPass",
				"confirm_password": "defaultPass"
			}`,
			http.StatusBadRequest,
			true,
		},
	}

	for _, e := range testCases {
		rr := helpers.FakePostRequest(e.requestBody, "/register", http.HandlerFunc(HandlerRepo.Register))
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d, %v", e.name, e.expectedStatusCode, rr.Code, rr.Body)
		}
	}
}

func TestHandlerRepository_Login(t *testing.T) {
	var testCases = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{
			"valid data",
			`{"username": "defaultUser", "password": "defaultPass"}`,
			http.StatusOK,
		},
		{
			"no user row",
			`{"username": "NotExistsUser", "password": "NotExistsUserPass"}`,
			http.StatusUnauthorized,
		},
		{
			"wrong password",
			`{"username": "defaultUser", "password": "WrongDefaultPass"}`,
			http.StatusUnauthorized,
		},
	}

	for _, e := range testCases {
		rr := helpers.FakePostRequest(e.requestBody, "/login", http.HandlerFunc(HandlerRepo.Login))
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}
