package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // we dont send back password
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// not null means deleted, * type as time cannot be nil, go will treat it as nil pointer
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type Todo struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      bool      `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	// not null means deleted, * type as time cannot be nil, go will treat it as nil pointer
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"` // for soft deleting, not null = deleted
}

type Session struct {
	ID         string     `json:"id" db:"id"` //radomly generated unique id, stored as string
	UserID     int        `json:"user_id" db:"user_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"` //omit this field if null or empty
}
