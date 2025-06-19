package dbhelper

import (
	db "TodoApp/database"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// for creating new session
func CreateSession(userID int) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err := db.DB.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, sessionID, userID, expiresAt)

	if err != nil {
		return "", err
	}
	return sessionID, nil
} //CreateSession

// find and archive session / soft delete
func RemoveSession(sessionID string) error {
	res, err := db.DB.Exec(`
		UPDATE sessions 
		SET archived_at = $1, expires_at = $2 
		WHERE id = $3 AND archived_at IS NULL
	`, time.Now(), time.Now(), sessionID)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("custom_error_id")
	}
	return nil
}
