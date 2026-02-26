package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Sacs616/streaming-app/services/auth/internal/domain"
	grpcServer "github.com/Sacs616/streaming-app/services/auth/internal/grpc"
	"github.com/Sacs616/streaming-app/services/auth/internal/repository"
	"github.com/Sacs616/streaming-app/services/auth/internal/service"
	pb "github.com/Sacs616/streaming-app/services/auth/proto"
)

func main() {
	// Load config from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://netflix:netflix@localhost:5432/auth_db?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-production"
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authServer := grpcServer.NewAuthServer(authService)

	// Create gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, authServer)
	reflection.Register(s)

	log.Printf("Auth Service listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
