package main

import (
	"context"
	"inventory-servicee/config"
	"inventory-servicee/internal/proto"
	"inventory-servicee/internal/repository"
	"inventory-servicee/internal/usecase"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Подключение к MongoDB
	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database(cfg.MongoDatabase)

	// 3. Инициализация зависимостей
	productRepo := repository.NewProductMongoRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	grpcHandler := grpc.NewInventoryGRPCHandler(productUsecase)

	// 4. Создание gRPC сервера
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	proto.RegisterInventoryServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer) // Для отладки через grpcurl

	// 5. Запуск сервера с graceful shutdown
	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC server on %s", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидание сигналов завершения
	waitForShutdown(grpcServer)
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Method: %s, Duration: %v, Error: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func waitForShutdown(server *grpc.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	server.GracefulStop()
	log.Println("Server stopped")
}
