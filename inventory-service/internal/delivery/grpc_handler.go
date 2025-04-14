package delivery

import (
	"context"
	"errors"
	"inventory-servicee/internal/domain"
	"inventory-servicee/internal/proto"
	"inventory-servicee/internal/usecase"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type InventoryGRPCHandler struct {
	proto.UnimplementedInventoryServiceServer
	productUsecase *usecase.ProductUsecase
}

func NewInventoryGRPCHandler(uc *usecase.ProductUsecase) *InventoryGRPCHandler {
	return &InventoryGRPCHandler{productUsecase: uc}
}

func (h *InventoryGRPCHandler) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.ProductResponse, error) {
	if req.Name == "" || req.Price <= 0 || req.Category == "" {
		return nil, status.Error(codes.InvalidArgument, "name, price, and category are required")
	}

	product := &domain.Product{
		Name:     req.Name,
		Category: req.Category,
		Price:    float64(req.Price),
		Stock:    req.Stock,
	}

	if err := h.productUsecase.CreateProduct(ctx, product); err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.ProductResponse{
		Product: &proto.Product{
			Id:       product.ID, 
			Name:     product.Name,
			Category: product.Category,
			Price:    float32(product.Price),
			Stock:    product.Stock,
		},
	}, nil
}

func (h *InventoryGRPCHandler) GetProductByID(ctx context.Context, req *proto.GetProductRequest) (*proto.ProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	product, err := h.productUsecase.GetProductByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, status.Error(codes.InvalidArgument, "invalid product ID format")
		}
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &proto.ProductResponse{
		Product: &proto.Product{
			Id:       product.ID, 
			Name:     product.Name,
			Category: product.Category,
			Price:    float32(product.Price),
			Stock:    product.Stock,
		},
	}, nil
}

func (h *InventoryGRPCHandler) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.ProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	product := &domain.Product{
		Name:     req.Name,
		Category: req.Category,
		Price:    float64(req.Price),
		Stock:    req.Stock,
	}

	if err := h.productUsecase.UpdateProduct(ctx, req.Id, product); err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, status.Error(codes.InvalidArgument, "invalid product ID format")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	updatedProduct, err := h.productUsecase.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to fetch updated product")
	}

	return &proto.ProductResponse{
		Product: &proto.Product{
			Id:       updatedProduct.ID, 
			Name:     updatedProduct.Name,
			Category: updatedProduct.Category,
			Price:    float32(updatedProduct.Price),
			Stock:    updatedProduct.Stock,
		},
	}, nil
}

func (h *InventoryGRPCHandler) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	if err := h.productUsecase.DeleteProduct(ctx, req.Id); err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, status.Error(codes.InvalidArgument, "invalid product ID format")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (h *InventoryGRPCHandler) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	products, err := h.productUsecase.ListProducts(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoProducts []*proto.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &proto.Product{
			Id:       p.ID, 
			Name:     p.Name,
			Category: p.Category,
			Price:    float32(p.Price),
			Stock:    p.Stock,
		})
	}

	return &proto.ListProductsResponse{
		Products: protoProducts,
	}, nil
}
