package handler

import (
	"net/http"
	"strconv"

	apierrors "github.com/azizjon-top/techno-re-ecosystem/internal/errors"
	"github.com/azizjon-top/techno-re-ecosystem/internal/middleware"
	"github.com/azizjon-top/techno-re-ecosystem/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductHandler handles product and order endpoints.
type ProductHandler struct {
	productSvc *service.ProductService
}

// NewProductHandler creates a ProductHandler.
func NewProductHandler(productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

// ListProducts godoc
// GET /api/v1/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.productSvc.ListProducts(c.Request.Context(), limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"limit":    limit,
		"offset":   offset,
	})
}

// createProductRequest is the JSON body for POST /products.
type createProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	PriceUSD    string  `json:"price_usd" binding:"required"`
	Category    *string `json:"category"`
	Stock       int     `json:"stock"`
}

// CreateProduct godoc
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeValidationFailed,
			"message": err.Error(),
		})
		return
	}

	price, err := decimal.NewFromString(req.PriceUSD)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeBadRequest,
			"message": "invalid price format",
		})
		return
	}

	sellerID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	product, svcErr := h.productSvc.CreateProduct(c.Request.Context(), service.CreateProductRequest{
		SellerID:    sellerID,
		Name:        req.Name,
		Description: req.Description,
		PriceUSD:    price,
		Category:    req.Category,
		Stock:       req.Stock,
	})
	if svcErr != nil {
		respondError(c, svcErr)
		return
	}
	c.JSON(http.StatusCreated, product)
}

// createOrderRequest is the JSON body for POST /orders.
type createOrderRequest struct {
	ProductID  string `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	TKDiscount string `json:"tk_discount"`
}

// CreateOrder godoc
// POST /api/v1/orders
func (h *ProductHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeValidationFailed,
			"message": err.Error(),
		})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apierrors.ErrCodeBadRequest,
			"message": "invalid product_id",
		})
		return
	}

	tkDiscount := decimal.Zero
	if req.TKDiscount != "" {
		tkDiscount, err = decimal.NewFromString(req.TKDiscount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    apierrors.ErrCodeBadRequest,
				"message": "invalid tk_discount format",
			})
			return
		}
	}

	buyerID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	order, svcErr := h.productSvc.CreateOrder(c.Request.Context(), service.CreateOrderRequest{
		BuyerID:    buyerID,
		ProductID:  productID,
		Quantity:   req.Quantity,
		TKDiscount: tkDiscount,
	})
	if svcErr != nil {
		respondError(c, svcErr)
		return
	}
	c.JSON(http.StatusCreated, order)
}

// ListOrders godoc
// GET /api/v1/orders
func (h *ProductHandler) ListOrders(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	buyerID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	orders, err := h.productSvc.ListOrders(c.Request.Context(), buyerID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"limit":  limit,
		"offset": offset,
	})
}
