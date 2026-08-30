package app

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/azizjon-top/techno-re-ecosystem/internal/config"
	apierrors "github.com/azizjon-top/techno-re-ecosystem/internal/errors"
	"github.com/azizjon-top/techno-re-ecosystem/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type App struct {
	mu sync.RWMutex

	cfg *config.Config

	users        map[uuid.UUID]*models.User
	usersByEmail map[string]uuid.UUID
	wallets      map[uuid.UUID]*models.Wallet
	products     map[uuid.UUID]*models.Product
	orders       map[uuid.UUID]*OrderRecord
	chats        map[uuid.UUID]*models.Chat
	messages     map[uuid.UUID][]*models.Message
	channels     map[uuid.UUID]*models.Channel
	videos       map[uuid.UUID]*models.Video
	mining       map[uuid.UUID]*models.MiningSession
	aiRequests   map[uuid.UUID]*models.AIRequest
	validations  map[uuid.UUID][]models.ConsensusValidation
	campaigns    map[uuid.UUID]*models.AdCampaign
	transactions map[uuid.UUID]*models.Transaction
	rewardedAI   map[uuid.UUID]bool
}

type OrderRecord struct {
	Order *models.Order       `json:"order"`
	Items []models.OrderItem  `json:"items"`
}

type ChatView struct {
	Chat     *models.Chat   `json:"chat"`
	Messages []MessageView  `json:"messages"`
}

type MessageView struct {
	MessageID        uuid.UUID        `json:"message_id"`
	ChatID           uuid.UUID        `json:"chat_id"`
	SenderID         uuid.UUID        `json:"sender_id"`
	Text             string           `json:"text"`
	TKTransferAmount *decimal.Decimal `json:"tk_transfer_amount,omitempty"`
	IsRead           bool             `json:"is_read"`
	CreatedAt        time.Time        `json:"created_at"`
}

type WalletView struct {
	Wallet           *models.Wallet   `json:"wallet"`
	AvailableBalance decimal.Decimal  `json:"available_balance"`
}

type ProductInput struct {
	SellerID     uuid.UUID
	Name         string
	Description  *string
	PriceUSD     decimal.Decimal
	Category     *string
	Stock        int
	ImageURL     *string
	AISearchTags []string
}

type OrderItemInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type VideoInput struct {
	OwnerID         uuid.UUID
	ChannelID       *uuid.UUID
	ChannelName     *string
	Title           string
	Description     *string
	FilePath        string
	FileHash        *string
	DurationSeconds *int
	P2PEnabled      bool
}

type ValidateFactInput struct {
	UserID           uuid.UUID
	AIRequestID      *uuid.UUID
	RequestType      models.AIRequestType
	QueryText        string
	FactClaim        string
	ValidationResult models.ValidationResult
	ConfidenceScore  decimal.Decimal
}

type CampaignInput struct {
	CreatorID   uuid.UUID
	Title       string
	Description *string
	BannerURL   *string
	StartDate   time.Time
	EndDate     time.Time
	BudgetTK    decimal.Decimal
}

type Stats struct {
	Users          int `json:"users"`
	Products       int `json:"products"`
	Orders         int `json:"orders"`
	Chats          int `json:"chats"`
	Videos         int `json:"videos"`
	MiningSessions int `json:"mining_sessions"`
	Campaigns      int `json:"campaigns"`
}

