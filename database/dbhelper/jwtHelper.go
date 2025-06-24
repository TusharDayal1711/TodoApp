package dbhelper

import (
	db "TodoApp/database"
	"github.com/google/uuid"
	"time"
)

func DeleteRecord(userID int) error {
	//delete_session_Record
	_, err := db.DB.Exec(`
		DELETE FROM sessions
		WHERE user_id = $1
	`, userID)
	return err

}

func StoreRefreshTokenToDB(userID int, token string, expiresAt time.Time) error {
	sessionID := uuid.New().String()
	_, err := db.DB.Exec(`
		INSERT INTO sessions (user_id, id, refresh_token, expires_at)
		VALUES ($1, $2, $3, $4)
	`, userID, sessionID, token, expiresAt)
	return err
}
