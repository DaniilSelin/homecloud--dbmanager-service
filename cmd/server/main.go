package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"homecloud--dbmanager-service/config"
	"homecloud--dbmanager-service/internal/logger"
	"homecloud--dbmanager-service/internal/repository"
	grpcServer "homecloud--dbmanager-service/internal/transport/grpc/dbManagerServer"
	protos "homecloud--dbmanager-service/internal/transport/grpc/protos"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig("config/config.local.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logr, err := logger.New("debug")
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logr.Error(context.Background(), "failed to connect to db", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logr.Error(context.Background(), "db ping failed", zap.Error(err))
		os.Exit(1)
	}

	repo := repository.NewDBRepository(db)

	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	logr.Info(context.Background(), "Starting gRPC server", zap.String("address", addr))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logr.Error(context.Background(), "failed to listen", zap.Error(err))
		os.Exit(1)
	}

	s := grpc.NewServer()
	protos.RegisterDBServiceServer(s, &grpcServer.Server{Repo: repo, Logger: logr})

	// Graceful shutdown
	go func() {
		logr.Info(context.Background(), "gRPC server listening", zap.String("address", addr))
		if err := s.Serve(lis); err != nil {
			logr.Error(context.Background(), "failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logr.Info(context.Background(), "Shutting down server...")
	s.GracefulStop()
	logr.Info(context.Background(), "Server stopped")
}
