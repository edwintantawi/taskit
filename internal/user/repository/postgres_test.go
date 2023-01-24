package repository

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/entity"
	"github.com/edwintantawi/taskit/pkg/database"
	"github.com/edwintantawi/taskit/test"
	"github.com/edwintantawi/taskit/test/testdb"
)

type UserRepositoryPostgresTestSuite struct {
	suite.Suite
	db              *sql.DB
	userTableHelper testdb.UserTableHelper
	cleanUp         func() error
}

func TestUserRepositoryPostgresSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryPostgresTestSuite))
}

func (s *UserRepositoryPostgresTestSuite) SetupSuite() {
	resource := test.NewPostgresResource()
	if err := database.Migration(resource.DSN, "../../../migrations"); err != nil {
		s.Fail("Could not migrate database", err)
	}

	s.userTableHelper = testdb.NewUserTableHelper(resource.DB)

	s.db = resource.DB
	s.cleanUp = resource.CleanUp
}

func (s *UserRepositoryPostgresTestSuite) TearDownSuite() {
	if err := s.cleanUp(); err != nil {
		log.Fatalf("Could not cleanup resource: %s", err)
	}
}

func (s *UserRepositoryPostgresTestSuite) TestSave() {
	s.Run("it should return an error if database fail or operation canceled", func() {
		ctx := context.Background()
		repository := NewPostgres(s.db)
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := repository.Save(ctx, entity.NewUser{})

		s.EqualError(err, context.Canceled.Error())
	})

	s.Run("it should save a new user to the database", func() {
		ctx := context.Background()
		repository := NewPostgres(s.db)
		newUser := entity.NewUser{
			Name:     "Gopher",
			Email:    "gopher@go.dev",
			Password: "secret_password",
		}

		addedUser, err := repository.Save(ctx, newUser)

		s.NoError(err)
		s.NotEmpty(addedUser.ID)
		s.Equal(newUser.Email, addedUser.Email)

		userInDB := s.userTableHelper.GetByID(addedUser.ID)
		s.Equal(addedUser.ID, userInDB.ID)
		s.Equal(newUser.Name, userInDB.Name)
		s.Equal(newUser.Email, userInDB.Email)
		s.Equal(newUser.Password, userInDB.Password)
		s.NotEmpty(userInDB.CreatedAt)
		s.Equal(userInDB.CreatedAt, userInDB.UpdatedAt)
	})
}
