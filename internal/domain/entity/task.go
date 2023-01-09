package entity

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"
)

type TaskID string

// Task represents a task in the system.
type Task struct {
	ID          TaskID
	UserID      UserID
	Content     string
	Description string
	IsCompleted bool
	DueDate     TaskDueDate
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TaskDueDate is a type for due date time and handling null values.
type TaskDueDate struct {
	sql.NullTime
}

func (t *TaskDueDate) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &t.Time); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

func (t TaskDueDate) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}
