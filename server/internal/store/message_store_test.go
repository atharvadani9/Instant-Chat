package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Message struct JSON marshaling
func TestMessageJSONMarshaling(t *testing.T) {
	message := &Message{
		ID:               1,
		SenderID:         123,
		ReceiverID:       456,
		EncryptedContent: "encrypted_secret_content",
		Content:          "Hello, World!",
		CreatedAt:        "2023-01-01T00:00:00Z",
	}

	// Test that encrypted content is excluded from JSON (due to json:"-" tag)
	jsonData, err := json.Marshal(message)
	require.NoError(t, err)

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, "Hello, World!")
	assert.Contains(t, jsonStr, "123")
	assert.Contains(t, jsonStr, "456")
	assert.Contains(t, jsonStr, "2023-01-01T00:00:00Z")
	assert.NotContains(t, jsonStr, "encrypted_secret_content")
	assert.NotContains(t, jsonStr, "encrypted_content")
}

func TestMessageStructFields(t *testing.T) {
	message := &Message{
		ID:               42,
		SenderID:         100,
		ReceiverID:       200,
		EncryptedContent: "encrypted_data",
		Content:          "decrypted_data",
		CreatedAt:        "2023-12-01T10:30:00Z",
	}

	assert.Equal(t, 42, message.ID)
	assert.Equal(t, 100, message.SenderID)
	assert.Equal(t, 200, message.ReceiverID)
	assert.Equal(t, "encrypted_data", message.EncryptedContent)
	assert.Equal(t, "decrypted_data", message.Content)
	assert.Equal(t, "2023-12-01T10:30:00Z", message.CreatedAt)
}

func TestMessageJSONUnmarshaling(t *testing.T) {
	jsonStr := `{
		"id": 1,
		"sender_id": 123,
		"receiver_id": 456,
		"content": "Hello, World!",
		"created_at": "2023-01-01T00:00:00Z"
	}`

	var message Message
	err := json.Unmarshal([]byte(jsonStr), &message)
	require.NoError(t, err)

	assert.Equal(t, 1, message.ID)
	assert.Equal(t, 123, message.SenderID)
	assert.Equal(t, 456, message.ReceiverID)
	assert.Equal(t, "Hello, World!", message.Content)
	assert.Equal(t, "2023-01-01T00:00:00Z", message.CreatedAt)
	assert.Empty(t, message.EncryptedContent) // Should not be set from JSON
}

func TestMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		isValid bool
	}{
		{
			name: "valid message",
			message: Message{
				ID:         1,
				SenderID:   100,
				ReceiverID: 200,
				Content:    "Hello",
				CreatedAt:  "2023-01-01T00:00:00Z",
			},
			isValid: true,
		},
		{
			name: "missing sender",
			message: Message{
				ID:         1,
				SenderID:   0,
				ReceiverID: 200,
				Content:    "Hello",
				CreatedAt:  "2023-01-01T00:00:00Z",
			},
			isValid: false,
		},
		{
			name: "missing receiver",
			message: Message{
				ID:         1,
				SenderID:   100,
				ReceiverID: 0,
				Content:    "Hello",
				CreatedAt:  "2023-01-01T00:00:00Z",
			},
			isValid: false,
		},
		{
			name: "empty content allowed",
			message: Message{
				ID:         1,
				SenderID:   100,
				ReceiverID: 200,
				Content:    "",
				CreatedAt:  "2023-01-01T00:00:00Z",
			},
			isValid: true,
		},
		{
			name: "same sender and receiver not allowed",
			message: Message{
				ID:         1,
				SenderID:   100,
				ReceiverID: 100,
				Content:    "Hello",
				CreatedAt:  "2023-01-01T00:00:00Z",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateMessage(&tt.message)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

// Helper function to validate message (this would be part of the actual implementation)
func validateMessage(msg *Message) bool {
	if msg.SenderID <= 0 || msg.ReceiverID <= 0 {
		return false
	}
	if msg.SenderID == msg.ReceiverID {
		return false
	}
	return true
}

func TestMessageContentTypes(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "simple text",
			content: "Hello, World!",
		},
		{
			name:    "unicode content",
			content: "Hello 世界! 🌍",
		},
		{
			name:    "empty content",
			content: "",
		},
		{
			name:    "special characters",
			content: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:    "multiline content",
			content: "Line 1\nLine 2\nLine 3",
		},
		{
			name:    "long content",
			content: "This is a very long message that contains multiple sentences and should test how the system handles larger text blocks to ensure everything works correctly with various message sizes and content types.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{
				ID:         1,
				SenderID:   100,
				ReceiverID: 200,
				Content:    tt.content,
				CreatedAt:  "2023-01-01T00:00:00Z",
			}

			// Test JSON marshaling preserves content
			jsonData, err := json.Marshal(message)
			require.NoError(t, err)

			var unmarshaled Message
			err = json.Unmarshal(jsonData, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, tt.content, unmarshaled.Content)
		})
	}
}

func BenchmarkMessageJSONMarshaling(b *testing.B) {
	message := &Message{
		ID:         1,
		SenderID:   123,
		ReceiverID: 456,
		Content:    "This is a test message for benchmarking JSON marshaling performance",
		CreatedAt:  "2023-01-01T00:00:00Z",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(message)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMessageJSONUnmarshaling(b *testing.B) {
	jsonStr := `{
		"id": 1,
		"sender_id": 123,
		"receiver_id": 456,
		"content": "This is a test message for benchmarking JSON unmarshaling performance",
		"created_at": "2023-01-01T00:00:00Z"
	}`
	jsonData := []byte(jsonStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var message Message
		err := json.Unmarshal(jsonData, &message)
		if err != nil {
			b.Fatal(err)
		}
	}
}
