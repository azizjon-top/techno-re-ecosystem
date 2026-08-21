package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product represents a product in the e-commerce store
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Image       string    `json:"image"`
	Category    string    `json:"category"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Order represents a customer order
type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Items     []OrderItem `json:"items"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Wallet represents a user's wallet
type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Chat represents a chat conversation
type Chat struct {
	ID          string    `json:"id"`
	Participants []string  `json:"participants"`
	LastMessage string    `json:"last_message"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// Video represents a video in the platform
type Video struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    int       `json:"duration"`
	Views       int       `json:"views"`
	Likes       int       `json:"likes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MiningSession represents a mining session
type MiningSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	Reward    float64   `json:"reward"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

// Fact represents a validated fact in the mining network
type Fact struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Validator string    `json:"validator"`
	Status    string    `json:"status"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

// Campaign represents an advertising campaign
type Campaign struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Budget    float64   `json:"budget"`
	Spent     float64   `json:"spent"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}

// Impression represents a campaign impression
type Impression struct {
	ID        string    `json:"id"`
	CampaignID string   `json:"campaign_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}