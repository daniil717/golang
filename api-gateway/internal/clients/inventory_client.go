package clients

import (
	"api-gateway/internal/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryClient struct {
	conn    *grpc.ClientConn
	Service proto.InventoryServiceClient
}

func NewInventoryClient(addr string) (*InventoryClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	service := proto.NewInventoryServiceClient(conn)
	log.Printf("Connected to Inventory Service at %s", addr)

	return &InventoryClient{
		conn:    conn,
		Service: service,
	}, nil
}

func (c *InventoryClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing Inventory client: %v", err)
	}
}
