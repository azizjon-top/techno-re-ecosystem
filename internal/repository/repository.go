package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/azizjon-top/techno-re-ecosystem/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("record not found")

// ErrConflict is returned when a unique constraint is violated.
var ErrConflict = errors.New("record already exists")

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

// ProductRepository defines the interface for product persistence.
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	List(ctx context.Context, limit, offset int) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
}

// OrderRepository defines the interface for order persistence.
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order, items []*models.OrderItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, []*models.OrderItem, error)
	ListByBuyer(ctx context.Context, buyerID uuid.UUID, limit, offset int) ([]*models.Order, error)
}

// WalletRepository defines the interface for wallet persistence.
type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Wallet, error)
	Create(ctx context.Context, wallet *models.Wallet) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, newBalance, newFrozen decimal.Decimal) error
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
}

// ---- In-memory implementations ----

// InMemoryUserRepository is a thread-safe in-memory user store (for development / testing).
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	byID   map[uuid.UUID]*models.User
	byEmail map[string]*models.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		byID:    make(map[uuid.UUID]*models.User),
		byEmail: make(map[string]*models.User),
	}
}

func (r *InMemoryUserRepository) Create(_ context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byEmail[user.Email]; exists {
		return ErrConflict
	}
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	cp := *user
	r.byID[user.UserID] = &cp
	r.byEmail[user.Email] = &cp
	return nil
}

func (r *InMemoryUserRepository) GetByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (r *InMemoryUserRepository) GetByEmail(_ context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (r *InMemoryUserRepository) Update(_ context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byID[user.UserID]
	if !ok {
		return ErrNotFound
	}
	user.UpdatedAt = time.Now()
	cp := *user
	r.byID[user.UserID] = &cp
	// update email index
	delete(r.byEmail, existing.Email)
	r.byEmail[user.Email] = &cp
	return nil
}

// InMemoryProductRepository is a thread-safe in-memory product store.
type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products []*models.Product
	byID     map[uuid.UUID]*models.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		byID: make(map[uuid.UUID]*models.Product),
	}
}

func (r *InMemoryProductRepository) Create(_ context.Context, product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if product.ProductID == uuid.Nil {
		product.ProductID = uuid.New()
	}
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	cp := *product
	r.byID[product.ProductID] = &cp
	r.products = append(r.products, &cp)
	return nil
}

func (r *InMemoryProductRepository) GetByID(_ context.Context, id uuid.UUID) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *InMemoryProductRepository) List(_ context.Context, limit, offset int) ([]*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.products)
	if offset >= total {
		return []*models.Product{}, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	result := make([]*models.Product, end-offset)
	for i, p := range r.products[offset:end] {
		cp := *p
		result[i] = &cp
	}
	return result, nil
}

func (r *InMemoryProductRepository) Update(_ context.Context, product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[product.ProductID]; !ok {
		return ErrNotFound
	}
	product.UpdatedAt = time.Now()
	cp := *product
	r.byID[product.ProductID] = &cp
	for i, p := range r.products {
		if p.ProductID == product.ProductID {
			r.products[i] = &cp
			break
		}
	}
	return nil
}

// InMemoryOrderRepository is a thread-safe in-memory order store.
type InMemoryOrderRepository struct {
	mu    sync.RWMutex
	byID  map[uuid.UUID]*models.Order
	items map[uuid.UUID][]*models.OrderItem
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		byID:  make(map[uuid.UUID]*models.Order),
		items: make(map[uuid.UUID][]*models.OrderItem),
	}
}

func (r *InMemoryOrderRepository) Create(_ context.Context, order *models.Order, items []*models.OrderItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if order.OrderID == uuid.Nil {
		order.OrderID = uuid.New()
	}
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	cp := *order
	r.byID[order.OrderID] = &cp
	its := make([]*models.OrderItem, len(items))
	for i, item := range items {
		if item.ItemID == uuid.Nil {
			item.ItemID = uuid.New()
		}
		item.OrderID = order.OrderID
		ic := *item
		its[i] = &ic
	}
	r.items[order.OrderID] = its
	return nil
}

func (r *InMemoryOrderRepository) GetByID(_ context.Context, id uuid.UUID) (*models.Order, []*models.OrderItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.byID[id]
	if !ok {
		return nil, nil, ErrNotFound
	}
	cp := *o
	its := r.items[id]
	cpItems := make([]*models.OrderItem, len(its))
	for i, item := range its {
		ic := *item
		cpItems[i] = &ic
	}
	return &cp, cpItems, nil
}

func (r *InMemoryOrderRepository) ListByBuyer(_ context.Context, buyerID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []*models.Order
	for _, o := range r.byID {
		if o.BuyerID == buyerID {
			cp := *o
			all = append(all, &cp)
		}
	}
	if offset >= len(all) {
		return []*models.Order{}, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

// InMemoryWalletRepository is a thread-safe in-memory wallet store.
type InMemoryWalletRepository struct {
	mu           sync.RWMutex
	byID         map[uuid.UUID]*models.Wallet
	byUserID     map[uuid.UUID]*models.Wallet
	transactions []*models.Transaction
}

func NewInMemoryWalletRepository() *InMemoryWalletRepository {
	return &InMemoryWalletRepository{
		byID:     make(map[uuid.UUID]*models.Wallet),
		byUserID: make(map[uuid.UUID]*models.Wallet),
	}
}

func (r *InMemoryWalletRepository) GetByUserID(_ context.Context, userID uuid.UUID) (*models.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.byUserID[userID]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *w
	return &cp, nil
}

func (r *InMemoryWalletRepository) Create(_ context.Context, wallet *models.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if wallet.WalletID == uuid.Nil {
		wallet.WalletID = uuid.New()
	}
	wallet.LastUpdated = time.Now()
	cp := *wallet
	r.byID[wallet.WalletID] = &cp
	r.byUserID[wallet.UserID] = &cp
	return nil
}

func (r *InMemoryWalletRepository) UpdateBalance(_ context.Context, walletID uuid.UUID, newBalance, newFrozen decimal.Decimal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.byID[walletID]
	if !ok {
		return ErrNotFound
	}
	w.BalanceTK = newBalance
	w.FrozenBalanceTK = newFrozen
	w.LastUpdated = time.Now()
	r.byUserID[w.UserID] = w
	return nil
}

func (r *InMemoryWalletRepository) CreateTransaction(_ context.Context, tx *models.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if tx.TransactionID == uuid.Nil {
		tx.TransactionID = uuid.New()
	}
	tx.CreatedAt = time.Now()
	cp := *tx
	r.transactions = append(r.transactions, &cp)
	return nil
}
