package usecase

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Shuhrat55/auth/internal/entity"
	"testing"
)

// MockUserRepository - мок для репозитория пользователей
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]entity.User, error) {
	args := m.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) Create(user entity.User) (entity.User, error) {
	args := m.Called(user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(id int, user entity.User) (entity.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) CheckPassword(id int, password string) bool {
	args := m.Called(id, password)
	return args.Bool(0)
}

func TestUserUseCase_GetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	mockUsers := []entity.User{
		{ID: 1, Name: "User1", Email: "user1@test.com"},
		{ID: 2, Name: "User2", Email: "user2@test.com"},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAll").Return(mockUsers, nil).Once()

		users, err := uc.GetAllUsers()

		assert.NoError(t, err)
		assert.Equal(t, mockUsers, users)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetAll").Return([]entity.User(nil), errors.New("db error")).Once()

		users, err := uc.GetAllUsers()

		assert.Error(t, err)
		assert.Nil(t, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	mockUser := entity.User{ID: 1, Name: "Test User", Email: "test@test.com"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByID", 1).Return(mockUser, nil).Once()

		user, err := uc.GetUserByID(1)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", 999).Return(entity.User{}, entity.ErrorUserNotFound).Once()

		user, err := uc.GetUserByID(999)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrorUserNotFound, err)
		assert.Equal(t, entity.User{}, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_GetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	mockUser := entity.User{ID: 1, Name: "Test User", Email: "test@test.com"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByEmail", "test@test.com").Return(mockUser, nil).Once()

		user, err := uc.GetUserByEmail("test@test.com")

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByEmail", "notfound@test.com").Return(entity.User{}, entity.ErrorUserNotFound).Once()

		user, err := uc.GetUserByEmail("notfound@test.com")

		assert.Error(t, err)
		assert.Equal(t, entity.ErrorUserNotFound, err)
		assert.Equal(t, entity.User{}, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		inputUser := entity.User{
			Name:     "New User",
			Email:    "new@test.com",
			Password: "password123",
		}

		expectedUser := entity.User{
			ID:       1,
			Name:     "New User",
			Email:    "new@test.com",
			Role:     "user",
		}

		mockRepo.On("Create", mock.MatchedBy(func(u entity.User) bool {
			return u.Name == inputUser.Name && u.Email == inputUser.Email && u.Role == "user" && u.Password != ""
		})).Return(expectedUser, nil).Once()

		user, err := uc.CreateUser(inputUser)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty password", func(t *testing.T) {
		inputUser := entity.User{
			Name:     "New User",
			Email:    "new@test.com",
			Password: "",
		}

		user, err := uc.CreateUser(inputUser)

		assert.Error(t, err)
		assert.Equal(t, "пароль не может быть пустым", err.Error())
		assert.Equal(t, entity.User{}, user)
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	updateUser := entity.User{
		Name:  "Updated User",
		Email: "updated@test.com",
	}

	expectedUser := entity.User{
		ID:    1,
		Name:  "Updated User",
		Email: "updated@test.com",
		Role:  "user",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Update", 1, updateUser).Return(expectedUser, nil).Once()

		user, err := uc.UpdateUser(1, updateUser)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("Update", 999, updateUser).Return(entity.User{}, errors.New("update error")).Once()

		user, err := uc.UpdateUser(999, updateUser)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()

		err := uc.DeleteUser(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("Delete", 999).Return(errors.New("delete error")).Once()

		err := uc.DeleteUser(999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_CheckPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	t.Run("valid password", func(t *testing.T) {
		mockRepo.On("CheckPassword", 1, "correctpass").Return(true).Once()

		valid := uc.CheckPassword(1, "correctpass")

		assert.True(t, valid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo.On("CheckPassword", 1, "wrongpass").Return(false).Once()

		valid := uc.CheckPassword(1, "wrongpass")

		assert.False(t, valid)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_Authenticate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	testUser := entity.User{
		ID:       1,
		Email:    "test@test.com",
		Password: "$2a$10$validhash",
	}

	t.Run("success authentication", func(t *testing.T) {
		mockRepo.On("GetByEmail", "test@test.com").Return(testUser, nil).Once()

		accessToken, refreshToken, expiresIn, err := uc.Authenticate("test@test.com", "correctpass")

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.Greater(t, expiresIn, int64(0))
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByEmail", "notfound@test.com").Return(entity.User{}, entity.ErrorUserNotFound).Once()

		accessToken, refreshToken, expiresIn, err := uc.Authenticate("notfound@test.com", "anypass")

		assert.Error(t, err)
		assert.Equal(t, entity.ErrorWrongPassword, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, int64(0), expiresIn)
		mockRepo.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockRepo.On("GetByEmail", "test@test.com").Return(entity.User{}, errors.New("db error")).Once()

		accessToken, refreshToken, expiresIn, err := uc.Authenticate("test@test.com", "correctpass")

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, int64(0), expiresIn)
		mockRepo.AssertExpectations(t)
	})
}