func New(cfg *config.Config) (*App, error) {
	a := &App{
		cfg:          cfg,
		users:        make(map[uuid.UUID]*models.User),
		usersByEmail: make(map[string]uuid.UUID),
		wallets:      make(map[uuid.UUID]*models.Wallet),
		products:     make(map[uuid.UUID]*models.Product),
		orders:       make(map[uuid.UUID]*OrderRecord),
		chats:        make(map[uuid.UUID]*models.Chat),
		messages:     make(map[uuid.UUID][]*models.Message),
		channels:     make(map[uuid.UUID]*models.Channel),
		videos:       make(map[uuid.UUID]*models.Video),
		mining:       make(map[uuid.UUID]*models.MiningSession),
		aiRequests:   make(map[uuid.UUID]*models.AIRequest),
		validations:  make(map[uuid.UUID][]models.ConsensusValidation),
		campaigns:    make(map[uuid.UUID]*models.AdCampaign),
		transactions: make(map[uuid.UUID]*models.Transaction),
		rewardedAI:   make(map[uuid.UUID]bool),
	}

	if err := a.seed(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) seed() error {
	user, err := a.seedUser("user1@techno.re", "+1234567890", "user12345", models.RoleUser, true, 5)
	if err != nil {
		return err
	}
	seller, err := a.seedUser("seller1@techno.re", "+1234567891", "seller12345", models.RoleSeller, false, 4)
	if err != nil {
		return err
	}
	admin, err := a.seedUser("admin@techno.re", "+1234567892", "admin12345", models.RoleAdmin, false, 1)
	if err != nil {
		return err
	}

	a.wallets[user.UserID] = &models.Wallet{
		WalletID:        uuid.New(),
		UserID:          user.UserID,
		BalanceTK:       decimal.RequireFromString("1000"),
		FrozenBalanceTK: decimal.Zero,
		LastUpdated:     time.Now().UTC(),
	}
	a.wallets[seller.UserID] = &models.Wallet{
		WalletID:        uuid.New(),
		UserID:          seller.UserID,
		BalanceTK:       decimal.RequireFromString("500"),
		FrozenBalanceTK: decimal.Zero,
		LastUpdated:     time.Now().UTC(),
	}
	a.wallets[admin.UserID] = &models.Wallet{
		WalletID:        uuid.New(),
		UserID:          admin.UserID,
		BalanceTK:       decimal.Zero,
		FrozenBalanceTK: decimal.Zero,
		LastUpdated:     time.Now().UTC(),
	}

	channel := &models.Channel{
		ChannelID: uuid.New(),
		OwnerID:   seller.UserID,
		Name:      "Techno Seller Channel",
		CreatedAt: time.Now().UTC(),
	}
	a.channels[channel.ChannelID] = channel

	description := "Starter product for the ecosystem MVP"
	category := "electronics"
	imageURL := "https://example.com/images/techno-router.png"
	product := &models.Product{
		ProductID:    uuid.New(),
		SellerID:     seller.UserID,
		Name:         "Techno Router",
		Description:  &description,
		PriceUSD:     decimal.RequireFromString("249.99"),
		Category:     &category,
		Stock:        15,
		ImageURL:     &imageURL,
		AISearchTags: []string{"network", "router", "smart-home"},
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	a.products[product.ProductID] = product

	videoDescription := "Launch overview for the Techno Router"
	duration := 420
	fileHash := "demo-file-hash"
	video := &models.Video{
		VideoID:         uuid.New(),
		ChannelID:       channel.ChannelID,
		Title:           "Techno Router Demo",
		Description:     &videoDescription,
		FilePath:        "/videos/techno-router-demo.mp4",
		FileHash:        &fileHash,
		DurationSeconds: &duration,
		ViewCount:       128,
		P2PEnabled:      true,
		CreatedAt:       time.Now().UTC(),
	}
	a.videos[video.VideoID] = video

	chat := &models.Chat{
		ChatID:         uuid.New(),
		Participant1ID: user.UserID,
		Participant2ID: seller.UserID,
		CreatedAt:      time.Now().UTC(),
	}
	if chat.Participant1ID.String() > chat.Participant2ID.String() {
		chat.Participant1ID, chat.Participant2ID = chat.Participant2ID, chat.Participant1ID
	}
	a.chats[chat.ChatID] = chat

	encrypted, err := models.EncryptMessage("Welcome to Techno RE!", nil)
	if err != nil {
		return err
	}
	a.messages[chat.ChatID] = []*models.Message{
		{
			MessageID:     uuid.New(),
			ChatID:        chat.ChatID,
			SenderID:      seller.UserID,
			EncryptedText: encrypted,
			IsRead:        false,
			CreatedAt:     time.Now().UTC(),
		},
	}

	campaignDescription := "Awareness campaign for the Techno Router"
	bannerURL := "https://example.com/banners/router-campaign.png"
	a.campaigns[uuid.New()] = &models.AdCampaign{
		CampaignID:  uuid.New(),
		CreatorID:   seller.UserID,
		Title:       "Router Launch Campaign",
		Description: &campaignDescription,
		BannerURL:   &bannerURL,
		StartDate:   time.Now().UTC().Add(-24 * time.Hour),
		EndDate:     time.Now().UTC().Add(7 * 24 * time.Hour),
		BudgetTK:    decimal.RequireFromString("250"),
		SpentTK:     decimal.RequireFromString("25"),
		Status:      models.CampaignStatusActive,
		CreatedAt:   time.Now().UTC(),
	}

	return nil
}

func (a *App) seedUser(email, phone, password string, role models.UserRole, miningEnabled bool, cpuLimit int) (*models.User, error) {
	now := time.Now().UTC()
	user := &models.User{
		UserID:           uuid.New(),
		Email:            email,
		Phone:            stringPtr(phone),
		Role:             role,
		MiningEnabled:    miningEnabled,
		MiningCPULimit:   cpuLimit,
		SecuritySettings: map[string]interface{}{"mfa_enabled": false},
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	hash, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hash
	a.users[user.UserID] = user
	a.usersByEmail[strings.ToLower(user.Email)] = user.UserID
	return user, nil
}

func (a *App) Stats() Stats {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return Stats{
		Users:          len(a.users),
		Products:       len(a.products),
		Orders:         len(a.orders),
		Chats:          len(a.chats),
		Videos:         len(a.videos),
		MiningSessions: len(a.mining),
		Campaigns:      len(a.campaigns),
	}
}

func (a *App) RegisterUser(email string, password string, phone *string, role models.UserRole) (*models.User, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "email and password are required", 400)
	}
	if _, exists := a.usersByEmail[email]; exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeConflict, "user already exists", 409)
	}
	if role == "" {
		role = models.RoleUser
	}

	user := &models.User{
		UserID:           uuid.New(),
		Email:            email,
		Phone:            phone,
		Role:             role,
		MiningEnabled:    role != models.RoleAdmin,
		MiningCPULimit:   5,
		SecuritySettings: map[string]interface{}{"mfa_enabled": false},
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	if !user.IsValidRole() {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "invalid role", 400)
	}
	hash, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hash
	a.users[user.UserID] = user
	a.usersByEmail[email] = user.UserID
	a.wallets[user.UserID] = &models.Wallet{
		WalletID:        uuid.New(),
		UserID:          user.UserID,
		BalanceTK:       decimal.Zero,
		FrozenBalanceTK: decimal.Zero,
		LastUpdated:     time.Now().UTC(),
	}

	return user, nil
}

