package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type AIRequestType string
type ValidationResult string

const (
	AIRequestProductSearch       AIRequestType = "product_search"
	AIRequestChatAssistance      AIRequestType = "chat_assistance"
	AIRequestContentModeration   AIRequestType = "content_moderation"
	AIRequestVideoRecommendation AIRequestType = "video_recommendation"

	ValidationTrue        ValidationResult = "true"
	ValidationFalse       ValidationResult = "false"
	ValidationInconcluisve ValidationResult = "inconclusive"
)

const ConsensusThreshold = 0.8

type AIRequest struct {
	RequestID    uuid.UUID     `db:"request_id" json:"request_id"`
	UserID       uuid.UUID     `db:"user_id" json:"user_id"`
	RequestType  AIRequestType `db:"request_type" json:"request_type"`
	QueryText    string        `db:"query_text" json:"query_text"`
	ResponseText *string       `db:"response_text" json:"response_text"`
	CreatedAt    time.Time     `db:"created_at" json:"created_at"`
}

type ConsensusValidation struct {
	ValidationID    uuid.UUID        `db:"validation_id" json:"validation_id"`
	AIRequestID     uuid.UUID        `db:"ai_request_id" json:"ai_request_id"`
	ValidatorNodeID uuid.UUID        `db:"validator_node_id" json:"validator_node_id"`
	FactClaim       string           `db:"fact_claim" json:"fact_claim"`
	ValidationResult ValidationResult `db:"validation_result" json:"validation_result"`
	ConfidenceScore decimal.Decimal   `db:"confidence_score" json:"confidence_score"`
	CreatedAt       time.Time        `db:"created_at" json:"created_at"`
}

// AverageConfidence calculates average confidence score from validations
func AverageConfidence(validations []ConsensusValidation) decimal.Decimal {
	if len(validations) == 0 {
		return decimal.Zero
	}
	sum := decimal.Zero
	for _, v := range validations {
		sum = sum.Add(v.ConfidenceScore)
	}
	return sum.Div(decimal.NewFromInt(int64(len(validations))))
}

// IsConsensusReached checks if consensus threshold is met (V >= 0.8)
func IsConsensusReached(validations []ConsensusValidation) bool {
	avg := AverageConfidence(validations)
	threshold := decimal.NewFromString("0.8")
	return avg.GreaterThanOrEqual(threshold)
}
