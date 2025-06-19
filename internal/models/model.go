package models

import "time"

type User struct {
	ID                  string
	Email               string
	Username            string
	PasswordHash        string
	IsActive            bool
	IsEmailVerified     bool
	Role                string
	StorageQuota        int64
	UsedSpace           int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
	FailedLoginAttempts int32
	LockedUntil         *time.Time
	LastLogin           *time.Time
}
