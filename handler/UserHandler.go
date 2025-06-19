package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/models"
	"TodoApp/utils"
	"errors"
	"strings"

	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Problem in hashing password", http.StatusInternalServerError)
		return
	}
	newUser := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	if err := dbhelper.CreateUser(newUser); err != nil {
		if strings.Contains(err.Error(), "username") || strings.Contains(err.Error(), "email") == true {
			http.Error(w, "Username or email already exist", http.StatusConflict)
			return
		} else {
			http.Error(w, errors.New("unable to process").Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully...",
	})
}

// login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}
	user, err := dbhelper.GetUserByEmail(input.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	isSessionExist, err := dbhelper.CheckIfExist(user.ID)
	if err != nil {
		http.Error(w, "Internal server error ", http.StatusInternalServerError)
		return
	}

	if isSessionExist != "" {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Already logged in",
			"session": isSessionExist,
		})
		return
	}

	sessionID, err := dbhelper.CreateSession(user.ID)
	if err != nil {
		http.Error(w, "could not create session", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful...",
		"session": sessionID,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("Authorization")
	userID, err := utils.AuthHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_ = userID //no use in this func

	err = dbhelper.RemoveSession(sessionID)
	if err != nil {
		if strings.Contains(err.Error(), "custom_error_id") {
			http.Error(w, "No session removed...", http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
} //
