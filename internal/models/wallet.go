package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionMiningReward           TransactionType = "mining_reward"
	TransactionPurchaseDiscount       TransactionType = "purchase_discount"
	TransactionSubscriptionPayment    TransactionType = "subscription_payment"
	TransactionTransferSent           TransactionType = "transfer_sent"
	TransactionTransferReceived       TransactionType = "transfer_received"
	TransactionAdsPayment             TransactionType = "ads_payment"
	TransactionConsensusValidationRew TransactionType = "consensus_validation_reward"

	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Wallet struct {
	WalletID        uuid.UUID       `db:"wallet_id" json:"wallet_id"`
	UserID          uuid.UUID       `db:"user_id" json:"user_id"`
	BalanceTK       decimal.Decimal `db:"balance_tk" json:"balance_tk"`
	FrozenBalanceTK decimal.Decimal `db:"frozen_balance_tk" json:"frozen_balance_tk"`
	LastUpdated     time.Time       `db:"last_updated" json:"last_updated"`
}

type Transaction struct {
	TransactionID   uuid.UUID         `db:"transaction_id" json:"transaction_id"`
	UserID          uuid.UUID         `db:"user_id" json:"user_id"`
	TransactionType TransactionType   `db:"transaction_type" json:"transaction_type"`
	AmountTK        decimal.Decimal   `db:"amount_tk" json:"amount_tk"`
	Description     *string           `db:"description" json:"description"`
	RelatedEntityID *uuid.UUID        `db:"related_entity_id" json:"related_entity_id"`
	Status          TransactionStatus `db:"status" json:"status"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
}

// HasSufficientBalance checks if wallet has enough available balance
func (w *Wallet) HasSufficientBalance(amount decimal.Decimal) bool {
	availableBalance := w.BalanceTK.Sub(w.FrozenBalanceTK)
	return availableBalance.GreaterThanOrEqual(amount)
}

// AvailableBalance returns the available balance (not frozen)
func (w *Wallet) AvailableBalance() decimal.Decimal {
	return w.BalanceTK.Sub(w.FrozenBalanceTK)
}

// CanApplyTKDiscount calculates max TK discount (50% of order total)
func (w *Wallet) CanApplyTKDiscount(orderTotal decimal.Decimal) decimal.Decimal {
	maxDiscount := orderTotal.Div(decimal.NewFromInt(2))
	if w.AvailableBalance().LessThan(maxDiscount) {
		return w.AvailableBalance()
	}
	return maxDiscount
}
