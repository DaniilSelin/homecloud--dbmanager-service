package test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"homecloud--dbmanager-service/internal/interfaces"
	"homecloud--dbmanager-service/internal/repository"
	grpcServer "homecloud--dbmanager-service/internal/transport/grpc/dbManagerServer"
	protos "homecloud--dbmanager-service/internal/transport/grpc/protos"
)

const (
	testDBName     = "homecloud_test"
	testDBUser     = "postgres"
	testDBPassword = "changeme"
	testDBHost     = "localhost"
	testDBPort     = 5432
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", testDBHost, testDBPort, testDBUser, testDBPassword)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	_, _ = db.Exec("DROP DATABASE IF EXISTS " + testDBName)
	_, err = db.Exec("CREATE DATABASE " + testDBName)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	db.Close()

	dsnTest := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", testDBHost, testDBPort, testDBUser, testDBPassword, testDBName)
	dbTest, err := sql.Open("postgres", dsnTest)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	_, err = dbTest.Exec(`CREATE SCHEMA IF NOT EXISTS dbmanager;`)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	_, err = dbTest.Exec(`CREATE TABLE IF NOT EXISTS dbmanager.users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
		last_login_at TIMESTAMP,
		failed_login_attempts INTEGER NOT NULL DEFAULT 0,
		locked_until TIMESTAMP,
		two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
		storage_quota BIGINT NOT NULL DEFAULT 10737418240,
		used_space BIGINT NOT NULL DEFAULT 0,
		role TEXT NOT NULL DEFAULT 'user',
		is_admin BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP NOT NULL DEFAULT now()
	);`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}
	return dbTest
}

func startTestGRPCServer(t *testing.T, repo interfaces.DBRepository) (addr string, stop func()) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	protos.RegisterDBServiceServer(server, &grpcServer.Server{Repo: repo})
	go server.Serve(lis)
	return lis.Addr().String(), func() { server.Stop(); lis.Close() }
}

func getClient(t *testing.T, addr string) protos.DBServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial grpc: %v", err)
	}
	return protos.NewDBServiceClient(conn)
}

func TestDBService_AllMethods(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewDBRepository(db)
	addr, stop := startTestGRPCServer(t, repo)
	defer stop()
	client := getClient(t, addr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// CreateUser
	user := &protos.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hash",
	}
	userID, err := client.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if userID.Id == "" {
		t.Fatalf("CreateUser returned empty id")
	}

	// GetUserByID
	gotUser, err := client.GetUserByID(ctx, &protos.UserID{Id: userID.Id})
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if gotUser.Email != user.Email {
		t.Errorf("GetUserByID: expected email %s, got %s", user.Email, gotUser.Email)
	}

	// GetUserByEmail
	gotUser2, err := client.GetUserByEmail(ctx, &protos.EmailRequest{Email: user.Email})
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if gotUser2.Username != user.Username {
		t.Errorf("GetUserByEmail: expected username %s, got %s", user.Username, gotUser2.Username)
	}

	// UpdateUser
	gotUser2.Username = "updateduser"
	_, err = client.UpdateUser(ctx, gotUser2)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	gotUser3, _ := client.GetUserByID(ctx, &protos.UserID{Id: userID.Id})
	if gotUser3.Username != "updateduser" {
		t.Errorf("UpdateUser: username not updated")
	}

	// UpdatePassword
	_, err = client.UpdatePassword(ctx, &protos.UpdatePasswordRequest{Id: userID.Id, PasswordHash: "newhash"})
	if err != nil {
		t.Fatalf("UpdatePassword failed: %v", err)
	}

	// UpdateUsername
	_, err = client.UpdateUsername(ctx, &protos.UpdateUsernameRequest{Id: userID.Id, Username: "newname"})
	if err != nil {
		t.Fatalf("UpdateUsername failed: %v", err)
	}

	// UpdateEmailVerification
	_, err = client.UpdateEmailVerification(ctx, &protos.UpdateEmailVerificationRequest{Id: userID.Id, IsVerified: true})
	if err != nil {
		t.Fatalf("UpdateEmailVerification failed: %v", err)
	}

	// UpdateLastLogin
	_, err = client.UpdateLastLogin(ctx, &protos.UserID{Id: userID.Id})
	if err != nil {
		t.Fatalf("UpdateLastLogin failed: %v", err)
	}

	// UpdateFailedLoginAttempts
	_, err = client.UpdateFailedLoginAttempts(ctx, &protos.UpdateFailedLoginAttemptsRequest{Id: userID.Id, Attempts: 3})
	if err != nil {
		t.Fatalf("UpdateFailedLoginAttempts failed: %v", err)
	}

	// UpdateLockedUntil
	_, err = client.UpdateLockedUntil(ctx, &protos.UpdateLockedUntilRequest{Id: userID.Id})
	if err != nil {
		t.Fatalf("UpdateLockedUntil failed: %v", err)
	}

	// UpdateStorageUsage
	_, err = client.UpdateStorageUsage(ctx, &protos.UpdateStorageUsageRequest{Id: userID.Id, UsedSpace: 12345})
	if err != nil {
		t.Fatalf("UpdateStorageUsage failed: %v", err)
	}

	// CheckEmailExists
	emailExists, err := client.CheckEmailExists(ctx, &protos.EmailRequest{Email: user.Email})
	if err != nil {
		t.Fatalf("CheckEmailExists failed: %v", err)
	}
	if !emailExists.Exists {
		t.Errorf("CheckEmailExists: expected true, got false")
	}

	// CheckUsernameExists
	usernameExists, err := client.CheckUsernameExists(ctx, &protos.UsernameRequest{Username: "newname"})
	if err != nil {
		t.Fatalf("CheckUsernameExists failed: %v", err)
	}
	if !usernameExists.Exists {
		t.Errorf("CheckUsernameExists: expected true, got false")
	}
}
