package crypto

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Set environment variable for the main package
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
}

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "Hello, World!",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "unicode text",
			plaintext: "Hello 世界! 🌍",
		},
		{
			name:      "long text",
			plaintext: "This is a very long message that contains multiple sentences and should test the encryption and decryption of larger text blocks to ensure the system works correctly with various message sizes.",
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the plaintext
			encrypted, err := Encrypt(tt.plaintext)
			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)
			assert.NotEqual(t, tt.plaintext, encrypted)

			// Decrypt the ciphertext
			decrypted, err := Decrypt(encrypted)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestEncryptProducesUniqueResults(t *testing.T) {
	plaintext := "test message"

	// Encrypt the same message multiple times
	encrypted1, err := Encrypt(plaintext)
	require.NoError(t, err)

	encrypted2, err := Encrypt(plaintext)
	require.NoError(t, err)

	// Results should be different due to random nonce
	assert.NotEqual(t, encrypted1, encrypted2)

	// But both should decrypt to the same plaintext
	decrypted1, err := Decrypt(encrypted1)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted1)

	decrypted2, err := Decrypt(encrypted2)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted2)
}

func TestDecryptInvalidData(t *testing.T) {
	tests := []struct {
		name       string
		ciphertext string
		expectErr  bool
	}{
		{
			name:       "invalid base64",
			ciphertext: "invalid-base64!@#",
			expectErr:  true,
		},
		{
			name:       "too short data",
			ciphertext: "dGVzdA==", // "test" in base64, but too short for nonce
			expectErr:  true,
		},
		{
			name:       "empty string",
			ciphertext: "",
			expectErr:  true,
		},
		{
			name:       "corrupted data",
			ciphertext: "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY3ODkwYWJjZGVmZ2hpamtsbW5vcA==",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.ciphertext)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEncryptEmptyString(t *testing.T) {
	encrypted, err := Encrypt("")
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	decrypted, err := Decrypt(encrypted)
	require.NoError(t, err)
	assert.Equal(t, "", decrypted)
}

func BenchmarkEncrypt(b *testing.B) {
	plaintext := "This is a test message for benchmarking encryption performance"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Encrypt(plaintext)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	plaintext := "This is a test message for benchmarking decryption performance"
	encrypted, err := Encrypt(plaintext)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decrypt(encrypted)
		if err != nil {
			b.Fatal(err)
		}
	}
}
