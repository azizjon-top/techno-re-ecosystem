package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type MiningStatus string

const (
	MiningStatusActive    MiningStatus = "active"
	MiningStatusPaused    MiningStatus = "paused"
	MiningStatusCompleted MiningStatus = "completed"
)

type MiningSession struct {
	SessionID       uuid.UUID       `db:"session_id" json:"session_id"`
	UserID          uuid.UUID       `db:"user_id" json:"user_id"`
	StartedAt       time.Time       `db:"started_at" json:"started_at"`
	PausedAt        *time.Time      `db:"paused_at" json:"paused_at"`
	DurationMinutes int             `db:"duration_minutes" json:"duration_minutes"`
	FactsVerified   int             `db:"facts_verified" json:"facts_verified"`
	RewardTKEarned  decimal.Decimal `db:"reward_tk_earned" json:"reward_tk_earned"`
	Status          MiningStatus    `db:"status" json:"status"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
}

// CalculateDuration calculates session duration
func (ms *MiningSession) CalculateDuration() time.Duration {
	if ms.PausedAt != nil {
		return ms.PausedAt.Sub(ms.StartedAt)
	}
	return time.Since(ms.StartedAt)
}

// IsActive checks if mining session is currently active
func (ms *MiningSession) IsActive() bool {
	return ms.Status == MiningStatusActive
}

// CalculateReward calculates reward based on facts verified
func (ms *MiningSession) CalculateReward(ratePerFact decimal.Decimal) decimal.Decimal {
	return ratePerFact.Mul(decimal.NewFromInt(int64(ms.FactsVerified)))
}
