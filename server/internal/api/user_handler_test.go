package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"chat/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserStore implements the UserStore interface for testing
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUser(user *store.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) GetUserByID(id int) (*store.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockUserStore) GetUserByUsername(username string) (*store.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockUserStore) GetUsersExcept(excludeUserID int) ([]*store.User, error) {
	args := m.Called(excludeUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*store.User), args.Error(1)
}

func (m *MockUserStore) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockUserStore) CheckPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockUserStore) AuthenticateUser(username, password string) (*store.User, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func TestNewUserHandler(t *testing.T) {
	mockStore := &MockUserStore{}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	handler := NewUserHandler(mockStore, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockStore, handler.Store)
	assert.Equal(t, logger, handler.logger)
}

func TestUserHandler_Register_Success(t *testing.T) {
	mockStore := &MockUserStore{}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	handler := NewUserHandler(mockStore, logger)

	// Setup mock expectations
	mockStore.On("GetUserByUsername", "testuser").Return(nil, sql.ErrNoRows)
	mockStore.On("HashPassword", "password123").Return("hashed_password", nil)
	mockStore.On("CreateUser", mock.AnythingOfType("*store.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*store.User)
		user.ID = 1
		user.CreatedAt = "2023-01-01T00:00:00Z"
	})

	// Create request
	reqBody := UserRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/user.register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.Register(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "User created successfully", response["message"])
	assert.NotNil(t, response["user"])

	mockStore.AssertExpectations(t)
}

func TestUserHandler_Register_UserExists(t *testing.T) {
	mockStore := &MockUserStore{}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	handler := NewUserHandler(mockStore, logger)

	// Setup mock expectations - user already exists
	existingUser := &store.User{ID: 1, Username: "testuser"}
	mockStore.On("GetUserByUsername", "testuser").Return(existingUser, nil)

	// Create request
	reqBody := UserRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/user.register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.Register(w, req)

	// Assertions
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Username already exists", response["error"])

	mockStore.AssertExpectations(t)
}

func TestUserHandler_Register_InvalidInput(t *testing.T) {
	mockStore := &MockUserStore{}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	handler := NewUserHandler(mockStore, logger)

	tests := []struct {
		name     string
		request  UserRequest
		expected string
	}{
		{
			name:     "empty username",
			request:  UserRequest{Username: "", Password: "password123"},
			expected: "Username and password are required",
		},
		{
			name:     "empty password",
			request:  UserRequest{Username: "testuser", Password: ""},
			expected: "Username and password are required",
		},
		{
			name:     "both empty",
			request:  UserRequest{Username: "", Password: ""},
			expected: "Username and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/user.register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Register(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, response["error"])
		})
	}
}

func TestUserHandler_Login_Success(t *testing.T) {
	mockStore := &MockUserStore{}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	handler := NewUserHandler(mockStore, logger)

	// Setup mock expectations
	user := &store.User{
		ID:       1,
		Username: "testuser",
	}
	mockStore.On("AuthenticateUser", "testuser", "password123").Return(user, nil)

	// Create request
	reqBody := UserRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/user.login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.Login(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Login successful", response["message"])
	assert.NotNil(t, response["user"])

	userMap := response["user"].(map[string]interface{})
	assert.Equal(t, float64(1), userMap["id"]) // JSON numbers are float64
	assert.Equal(t, "testuser", userMap["username"])

	mockStore.AssertExpectations(t)
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	mockStore := &MockUserStore{}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	handler := NewUserHandler(mockStore, logger)

	// Setup mock expectations - user not found
	mockStore.On("AuthenticateUser", "testuser", "wrongpassword").Return(nil, sql.ErrNoRows)

	// Create request
	reqBody := UserRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/user.login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.Login(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Invalid username or password", response["error"])

	mockStore.AssertExpectations(t)
}

func TestUserHandler_Login_DatabaseError(t *testing.T) {
	mockStore := &MockUserStore{}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	handler := NewUserHandler(mockStore, logger)

	// Setup mock expectations - database error
	mockStore.On("AuthenticateUser", "testuser", "password123").Return(nil, errors.New("database connection failed"))

	// Create request
	reqBody := UserRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/user.login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.Login(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Invalid username or password", response["error"])

	mockStore.AssertExpectations(t)
}
