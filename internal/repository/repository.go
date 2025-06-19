package repository

import (
	"context"
	"database/sql"
	"time"

	"homecloud--dbmanager-service/internal/interfaces"
	"homecloud--dbmanager-service/internal/models"
)

type dbRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) interfaces.DBRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	query := `INSERT INTO dbmanager.users (email, username, password_hash, is_active, is_email_verified, role, storage_quota, used_space, created_at, updated_at, failed_login_attempts, locked_until, last_login_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW(),$9,$10,$11) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		user.Email, user.Username, user.PasswordHash, user.IsActive, user.IsEmailVerified, user.Role, user.StorageQuota, user.UsedSpace, user.FailedLoginAttempts, user.LockedUntil, user.LastLogin,
	).Scan(&id)
	return id, err
}

func (r *dbRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, is_active, is_email_verified, role, storage_quota, used_space, created_at, updated_at, failed_login_attempts, locked_until, last_login_at FROM dbmanager.users WHERE id=$1`
	user := &models.User{}
	var lockedUntil, lastLogin sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsActive, &user.IsEmailVerified, &user.Role, &user.StorageQuota, &user.UsedSpace, &user.CreatedAt, &user.UpdatedAt, &user.FailedLoginAttempts, &lockedUntil, &lastLogin,
	)
	if err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		user.LockedUntil = &lockedUntil.Time
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	return user, nil
}

func (r *dbRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, is_active, is_email_verified, role, storage_quota, used_space, created_at, updated_at, failed_login_attempts, locked_until, last_login_at FROM dbmanager.users WHERE email=$1`
	user := &models.User{}
	var lockedUntil, lastLogin sql.NullTime
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsActive, &user.IsEmailVerified, &user.Role, &user.StorageQuota, &user.UsedSpace, &user.CreatedAt, &user.UpdatedAt, &user.FailedLoginAttempts, &lockedUntil, &lastLogin,
	)
	if err != nil {
		return nil, err
	}
	if lockedUntil.Valid {
		user.LockedUntil = &lockedUntil.Time
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	return user, nil
}

func (r *dbRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE dbmanager.users SET email=$1, username=$2, password_hash=$3, is_active=$4, is_email_verified=$5, role=$6, storage_quota=$7, used_space=$8, updated_at=NOW(), failed_login_attempts=$9, locked_until=$10, last_login_at=$11 WHERE id=$12`
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.Username, user.PasswordHash, user.IsActive, user.IsEmailVerified, user.Role, user.StorageQuota, user.UsedSpace, user.FailedLoginAttempts, user.LockedUntil, user.LastLogin, user.ID,
	)
	return err
}

func (r *dbRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbmanager.users SET password_hash=$1, updated_at=NOW() WHERE id=$2`, passwordHash, id)
	return err
}

func (r *dbRepository) UpdateUsername(ctx context.Context, id, username string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbmanager.users SET username=$1, updated_at=NOW() WHERE id=$2`, username, id)
	return err
}

func (r *dbRepository) UpdateEmailVerification(ctx context.Context, id string, isVerified bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbmanager.users SET is_email_verified=$1, updated_at=NOW() WHERE id=$2`, isVerified, id)
	return err
}

func (r *dbRepository) UpdateLastLogin(ctx context.Context, id string, lastLogin time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbmanager.users SET last_login_at=$1, updated_at=NOW() WHERE id=$2`, lastLogin, id)
	return err
}

func (r *dbRepository) UpdateFailedLoginAttempts(ctx context.Context, id string, attempts int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbmanager.users SET failed_login_attempts=$1, updated_at=NOW() WHERE id=$2`, attempts, id)
	return err
}

func (r *dbRepository) UpdateLockedUntil(ctx context.Context, id string, lockedUntil time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbmanager.users SET locked_until=$1, updated_at=NOW() WHERE id=$2`, lockedUntil, id)
	return err
}

func (r *dbRepository) UpdateStorageUsage(ctx context.Context, id string, usedSpace int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE dbmanager.users SET used_space=$1, updated_at=NOW() WHERE id=$2`, usedSpace, id)
	return err
}

func (r *dbRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM dbmanager.users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

func (r *dbRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM dbmanager.users WHERE username=$1)`, username).Scan(&exists)
	return exists, err
}
