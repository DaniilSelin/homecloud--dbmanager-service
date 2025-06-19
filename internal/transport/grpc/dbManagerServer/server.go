package dbManagerServer

import (
	"context"
	"database/sql"
	"time"

	"homecloud--dbmanager-service/internal/interfaces"
	"homecloud--dbmanager-service/internal/models"
	protos "homecloud--dbmanager-service/internal/transport/grpc/protos"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	protos.UnimplementedDBServiceServer
	Repo interfaces.DBRepository
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
	u, err := s.Repo.GetUserByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return userModelToProto(u), nil
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
