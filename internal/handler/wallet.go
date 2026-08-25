package handler

import (
	"net/http"

	apierrors "github.com/azizjon-top/techno-re-ecosystem/internal/errors"
	"github.com/azizjon-top/techno-re-ecosystem/internal/middleware"
	"github.com/azizjon-top/techno-re-ecosystem/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletHandler handles wallet endpoints.
type WalletHandler struct {
	walletSvc *service.WalletService
}

// NewWalletHandler creates a WalletHandler.
func NewWalletHandler(walletSvc *service.WalletService) *WalletHandler {
	return &WalletHandler{walletSvc: walletSvc}
}

// GetBalance godoc
// GET /api/v1/wallet/balance
func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)
	wallet, err := h.walletSvc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, wallet)
}

// transferRequest is the JSON body for POST /wallet/transfer.
type transferRequest struct {
	ReceiverID string `json:"receiver_id" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

// Transfer godoc
// POST /api/v1/wallet/transfer
func (h *WalletHandler) Transfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeValidationFailed,
			"message": err.Error(),
		})
		return
	}

	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeBadRequest,
			"message": "invalid receiver_id",
		})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeBadRequest,
			"message": "invalid amount format",
		})
		return
	}

	senderID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	if svcErr := h.walletSvc.Transfer(c.Request.Context(), service.TransferRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
	}); svcErr != nil {
		respondError(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}
