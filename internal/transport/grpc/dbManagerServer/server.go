package dbManagerServer

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"homecloud--dbmanager-service/internal/interfaces"
	"homecloud--dbmanager-service/internal/logger"
	"homecloud--dbmanager-service/internal/models"
	protos "homecloud--dbmanager-service/internal/transport/grpc/protos"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	protos.UnimplementedDBServiceServer
	Repo   interfaces.DBRepository
	Logger *logger.Logger
}

func userModelToProto(u *models.User) *protos.User {
	return &protos.User{
		Id:                  u.ID,
		Email:               u.Email,
		Username:            u.Username,
		PasswordHash:        u.PasswordHash,
		IsActive:            u.IsActive,
		IsEmailVerified:     u.IsEmailVerified,
		Role:                u.Role,
		StorageQuota:        u.StorageQuota,
		UsedSpace:           u.UsedSpace,
		CreatedAt:           timestamppb.New(u.CreatedAt),
		UpdatedAt:           timestamppb.New(u.UpdatedAt),
		FailedLoginAttempts: u.FailedLoginAttempts,
		LockedUntil:         timeToProto(u.LockedUntil),
		LastLogin:           timeToProto(u.LastLogin),
	}
}

func protoToUserModel(u *protos.User) *models.User {
	return &models.User{
		ID:                  u.Id,
		Email:               u.Email,
		Username:            u.Username,
		PasswordHash:        u.PasswordHash,
		IsActive:            u.IsActive,
		IsEmailVerified:     u.IsEmailVerified,
		Role:                u.Role,
		StorageQuota:        u.StorageQuota,
		UsedSpace:           u.UsedSpace,
		CreatedAt:           u.CreatedAt.AsTime(),
		UpdatedAt:           u.UpdatedAt.AsTime(),
		FailedLoginAttempts: u.FailedLoginAttempts,
		LockedUntil:         protoToTime(u.LockedUntil),
		LastLogin:           protoToTime(u.LastLogin),
	}
}

