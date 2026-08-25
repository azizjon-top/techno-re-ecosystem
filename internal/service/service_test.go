package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/azizjon-top/techno-re-ecosystem/internal/models"
	"github.com/azizjon-top/techno-re-ecosystem/internal/repository"
	"github.com/azizjon-top/techno-re-ecosystem/internal/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func newTestJWT() *service.JWTService {
	return service.NewJWTService("test-secret-key-1234", 15*time.Minute, 7*24*time.Hour, "test-issuer", "test-audience")
}

func TestJWTService_GenerateAndValidate(t *testing.T) {
	svc := newTestJWT()
	userID := uuid.New()

	token, err := svc.GenerateAccessToken(userID, "user")
	if err != nil {
		t.Fatalf("GenerateAccessToken error: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, claims.UserID)
	}
	if claims.Role != "user" {
		t.Errorf("expected role 'user', got %q", claims.Role)
	}
}

func TestJWTService_InvalidToken(t *testing.T) {
	svc := newTestJWT()
	_, err := svc.ValidateToken("invalid.token.here")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestUserService_RegisterAndLogin(t *testing.T) {
	userRepo := repository.NewInMemoryUserRepository()
	walletRepo := repository.NewInMemoryWalletRepository()
	jwtSvc := newTestJWT()
	svc := service.NewUserService(userRepo, walletRepo, jwtSvc)
	ctx := context.Background()

	user, tokens, err := svc.Register(ctx, service.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %q", user.Email)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Error("expected non-empty tokens")
	}

	// Duplicate registration should fail
	_, _, err = svc.Register(ctx, service.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected conflict error for duplicate registration")
	}

	// Login with correct credentials
	_, loginTokens, err := svc.Login(ctx, service.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	if loginTokens.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	// Login with wrong password
	_, _, err = svc.Login(ctx, service.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestWalletService_Transfer(t *testing.T) {
	walletRepo := repository.NewInMemoryWalletRepository()
	ctx := context.Background()

	senderID := uuid.New()
	receiverID := uuid.New()

	senderWallet := &models.Wallet{UserID: senderID, BalanceTK: decimal.NewFromInt(100)}
	receiverWallet := &models.Wallet{UserID: receiverID, BalanceTK: decimal.Zero}
	_ = walletRepo.Create(ctx, senderWallet)
	_ = walletRepo.Create(ctx, receiverWallet)

	walletSvc := service.NewWalletService(walletRepo)

	// Successful transfer
	err := walletSvc.Transfer(ctx, service.TransferRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     decimal.NewFromInt(30),
	})
	if err != nil {
		t.Fatalf("Transfer error: %v", err)
	}

	sender, _ := walletRepo.GetByUserID(ctx, senderID)
	receiver, _ := walletRepo.GetByUserID(ctx, receiverID)
	if !sender.BalanceTK.Equal(decimal.NewFromInt(70)) {
		t.Errorf("expected sender balance 70, got %v", sender.BalanceTK)
	}
	if !receiver.BalanceTK.Equal(decimal.NewFromInt(30)) {
		t.Errorf("expected receiver balance 30, got %v", receiver.BalanceTK)
	}

	// Insufficient balance
	err = walletSvc.Transfer(ctx, service.TransferRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     decimal.NewFromInt(1000),
	})
	if err == nil {
		t.Fatal("expected insufficient balance error")
	}
}

func TestProductService_CreateAndList(t *testing.T) {
	productRepo := repository.NewInMemoryProductRepository()
	orderRepo := repository.NewInMemoryOrderRepository()
	walletRepo := repository.NewInMemoryWalletRepository()
	svc := service.NewProductService(productRepo, orderRepo, walletRepo)
	ctx := context.Background()

	sellerID := uuid.New()
	desc := "A great product"
	product, err := svc.CreateProduct(ctx, service.CreateProductRequest{
		SellerID:    sellerID,
		Name:        "Widget",
		Description: &desc,
		PriceUSD:    decimal.NewFromFloat(9.99),
		Stock:       10,
	})
	if err != nil {
		t.Fatalf("CreateProduct error: %v", err)
	}
	if product.Name != "Widget" {
		t.Errorf("expected name Widget, got %q", product.Name)
	}

	products, err := svc.ListProducts(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListProducts error: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("expected 1 product, got %d", len(products))
	}
}

func TestProductService_CreateOrder_WithTKDiscount(t *testing.T) {
	productRepo := repository.NewInMemoryProductRepository()
	orderRepo := repository.NewInMemoryOrderRepository()
	walletRepo := repository.NewInMemoryWalletRepository()
	svc := service.NewProductService(productRepo, orderRepo, walletRepo)
	ctx := context.Background()

	sellerID := uuid.New()
	buyerID := uuid.New()

	// Create buyer wallet with 50 TK
	bWallet := &models.Wallet{UserID: buyerID, BalanceTK: decimal.NewFromInt(50)}
	_ = walletRepo.Create(ctx, bWallet)

	// Create a product ($20, 5 in stock)
	product, _ := svc.CreateProduct(ctx, service.CreateProductRequest{
		SellerID: sellerID,
		Name:     "Gadget",
		PriceUSD: decimal.NewFromInt(20),
		Stock:    5,
	})

	// Order 1 unit with 5 TK discount (should reduce price to $15)
	order, err := svc.CreateOrder(ctx, service.CreateOrderRequest{
		BuyerID:    buyerID,
		ProductID:  product.ProductID,
		Quantity:   1,
		TKDiscount: decimal.NewFromInt(5),
	})
	if err != nil {
		t.Fatalf("CreateOrder error: %v", err)
	}
	if !order.FinalPriceUSD.Equal(decimal.NewFromInt(15)) {
		t.Errorf("expected final price 15, got %v", order.FinalPriceUSD)
	}

	// Buyer wallet should be reduced by 5 TK
	wallet, _ := walletRepo.GetByUserID(ctx, buyerID)
	if !wallet.BalanceTK.Equal(decimal.NewFromInt(45)) {
		t.Errorf("expected buyer wallet 45 TK, got %v", wallet.BalanceTK)
	}
}
