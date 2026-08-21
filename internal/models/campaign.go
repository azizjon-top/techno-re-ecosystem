package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type CampaignStatus string

const (
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusCancelled CampaignStatus = "cancelled"
)

type AdCampaign struct {
	CampaignID  uuid.UUID       `db:"campaign_id" json:"campaign_id"`
	CreatorID   uuid.UUID       `db:"creator_id" json:"creator_id"`
	Title       string          `db:"title" json:"title"`
	Description *string         `db:"description" json:"description"`
	BannerURL   *string         `db:"banner_url" json:"banner_url"`
	StartDate   time.Time       `db:"start_date" json:"start_date"`
	EndDate     time.Time       `db:"end_date" json:"end_date"`
	BudgetTK    decimal.Decimal `db:"budget_tk" json:"budget_tk"`
	SpentTK     decimal.Decimal `db:"spent_tk" json:"spent_tk"`
	Status      CampaignStatus  `db:"status" json:"status"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
}

type AdImpression struct {
	ImpressionID uuid.UUID `db:"impression_id" json:"impression_id"`
	CampaignID   uuid.UUID `db:"campaign_id" json:"campaign_id"`
	ViewerID     uuid.UUID `db:"viewer_id" json:"viewer_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// IsActive checks if campaign is currently active
func (ac *AdCampaign) IsActive() bool {
	now := time.Now()
	return ac.Status == CampaignStatusActive &&
		ac.StartDate.Before(now) &&
		ac.EndDate.After(now)
}

// RemainingBudget returns remaining budget for campaign
func (ac *AdCampaign) RemainingBudget() decimal.Decimal {
	return ac.BudgetTK.Sub(ac.SpentTK)
}

// CanSpend checks if campaign has enough budget to spend
func (ac *AdCampaign) CanSpend(amount decimal.Decimal) bool {
	return ac.RemainingBudget().GreaterThanOrEqual(amount)
}