// Helper functions for preferences conversion
func convertPreferencesToMap(prefs map[string]interface{}) map[string]string {
	if prefs == nil {
		return nil
	}
	result := make(map[string]string)
	for k, v := range prefs {
		if str, ok := v.(string); ok {
			result[k] = str
		} else {
			// Convert other types to string representation
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}

func convertMapToPreferences(prefs map[string]string) map[string]interface{} {
	if prefs == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range prefs {
		result[k] = v
	}
	return result
}

func timeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func protoToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func (s *Server) CreateUser(ctx context.Context, req *protos.User) (*protos.UserID, error) {
	id, err := s.Repo.CreateUser(ctx, protoToUserModel(req))
	if err != nil {
		return nil, err
	}
	return &protos.UserID{Id: id}, nil
}

func (s *Server) GetUserByID(ctx context.Context, req *protos.UserID) (*protos.User, error) {
	s.Logger.Info(ctx, "GetUserByID called", zap.String("user_id", req.Id))
	u, err := s.Repo.GetUserByID(ctx, req.Id)
	if err != nil {
		s.Logger.Error(ctx, "GetUserByID failed", zap.Error(err))
		return nil, err
	}
	if u == nil {
		s.Logger.Info(ctx, "User not found", zap.String("user_id", req.Id))
		return nil, nil
	}
	resp := userModelToProto(u)
	s.Logger.Debug(ctx, "GetUserByID response", zap.Any("response", resp))
	return resp, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, req *protos.EmailRequest) (*protos.User, error) {
	u, err := s.Repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return userModelToProto(u), nil
}

func (s *Server) UpdateUser(ctx context.Context, req *protos.User) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateUser(ctx, protoToUserModel(req)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdatePassword(ctx context.Context, req *protos.UpdatePasswordRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdatePassword(ctx, req.Id, req.PasswordHash); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateUsername(ctx context.Context, req *protos.UpdateUsernameRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateUsername(ctx, req.Id, req.Username); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateEmailVerification(ctx context.Context, req *protos.UpdateEmailVerificationRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateEmailVerification(ctx, req.Id, req.IsVerified); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateLastLogin(ctx context.Context, req *protos.UserID) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateLastLogin(ctx, req.Id, time.Now()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateFailedLoginAttempts(ctx context.Context, req *protos.UpdateFailedLoginAttemptsRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateFailedLoginAttempts(ctx, req.Id, int(req.Attempts)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateLockedUntil(ctx context.Context, req *protos.UpdateLockedUntilRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateLockedUntil(ctx, req.Id, req.LockedUntil.AsTime()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateStorageUsage(ctx context.Context, req *protos.UpdateStorageUsageRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateStorageUsage(ctx, req.Id, req.UsedSpace); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) CheckEmailExists(ctx context.Context, req *protos.EmailRequest) (*protos.ExistsResponse, error) {
	exists, err := s.Repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &protos.ExistsResponse{Exists: exists}, nil
}

func (s *Server) CheckUsernameExists(ctx context.Context, req *protos.UsernameRequest) (*protos.ExistsResponse, error) {
	exists, err := s.Repo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &protos.ExistsResponse{Exists: exists}, nil
}

// File operations
func (s *Server) CreateFile(ctx context.Context, req *protos.File) (*protos.FileID, error) {
	file := protoToFileModel(req)
	id, err := s.Repo.CreateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	return &protos.FileID{Id: id}, nil
}

func (s *Server) GetFileByID(ctx context.Context, req *protos.FileID) (*protos.File, error) {
	file, err := s.Repo.GetFileByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return fileModelToProto(file), nil
}

func (s *Server) GetFileByPath(ctx context.Context, req *protos.GetFileByPathRequest) (*protos.File, error) {
	file, err := s.Repo.GetFileByPath(ctx, req.OwnerId, req.Path)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return fileModelToProto(file), nil
}

func (s *Server) UpdateFile(ctx context.Context, req *protos.File) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateFile(ctx, protoToFileModel(req)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteFile(ctx context.Context, req *protos.FileID) (*emptypb.Empty, error) {
	if err := s.Repo.DeleteFile(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) SoftDeleteFile(ctx context.Context, req *protos.FileID) (*emptypb.Empty, error) {
	if err := s.Repo.SoftDeleteFile(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) RestoreFile(ctx context.Context, req *protos.FileID) (*emptypb.Empty, error) {
	if err := s.Repo.RestoreFile(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListFiles(ctx context.Context, req *protos.ListFilesRequest) (*protos.ListFilesResponse, error) {
	files, total, err := s.Repo.ListFiles(ctx, req.ParentId, req.OwnerId, req.IsTrashed, req.Starred, int(req.Limit), int(req.Offset), req.OrderBy, req.OrderDir)
	if err != nil {
		return nil, err
	}

	protoFiles := make([]*protos.File, len(files))
	for i, file := range files {
		protoFiles[i] = fileModelToProto(file)
	}

	return &protos.ListFilesResponse{
		Files:  protoFiles,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

func (s *Server) ListFilesByParent(ctx context.Context, req *protos.ListFilesByParentRequest) (*protos.ListFilesResponse, error) {
	files, err := s.Repo.ListFilesByParent(ctx, req.OwnerId, req.ParentId)
	if err != nil {
		return nil, err
	}

	protoFiles := make([]*protos.File, len(files))
	for i, file := range files {
		protoFiles[i] = fileModelToProto(file)
	}

	return &protos.ListFilesResponse{
		Files: protoFiles,
		Total: int64(len(files)),
	}, nil
}

func (s *Server) ListStarredFiles(ctx context.Context, req *protos.ListStarredFilesRequest) (*protos.ListFilesResponse, error) {
	files, err := s.Repo.ListStarredFiles(ctx, req.OwnerId)
	if err != nil {
		return nil, err
	}

	protoFiles := make([]*protos.File, len(files))
	for i, file := range files {
		protoFiles[i] = fileModelToProto(file)
	}

	return &protos.ListFilesResponse{
		Files: protoFiles,
		Total: int64(len(files)),
	}, nil
}

func (s *Server) ListTrashedFiles(ctx context.Context, req *protos.ListTrashedFilesRequest) (*protos.ListFilesResponse, error) {
	files, err := s.Repo.ListTrashedFiles(ctx, req.OwnerId)
	if err != nil {
		return nil, err
	}

	protoFiles := make([]*protos.File, len(files))
	for i, file := range files {
		protoFiles[i] = fileModelToProto(file)
	}

	return &protos.ListFilesResponse{
		Files: protoFiles,
		Total: int64(len(files)),
	}, nil
}

func (s *Server) SearchFiles(ctx context.Context, req *protos.SearchFilesRequest) (*protos.ListFilesResponse, error) {
	files, err := s.Repo.SearchFiles(ctx, req.OwnerId, req.Query)
	if err != nil {
		return nil, err
	}

	protoFiles := make([]*protos.File, len(files))
	for i, file := range files {
		protoFiles[i] = fileModelToProto(file)
	}

	return &protos.ListFilesResponse{
		Files: protoFiles,
		Total: int64(len(files)),
	}, nil
}

func (s *Server) GetFileSize(ctx context.Context, req *protos.FileID) (*protos.FileSizeResponse, error) {
	size, err := s.Repo.GetFileSize(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &protos.FileSizeResponse{Size: size}, nil
}

func (s *Server) UpdateFileSize(ctx context.Context, req *protos.UpdateFileSizeRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateFileSize(ctx, req.Id, req.Size); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateLastViewed(ctx context.Context, req *protos.FileID) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateLastViewed(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetFileTree(ctx context.Context, req *protos.GetFileTreeRequest) (*protos.ListFilesResponse, error) {
	files, err := s.Repo.GetFileTree(ctx, req.OwnerId, req.RootId)
	if err != nil {
		return nil, err
	}

	protoFiles := make([]*protos.File, len(files))
	for i, file := range files {
		protoFiles[i] = fileModelToProto(file)
	}

	return &protos.ListFilesResponse{
		Files: protoFiles,
		Total: int64(len(files)),
	}, nil
}

// File revision operations
func (s *Server) CreateRevision(ctx context.Context, req *protos.FileRevision) (*protos.RevisionID, error) {
	revision := protoToFileRevisionModel(req)
	id, err := s.Repo.CreateRevision(ctx, revision)
	if err != nil {
		return nil, err
	}
	return &protos.RevisionID{Id: id}, nil
}

func (s *Server) GetRevisions(ctx context.Context, req *protos.FileID) (*protos.ListRevisionsResponse, error) {
	revisions, err := s.Repo.GetRevisions(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	protoRevisions := make([]*protos.FileRevision, len(revisions))
	for i, revision := range revisions {
		protoRevisions[i] = fileRevisionModelToProto(revision)
	}

	return &protos.ListRevisionsResponse{Revisions: protoRevisions}, nil
}

func (s *Server) GetRevision(ctx context.Context, req *protos.GetRevisionRequest) (*protos.FileRevision, error) {
	revision, err := s.Repo.GetRevision(ctx, req.FileId, req.RevisionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return fileRevisionModelToProto(revision), nil
}

func (s *Server) DeleteRevision(ctx context.Context, req *protos.RevisionID) (*emptypb.Empty, error) {
	if err := s.Repo.DeleteRevision(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// File permission operations
func (s *Server) CreatePermission(ctx context.Context, req *protos.FilePermission) (*protos.PermissionID, error) {
	permission := protoToFilePermissionModel(req)
	id, err := s.Repo.CreatePermission(ctx, permission)
	if err != nil {
		return nil, err
	}
	return &protos.PermissionID{Id: id}, nil
}

func (s *Server) GetPermissions(ctx context.Context, req *protos.FileID) (*protos.ListPermissionsResponse, error) {
	permissions, err := s.Repo.GetPermissions(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	protoPermissions := make([]*protos.FilePermission, len(permissions))
	for i, permission := range permissions {
		protoPermissions[i] = filePermissionModelToProto(permission)
	}

	return &protos.ListPermissionsResponse{Permissions: protoPermissions}, nil
}

func (s *Server) UpdatePermission(ctx context.Context, req *protos.FilePermission) (*emptypb.Empty, error) {
	if err := s.Repo.UpdatePermission(ctx, protoToFilePermissionModel(req)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeletePermission(ctx context.Context, req *protos.PermissionID) (*emptypb.Empty, error) {
	if err := s.Repo.DeletePermission(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) CheckPermission(ctx context.Context, req *protos.CheckPermissionRequest) (*protos.PermissionResponse, error) {
	hasPermission, err := s.Repo.CheckPermission(ctx, req.FileId, req.UserId, req.RequiredRole)
	if err != nil {
		return nil, err
	}
	return &protos.PermissionResponse{HasPermission: hasPermission}, nil
}

// File metadata operations
func (s *Server) UpdateFileMetadata(ctx context.Context, req *protos.UpdateFileMetadataRequest) (*emptypb.Empty, error) {
	if err := s.Repo.UpdateFileMetadata(ctx, req.FileId, req.Metadata); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetFileMetadata(ctx context.Context, req *protos.FileID) (*protos.FileMetadataResponse, error) {
	metadata, err := s.Repo.GetFileMetadata(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &protos.FileMetadataResponse{Metadata: metadata}, nil
}

// File operations (star, move, copy, rename)
func (s *Server) StarFile(ctx context.Context, req *protos.FileID) (*emptypb.Empty, error) {
	if err := s.Repo.StarFile(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UnstarFile(ctx context.Context, req *protos.FileID) (*emptypb.Empty, error) {
	if err := s.Repo.UnstarFile(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) MoveFile(ctx context.Context, req *protos.MoveFileRequest) (*emptypb.Empty, error) {
	if err := s.Repo.MoveFile(ctx, req.FileId, req.NewParentId); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) CopyFile(ctx context.Context, req *protos.CopyFileRequest) (*protos.File, error) {
	file, err := s.Repo.CopyFile(ctx, req.FileId, req.NewParentId, req.NewName)
	if err != nil {
		return nil, err
	}
	return fileModelToProto(file), nil
}

func (s *Server) RenameFile(ctx context.Context, req *protos.RenameFileRequest) (*emptypb.Empty, error) {
	if err := s.Repo.RenameFile(ctx, req.FileId, req.NewName); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// File integrity operations
func (s *Server) VerifyFileIntegrity(ctx context.Context, req *protos.FileID) (*protos.IntegrityResponse, error) {
	isVerified, err := s.Repo.VerifyFileIntegrity(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &protos.IntegrityResponse{IsIntegrityVerified: isVerified}, nil
}

func (s *Server) CalculateFileChecksums(ctx context.Context, req *protos.FileID) (*protos.ChecksumsResponse, error) {
	checksums, err := s.Repo.CalculateFileChecksums(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &protos.ChecksumsResponse{Checksums: checksums}, nil
}

// Helper functions for converting between models and protos
func fileModelToProto(f *models.File) *protos.File {
	if f == nil {
		return nil
	}

	// Helper function to safely convert *string to string
	safeString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	return &protos.File{
		Id:             f.ID,
		OwnerId:        f.OwnerID,
		ParentId:       safeString(f.ParentID),
		Name:           f.Name,
		FileExtension:  safeString(f.FileExtension),
		MimeType:       f.MimeType,
		StoragePath:    f.StoragePath,
		Size:           f.Size,
		Md5Checksum:    safeString(f.MD5Checksum),
		Sha256Checksum: safeString(f.SHA256Checksum),
		IsFolder:       f.IsFolder,
		IsTrashed:      f.IsTrashed,
		TrashedAt:      timeToProto(f.TrashedAt),
		Starred:        f.Starred,
		CreatedAt:      timestamppb.New(f.CreatedAt),
		UpdatedAt:      timestamppb.New(f.UpdatedAt),
		LastViewedAt:   timeToProto(f.LastViewedAt),
		ViewedByMe:     f.ViewedByMe,
		Version:        f.Version,
		RevisionId:     safeString(f.RevisionID),
		IndexableText:  safeString(f.IndexableText),
		ThumbnailLink:  safeString(f.ThumbnailLink),
		WebViewLink:    safeString(f.WebViewLink),
		WebContentLink: safeString(f.WebContentLink),
		IconLink:       safeString(f.IconLink),
	}
}

func protoToFileModel(f *protos.File) *models.File {
	if f == nil {
		return nil
	}

	// Helper function to safely convert string to *string
	safeStringPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	return &models.File{
		ID:             f.Id,
		OwnerID:        f.OwnerId,
		ParentID:       safeStringPtr(f.ParentId),
		Name:           f.Name,
		FileExtension:  safeStringPtr(f.FileExtension),
		MimeType:       f.MimeType,
		StoragePath:    f.StoragePath,
		Size:           f.Size,
		MD5Checksum:    safeStringPtr(f.Md5Checksum),
		SHA256Checksum: safeStringPtr(f.Sha256Checksum),
		IsFolder:       f.IsFolder,
		IsTrashed:      f.IsTrashed,
		TrashedAt:      protoToTime(f.TrashedAt),
		Starred:        f.Starred,
		CreatedAt:      f.CreatedAt.AsTime(),
		UpdatedAt:      f.UpdatedAt.AsTime(),
		LastViewedAt:   protoToTime(f.LastViewedAt),
		ViewedByMe:     f.ViewedByMe,
		Version:        f.Version,
		RevisionID:     safeStringPtr(f.RevisionId),
		IndexableText:  safeStringPtr(f.IndexableText),
		ThumbnailLink:  safeStringPtr(f.ThumbnailLink),
		WebViewLink:    safeStringPtr(f.WebViewLink),
		WebContentLink: safeStringPtr(f.WebContentLink),
		IconLink:       safeStringPtr(f.IconLink),
	}
}

func fileRevisionModelToProto(fr *models.FileRevision) *protos.FileRevision {
	if fr == nil {
		return nil
	}

	// Helper function to safely convert *string to string
	safeString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	return &protos.FileRevision{
		Id:          fr.ID,
		FileId:      fr.FileID,
		RevisionId:  fr.RevisionID,
		Md5Checksum: safeString(fr.MD5Checksum),
		Size:        fr.Size,
		CreatedAt:   timestamppb.New(fr.CreatedAt),
		StoragePath: fr.StoragePath,
		MimeType:    safeString(fr.MimeType),
		UserId:      safeString(fr.UserID),
	}
}

func protoToFileRevisionModel(fr *protos.FileRevision) *models.FileRevision {
	if fr == nil {
		return nil
	}

	// Helper function to safely convert string to *string
	safeStringPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	return &models.FileRevision{
		ID:          fr.Id,
		FileID:      fr.FileId,
		RevisionID:  fr.RevisionId,
		MD5Checksum: safeStringPtr(fr.Md5Checksum),
		Size:        fr.Size,
		CreatedAt:   fr.CreatedAt.AsTime(),
		StoragePath: fr.StoragePath,
		MimeType:    safeStringPtr(fr.MimeType),
		UserID:      safeStringPtr(fr.UserId),
	}
}

func filePermissionModelToProto(fp *models.FilePermission) *protos.FilePermission {
	if fp == nil {
		return nil
	}

	// Helper function to safely convert *string to string
	safeString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	return &protos.FilePermission{
		Id:          fp.ID,
		FileId:      fp.FileID,
		GranteeId:   safeString(fp.GranteeID),
		GranteeType: fp.GranteeType,
		Role:        fp.Role,
		AllowShare:  fp.AllowShare,
		CreatedAt:   timestamppb.New(fp.CreatedAt),
	}
}

func protoToFilePermissionModel(fp *protos.FilePermission) *models.FilePermission {
	if fp == nil {
		return nil
	}

	// Helper function to safely convert string to *string
	safeStringPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	return &models.FilePermission{
		ID:          fp.Id,
		FileID:      fp.FileId,
		GranteeID:   safeStringPtr(fp.GranteeId),
		GranteeType: fp.GranteeType,
		Role:        fp.Role,
		AllowShare:  fp.AllowShare,
		CreatedAt:   fp.CreatedAt.AsTime(),
	}
}

func (s *Server) GetUserExtendedInfo(ctx context.Context, req *protos.UserID) (*protos.UserExtendedInfo, error) {
	s.Logger.Info(ctx, "GetUserExtendedInfo called", zap.String("user_id", req.Id))
	u, err := s.Repo.GetUserByID(ctx, req.Id)
	if err != nil {
		s.Logger.Error(ctx, "GetUserByID failed", zap.Error(err))
		return nil, err
	}
	if u == nil {
		s.Logger.Info(ctx, "User not found", zap.String("user_id", req.Id))
		return nil, nil
	}

	// Вычисляемые поля
	var (
		usagePct                             float64
		usageFmt, remainFmt                  string
		isLocked, isQuotaExceeded            bool
		daysSinceLastLogin, daysSinceCreated int32
		accountStatus, securityStatus        string
		warnings, recommendations            []string
		metadata                             = map[string]string{}
	)

	if u.StorageQuota > 0 {
		usagePct = float64(u.UsedSpace) / float64(u.StorageQuota) * 100
		usageFmt = formatBytes(u.UsedSpace) + " / " + formatBytes(u.StorageQuota)
		remainFmt = formatBytes(u.StorageQuota - u.UsedSpace)
		isQuotaExceeded = u.UsedSpace > u.StorageQuota
		if usagePct > 90 {
			warnings = append(warnings, "Почти закончилось место")
		}
	} else {
		usageFmt = formatBytes(u.UsedSpace)
		remainFmt = "∞"
	}

	isLocked = u.LockedUntil != nil && u.LockedUntil.After(time.Now())
	if isLocked {
		accountStatus = "locked"
		warnings = append(warnings, "Аккаунт заблокирован")
	} else if !u.IsActive {
		accountStatus = "inactive"
	} else {
		accountStatus = "active"
	}

	if !u.IsEmailVerified {
		securityStatus = "needs_verification"
		recommendations = append(recommendations, "Подтвердите email")
	} else {
		securityStatus = "secure"
	}

	if u.FailedLoginAttempts > 3 {
		warnings = append(warnings, "Много неудачных попыток входа")
	}

	if u.LastLogin != nil {
		daysSinceLastLogin = int32(math.Round(time.Since(*u.LastLogin).Hours() / 24))
	} else {
		daysSinceLastLogin = -1
	}
	daysSinceCreated = int32(math.Round(time.Since(u.CreatedAt).Hours() / 24))

	if u.UsedSpace < u.StorageQuota/10 {
		recommendations = append(recommendations, "Загрузите больше файлов!")
	}

	metadata["role"] = u.Role
	metadata["email"] = u.Email
	metadata["username"] = u.Username

	resp := &protos.UserExtendedInfo{
		User:                      userModelToProto(u),
		StorageUsagePercentage:    usagePct,
		StorageUsageFormatted:     usageFmt,
		StorageRemainingFormatted: remainFmt,
		IsLocked:                  isLocked,
		IsQuotaExceeded:           isQuotaExceeded,
		DaysSinceLastLogin:        daysSinceLastLogin,
		DaysSinceCreated:          daysSinceCreated,
		AccountStatus:             accountStatus,
		SecurityStatus:            securityStatus,
		Warnings:                  warnings,
		Recommendations:           recommendations,
		Metadata:                  metadata,
	}
	s.Logger.Debug(ctx, "GetUserExtendedInfo response", zap.Any("response", resp))
	return resp, nil
}

// formatBytes - форматирует байты в строку (например, 1.5 GB)
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPEZY"[exp])
}
