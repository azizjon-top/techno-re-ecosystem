package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Product struct {
	ProductID    uuid.UUID       `db:"product_id" json:"product_id"`
	SellerID     uuid.UUID       `db:"seller_id" json:"seller_id"`
	Name         string          `db:"name" json:"name"`
	Description  *string         `db:"description" json:"description"`
	PriceUSD     decimal.Decimal `db:"price_usd" json:"price_usd"`
	Category     *string         `db:"category" json:"category"`
	Stock        int             `db:"stock" json:"stock"`
	ImageURL     *string         `db:"image_url" json:"image_url"`
	AISearchTags []string        `db:"ai_search_tags" json:"ai_search_tags"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

type Order struct {
	OrderID        uuid.UUID       `db:"order_id" json:"order_id"`
	BuyerID        uuid.UUID       `db:"buyer_id" json:"buyer_id"`
	SellerID       uuid.UUID       `db:"seller_id" json:"seller_id"`
	TotalPriceUSD  decimal.Decimal `db:"total_price_usd" json:"total_price_usd"`
	TKDiscountUsed decimal.Decimal `db:"tk_discount_used" json:"tk_discount_used"`
	FinalPriceUSD  decimal.Decimal `db:"final_price_usd" json:"final_price_usd"`
	Status         OrderStatus     `db:"status" json:"status"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ItemID       uuid.UUID       `db:"item_id" json:"item_id"`
	OrderID      uuid.UUID       `db:"order_id" json:"order_id"`
	ProductID    uuid.UUID       `db:"product_id" json:"product_id"`
	Quantity     int             `db:"quantity" json:"quantity"`
	PricePerUnit decimal.Decimal `db:"price_per_unit" json:"price_per_unit"`
}

// IsInStock checks if product has sufficient stock
func (p *Product) IsInStock(quantity int) bool {
	return p.Stock >= quantity
}

// CalculateFinalPrice calculates final price after TK discount (max 50%)
func (o *Order) CalculateFinalPrice() decimal.Decimal {
	maxDiscount := o.TotalPriceUSD.Div(decimal.NewFromInt(2))
	if o.TKDiscountUsed.GreaterThan(maxDiscount) {
		return o.TotalPriceUSD.Sub(maxDiscount)
	}
	return o.TotalPriceUSD.Sub(o.TKDiscountUsed)
}
