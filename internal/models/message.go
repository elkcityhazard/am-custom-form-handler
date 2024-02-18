package models

import "time"

type Message struct {
	ID        int64
	UserID    int64
	Email     string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}
