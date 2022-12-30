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
	tests := []struct {
		name     string
		input    Task
		expected error
	}{
		{name: "it should return error when content is empty", input: Task{}, expected: ErrContentEmpty},
		{name: "it should return nil when all fields are valid", input: Task{Content: "Task content"}, expected: nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}
