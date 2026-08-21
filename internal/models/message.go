package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Chat struct {
	ChatID         uuid.UUID `db:"chat_id" json:"chat_id"`
	Participant1ID uuid.UUID `db:"participant_1_id" json:"participant_1_id"`
	Participant2ID uuid.UUID `db:"participant_2_id" json:"participant_2_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type Message struct {
	MessageID        uuid.UUID          `db:"message_id" json:"message_id"`
	ChatID           uuid.UUID          `db:"chat_id" json:"chat_id"`
	SenderID         uuid.UUID          `db:"sender_id" json:"sender_id"`
	EncryptedText    []byte             `db:"encrypted_text" json:"-"`
	TKTransferAmount *decimal.Decimal   `db:"tk_transfer_amount" json:"tk_transfer_amount"`
	IsRead           bool               `db:"is_read" json:"is_read"`
	CreatedAt        time.Time          `db:"created_at" json:"created_at"`
}

// EncryptMessage placeholder for E2EE encryption
// In production, use proper encryption library (NaCl Box, etc.)
func EncryptMessage(plaintext string, key []byte) ([]byte, error) {
	// TODO: Implement proper E2EE encryption with NaCl
	return []byte(plaintext), nil
}

// DecryptMessage placeholder for E2EE decryption
// In production, use proper decryption library (NaCl Box, etc.)
func DecryptMessage(ciphertext []byte, key []byte) (string, error) {
	// TODO: Implement proper E2EE decryption with NaCl
	return string(ciphertext), nil
}
