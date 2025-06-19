package dbhelper

import (
	db "TodoApp/database"
	"TodoApp/models"
	"fmt"
	"time"
)

// CreateUser inserts a new user if username or email does not already exist.
func CreateUser(user models.User) error {
	res, err := db.DB.Exec(`
		INSERT INTO users (username, email, password)
		SELECT $1, $2, $3
		WHERE NOT EXISTS (
			SELECT 1 FROM users 
			WHERE (username = $1 OR email = $2) AND archived_at IS NULL
		)
	`, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("username or email already exists")
	}

	return nil
}

// GetUserByEmail returns a user record by email (if not archived).
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.DB.Get(&user, `
		SELECT id, password
		FROM users
		WHERE email = $1 AND archived_at IS NULL
	`, email)
	return user, err
}

// GetUserBySessionKey returns the user ID associated with a valid session.
func GetUserBySessionKey(sessionID string) (int, error) {
	var userID int
	err := db.DB.Get(&userID, `
		SELECT user_id FROM sessions
		WHERE id = $1 
		AND expires_at > $2 
		AND archived_at IS NULL
	`, sessionID, time.Now())
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// CheckIfExist verifies an active session exists for the given user ID.
func CheckIfExist(userID int) (string, error) {
	var sessionID string
	err := db.DB.Get(&sessionID, `
		SELECT id FROM sessions
		WHERE user_id = $1 
		AND expires_at > NOW() 
		AND archived_at IS NULL
	`, userID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

// GetUserByID returns user details by user ID.
func GetUserByID(id int) (models.User, error) {
	var user models.User
	err := db.DB.Get(&user, `
		SELECT id, username, email, created_at 
		FROM users 
		WHERE id = $1
	`, id)
	return user, err
}

// DeleteUserByID soft deletes a user by setting archived_at timestamp.
func DeleteUserByID(userID int) error {
	_, err := db.DB.Exec(`
		UPDATE users 
		SET archived_at = $1 
		WHERE id = $2
	`, time.Now(), userID)
	return err
}
