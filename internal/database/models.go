// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"
)

type Chat struct {
	ID           int64
	CreationDate time.Time
}

type ChatMember struct {
	ChatID   sql.NullInt64
	UserID   sql.NullInt64
	JoinDate time.Time
}

type Message struct {
	ID       int64
	SenderID interface{}
	ChatID   sql.NullInt64
	Content  string
	SentAt   time.Time
}

type User struct {
	ID        int64
	CreatedAt time.Time
	Password  string
	Username  string
}
