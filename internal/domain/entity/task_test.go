package entity

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TaskEntityTestSuite struct {
	suite.Suite
}

func TestTaskEntitySuite(t *testing.T) {
	suite.Run(t, new(TaskEntityTestSuite))
}

func (s *TaskEntityTestSuite) TestValidate() {
	s.Run("it should return error when token is empty", func() {
		t := Task{}
		err := t.Validate()

		s.Equal(ErrContentEmpty, err)
	})

	s.Run("it should return nil when all fields are valid", func() {
		t := Task{
			Content: "Task content",
		}
		err := t.Validate()

		s.Nil(err)
	})
}
