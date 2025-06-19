package interfaces

import (
	"context"
	"homecloud--dbmanager-service/internal/models"
	"time"
)

type DBRepository interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	UpdateUsername(ctx context.Context, id, username string) error
	UpdateEmailVerification(ctx context.Context, id string, isVerified bool) error
	UpdateLastLogin(ctx context.Context, id string, lastLogin time.Time) error
	UpdateFailedLoginAttempts(ctx context.Context, id string, attempts int) error
	UpdateLockedUntil(ctx context.Context, id string, lockedUntil time.Time) error
	UpdateStorageUsage(ctx context.Context, id string, usedSpace int64) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}
