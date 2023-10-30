package entity

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// nullBytes represent the bytes for null.
var nullBytes = []byte("null")

// NullTime that may be null. NullTime embed sql.NullTime and implement json Unmarshaler and Marshaler
type NullTime struct {
	sql.NullTime
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		t.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &t.Time); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullBytes, nil
	}
	return json.Marshal(t.Time)
}

type NullString struct {
	sql.NullString
}

func (t *NullString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		t.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &t.String); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

func (t NullString) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullBytes, nil
	}
	return json.Marshal(t.String)
}
