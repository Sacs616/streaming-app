package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Sacs616/streaming-app/services/video/internal/domain"
	grpcServer "github.com/Sacs616/streaming-app/services/video/internal/grpc"
	"github.com/Sacs616/streaming-app/services/video/internal/repository"
	"github.com/Sacs616/streaming-app/services/video/internal/service"
	videopb "github.com/Sacs616/streaming-app/services/video/proto"
)

func main() {
	// Load config
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://netflix:netflix@localhost:5434/video_db?sslmode=disable"
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(&domain.Video{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize layers
	videoRepo := repository.NewVideoRepository(db)
	videoService := service.NewVideoService(videoRepo)
	videoServer := grpcServer.NewVideoServer(videoService)

	// Create gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	videopb.RegisterVideoServiceServer(s, videoServer)
	reflection.Register(s)

	log.Printf("Video Service listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
