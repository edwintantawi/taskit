package entity

import "time"

type ProjectID string

type Project struct {
	ID        ProjectID
	UserID    UserID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