func (a *App) Authenticate(email string, password string) (string, *models.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	userID, exists := a.usersByEmail[strings.ToLower(strings.TrimSpace(email))]
	if !exists {
		return "", nil, apierrors.NewAPIError(apierrors.ErrCodeInvalidCredentials, "invalid credentials", 401)
	}

	user := a.users[userID]
	if !user.VerifyPassword(password) {
		return "", nil, apierrors.NewAPIError(apierrors.ErrCodeInvalidCredentials, "invalid credentials", 401)
	}

	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.UserID.String(),
		"email": user.Email,
		"role":  string(user.Role),
		"iss":   a.cfg.JWT.TokenIssuer,
		"aud":   a.cfg.JWT.TokenAudience,
		"iat":   now.Unix(),
		"exp":   now.Add(a.cfg.JWT.AccessTokenTTL).Unix(),
	})

	signed, err := token.SignedString([]byte(a.cfg.JWT.Secret))
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}

func (a *App) GetUser(userID uuid.UUID) (*models.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	user, exists := a.users[userID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "user not found", 404)
	}
	return user, nil
}

func (a *App) ListProducts() []*models.Product {
	a.mu.RLock()
	defer a.mu.RUnlock()

	products := make([]*models.Product, 0, len(a.products))
	for _, product := range a.products {
		products = append(products, product)
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].CreatedAt.Before(products[j].CreatedAt)
	})
	return products
}

func (a *App) CreateProduct(input ProductInput) (*models.Product, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user, exists := a.users[input.SellerID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "seller not found", 404)
	}
	if user.Role != models.RoleSeller && user.Role != models.RoleAdmin {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeForbidden, "user cannot create products", 403)
	}
	if input.Name == "" || input.Stock < 0 || input.PriceUSD.LessThanOrEqual(decimal.Zero) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "invalid product payload", 400)
	}

	now := time.Now().UTC()
	product := &models.Product{
		ProductID:    uuid.New(),
		SellerID:     input.SellerID,
		Name:         input.Name,
		Description:  input.Description,
		PriceUSD:     input.PriceUSD,
		Category:     input.Category,
		Stock:        input.Stock,
		ImageURL:     input.ImageURL,
		AISearchTags: input.AISearchTags,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	a.products[product.ProductID] = product
	return product, nil
}

