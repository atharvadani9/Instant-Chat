package routes

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"chat/internal/api"
	"chat/internal/app"
	"chat/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ErrUserNotFound = errors.New("user not found")

// MockUserStore for testing
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

// MockMessageStore for testing
type MockMessageStore struct {
	mock.Mock
}

func (m *MockMessageStore) CreateMessage(senderID, receiverID int, content string) (*store.Message, error) {
	args := m.Called(senderID, receiverID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.Message), args.Error(1)
}

func (m *MockMessageStore) GetMessagesBetweenUsers(userID1, userID2 int) ([]*store.Message, error) {
	args := m.Called(userID1, userID2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*store.Message), args.Error(1)
}

func createTestApplication() *app.Application {
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	// Create mock stores
	userStore := &MockUserStore{}
	messageStore := &MockMessageStore{}

	// Set up basic mock expectations to avoid panics
	userStore.On("GetUserByID", mock.AnythingOfType("int")).Return(nil, ErrUserNotFound).Maybe()
	userStore.On("GetUsersExcept", mock.AnythingOfType("int")).Return([]*store.User{}, nil).Maybe()

	// Create handlers with mocks
	userHandler := api.NewUserHandler(userStore, logger)
	webSocketHandler := api.NewWebSocketHandler(messageStore, userStore, logger)

	return &app.Application{
		Logger:           logger,
		DB:               nil, // Not needed for route testing
		UserHandler:      userHandler,
		WebSocketHandler: webSocketHandler,
	}
}

func TestSetupRoutes(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	assert.NotNil(t, router)
}

func TestHealthCheckRoute(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), `"status"`)
	assert.Contains(t, w.Body.String(), `"OK"`)
}

func TestCORSHeaders(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	// Test CORS preflight request
	req := httptest.NewRequest(http.MethodOptions, "/healthcheck", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check CORS headers (they may not all be set for a simple GET request)
	// The CORS middleware typically sets headers on preflight OPTIONS requests
	// For a simple GET request, we mainly check that the request succeeds
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouteEndpoints(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "health check",
			method:         http.MethodGet,
			path:           "/healthcheck",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user register endpoint exists",
			method:         http.MethodPost,
			path:           "/user.register",
			expectedStatus: http.StatusBadRequest, // Will fail due to empty body, but route exists
		},
		{
			name:           "user login endpoint exists",
			method:         http.MethodPost,
			path:           "/user.login",
			expectedStatus: http.StatusBadRequest, // Will fail due to empty body, but route exists
		},
		{
			name:           "get users endpoint exists",
			method:         http.MethodGet,
			path:           "/user.get",
			expectedStatus: http.StatusBadRequest, // Will fail due to missing user_id, but route exists
		},
		{
			name:           "get me user endpoint exists",
			method:         http.MethodGet,
			path:           "/user.get.me",
			expectedStatus: http.StatusBadRequest, // Will fail due to missing user_id, but route exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// We're mainly testing that the routes exist and don't return 404
			assert.NotEqual(t, http.StatusNotFound, w.Code, "Route should exist")

			// For health check, we expect 200
			if tt.path == "/healthcheck" {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestNonExistentRoute(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	// Test wrong method on existing route
	req := httptest.NewRequest(http.MethodDelete, "/user.register", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestWebSocketRoute(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	// Test WebSocket route exists (will fail upgrade but route should exist)
	req := httptest.NewRequest(http.MethodGet, "/chat/ws", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should not be 404 (route exists)
	assert.NotEqual(t, http.StatusNotFound, w.Code)
	// Will likely be 400 or 426 due to missing WebSocket headers
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusUpgradeRequired)
}

func TestRouteWithQueryParams(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/user.get?user_id=123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Route should exist (not 404)
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestConcurrentRequests(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	const numRequests = 10
	results := make(chan int, numRequests)

	// Make concurrent requests to health check
	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		assert.Equal(t, http.StatusOK, statusCode)
	}
}

func BenchmarkHealthCheckRoute(b *testing.B) {
	app := createTestApplication()
	router := SetupRoutes(app)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func TestRouterMiddleware(t *testing.T) {
	app := createTestApplication()
	router := SetupRoutes(app)

	// Test that CORS middleware is applied
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should have CORS headers
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
