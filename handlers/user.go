package handlers

import (
	"errors"
	"github.com/mr-time2028/WebChat/helpers"
	"github.com/mr-time2028/WebChat/models"
	"github.com/mr-time2028/WebChat/validators"
	"log"
	"net/http"
)

func (h *HandlerRepository) Register(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username        string `json:"username" required:"true" min:"5" max:"255"`
		Password        string `json:"password" required:"true" min:"8" max:"30"`
		ConfirmPassword string `json:"confirm_password" required:"true" min:"8" max:"30"`
	}

	var responseBody struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	// get user data and json validation
	if validator := helpers.ReadJSON(w, r, &requestBody); !validator.Valid() {
		if err := helpers.ErrorMapJSON(w, validator.Errors); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// validate password
	validator := validators.New()
	password1 := requestBody.Password
	password2 := requestBody.ConfirmPassword

	validator.UserPasswordValidation(password1, password2)
	if !validator.Valid() {
		if err := helpers.ErrorMapJSON(w, validator.Errors); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// check if user with this information already exists
	isExistsUser, err := h.App.Models.User.CheckIfExistsUser(requestBody.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if isExistsUser {
		if err = helpers.ErrorStrJSON(w, errors.New("user with this username already exists")); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// insert user to the database
	user := &models.User{
		Username: requestBody.Username,
		Password: password1,
		IsActive: true, // TODO: account confirmation
	}

	_, err = h.App.Models.User.InsertOneUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// write tokens to the output
	responseBody.Error = false
	responseBody.Message = "registration was successful"
	if err = helpers.WriteJSON(w, http.StatusOK, responseBody); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
