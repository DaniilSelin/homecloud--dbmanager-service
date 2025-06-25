package interfaces

import (
	"context"
	"homecloud--dbmanager-service/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	UpdateUsername(ctx context.Context, id, username string) error
	UpdateEmailVerification(ctx context.Context, id string, isVerified bool) error
	UpdateLastLogin(ctx context.Context, id string) error
	UpdateFailedLoginAttempts(ctx context.Context, id string, attempts int) error
	UpdateLockedUntil(ctx context.Context, id string) error
	UpdateStorageUsage(ctx context.Context, id string, usedSpace int64) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}

type FileService interface {
	// File operations
	CreateFile(ctx context.Context, file *models.File) (string, error)
	GetFileByID(ctx context.Context, id string) (*models.File, error)
	GetFileByPath(ctx context.Context, ownerID, path string) (*models.File, error)
	UpdateFile(ctx context.Context, file *models.File) error
	DeleteFile(ctx context.Context, id string) error
	SoftDeleteFile(ctx context.Context, id string) error
	RestoreFile(ctx context.Context, id string) error
	ListFiles(ctx context.Context, parentID, ownerID string, isTrashed, starred bool, limit, offset int, orderBy, orderDir string) ([]*models.File, int64, error)
	ListFilesByParent(ctx context.Context, ownerID, parentID string) ([]*models.File, error)
	ListStarredFiles(ctx context.Context, ownerID string) ([]*models.File, error)
	ListTrashedFiles(ctx context.Context, ownerID string) ([]*models.File, error)
	SearchFiles(ctx context.Context, ownerID, query string) ([]*models.File, error)
	GetFileSize(ctx context.Context, id string) (int64, error)
	UpdateFileSize(ctx context.Context, id string, size int64) error
	UpdateLastViewed(ctx context.Context, id string) error
	GetFileTree(ctx context.Context, ownerID, rootID string) ([]*models.File, error)

	// File revision operations
	CreateRevision(ctx context.Context, revision *models.FileRevision) (string, error)
	GetRevisions(ctx context.Context, fileID string) ([]*models.FileRevision, error)
	GetRevision(ctx context.Context, fileID string, revisionID int64) (*models.FileRevision, error)
	DeleteRevision(ctx context.Context, id string) error

	// File permission operations
	CreatePermission(ctx context.Context, permission *models.FilePermission) (string, error)
	GetPermissions(ctx context.Context, fileID string) ([]*models.FilePermission, error)
	UpdatePermission(ctx context.Context, permission *models.FilePermission) error
	DeletePermission(ctx context.Context, id string) error
	CheckPermission(ctx context.Context, fileID, userID, requiredRole string) (bool, error)

	// File metadata operations
	UpdateFileMetadata(ctx context.Context, fileID, metadata string) error
	GetFileMetadata(ctx context.Context, fileID string) (string, error)

	// File operations (star, move, copy, rename)
	StarFile(ctx context.Context, id string) error
	UnstarFile(ctx context.Context, id string) error
	MoveFile(ctx context.Context, fileID, newParentID string) error
	CopyFile(ctx context.Context, fileID, newParentID, newName string) (*models.File, error)
	RenameFile(ctx context.Context, fileID, newName string) error

	// File integrity operations
	VerifyFileIntegrity(ctx context.Context, id string) (bool, error)
	CalculateFileChecksums(ctx context.Context, id string) (map[string]string, error)
}
