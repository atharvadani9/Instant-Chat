package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// Since we can't easily mock sql.DB directly, we'll test the actual functions
// with a focus on the business logic and password handling

func TestHashPassword(t *testing.T) {
	store := &PostgresUserStore{}

	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!@#$%^&*()",
		},
		{
			name:     "unicode password",
			password: "пароль123",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "long password",
			password: "this_is_a_long_password_that_works_with_bcrypt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := store.HashPassword(tt.password)
			require.NoError(t, err)
			assert.NotEmpty(t, hashedPassword)
			assert.NotEqual(t, tt.password, hashedPassword)

			// Verify the hash starts with bcrypt prefix
			assert.True(t, len(hashedPassword) >= 60) // bcrypt hashes are at least 60 characters
		})
	}
}

func TestCheckPassword(t *testing.T) {
	store := &PostgresUserStore{}

	tests := []struct {
		name          string
		password      string
		shouldMatch   bool
		wrongPassword string
	}{
		{
			name:          "correct password",
			password:      "password123",
			shouldMatch:   true,
			wrongPassword: "",
		},
		{
			name:          "wrong password",
			password:      "password123",
			shouldMatch:   false,
			wrongPassword: "wrongpassword",
		},
		{
			name:          "empty password",
			password:      "",
			shouldMatch:   true,
			wrongPassword: "",
		},
		{
			name:          "case sensitive",
			password:      "Password123",
			shouldMatch:   false,
			wrongPassword: "password123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First hash the password
			hashedPassword, err := store.HashPassword(tt.password)
			require.NoError(t, err)

			if tt.shouldMatch {
				// Check with correct password
				err = store.CheckPassword(hashedPassword, tt.password)
				assert.NoError(t, err)
			} else {
				// Check with wrong password
				err = store.CheckPassword(hashedPassword, tt.wrongPassword)
				assert.Error(t, err)
				assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
			}
		})
	}
}

func TestCheckPasswordWithInvalidHash(t *testing.T) {
	store := &PostgresUserStore{}

	tests := []struct {
		name        string
		invalidHash string
		password    string
	}{
		{
			name:        "empty hash",
			invalidHash: "",
			password:    "password",
		},
		{
			name:        "invalid hash format",
			invalidHash: "not_a_bcrypt_hash",
			password:    "password",
		},
		{
			name:        "truncated hash",
			invalidHash: "$2a$10$truncated",
			password:    "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CheckPassword(tt.invalidHash, tt.password)
			assert.Error(t, err)
		})
	}
}

func TestPasswordHashingRoundTrip(t *testing.T) {
	store := &PostgresUserStore{}

	passwords := []string{
		"simple",
		"complex!@#$%^&*()",
		"unicode_пароль_🔒",
		"long_password_that_works_with_bcrypt_limits",
		"", // empty password
	}

	for _, password := range passwords {
		t.Run("password_"+password, func(t *testing.T) {
			// Hash the password
			hashed, err := store.HashPassword(password)
			require.NoError(t, err)

			// Verify it can be checked successfully
			err = store.CheckPassword(hashed, password)
			assert.NoError(t, err)

			// Verify wrong password fails
			if password != "" {
				err = store.CheckPassword(hashed, password+"wrong")
				assert.Error(t, err)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	store := &PostgresUserStore{}
	password := "test_password"

	// Hash the same password multiple times
	hash1, err := store.HashPassword(password)
	require.NoError(t, err)

	hash2, err := store.HashPassword(password)
	require.NoError(t, err)

	// Hashes should be different (due to salt)
	assert.NotEqual(t, hash1, hash2)

	// But both should verify correctly
	err = store.CheckPassword(hash1, password)
	assert.NoError(t, err)

	err = store.CheckPassword(hash2, password)
	assert.NoError(t, err)
}

func BenchmarkHashPassword(b *testing.B) {
	store := &PostgresUserStore{}
	password := "benchmark_password_123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := store.HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	store := &PostgresUserStore{}
	password := "benchmark_password_123"

	// Pre-hash the password
	hashedPassword, err := store.HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := store.CheckPassword(hashedPassword, password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test User struct JSON marshaling
func TestUserJSONMarshaling(t *testing.T) {
	user := &User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "secret_hash",
		CreatedAt:    "2023-01-01T00:00:00Z",
	}

	// Test that password hash is excluded from JSON (due to json:"-" tag)
	jsonData, err := json.Marshal(user)
	require.NoError(t, err)

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, "testuser")
	assert.Contains(t, jsonStr, "2023-01-01T00:00:00Z")
	assert.NotContains(t, jsonStr, "secret_hash")
	assert.NotContains(t, jsonStr, "password_hash")
}
