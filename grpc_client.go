package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "homecloud--dbmanager-service/internal/transport/grpc/protos"
)

func main() {
	userID := flag.String("id", "", "User ID to query")
	addr := flag.String("addr", "127.0.0.1:50051", "gRPC server address")
	flag.Parse()

	if *userID == "" {
		fmt.Println("Usage: go run test/grpc_client.go -id <USER_ID> [-addr <ADDR>]")
		os.Exit(1)
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDBServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUserByID(ctx, &pb.UserID{Id: *userID})
	if err != nil {
		log.Fatalf("gRPC error: %v", err)
	}

	fmt.Printf("User info:\n%+v\n", resp)
}
