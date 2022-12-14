package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PrimitiveTestSuite struct {
	suite.Suite
}

func TestPrimitiveSuite(t *testing.T) {
	suite.Run(t, new(PrimitiveTestSuite))
}

func (s *PrimitiveTestSuite) TestNullTimeUnmarshalJSON() {
	s.Run("it should return error when fail to unmarshal with invalid json", func() {
		rawJson := `[]`
		var dateTime NullTime

		err := json.Unmarshal([]byte(rawJson), &dateTime)

		s.Error(err)
		s.False(dateTime.Valid)
		s.Empty(dateTime.Time)
	})

	s.Run("it should successfully unmarshal and return valid false and time is zero value", func() {
		rawJson := `null`
		var dateTime NullTime

		err := json.Unmarshal([]byte(rawJson), &dateTime)

		s.NoError(err)
		s.False(dateTime.Valid)
		s.Empty(dateTime.Time)
	})

	s.Run("it should successfully unmarshal and return valid true and time is actual time form json", func() {
		rawJson := `"2022-12-25T00:00:00.000Z"`
		var dateTime NullTime

		err := json.Unmarshal([]byte(rawJson), &dateTime)

		s.NoError(err)
		s.True(dateTime.Valid)
		s.Equal("2022-12-25 00:00:00 +0000 UTC", dateTime.Time.String())
	})
}

func (s *PrimitiveTestSuite) TestNullTimeMarshalJSON() {
	s.Run("it should successfully marshal and return json null when not valid", func() {
		dateTime := NullTime{
			NullTime: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		}

		r, err := json.Marshal(dateTime)

		s.NoError(err)
		s.Equal("null", string(r))
	})

	s.Run("it should successfully marshal and return json time correctly", func() {
		currentTime := time.Now()

		dateTime := NullTime{
			NullTime: sql.NullTime{
				Time:  currentTime,
				Valid: true,
			},
		}

		r, err := json.Marshal(dateTime)

		s.NoError(err)
		s.Equal(fmt.Sprintf("\"%s\"", currentTime.Format(time.RFC3339Nano)), string(r))
	})
}