func (a *App) ListOrders(userID *uuid.UUID) []OrderRecord {
	a.mu.RLock()
	defer a.mu.RUnlock()

	records := make([]OrderRecord, 0, len(a.orders))
	for _, record := range a.orders {
		if userID != nil && record.Order.BuyerID != *userID && record.Order.SellerID != *userID {
			continue
		}
		records = append(records, OrderRecord{Order: record.Order, Items: append([]models.OrderItem(nil), record.Items...)})
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Order.CreatedAt.Before(records[j].Order.CreatedAt)
	})
	return records
}

func (a *App) CreateOrder(buyerID uuid.UUID, items []OrderItemInput, tkDiscount decimal.Decimal) (*OrderRecord, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.users[buyerID]; !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "buyer not found", 404)
	}
	if len(items) == 0 {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "order items are required", 400)
	}

	total := decimal.Zero
	var sellerID uuid.UUID
	orderItems := make([]models.OrderItem, 0, len(items))
	for _, item := range items {
		product, exists := a.products[item.ProductID]
		if !exists {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeProductNotFound, "product not found", 404)
		}
		if item.Quantity <= 0 || !product.IsInStock(item.Quantity) {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "product stock is insufficient", 400)
		}
		if sellerID == uuid.Nil {
			sellerID = product.SellerID
		}
		if sellerID != product.SellerID {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "all items must belong to one seller", 400)
		}

		lineTotal := product.PriceUSD.Mul(decimal.NewFromInt(int64(item.Quantity)))
		total = total.Add(lineTotal)
		orderItems = append(orderItems, models.OrderItem{
			ItemID:       uuid.New(),
			ProductID:    product.ProductID,
			Quantity:     item.Quantity,
			PricePerUnit: product.PriceUSD,
		})
	}

	wallet, exists := a.wallets[buyerID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "wallet not found", 404)
	}
	maxDiscount := wallet.CanApplyTKDiscount(total)
	if tkDiscount.GreaterThan(maxDiscount) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInsufficientBalance, "tk discount exceeds available or allowed balance", 400)
	}

	now := time.Now().UTC()
	order := &models.Order{
		OrderID:        uuid.New(),
		BuyerID:        buyerID,
		SellerID:       sellerID,
		TotalPriceUSD:  total,
		TKDiscountUsed: tkDiscount,
		FinalPriceUSD:  total.Sub(tkDiscount),
		Status:         models.OrderStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	for index := range orderItems {
		orderItems[index].OrderID = order.OrderID
		product := a.products[orderItems[index].ProductID]
		product.Stock -= orderItems[index].Quantity
		product.UpdatedAt = now
	}

	if tkDiscount.GreaterThan(decimal.Zero) {
		wallet.BalanceTK = wallet.BalanceTK.Sub(tkDiscount)
		wallet.LastUpdated = now
		a.addTransactionLocked(buyerID, models.TransactionPurchaseDiscount, tkDiscount, stringPtr("Applied TK discount to order"), &order.OrderID, models.TransactionStatusCompleted)
	}

	record := &OrderRecord{Order: order, Items: orderItems}
	a.orders[order.OrderID] = record

	return record, nil
}

func (a *App) GetWallet(userID uuid.UUID) (*WalletView, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	wallet, exists := a.wallets[userID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "wallet not found", 404)
	}

	return &WalletView{
		Wallet:           wallet,
		AvailableBalance: wallet.AvailableBalance(),
	}, nil
}

func (a *App) TransferTokens(fromUserID uuid.UUID, toUserID uuid.UUID, amount decimal.Decimal, description *string) ([]*models.Transaction, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "amount must be greater than zero", 400)
	}
	fromWallet, exists := a.wallets[fromUserID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "sender wallet not found", 404)
	}
	toWallet, exists := a.wallets[toUserID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "recipient wallet not found", 404)
	}
	if !fromWallet.HasSufficientBalance(amount) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInsufficientBalance, "insufficient balance", 400)
	}

	now := time.Now().UTC()
	fromWallet.BalanceTK = fromWallet.BalanceTK.Sub(amount)
	fromWallet.LastUpdated = now
	toWallet.BalanceTK = toWallet.BalanceTK.Add(amount)
	toWallet.LastUpdated = now

	sent := a.addTransactionLocked(fromUserID, models.TransactionTransferSent, amount, description, nil, models.TransactionStatusCompleted)
	received := a.addTransactionLocked(toUserID, models.TransactionTransferReceived, amount, description, nil, models.TransactionStatusCompleted)
	return []*models.Transaction{sent, received}, nil
}

