package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"database/sql"
	"homecloud--dbmanager-service/internal/repository"
	protos "homecloud--dbmanager-service/internal/transport/grpc/protos"

	"github.com/stretchr/testify/require"
)

func TestGetUserExtendedInfo(t *testing.T) {
	// Настройка тестовой БД (используйте свою строку подключения или мок)
	dsn := "host=localhost port=5432 user=postgres password=changeme dbname=homecloud_test sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	repo := repository.NewDBRepository(db)

	addr, stop := startTestGRPCServer(t, repo)
	defer stop()

	client := getClient(t, addr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаём пользователя для теста
	user := &protos.User{
		Email:        "test2@example.com",
		Username:     "testuser2",
		PasswordHash: "hash2",
	}
	userID, err := client.CreateUser(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, userID.Id)

	resp, err := client.GetUserExtendedInfo(ctx, &protos.UserID{Id: userID.Id})
	require.NoError(t, err)
	require.NotNil(t, resp)

	fmt.Printf("User: %+v\n", resp.User)
	fmt.Printf("Storage usage: %s\n", resp.StorageUsageFormatted)
	fmt.Printf("Warnings: %v\n", resp.Warnings)
	fmt.Printf("Recommendations: %v\n", resp.Recommendations)
}
