package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TaskEntityTestSuite struct {
	suite.Suite
}

func TestTaskEntitySuite(t *testing.T) {
	suite.Run(t, new(TaskEntityTestSuite))
}

func (s *TaskEntityTestSuite) TestTaskDueDateUnmarshalJSON() {
	s.Run("it should return error when invalid json", func() {
		rawJson := `[]`
		var taskDueDate TaskDueDate

		err := json.Unmarshal([]byte(rawJson), &taskDueDate)

		log.Println("HERE", err)
		s.Error(err)
		s.False(taskDueDate.Valid)
		s.Empty(taskDueDate.Time)
	})

	s.Run("it should return valid false and time is zero value", func() {
		rawJson := `null`
		var taskDueDate TaskDueDate

		err := json.Unmarshal([]byte(rawJson), &taskDueDate)

		s.NoError(err)
		s.False(taskDueDate.Valid)
		s.Empty(taskDueDate.Time)
	})

	s.Run("it should return valid false and time is zero value", func() {
		rawJson := `"2022-12-25T00:00:00.000Z"`
		var taskDueDate TaskDueDate

		err := json.Unmarshal([]byte(rawJson), &taskDueDate)

		s.NoError(err)
		s.True(taskDueDate.Valid)
		s.Equal("2022-12-25 00:00:00 +0000 UTC", taskDueDate.Time.String())
	})
}

func (s *TaskEntityTestSuite) TestTaskDueDateMarshalJSON() {
	s.Run("it should return json null when not valid", func() {
		taskDueDate := TaskDueDate{
			NullTime: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		}

		r, err := json.Marshal(taskDueDate)

		s.NoError(err)
		s.Equal("null", string(r))
	})

	s.Run("it should return string of time correctly", func() {
		currentTime := time.Now()

		taskDueDate := TaskDueDate{
			NullTime: sql.NullTime{
				Time:  currentTime,
				Valid: true,
			},
		}

		r, err := json.Marshal(taskDueDate)

		s.NoError(err)
		s.Equal(fmt.Sprintf("\"%s\"", currentTime.Format(time.RFC3339Nano)), string(r))
	})
}