func (a *App) ListChats(userID *uuid.UUID) ([]ChatView, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	chats := make([]ChatView, 0, len(a.chats))
	for _, chat := range a.chats {
		if userID != nil && chat.Participant1ID != *userID && chat.Participant2ID != *userID {
			continue
		}
		messages := make([]MessageView, 0, len(a.messages[chat.ChatID]))
		for _, message := range a.messages[chat.ChatID] {
			text, err := models.DecryptMessage(message.EncryptedText, nil)
			if err != nil {
				return nil, err
			}
			messages = append(messages, MessageView{
				MessageID:        message.MessageID,
				ChatID:           message.ChatID,
				SenderID:         message.SenderID,
				Text:             text,
				TKTransferAmount: message.TKTransferAmount,
				IsRead:           message.IsRead,
				CreatedAt:        message.CreatedAt,
			})
		}
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].CreatedAt.Before(messages[j].CreatedAt)
		})
		chats = append(chats, ChatView{Chat: chat, Messages: messages})
	}
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].Chat.CreatedAt.Before(chats[j].Chat.CreatedAt)
	})
	return chats, nil
}

func (a *App) PostMessage(chatID uuid.UUID, senderID uuid.UUID, text string, transferAmount *decimal.Decimal) (*MessageView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	chat, exists := a.chats[chatID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "chat not found", 404)
	}
	if chat.Participant1ID != senderID && chat.Participant2ID != senderID {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeForbidden, "sender is not part of this chat", 403)
	}
	if strings.TrimSpace(text) == "" {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "message text is required", 400)
	}

	if transferAmount != nil && transferAmount.GreaterThan(decimal.Zero) {
		recipientID := chat.Participant1ID
		if recipientID == senderID {
			recipientID = chat.Participant2ID
		}
		if _, err := a.transferTokensLocked(senderID, recipientID, *transferAmount, stringPtr("Chat TK transfer")); err != nil {
			return nil, err
		}
	}

	encrypted, err := models.EncryptMessage(text, nil)
	if err != nil {
		return nil, err
	}
	message := &models.Message{
		MessageID:        uuid.New(),
		ChatID:           chatID,
		SenderID:         senderID,
		EncryptedText:    encrypted,
		TKTransferAmount: transferAmount,
		IsRead:           false,
		CreatedAt:        time.Now().UTC(),
	}
	a.messages[chatID] = append(a.messages[chatID], message)

	return &MessageView{
		MessageID:        message.MessageID,
		ChatID:           message.ChatID,
		SenderID:         message.SenderID,
		Text:             text,
		TKTransferAmount: message.TKTransferAmount,
		IsRead:           message.IsRead,
		CreatedAt:        message.CreatedAt,
	}, nil
}

func (a *App) ListVideos(channelID *uuid.UUID) []*models.Video {
	a.mu.RLock()
	defer a.mu.RUnlock()

	videos := make([]*models.Video, 0, len(a.videos))
	for _, video := range a.videos {
		if channelID != nil && video.ChannelID != *channelID {
			continue
		}
		videos = append(videos, video)
	}
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].CreatedAt.Before(videos[j].CreatedAt)
	})
	return videos
}

func (a *App) UploadVideo(input VideoInput) (*models.Video, *models.Channel, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.users[input.OwnerID]; !exists {
		return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "owner not found", 404)
	}
	if input.Title == "" || input.FilePath == "" {
		return nil, nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "title and file_path are required", 400)
	}

	channel, err := a.resolveChannelLocked(input.OwnerID, input.ChannelID, input.ChannelName)
	if err != nil {
		return nil, nil, err
	}

	video := &models.Video{
		VideoID:         uuid.New(),
		ChannelID:       channel.ChannelID,
		Title:           input.Title,
		Description:     input.Description,
		FilePath:        input.FilePath,
		FileHash:        input.FileHash,
		DurationSeconds: input.DurationSeconds,
		ViewCount:       0,
		P2PEnabled:      input.P2PEnabled,
		CreatedAt:       time.Now().UTC(),
	}
	a.videos[video.VideoID] = video
	return video, channel, nil
}

