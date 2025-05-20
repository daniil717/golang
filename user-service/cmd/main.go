package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"user-service/config"
	"user-service/internal/handler"
	"user-service/internal/pb"
	"user-service/internal/repository"
	"user-service/internal/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// ✅ Принудительно подгружаем JWT_SECRET из окружения
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	client, err := mongo.Connect(cfg.Ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Mongo connect error: %v", err)
	}
	defer client.Disconnect(cfg.Ctx)

	col := client.Database(cfg.MongoDBName).Collection("users")
	userRepo := repository.NewMongoUserRepository(col)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUC, cfg.JWTSecret)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	fmt.Printf("UserService running on :%s\n", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Serve error: %v", err)
	}
}
