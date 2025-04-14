package handlers

import (
	"net/http"
	"strconv"

	"api-gateway/internal/clients"
	"api-gateway/internal/proto"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	client *clients.InventoryClient
}

func NewInventoryHandler(client *clients.InventoryClient) *InventoryHandler {
	return &InventoryHandler{client: client}
}

func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required"`
		Category string  `json:"category" binding:"required"`
		Price    float64 `json:"price" binding:"required,gt=0"`
		Stock    int32   `json:"stock" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcResp, err := h.client.Service.CreateProduct(c.Request.Context(), &proto.CreateProductRequest{
		Name:     req.Name,
		Category: req.Category,
		Price:    float32(req.Price),
		Stock:    req.Stock,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       grpcResp.Product.Id,
		"name":     grpcResp.Product.Name,
		"category": grpcResp.Product.Category,
		"price":    grpcResp.Product.Price,
		"stock":    grpcResp.Product.Stock,
	})
}

func (h *InventoryHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	grpcResp, err := h.client.Service.GetProductByID(c.Request.Context(), &proto.GetProductRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       grpcResp.Product.Id,
		"name":     grpcResp.Product.Name,
		"category": grpcResp.Product.Category,
		"price":    grpcResp.Product.Price,
		"stock":    grpcResp.Product.Stock,
	})
}

func (h *InventoryHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	var req struct {
		Name     string  `json:"name"`
		Category string  `json:"category"`
		Price    float64 `json:"price" binding:"gt=0"`
		Stock    int32   `json:"stock" binding:"gte=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcResp, err := h.client.Service.UpdateProduct(c.Request.Context(), &proto.UpdateProductRequest{
		Id:       id,
		Name:     req.Name,
		Category: req.Category,
		Price:    float32(req.Price),
		Stock:    req.Stock,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       grpcResp.Product.Id,
		"name":     grpcResp.Product.Name,
		"category": grpcResp.Product.Category,
		"price":    grpcResp.Product.Price,
		"stock":    grpcResp.Product.Stock,
	})
}

func (h *InventoryHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	_, err := h.client.Service.DeleteProduct(c.Request.Context(), &proto.DeleteProductRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *InventoryHandler) ListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	grpcResp, err := h.client.Service.ListProducts(c.Request.Context(), &proto.ListProductsRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var products []gin.H
	for _, p := range grpcResp.Products {
		products = append(products, gin.H{
			"id":       p.Id,
			"name":     p.Name,
			"category": p.Category,
			"price":    p.Price,
			"stock":    p.Stock,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}
