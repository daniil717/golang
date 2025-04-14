package handlers

import (
	"net/http"
	"strconv"
	"time"

	"api-gateway/internal/clients"
	"api-gateway/internal/proto"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	client *clients.OrderClient
}

func NewOrderHandler(client *clients.OrderClient) *OrderHandler {
	return &OrderHandler{client: client}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString("userID") 

	var req struct {
		Items []struct {
			ProductID string  `json:"product_id" binding:"required"`
			Quantity  int32   `json:"quantity" binding:"required,gt=0"`
			Price     float64 `json:"price" binding:"required,gt=0"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []*proto.OrderItem
	for _, item := range req.Items {
		items = append(items, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	grpcResp, err := h.client.Service.CreateOrder(c.Request.Context(), &proto.CreateOrderRequest{
		UserId: userID,
		Items:  items,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order_id":   grpcResp.Order.Id,
		"user_id":    grpcResp.Order.UserId,
		"status":     grpcResp.Order.Status,
		"total":      grpcResp.Order.Total,
		"created_at": grpcResp.Order.CreatedAt,
		"items":      grpcResp.Order.Items,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetString("userID")
	orderID := c.Param("id")

	grpcResp, err := h.client.Service.GetOrderByID(c.Request.Context(), &proto.GetOrderRequest{
		OrderId: orderID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if grpcResp.Order.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":   grpcResp.Order.Id,
		"user_id":    grpcResp.Order.UserId,
		"status":     grpcResp.Order.Status,
		"total":      grpcResp.Order.Total,
		"created_at": grpcResp.Order.CreatedAt,
		"items":      grpcResp.Order.Items,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	userID := c.GetString("userID")
	orderID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required,oneof=paid shipped completed cancelled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getResp, err := h.client.Service.GetOrderByID(c.Request.Context(), &proto.GetOrderRequest{
		OrderId: orderID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if getResp.Order.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	updateResp, err := h.client.Service.UpdateOrderStatus(c.Request.Context(), &proto.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  req.Status,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":   updateResp.Order.Id,
		"status":     updateResp.Order.Status,
		"updated_at": time.Now().Format(time.RFC3339),
	})
}

func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	grpcResp, err := h.client.Service.ListUserOrders(c.Request.Context(), &proto.ListOrdersRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var orders []gin.H
	for _, o := range grpcResp.Orders {
		orders = append(orders, gin.H{
			"order_id":   o.Id,
			"status":     o.Status,
			"total":      o.Total,
			"created_at": o.CreatedAt,
			"item_count": len(o.Items),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}
