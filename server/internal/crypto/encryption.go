package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/nacl/secretbox"
	"log"
	"os"
)

var encryptionKey [32]byte

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using default key")
	}

	keyHex := os.Getenv("ENCRYPTION_KEY")
	if keyHex == "" {
		log.Fatal("ENCRYPTION_KEY environment variable is required")
	}

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil || len(keyBytes) != 32 {
		log.Fatal("ENCRYPTION_KEY must be a 64-character hex string (32 bytes)")
	}

	copy(encryptionKey[:], keyBytes)
}

func Encrypt(plaintext string) (string, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	message := []byte(plaintext)
	encrypted := secretbox.Seal(nonce[:], message, &nonce, &encryptionKey)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	if len(data) < 24 {
		return "", fmt.Errorf("encrypted message is too short")
	}

	var nonce [24]byte
	copy(nonce[:], data[:24])
	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, &encryptionKey)
	if !ok {
		return "", fmt.Errorf("failed to decrypt message")
	}

	return string(decrypted), nil
}
