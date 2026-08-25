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

// ProductService handles product and order business logic.
type ProductService struct {
	products repository.ProductRepository
	orders   repository.OrderRepository
	wallets  repository.WalletRepository
}

// NewProductService creates a ProductService.
func NewProductService(products repository.ProductRepository, orders repository.OrderRepository, wallets repository.WalletRepository) *ProductService {
	return &ProductService{products: products, orders: orders, wallets: wallets}
}

// CreateProductRequest holds input for creating a product.
type CreateProductRequest struct {
	SellerID    uuid.UUID
	Name        string
	Description *string
	PriceUSD    decimal.Decimal
	Category    *string
	Stock       int
}

// CreateProduct creates a new product.
func (s *ProductService) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {
	if req.Name == "" {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeBadRequest, "product name is required", 400)
	}
	if req.PriceUSD.LessThanOrEqual(decimal.Zero) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeBadRequest, "price must be greater than zero", 400)
	}
	product := &models.Product{
		SellerID:    req.SellerID,
		Name:        req.Name,
		Description: req.Description,
		PriceUSD:    req.PriceUSD,
		Category:    req.Category,
		Stock:       req.Stock,
	}
	if err := s.products.Create(ctx, product); err != nil {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to create product", 500)
	}
	return product, nil
}

// ListProducts returns a paginated list of products.
func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.products.List(ctx, limit, offset)
}

// GetProduct returns a product by ID.
func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	p, err := s.products.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeProductNotFound, "product not found", 404)
		}
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to get product", 500)
	}
	return p, nil
}

// CreateOrderRequest holds input for creating an order.
type CreateOrderRequest struct {
	BuyerID      uuid.UUID
	ProductID    uuid.UUID
	Quantity     int
	TKDiscount   decimal.Decimal // amount of TK tokens to use as discount
}

// CreateOrder places an order, applying optional TK discount and updating wallet.
func (s *ProductService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*models.Order, error) {
	product, err := s.products.GetByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeProductNotFound, "product not found", 404)
		}
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to get product", 500)
	}
	if !product.IsInStock(req.Quantity) {
		return nil, apierrors.NewAPIError("PRODUCT_OUT_OF_STOCK", "product is not in stock", 400)
	}

	qty := decimal.NewFromInt(int64(req.Quantity))
	total := product.PriceUSD.Mul(qty)

	// Apply TK discount (max 50%)
	maxDiscount := total.Div(decimal.NewFromInt(2))
	tkUsed := req.TKDiscount
	if tkUsed.GreaterThan(maxDiscount) {
		tkUsed = maxDiscount
	}
	// TK discount converts 1:1 to USD for simplicity
	finalPrice := total.Sub(tkUsed)
	if finalPrice.LessThan(decimal.Zero) {
		finalPrice = decimal.Zero
	}

	// Deduct TK from buyer wallet if discount is used
	if tkUsed.GreaterThan(decimal.Zero) {
		wallet, wErr := s.wallets.GetByUserID(ctx, req.BuyerID)
		if wErr != nil {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to get wallet", 500)
		}
		if !wallet.HasSufficientBalance(tkUsed) {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeInsufficientBalance, "insufficient TK balance", 400)
		}
		newBalance := wallet.BalanceTK.Sub(tkUsed)
		if err := s.wallets.UpdateBalance(ctx, wallet.WalletID, newBalance, wallet.FrozenBalanceTK); err != nil {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to update wallet", 500)
		}
		desc := "TK discount on order"
		_ = s.wallets.CreateTransaction(ctx, &models.Transaction{
			UserID:          req.BuyerID,
			TransactionType: models.TransactionPurchaseDiscount,
			AmountTK:        tkUsed,
			Description:     &desc,
			Status:          models.TransactionStatusCompleted,
		})
	}

	order := &models.Order{
		BuyerID:        req.BuyerID,
		SellerID:       product.SellerID,
		TotalPriceUSD:  total,
		TKDiscountUsed: tkUsed,
		FinalPriceUSD:  finalPrice,
		Status:         models.OrderStatusPending,
	}
	items := []*models.OrderItem{
		{
			ProductID:    product.ProductID,
			Quantity:     req.Quantity,
			PricePerUnit: product.PriceUSD,
		},
	}
	if err := s.orders.Create(ctx, order, items); err != nil {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInternalError, "failed to create order", 500)
	}

	// Reduce product stock
	product.Stock -= req.Quantity
	_ = s.products.Update(ctx, product)

	return order, nil
}

// ListOrders returns orders for a buyer.
func (s *ProductService) ListOrders(ctx context.Context, buyerID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.orders.ListByBuyer(ctx, buyerID, limit, offset)
}
