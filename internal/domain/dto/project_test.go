package dto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProjectDTOTestSuite struct {
	suite.Suite
}

func TestProjectDTOSuite(t *testing.T) {
	suite.Run(t, new(ProjectDTOTestSuite))
}

func (s *ProjectDTOTestSuite) TestProjectCreateIn() {
	tests := []struct {
		name     string
		input    ProjectCreateIn
		expected error
	}{
		{
			name:     "it should return error when content is empty",
			input:    ProjectCreateIn{},
			expected: ErrTitleEmpty,
		},
		{
			name: "it should return nil when all fields are valid",
			input: ProjectCreateIn{
				Title: "project_title",
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
