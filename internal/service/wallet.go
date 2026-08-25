package service

import (
	"context"
	"errors"

	apierrors "github.com/azizjon-top/techno-re-ecosystem/internal/errors"
	"github.com/azizjon-top/techno-re-ecosystem/internal/models"
	"github.com/azizjon-top/techno-re-ecosystem/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletService handles wallet business logic.
type WalletService struct {
	wallets repository.WalletRepository
}

// NewWalletService creates a WalletService.
func NewWalletService(wallets repository.WalletRepository) *WalletService {
	return &WalletService{wallets: wallets}
}

// GetBalance returns the wallet for the given user.
func (s *WalletService) GetBalance(ctx context.Context, userID uuid.UUID) (*models.Wallet, error) {
	w, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "wallet not found", 404)
		}
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to get wallet", 500)
	}
	return w, nil
}

// TransferRequest holds input for a TK transfer.
type TransferRequest struct {
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	Amount     decimal.Decimal
}

// Transfer moves TK tokens between two wallets.
func (s *WalletService) Transfer(ctx context.Context, req TransferRequest) error {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return apierrors.NewAPIError(apierrors.ErrCodeBadRequest, "transfer amount must be greater than zero", 400)
	}
	if req.SenderID == req.ReceiverID {
		return apierrors.NewAPIError(apierrors.ErrCodeBadRequest, "cannot transfer to yourself", 400)
	}

	senderWallet, err := s.wallets.GetByUserID(ctx, req.SenderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apierrors.NewAPIError(apierrors.ErrCodeNotFound, "sender wallet not found", 404)
		}
		return apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to get sender wallet", 500)
	}
	if !senderWallet.HasSufficientBalance(req.Amount) {
		return apierrors.NewAPIError(apierrors.ErrCodeInsufficientBalance, "insufficient balance", 400)
	}

	receiverWallet, err := s.wallets.GetByUserID(ctx, req.ReceiverID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apierrors.NewAPIError(apierrors.ErrCodeNotFound, "receiver wallet not found", 404)
		}
		return apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to get receiver wallet", 500)
	}

	// Debit sender
	newSenderBalance := senderWallet.BalanceTK.Sub(req.Amount)
	if err := s.wallets.UpdateBalance(ctx, senderWallet.WalletID, newSenderBalance, senderWallet.FrozenBalanceTK); err != nil {
		return apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to update sender wallet", 500)
	}

	// Credit receiver
	newReceiverBalance := receiverWallet.BalanceTK.Add(req.Amount)
	if err := s.wallets.UpdateBalance(ctx, receiverWallet.WalletID, newReceiverBalance, receiverWallet.FrozenBalanceTK); err != nil {
		return apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to update receiver wallet", 500)
	}

	desc := "TK transfer"
	_ = s.wallets.CreateTransaction(ctx, &models.Transaction{
		UserID:          req.SenderID,
		TransactionType: models.TransactionTransferSent,
		AmountTK:        req.Amount,
		Description:     &desc,
		RelatedEntityID: &req.ReceiverID,
		Status:          models.TransactionStatusCompleted,
	})
	_ = s.wallets.CreateTransaction(ctx, &models.Transaction{
		UserID:          req.ReceiverID,
		TransactionType: models.TransactionTransferReceived,
		AmountTK:        req.Amount,
		Description:     &desc,
		RelatedEntityID: &req.SenderID,
		Status:          models.TransactionStatusCompleted,
	})
	return nil
}
