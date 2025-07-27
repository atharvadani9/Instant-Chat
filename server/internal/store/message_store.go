package store

import (
	"chat/internal/crypto"
	"database/sql"
)

type Message struct {
	ID               int    `json:"id"`
	SenderID         int    `json:"sender_id"`
	ReceiverID       int    `json:"receiver_id"`
	EncryptedContent string `json:"-"`
	Content          string `json:"content"`
	CreatedAt        string `json:"created_at"`
}

type PostgresMessageStore struct {
	db *sql.DB
}

func NewPostgresMessageStore(db *sql.DB) *PostgresMessageStore {
	return &PostgresMessageStore{db: db}
}

type MessageStore interface {
	CreateMessage(senderID, receiverID int, content string) (*Message, error)
	GetMessagesBetweenUsers(userID1, userID2 int) ([]*Message, error)
}

func (s *PostgresMessageStore) CreateMessage(senderID, receiverID int, content string) (*Message, error) {
	encryptedContent, err := crypto.Encrypt(content)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO messages (sender_id, receiver_id, encrypted_content) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`
	message := &Message{}
	err = s.db.QueryRow(query, senderID, receiverID, encryptedContent).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	message.SenderID = senderID
	message.ReceiverID = receiverID
	message.EncryptedContent = encryptedContent
	message.Content = content
	return message, nil
}

func (s *PostgresMessageStore) GetMessagesBetweenUsers(userID1, userID2 int) ([]*Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, encrypted_content, created_at 
		FROM messages 
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at
	`
	rows, err := s.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.EncryptedContent, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		message.Content, err = crypto.Decrypt(message.EncryptedContent)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
