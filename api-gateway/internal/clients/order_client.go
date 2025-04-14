package clients

import (
	"api-gateway/internal/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	conn    *grpc.ClientConn
	Service proto.OrderServiceClient
}

func NewOrderClient(addr string) (*OrderClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	
	service := proto.NewOrderServiceClient(conn)
	log.Printf("Connected to Order Service at %s", addr)
	
	return &OrderClient{
		conn:    conn,
		Service: service,
	}, nil
}

func (c *OrderClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing Order client: %v", err)
	}
}