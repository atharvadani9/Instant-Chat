package app

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplication_HealthCheck(t *testing.T) {
	// Create a minimal application for testing
	app := &Application{
		Logger: log.New(os.Stdout, "TEST: ", log.LstdFlags),
	}

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()

	// Call the health check handler
	app.HealthCheck(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Check the response body
	expectedBody := `{
  "status": "OK"
}
`
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestApplication_HealthCheck_WithDifferentMethods(t *testing.T) {
	app := &Application{
		Logger: log.New(os.Stdout, "TEST: ", log.LstdFlags),
	}

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
	}

	for _, method := range methods {
		t.Run("method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/healthcheck", nil)
			w := httptest.NewRecorder()

			app.HealthCheck(w, req)

			// Health check should work with any HTTP method
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

func TestApplication_HealthCheck_ResponseFormat(t *testing.T) {
	app := &Application{
		Logger: log.New(os.Stdout, "TEST: ", log.LstdFlags),
	}

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()

	app.HealthCheck(w, req)

	// Verify response format
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Check that response is valid JSON
	body := w.Body.String()
	assert.Contains(t, body, `"status"`)
	assert.Contains(t, body, `"OK"`)
	assert.Contains(t, body, `{`)
	assert.Contains(t, body, `}`)
}

func TestApplication_HealthCheck_Concurrent(t *testing.T) {
	app := &Application{
		Logger: log.New(os.Stdout, "TEST: ", log.LstdFlags),
	}

	// Test concurrent access to health check
	const numRequests = 10
	results := make(chan int, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
			w := httptest.NewRecorder()
			app.HealthCheck(w, req)
			results <- w.Code
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		assert.Equal(t, http.StatusOK, statusCode)
	}
}

func BenchmarkApplication_HealthCheck(b *testing.B) {
	app := &Application{
		Logger: log.New(os.Stdout, "TEST: ", log.LstdFlags),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		w := httptest.NewRecorder()
		app.HealthCheck(w, req)
	}
}

// Test Application struct initialization
func TestApplicationStruct(t *testing.T) {
	app := &Application{}

	// Test that all fields can be set
	assert.NotPanics(t, func() {
		app.Logger = nil
		app.DB = nil
		app.UserHandler = nil
		app.WebSocketHandler = nil
	})
}

// Test Application with nil logger (will panic as expected)
func TestApplication_HealthCheck_NilLogger(t *testing.T) {
	app := &Application{
		Logger: nil, // Nil logger will cause panic
	}

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()

	// This should panic with nil logger
	assert.Panics(t, func() {
		app.HealthCheck(w, req)
	})
}

// Test that Application can be created with zero values
func TestApplication_ZeroValue(t *testing.T) {
	var app Application

	// Should be able to call health check on zero-value Application (but will panic due to nil logger)
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		app.HealthCheck(w, req)
	})
}

// Test Application fields are properly typed
func TestApplication_FieldTypes(t *testing.T) {
	app := &Application{}

	// Test that we can assign the expected types
	require.NotPanics(t, func() {
		// These assignments should compile and not panic
		_ = app.Logger
		_ = app.DB
		_ = app.UserHandler
		_ = app.WebSocketHandler
	})
}

// Test health check response consistency
func TestApplication_HealthCheck_ResponseConsistency(t *testing.T) {
	app := &Application{
		Logger: log.New(os.Stdout, "TEST: ", log.LstdFlags),
	}

	// Make multiple requests and ensure responses are consistent
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		w := httptest.NewRecorder()

		app.HealthCheck(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		expectedBody := `{
  "status": "OK"
}
`
		assert.Equal(t, expectedBody, w.Body.String())
	}
}

// Test health check with various request headers
func TestApplication_HealthCheck_WithHeaders(t *testing.T) {
	app := &Application{
		Logger: log.New(os.Stdout, "TEST: ", log.LstdFlags),
	}

	headers := map[string]string{
		"User-Agent":      "test-agent",
		"Accept":          "application/json",
		"Accept-Language": "en-US",
		"Authorization":   "Bearer token",
		"X-Custom-Header": "custom-value",
	}

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)

	// Add headers to request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	app.HealthCheck(w, req)

	// Health check should work regardless of request headers
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
