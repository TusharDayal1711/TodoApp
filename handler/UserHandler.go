package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/middleware"
	"TodoApp/models"
	"TodoApp/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	w.WriteHeader(http.StatusOK) //sending status code 200
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully...",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	user, err := dbhelper.GetUserByEmail(input.Email)
	fmt.Println(user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		http.Error(w, "invalid email or password....", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "could not create access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "could not create refresh token", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour) //refresh last for 7 days
	err2 := dbhelper.StoreRefreshTokenToDB(user.ID, refreshToken, expiresAt)
	if err2 != nil {
		http.Error(w, "Could not create refresh token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Login successful",
		"acess_token":   accessToken,
		"refresh_token": refreshToken,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.AuthUserFromMiddleWare(r)
	if err != nil {
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}

	err = dbhelper.DeleteRecord(userID)
	if err != nil {
		http.Error(w, "failed to delete token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