func (a *App) StartMiningSession(userID uuid.UUID) (*models.MiningSession, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user, exists := a.users[userID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "user not found", 404)
	}
	if !user.MiningEnabled || !user.CanEnableMining() {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeForbidden, "mining is not enabled for user", 403)
	}
	for _, session := range a.mining {
		if session.UserID == userID && session.IsActive() {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeConflict, "active mining session already exists", 409)
		}
	}

	now := time.Now().UTC()
	session := &models.MiningSession{
		SessionID:       uuid.New(),
		UserID:          userID,
		StartedAt:       now,
		DurationMinutes: 0,
		FactsVerified:   0,
		RewardTKEarned:  decimal.Zero,
		Status:          models.MiningStatusActive,
		CreatedAt:       now,
	}
	a.mining[session.SessionID] = session
	return session, nil
}

func (a *App) ValidateFact(input ValidateFactInput) (*models.ConsensusValidation, decimal.Decimal, bool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.users[input.UserID]; !exists {
		return nil, decimal.Zero, false, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "validator not found", 404)
	}
	if input.FactClaim == "" {
		return nil, decimal.Zero, false, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "fact_claim is required", 400)
	}
	if input.ConfidenceScore.LessThan(decimal.Zero) || input.ConfidenceScore.GreaterThan(decimal.NewFromInt(1)) {
		return nil, decimal.Zero, false, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "confidence_score must be between 0 and 1", 400)
	}

	requestID := uuid.Nil
	if input.AIRequestID != nil {
		requestID = *input.AIRequestID
		if _, exists := a.aiRequests[requestID]; !exists {
			return nil, decimal.Zero, false, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "ai request not found", 404)
		}
	} else {
		if input.RequestType == "" {
			input.RequestType = models.AIRequestContentModeration
		}
		if input.QueryText == "" {
			input.QueryText = input.FactClaim
		}
		requestID = uuid.New()
		a.aiRequests[requestID] = &models.AIRequest{
			RequestID:   requestID,
			UserID:      input.UserID,
			RequestType: input.RequestType,
			QueryText:   input.QueryText,
			CreatedAt:   time.Now().UTC(),
		}
	}

	validation := models.ConsensusValidation{
		ValidationID:     uuid.New(),
		AIRequestID:      requestID,
		ValidatorNodeID:  input.UserID,
		FactClaim:        input.FactClaim,
		ValidationResult: input.ValidationResult,
		ConfidenceScore:  input.ConfidenceScore,
		CreatedAt:        time.Now().UTC(),
	}
	a.validations[requestID] = append(a.validations[requestID], validation)

	average := models.AverageConfidence(a.validations[requestID])
	consensusReached := models.IsConsensusReached(a.validations[requestID])
	if consensusReached && !a.rewardedAI[requestID] {
		if err := a.rewardConsensusLocked(input.UserID, requestID); err != nil {
			return nil, decimal.Zero, false, err
		}
	}

	return &validation, average, consensusReached, nil
}

func (a *App) ListCampaigns(creatorID *uuid.UUID) []*models.AdCampaign {
	a.mu.RLock()
	defer a.mu.RUnlock()

	campaigns := make([]*models.AdCampaign, 0, len(a.campaigns))
	for _, campaign := range a.campaigns {
		if creatorID != nil && campaign.CreatorID != *creatorID {
			continue
		}
		campaigns = append(campaigns, campaign)
	}
	sort.Slice(campaigns, func(i, j int) bool {
		return campaigns[i].CreatedAt.Before(campaigns[j].CreatedAt)
	})
	return campaigns
}

func (a *App) CreateCampaign(input CampaignInput) (*models.AdCampaign, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user, exists := a.users[input.CreatorID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeUserNotFound, "creator not found", 404)
	}
	if user.Role != models.RoleSeller && user.Role != models.RoleAdmin {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeForbidden, "user cannot create campaigns", 403)
	}
	if input.Title == "" || input.BudgetTK.LessThanOrEqual(decimal.Zero) || input.EndDate.Before(input.StartDate) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "invalid campaign payload", 400)
	}

	campaign := &models.AdCampaign{
		CampaignID:  uuid.New(),
		CreatorID:   input.CreatorID,
		Title:       input.Title,
		Description: input.Description,
		BannerURL:   input.BannerURL,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		BudgetTK:    input.BudgetTK,
		SpentTK:     decimal.Zero,
		Status:      models.CampaignStatusActive,
		CreatedAt:   time.Now().UTC(),
	}
	a.campaigns[campaign.CampaignID] = campaign
	return campaign, nil
}

