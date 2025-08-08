package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		data           Envelope
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success response",
			status:         http.StatusOK,
			data:           Envelope{"message": "success", "data": "test"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":"test","message":"success"}`,
		},
		{
			name:           "error response",
			status:         http.StatusBadRequest,
			data:           Envelope{"error": "invalid request"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request"}`,
		},
		{
			name:           "empty data",
			status:         http.StatusNoContent,
			data:           Envelope{},
			expectedStatus: http.StatusNoContent,
			expectedBody:   `{}`,
		},
		{
			name:           "complex data",
			status:         http.StatusOK,
			data:           Envelope{"user": map[string]interface{}{"id": 1, "name": "test"}},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"user":{"id":1,"name":"test"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Call WriteJSON
			err := WriteJSON(w, tt.status, tt.data)
			require.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check content type
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			// Check body (parse JSON to avoid formatting differences)
			var actualData map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &actualData)
			require.NoError(t, err)

			var expectedData map[string]interface{}
			err = json.Unmarshal([]byte(tt.expectedBody), &expectedData)
			require.NoError(t, err)

			assert.Equal(t, expectedData, actualData)
		})
	}
}

func TestWriteJSONWithInvalidData(t *testing.T) {
	w := httptest.NewRecorder()

	// Create data that cannot be marshaled to JSON
	invalidData := Envelope{
		"invalid": make(chan int), // channels cannot be marshaled to JSON
	}

	err := WriteJSON(w, http.StatusOK, invalidData)
	assert.Error(t, err)
}

func TestWriteWebsocketMessage(t *testing.T) {
	// Create a simple test that doesn't require actual websocket connection
	// We'll test the function with a mock connection

	// Create a test logger
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	// Test data to send
	testData := map[string]interface{}{
		"type":    "test",
		"message": "hello world",
		"id":      123,
	}

	// Create a test server for websocket
	done := make(chan bool)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Write message using our function
		err = WriteWebsocketMessage(conn, testData, logger)
		if err != nil {
			t.Errorf("WriteWebsocketMessage failed: %v", err)
		}

		done <- true
	}))
	defer server.Close()

	// Connect to the test server
	wsURL := "ws" + server.URL[4:] // Replace http with ws
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Wait for the test to complete
	<-done
}

func TestWriteWebsocketMessageError(t *testing.T) {
	// Create a test logger
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	// Create a test server for websocket
	done := make(chan bool)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// Close the connection immediately
		conn.Close()

		// Try to write to closed connection
		testData := map[string]interface{}{"test": "data"}
		err = WriteWebsocketMessage(conn, testData, logger)
		if err == nil {
			t.Error("Expected error when writing to closed connection")
		}

		done <- true
	}))
	defer server.Close()

	// Connect to the test server
	wsURL := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Wait for the test to complete
	<-done
}

func BenchmarkWriteJSON(b *testing.B) {
	data := Envelope{
		"message": "test message",
		"user": map[string]interface{}{
			"id":       123,
			"username": "testuser",
			"active":   true,
		},
		"timestamp": "2023-01-01T00:00:00Z",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		err := WriteJSON(w, http.StatusOK, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
