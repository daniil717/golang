package delivery

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/proto"
	"order-service/internal/usecase"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGRPCHandler struct {
	proto.UnimplementedOrderServiceServer
	orderUsecase *usecase.OrderUsecase
}

func NewOrderGRPCHandler(uc *usecase.OrderUsecase) *OrderGRPCHandler {
	return &OrderGRPCHandler{orderUsecase: uc}
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	var items []domain.OrderItem
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	order, err := h.orderUsecase.CreateOrder(ctx, req.UserId, items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.OrderResponse{
		Order: toProtoOrder(order),
	}, nil
}

func (h *OrderGRPCHandler) GetOrderByID(ctx context.Context, req *proto.GetOrderRequest) (*proto.OrderResponse, error) {
	order, err := h.orderUsecase.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	return &proto.OrderResponse{Order: toProtoOrder(order)}, nil
}

func (h *OrderGRPCHandler) UpdateOrderStatus(ctx context.Context, req *proto.UpdateOrderStatusRequest) (*proto.OrderResponse, error) {
	if err := h.orderUsecase.UpdateOrderStatus(ctx, req.OrderId, req.Status); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	order, err := h.orderUsecase.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	return &proto.OrderResponse{Order: toProtoOrder(order)}, nil
}

func (h *OrderGRPCHandler) ListUserOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	orders, err := h.orderUsecase.ListUserOrders(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoOrders []*proto.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, toProtoOrder(order))
	}
	return &proto.ListOrdersResponse{Orders: protoOrders}, nil
}

func toProtoOrder(order *domain.Order) *proto.Order {
	var items []*proto.OrderItem
	for _, item := range order.Items {
		items = append(items, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &proto.Order{
		Id:        order.ID,
		UserId:    order.UserID,
		Items:     items,
		Total:     order.Total,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}
}