func (a *App) transferTokensLocked(fromUserID uuid.UUID, toUserID uuid.UUID, amount decimal.Decimal, description *string) ([]*models.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeValidationFailed, "amount must be greater than zero", 400)
	}
	fromWallet, exists := a.wallets[fromUserID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "sender wallet not found", 404)
	}
	toWallet, exists := a.wallets[toUserID]
	if !exists {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "recipient wallet not found", 404)
	}
	if !fromWallet.HasSufficientBalance(amount) {
		return nil, apierrors.NewAPIError(apierrors.ErrCodeInsufficientBalance, "insufficient balance", 400)
	}

	now := time.Now().UTC()
	fromWallet.BalanceTK = fromWallet.BalanceTK.Sub(amount)
	fromWallet.LastUpdated = now
	toWallet.BalanceTK = toWallet.BalanceTK.Add(amount)
	toWallet.LastUpdated = now

	sent := a.addTransactionLocked(fromUserID, models.TransactionTransferSent, amount, description, nil, models.TransactionStatusCompleted)
	received := a.addTransactionLocked(toUserID, models.TransactionTransferReceived, amount, description, nil, models.TransactionStatusCompleted)
	return []*models.Transaction{sent, received}, nil
}

func (a *App) addTransactionLocked(userID uuid.UUID, txType models.TransactionType, amount decimal.Decimal, description *string, relatedEntityID *uuid.UUID, status models.TransactionStatus) *models.Transaction {
	tx := &models.Transaction{
		TransactionID:   uuid.New(),
		UserID:          userID,
		TransactionType: txType,
		AmountTK:        amount,
		Description:     description,
		RelatedEntityID: relatedEntityID,
		Status:          status,
		CreatedAt:       time.Now().UTC(),
	}
	a.transactions[tx.TransactionID] = tx
	return tx
}

func (a *App) resolveChannelLocked(ownerID uuid.UUID, channelID *uuid.UUID, channelName *string) (*models.Channel, error) {
	if channelID != nil {
		channel, exists := a.channels[*channelID]
		if !exists {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeNotFound, "channel not found", 404)
		}
		if channel.OwnerID != ownerID {
			return nil, apierrors.NewAPIError(apierrors.ErrCodeForbidden, "channel does not belong to owner", 403)
		}
		return channel, nil
	}

	name := "Default Channel"
	if channelName != nil && strings.TrimSpace(*channelName) != "" {
		name = strings.TrimSpace(*channelName)
	} else {
		for _, channel := range a.channels {
			if channel.OwnerID == ownerID {
				return channel, nil
			}
		}
	}

	channel := &models.Channel{
		ChannelID: uuid.New(),
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
	a.channels[channel.ChannelID] = channel
	return channel, nil
}

func (a *App) rewardConsensusLocked(userID uuid.UUID, requestID uuid.UUID) error {
	reward, err := decimal.NewFromString(a.cfg.Mining.RewardPerFact)
	if err != nil {
		return fmt.Errorf("parse mining reward: %w", err)
	}

	wallet, exists := a.wallets[userID]
	if !exists {
		return apierrors.NewAPIError(apierrors.ErrCodeNotFound, "wallet not found", 404)
	}
	wallet.BalanceTK = wallet.BalanceTK.Add(reward)
	wallet.LastUpdated = time.Now().UTC()
	a.addTransactionLocked(userID, models.TransactionConsensusValidationReward, reward, stringPtr("Consensus validation reward"), &requestID, models.TransactionStatusCompleted)

	for _, session := range a.mining {
		if session.UserID == userID && session.IsActive() {
			session.FactsVerified++
			session.RewardTKEarned = session.RewardTKEarned.Add(reward)
			session.DurationMinutes = int(time.Since(session.StartedAt).Minutes())
			break
		}
	}
	a.rewardedAI[requestID] = true
	return nil
}

func stringPtr(value string) *string {
	return &value
}
