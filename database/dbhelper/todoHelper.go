package dbhelper

import (
	db "TodoApp/database"
	"TodoApp/models"
	"fmt"
	"time"
)

// CreateTodo inserts a new todo record into the database.
func CreateTodo(todoObj models.Todo) error {
	_, err := db.DB.Exec(`
		INSERT INTO todos (user_id, title, description) 
		VALUES ($1, $2, $3)
	`, todoObj.UserID, todoObj.Title, todoObj.Description)
	return err
}

// UpdateTodoDetails updates a user's todo entry with new details.
func UpdateTodoDetails(todo models.Todo) error {
	res, err := db.DB.Exec(`
		UPDATE todos
		SET title = $1, description = $2, status = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5 AND archived_at IS NULL
	`, todo.Title, todo.Description, todo.Status, todo.ID, todo.UserID)

	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("invalid todo id")
	}

	return nil
}

// MarkCompleted sets the status of a todo to true (completed).
func MarkCompleted(todoID int, userID int) error {
	_, err := db.DB.Exec(`
		UPDATE todos 
		SET status = TRUE, updated_at = $3
		WHERE id = $1 AND user_id = $2
	`, todoID, userID, time.Now())
	return err
}

// DeleteTodoByID soft deletes a todo item by setting archived_at.
func DeleteTodoByID(userID int, itemID int) error {
	_, err := db.DB.Exec(`
		UPDATE todos
		SET archived_at = $1
		WHERE user_id = $2 AND id = $3 AND archived_at IS NULL
	`, time.Now(), userID, itemID)
	return err
}

func GetTodosStatus(userID int, status string) ([]models.Todo, error) {
	var todos []models.Todo
	var err error
	// completed means true, incomplete means false else display all as default
	switch status {
	case "completed":
		err = db.DB.Select(&todos, `
			SELECT id, title, description, status, created_at, updated_at
			FROM todos
			WHERE user_id = $1 AND status = true AND archived_at IS NULL
			ORDER BY updated_at DESC
		`, userID)

	case "incompleted":
		err = db.DB.Select(&todos, `
			SELECT id, title, description, status, created_at, updated_at
			FROM todos
			WHERE user_id = $1 AND status = false AND archived_at IS NULL
			ORDER BY updated_at DESC
		`, userID)

	default:
		err = db.DB.Select(&todos, `
			SELECT id, title, description, status, created_at, updated_at
			FROM todos
			WHERE user_id = $1 AND archived_at IS NULL
			ORDER BY updated_at DESC
		`, userID)
	}

	return todos, err
}

////GetTodosByID fetches all non-archived todos for a user.
// func GetTodosByID(userID int) ([]models.Todo, error) {
// 	var todos []models.Todo
// 	err := db.DB.Select(&todos, `
// 		SELECT id, user_id, title, description, status, created_at, updated_at
// 		FROM todos
// 		WHERE user_id = $1 AND archived_at IS NULL
// 		ORDER BY created_at DESC
// 	`, userID)
// 	return todos, err
// }
