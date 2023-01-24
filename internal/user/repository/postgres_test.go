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
)

type UserRepositoryPostgresTestSuite struct {
	suite.Suite
	db              *sql.DB
	userTableHelper test.UserTableHelper
	cleanUpSuite    func() error
	cleanUpTest     func()
}

func TestUserRepositoryPostgresSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryPostgresTestSuite))
}

func (s *UserRepositoryPostgresTestSuite) SetupSuite() {
	resource := test.NewPostgresResource()
	if err := database.Migration(resource.DSN, "../../../migrations"); err != nil {
		s.Fail("Could not migrate database", err)
	}

	s.userTableHelper = test.NewUserTableHelper(resource.DB)

	s.db = resource.DB
	s.cleanUpSuite = resource.CleanUp
	s.cleanUpTest = test.TruncateTables(resource.DB)
}

func (s *UserRepositoryPostgresTestSuite) TearDownSuite() {
	if err := s.cleanUpSuite(); err != nil {
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
		defer s.cleanUpTest()

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

func (s *UserRepositoryPostgresTestSuite) TestFindByEmail() {
	s.Run("it should return an error if database fail or operation canceled", func() {
		ctx := context.Background()
		repository := NewPostgres(s.db)
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		user, err := repository.FindByEmail(ctx, "gopher@gmail.com")

		s.EqualError(err, context.Canceled.Error())
		s.Empty(user)
	})

	s.Run("it should return error when user not found", func() {
		ctx := context.Background()
		repository := NewPostgres(s.db)
		email := "gopher@go.dev"

		user, err := repository.FindByEmail(ctx, email)

		s.Equal(ErrUserNotFound, err)
		s.Empty(user)
	})

	s.Run("it should return user when user found", func() {
		defer s.cleanUpTest()

		ctx := context.Background()
		repository := NewPostgres(s.db)
		email := "gopher@go.dev"

		existingUserInDB := test.User{
			ID:        "user-xxxxx",
			Name:      "Gopher",
			Email:     email,
			Password:  "secret_password",
			CreatedAt: test.TimeBeforeNow,
			UpdatedAt: test.TimeAfterNow,
		}
		s.userTableHelper.Add(existingUserInDB)

		user, err := repository.FindByEmail(ctx, email)

		s.NoError(err)
		s.Equal(existingUserInDB.ID, user.ID)
		s.Equal(existingUserInDB.Name, user.Name)
		s.Equal(existingUserInDB.Email, user.Email)
		s.Equal(existingUserInDB.Password, user.Password)
		s.Equal(existingUserInDB.CreatedAt, user.CreatedAt)
		s.Equal(existingUserInDB.UpdatedAt, user.UpdatedAt)
	})
}
