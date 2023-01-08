package dto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TaskDTOTestSuite struct {
	suite.Suite
}

func TestTaskDTOSuite(t *testing.T) {
	suite.Run(t, new(TaskDTOTestSuite))
}

func (s *TaskDTOTestSuite) TestTaskCreateIn() {
	tests := []struct {
		name     string
		input    TaskCreateIn
		expected error
	}{
		{
			name:     "it should return error when content is empty",
			input:    TaskCreateIn{},
			expected: ErrContentEmpty,
		},
		{
			name: "it should return nil when all fields are valid",
			input: TaskCreateIn{
				Content: "content",
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}

func (s *TaskDTOTestSuite) TestTaskUpdateIn() {
	tests := []struct {
		name     string
		input    TaskUpdateIn
		expected error
	}{
		{
			name:     "it should return error when content is empty",
			input:    TaskUpdateIn{},
			expected: ErrContentEmpty,
		},
		{
			name: "it should return nil when all fields are valid",
			input: TaskUpdateIn{
				Content: "content",
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}
