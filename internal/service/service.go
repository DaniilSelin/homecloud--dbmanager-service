package service

import (
	"context"
	"time"

	"homecloud--dbmanager-service/internal/interfaces"
	"homecloud--dbmanager-service/internal/models"
)

type userService struct {
	repo interfaces.DBRepository
}

func NewUserService(repo interfaces.DBRepository) interfaces.UserService {
	return &userService{repo: repo}
}

// UserService implementations
func (s *userService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	return s.repo.UpdatePassword(ctx, id, passwordHash)
}

func (s *userService) UpdateUsername(ctx context.Context, id, username string) error {
	return s.repo.UpdateUsername(ctx, id, username)
}

func (s *userService) UpdateEmailVerification(ctx context.Context, id string, isVerified bool) error {
	return s.repo.UpdateEmailVerification(ctx, id, isVerified)
}

func (s *userService) UpdateLastLogin(ctx context.Context, id string) error {
	return s.repo.UpdateLastLogin(ctx, id, time.Now())
}

func (s *userService) UpdateFailedLoginAttempts(ctx context.Context, id string, attempts int) error {
	return s.repo.UpdateFailedLoginAttempts(ctx, id, attempts)
}

func (s *userService) UpdateLockedUntil(ctx context.Context, id string) error {
	return s.repo.UpdateLockedUntil(ctx, id, time.Now())
}

func (s *userService) UpdateStorageUsage(ctx context.Context, id string, usedSpace int64) error {
	return s.repo.UpdateStorageUsage(ctx, id, usedSpace)
}

func (s *userService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return s.repo.CheckEmailExists(ctx, email)
}

func (s *userService) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	return s.repo.CheckUsernameExists(ctx, username)
}

// FileService implementation
type fileService struct {
	repo interfaces.DBRepository
}

func NewFileService(repo interfaces.DBRepository) interfaces.FileService {
	return &fileService{repo: repo}
}

// File operations
func (s *fileService) CreateFile(ctx context.Context, file *models.File) (string, error) {
	return s.repo.CreateFile(ctx, file)
}

func (s *fileService) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	return s.repo.GetFileByID(ctx, id)
}

func (s *fileService) GetFileByPath(ctx context.Context, ownerID, path string) (*models.File, error) {
	return s.repo.GetFileByPath(ctx, ownerID, path)
}

func (s *fileService) UpdateFile(ctx context.Context, file *models.File) error {
	return s.repo.UpdateFile(ctx, file)
}

func (s *fileService) DeleteFile(ctx context.Context, id string) error {
	return s.repo.DeleteFile(ctx, id)
}

func (s *fileService) SoftDeleteFile(ctx context.Context, id string) error {
	return s.repo.SoftDeleteFile(ctx, id)
}

func (s *fileService) RestoreFile(ctx context.Context, id string) error {
	return s.repo.RestoreFile(ctx, id)
}

func (s *fileService) ListFiles(ctx context.Context, parentID, ownerID string, isTrashed, starred bool, limit, offset int, orderBy, orderDir string) ([]*models.File, int64, error) {
	return s.repo.ListFiles(ctx, parentID, ownerID, isTrashed, starred, limit, offset, orderBy, orderDir)
}

func (s *fileService) ListFilesByParent(ctx context.Context, ownerID, parentID string) ([]*models.File, error) {
	return s.repo.ListFilesByParent(ctx, ownerID, parentID)
}

func (s *fileService) ListStarredFiles(ctx context.Context, ownerID string) ([]*models.File, error) {
	return s.repo.ListStarredFiles(ctx, ownerID)
}

func (s *fileService) ListTrashedFiles(ctx context.Context, ownerID string) ([]*models.File, error) {
	return s.repo.ListTrashedFiles(ctx, ownerID)
}

func (s *fileService) SearchFiles(ctx context.Context, ownerID, query string) ([]*models.File, error) {
	return s.repo.SearchFiles(ctx, ownerID, query)
}

func (s *fileService) GetFileSize(ctx context.Context, id string) (int64, error) {
	return s.repo.GetFileSize(ctx, id)
}

func (s *fileService) UpdateFileSize(ctx context.Context, id string, size int64) error {
	return s.repo.UpdateFileSize(ctx, id, size)
}

func (s *fileService) UpdateLastViewed(ctx context.Context, id string) error {
	return s.repo.UpdateLastViewed(ctx, id)
}

func (s *fileService) GetFileTree(ctx context.Context, ownerID, rootID string) ([]*models.File, error) {
	return s.repo.GetFileTree(ctx, ownerID, rootID)
}

// File revision operations
func (s *fileService) CreateRevision(ctx context.Context, revision *models.FileRevision) (string, error) {
	return s.repo.CreateRevision(ctx, revision)
}

func (s *fileService) GetRevisions(ctx context.Context, fileID string) ([]*models.FileRevision, error) {
	return s.repo.GetRevisions(ctx, fileID)
}

func (s *fileService) GetRevision(ctx context.Context, fileID string, revisionID int64) (*models.FileRevision, error) {
	return s.repo.GetRevision(ctx, fileID, revisionID)
}

func (s *fileService) DeleteRevision(ctx context.Context, id string) error {
	return s.repo.DeleteRevision(ctx, id)
}

// File permission operations
func (s *fileService) CreatePermission(ctx context.Context, permission *models.FilePermission) (string, error) {
	return s.repo.CreatePermission(ctx, permission)
}

func (s *fileService) GetPermissions(ctx context.Context, fileID string) ([]*models.FilePermission, error) {
	return s.repo.GetPermissions(ctx, fileID)
}

func (s *fileService) UpdatePermission(ctx context.Context, permission *models.FilePermission) error {
	return s.repo.UpdatePermission(ctx, permission)
}

func (s *fileService) DeletePermission(ctx context.Context, id string) error {
	return s.repo.DeletePermission(ctx, id)
}

func (s *fileService) CheckPermission(ctx context.Context, fileID, userID, requiredRole string) (bool, error) {
	return s.repo.CheckPermission(ctx, fileID, userID, requiredRole)
}

// File metadata operations
func (s *fileService) UpdateFileMetadata(ctx context.Context, fileID, metadata string) error {
	return s.repo.UpdateFileMetadata(ctx, fileID, metadata)
}

func (s *fileService) GetFileMetadata(ctx context.Context, fileID string) (string, error) {
	return s.repo.GetFileMetadata(ctx, fileID)
}

// File operations (star, move, copy, rename)
func (s *fileService) StarFile(ctx context.Context, id string) error {
	return s.repo.StarFile(ctx, id)
}

func (s *fileService) UnstarFile(ctx context.Context, id string) error {
	return s.repo.UnstarFile(ctx, id)
}

func (s *fileService) MoveFile(ctx context.Context, fileID, newParentID string) error {
	return s.repo.MoveFile(ctx, fileID, newParentID)
}

func (s *fileService) CopyFile(ctx context.Context, fileID, newParentID, newName string) (*models.File, error) {
	return s.repo.CopyFile(ctx, fileID, newParentID, newName)
}

func (s *fileService) RenameFile(ctx context.Context, fileID, newName string) error {
	return s.repo.RenameFile(ctx, fileID, newName)
}

// File integrity operations
func (s *fileService) VerifyFileIntegrity(ctx context.Context, id string) (bool, error) {
	return s.repo.VerifyFileIntegrity(ctx, id)
}

func (s *fileService) CalculateFileChecksums(ctx context.Context, id string) (map[string]string, error) {
	return s.repo.CalculateFileChecksums(ctx, id)
}
