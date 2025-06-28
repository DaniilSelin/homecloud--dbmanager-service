package repository

import (
	"context"
	"database/sql"
	"fmt"
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
	query := `INSERT INTO homecloud.users (id, email, username, password_hash, is_active, is_email_verified, role, storage_quota, used_space, created_at, updated_at, failed_login_attempts, locked_until, last_login_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW(),NOW(),$10,$11,$12) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.Email, user.Username, user.PasswordHash, user.IsActive, user.IsEmailVerified, user.Role, user.StorageQuota, user.UsedSpace, user.FailedLoginAttempts, user.LockedUntil, user.LastLogin,
	).Scan(&id)
	return id, err
}

func (r *dbRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, is_active, is_email_verified, role, storage_quota, used_space, created_at, updated_at, failed_login_attempts, locked_until, last_login_at FROM homecloud.users WHERE id=$1`
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
	query := `SELECT id, email, username, password_hash, is_active, is_email_verified, role, storage_quota, used_space, created_at, updated_at, failed_login_attempts, locked_until, last_login_at FROM homecloud.users WHERE email=$1`
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
	query := `UPDATE homecloud.users SET email=$1, username=$2, password_hash=$3, is_active=$4, is_email_verified=$5, role=$6, storage_quota=$7, used_space=$8, updated_at=NOW(), failed_login_attempts=$9, locked_until=$10, last_login_at=$11 WHERE id=$12`
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.Username, user.PasswordHash, user.IsActive, user.IsEmailVerified, user.Role, user.StorageQuota, user.UsedSpace, user.FailedLoginAttempts, user.LockedUntil, user.LastLogin, user.ID,
	)
	return err
}

func (r *dbRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.users SET password_hash=$1, updated_at=NOW() WHERE id=$2`, passwordHash, id)
	return err
}

func (r *dbRepository) UpdateUsername(ctx context.Context, id, username string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.users SET username=$1, updated_at=NOW() WHERE id=$2`, username, id)
	return err
}

func (r *dbRepository) UpdateEmailVerification(ctx context.Context, id string, isVerified bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.users SET is_email_verified=$1, updated_at=NOW() WHERE id=$2`, isVerified, id)
	return err
}

func (r *dbRepository) UpdateLastLogin(ctx context.Context, id string, lastLogin time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.users SET last_login_at=$1, updated_at=NOW() WHERE id=$2`, lastLogin, id)
	return err
}

func (r *dbRepository) UpdateFailedLoginAttempts(ctx context.Context, id string, attempts int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.users SET failed_login_attempts=$1, updated_at=NOW() WHERE id=$2`, attempts, id)
	return err
}

func (r *dbRepository) UpdateLockedUntil(ctx context.Context, id string, lockedUntil time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.users SET locked_until=$1, updated_at=NOW() WHERE id=$2`, lockedUntil, id)
	return err
}

func (r *dbRepository) UpdateStorageUsage(ctx context.Context, id string, usedSpace int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.users SET used_space=$1, updated_at=NOW() WHERE id=$2`, usedSpace, id)
	return err
}

func (r *dbRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM homecloud.users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

func (r *dbRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM homecloud.users WHERE username=$1)`, username).Scan(&exists)
	return exists, err
}

// File operations
func (r *dbRepository) CreateFile(ctx context.Context, file *models.File) (string, error) {
	query := `INSERT INTO homecloud.files (owner_id, parent_id, name, file_extension, mime_type, storage_path, size, md5_checksum, sha256_checksum, is_folder, is_trashed, trashed_at, starred, created_at, updated_at, last_viewed_at, viewed_by_me, version, revision_id, indexable_text, thumbnail_link, web_view_link, web_content_link, icon_link)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW(), $14, $15, $16, $17, $18, $19, $20, $21, $22) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		file.OwnerID, file.ParentID, file.Name, file.FileExtension, file.MimeType, file.StoragePath, file.Size, file.MD5Checksum, file.SHA256Checksum, file.IsFolder, file.IsTrashed, file.TrashedAt, file.Starred, file.LastViewedAt, file.ViewedByMe, file.Version, file.RevisionID, file.IndexableText, file.ThumbnailLink, file.WebViewLink, file.WebContentLink, file.IconLink,
	).Scan(&id)
	return id, err
}

func (r *dbRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	query := `SELECT id, owner_id, parent_id, name, file_extension, mime_type, storage_path, size, md5_checksum, sha256_checksum, is_folder, is_trashed, trashed_at, starred, created_at, updated_at, last_viewed_at, viewed_by_me, version, revision_id, indexable_text, thumbnail_link, web_view_link, web_content_link, icon_link FROM homecloud.files WHERE id=$1`
	file := &models.File{}
	var parentID, fileExtension, md5Checksum, sha256Checksum, revisionID, indexableText, thumbnailLink, webViewLink, webContentLink, iconLink sql.NullString
	var trashedAt, lastViewedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID, &file.OwnerID, &parentID, &file.Name, &fileExtension, &file.MimeType, &file.StoragePath, &file.Size, &md5Checksum, &sha256Checksum, &file.IsFolder, &file.IsTrashed, &trashedAt, &file.Starred, &file.CreatedAt, &file.UpdatedAt, &lastViewedAt, &file.ViewedByMe, &file.Version, &revisionID, &indexableText, &thumbnailLink, &webViewLink, &webContentLink, &iconLink,
	)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		file.ParentID = &parentID.String
	}
	if fileExtension.Valid {
		file.FileExtension = &fileExtension.String
	}
	if md5Checksum.Valid {
		file.MD5Checksum = &md5Checksum.String
	}
	if sha256Checksum.Valid {
		file.SHA256Checksum = &sha256Checksum.String
	}
	if trashedAt.Valid {
		file.TrashedAt = &trashedAt.Time
	}
	if lastViewedAt.Valid {
		file.LastViewedAt = &lastViewedAt.Time
	}
	if revisionID.Valid {
		file.RevisionID = &revisionID.String
	}
	if indexableText.Valid {
		file.IndexableText = &indexableText.String
	}
	if thumbnailLink.Valid {
		file.ThumbnailLink = &thumbnailLink.String
	}
	if webViewLink.Valid {
		file.WebViewLink = &webViewLink.String
	}
	if webContentLink.Valid {
		file.WebContentLink = &webContentLink.String
	}
	if iconLink.Valid {
		file.IconLink = &iconLink.String
	}
	return file, nil
}

func (r *dbRepository) GetFileByPath(ctx context.Context, ownerID, path string) (*models.File, error) {
	// Если путь пустой, ищем корневую папку
	searchName := path
	if path == "" {
		searchName = "root"
	}

	// Простая реализация - поиск по имени файла в корне
	query := `SELECT id, owner_id, parent_id, name, file_extension, mime_type, storage_path, size, md5_checksum, sha256_checksum, is_folder, is_trashed, trashed_at, starred, created_at, updated_at, last_viewed_at, viewed_by_me, version, revision_id, indexable_text, thumbnail_link, web_view_link, web_content_link, icon_link FROM homecloud.files WHERE owner_id=$1 AND name=$2 AND parent_id IS NULL`
	file := &models.File{}
	var parentID, fileExtension, md5Checksum, sha256Checksum, revisionID, indexableText, thumbnailLink, webViewLink, webContentLink, iconLink sql.NullString
	var trashedAt, lastViewedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, ownerID, searchName).Scan(
		&file.ID, &file.OwnerID, &parentID, &file.Name, &fileExtension, &file.MimeType, &file.StoragePath, &file.Size, &md5Checksum, &sha256Checksum, &file.IsFolder, &file.IsTrashed, &trashedAt, &file.Starred, &file.CreatedAt, &file.UpdatedAt, &lastViewedAt, &file.ViewedByMe, &file.Version, &revisionID, &indexableText, &thumbnailLink, &webViewLink, &webContentLink, &iconLink,
	)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		file.ParentID = &parentID.String
	}
	if fileExtension.Valid {
		file.FileExtension = &fileExtension.String
	}
	if md5Checksum.Valid {
		file.MD5Checksum = &md5Checksum.String
	}
	if sha256Checksum.Valid {
		file.SHA256Checksum = &sha256Checksum.String
	}
	if trashedAt.Valid {
		file.TrashedAt = &trashedAt.Time
	}
	if lastViewedAt.Valid {
		file.LastViewedAt = &lastViewedAt.Time
	}
	if revisionID.Valid {
		file.RevisionID = &revisionID.String
	}
	if indexableText.Valid {
		file.IndexableText = &indexableText.String
	}
	if thumbnailLink.Valid {
		file.ThumbnailLink = &thumbnailLink.String
	}
	if webViewLink.Valid {
		file.WebViewLink = &webViewLink.String
	}
	if webContentLink.Valid {
		file.WebContentLink = &webContentLink.String
	}
	if iconLink.Valid {
		file.IconLink = &iconLink.String
	}
	return file, nil
}

func (r *dbRepository) UpdateFile(ctx context.Context, file *models.File) error {
	query := `UPDATE homecloud.files SET owner_id=$1, parent_id=$2, name=$3, file_extension=$4, mime_type=$5, storage_path=$6, size=$7, md5_checksum=$8, sha256_checksum=$9, is_folder=$10, is_trashed=$11, trashed_at=$12, starred=$13, updated_at=NOW(), last_viewed_at=$14, viewed_by_me=$15, version=$16, revision_id=$17, indexable_text=$18, thumbnail_link=$19, web_view_link=$20, web_content_link=$21, icon_link=$22 WHERE id=$23`
	_, err := r.db.ExecContext(ctx, query,
		file.OwnerID, file.ParentID, file.Name, file.FileExtension, file.MimeType, file.StoragePath, file.Size, file.MD5Checksum, file.SHA256Checksum, file.IsFolder, file.IsTrashed, file.TrashedAt, file.Starred, file.LastViewedAt, file.ViewedByMe, file.Version, file.RevisionID, file.IndexableText, file.ThumbnailLink, file.WebViewLink, file.WebContentLink, file.IconLink, file.ID,
	)
	return err
}

func (r *dbRepository) DeleteFile(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM homecloud.files WHERE id=$1`, id)
	return err
}

func (r *dbRepository) SoftDeleteFile(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET is_trashed=true, trashed_at=NOW(), updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *dbRepository) RestoreFile(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET is_trashed=false, trashed_at=NULL, updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *dbRepository) ListFiles(ctx context.Context, parentID, ownerID string, isTrashed, starred bool, limit, offset int, orderBy, orderDir string) ([]*models.File, int64, error) {
	baseQuery := `FROM homecloud.files WHERE owner_id=$1`
	args := []interface{}{ownerID}
	argIndex := 2

	// Добавляем фильтры
	if parentID != "" {
		baseQuery += fmt.Sprintf(" AND parent_id=$%d", argIndex)
		args = append(args, parentID)
		argIndex++
	} else {
		baseQuery += " AND parent_id IS NULL"
	}

	if isTrashed {
		baseQuery += fmt.Sprintf(" AND is_trashed=$%d", argIndex)
		args = append(args, true)
		argIndex++
	} else {
		baseQuery += " AND is_trashed=false"
	}

	if starred {
		baseQuery += fmt.Sprintf(" AND starred=$%d", argIndex)
		args = append(args, true)
		argIndex++
	}

	// Получаем общее количество
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Получаем файлы
	selectQuery := `SELECT id, owner_id, parent_id, name, file_extension, mime_type, storage_path, size, md5_checksum, sha256_checksum, is_folder, is_trashed, trashed_at, starred, created_at, updated_at, last_viewed_at, viewed_by_me, version, revision_id, indexable_text, thumbnail_link, web_view_link, web_content_link, icon_link ` + baseQuery

	// Добавляем сортировку
	if orderBy != "" {
		selectQuery += " ORDER BY " + orderBy
		if orderDir != "" {
			selectQuery += " " + orderDir
		}
	} else {
		selectQuery += " ORDER BY updated_at DESC"
	}

	// Добавляем лимит и оффсет
	if limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}
	if offset > 0 {
		selectQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		file := &models.File{}
		var parentID, fileExtension, md5Checksum, sha256Checksum, revisionID, indexableText, thumbnailLink, webViewLink, webContentLink, iconLink sql.NullString
		var trashedAt, lastViewedAt sql.NullTime
		err := rows.Scan(
			&file.ID, &file.OwnerID, &parentID, &file.Name, &fileExtension, &file.MimeType, &file.StoragePath, &file.Size, &md5Checksum, &sha256Checksum, &file.IsFolder, &file.IsTrashed, &trashedAt, &file.Starred, &file.CreatedAt, &file.UpdatedAt, &lastViewedAt, &file.ViewedByMe, &file.Version, &revisionID, &indexableText, &thumbnailLink, &webViewLink, &webContentLink, &iconLink,
		)
		if err != nil {
			return nil, 0, err
		}
		if parentID.Valid {
			file.ParentID = &parentID.String
		}
		if fileExtension.Valid {
			file.FileExtension = &fileExtension.String
		}
		if md5Checksum.Valid {
			file.MD5Checksum = &md5Checksum.String
		}
		if sha256Checksum.Valid {
			file.SHA256Checksum = &sha256Checksum.String
		}
		if trashedAt.Valid {
			file.TrashedAt = &trashedAt.Time
		}
		if lastViewedAt.Valid {
			file.LastViewedAt = &lastViewedAt.Time
		}
		if revisionID.Valid {
			file.RevisionID = &revisionID.String
		}
		if indexableText.Valid {
			file.IndexableText = &indexableText.String
		}
		if thumbnailLink.Valid {
			file.ThumbnailLink = &thumbnailLink.String
		}
		if webViewLink.Valid {
			file.WebViewLink = &webViewLink.String
		}
		if webContentLink.Valid {
			file.WebContentLink = &webContentLink.String
		}
		if iconLink.Valid {
			file.IconLink = &iconLink.String
		}
		files = append(files, file)
	}

	return files, total, nil
}

func (r *dbRepository) ListFilesByParent(ctx context.Context, ownerID, parentID string) ([]*models.File, error) {
	files, _, err := r.ListFiles(ctx, parentID, ownerID, false, false, 1000, 0, "name", "ASC")
	return files, err
}

func (r *dbRepository) ListStarredFiles(ctx context.Context, ownerID string) ([]*models.File, error) {
	files, _, err := r.ListFiles(ctx, "", ownerID, false, true, 1000, 0, "updated_at", "DESC")
	return files, err
}

func (r *dbRepository) ListTrashedFiles(ctx context.Context, ownerID string) ([]*models.File, error) {
	files, _, err := r.ListFiles(ctx, "", ownerID, true, false, 1000, 0, "trashed_at", "DESC")
	return files, err
}

func (r *dbRepository) SearchFiles(ctx context.Context, ownerID, query string) ([]*models.File, error) {
	searchQuery := `SELECT id, owner_id, parent_id, name, file_extension, mime_type, storage_path, size, md5_checksum, sha256_checksum, is_folder, is_trashed, trashed_at, starred, created_at, updated_at, last_viewed_at, viewed_by_me, version, revision_id, indexable_text, thumbnail_link, web_view_link, web_content_link, icon_link FROM homecloud.files WHERE owner_id=$1 AND is_trashed=false AND (name ILIKE $2 OR indexable_text ILIKE $2) ORDER BY updated_at DESC`
	rows, err := r.db.QueryContext(ctx, searchQuery, ownerID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		file := &models.File{}
		var parentID, fileExtension, md5Checksum, sha256Checksum, revisionID, indexableText, thumbnailLink, webViewLink, webContentLink, iconLink sql.NullString
		var trashedAt, lastViewedAt sql.NullTime
		err := rows.Scan(
			&file.ID, &file.OwnerID, &parentID, &file.Name, &fileExtension, &file.MimeType, &file.StoragePath, &file.Size, &md5Checksum, &sha256Checksum, &file.IsFolder, &file.IsTrashed, &trashedAt, &file.Starred, &file.CreatedAt, &file.UpdatedAt, &lastViewedAt, &file.ViewedByMe, &file.Version, &revisionID, &indexableText, &thumbnailLink, &webViewLink, &webContentLink, &iconLink,
		)
		if err != nil {
			return nil, err
		}
		if parentID.Valid {
			file.ParentID = &parentID.String
		}
		if fileExtension.Valid {
			file.FileExtension = &fileExtension.String
		}
		if md5Checksum.Valid {
			file.MD5Checksum = &md5Checksum.String
		}
		if sha256Checksum.Valid {
			file.SHA256Checksum = &sha256Checksum.String
		}
		if trashedAt.Valid {
			file.TrashedAt = &trashedAt.Time
		}
		if lastViewedAt.Valid {
			file.LastViewedAt = &lastViewedAt.Time
		}
		if revisionID.Valid {
			file.RevisionID = &revisionID.String
		}
		if indexableText.Valid {
			file.IndexableText = &indexableText.String
		}
		if thumbnailLink.Valid {
			file.ThumbnailLink = &thumbnailLink.String
		}
		if webViewLink.Valid {
			file.WebViewLink = &webViewLink.String
		}
		if webContentLink.Valid {
			file.WebContentLink = &webContentLink.String
		}
		if iconLink.Valid {
			file.IconLink = &iconLink.String
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *dbRepository) GetFileSize(ctx context.Context, id string) (int64, error) {
	var size int64
	err := r.db.QueryRowContext(ctx, `SELECT size FROM homecloud.files WHERE id=$1`, id).Scan(&size)
	return size, err
}

func (r *dbRepository) UpdateFileSize(ctx context.Context, id string, size int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET size=$1, updated_at=NOW() WHERE id=$2`, size, id)
	return err
}

func (r *dbRepository) UpdateLastViewed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET last_viewed_at=NOW(), viewed_by_me=true, updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *dbRepository) GetFileTree(ctx context.Context, ownerID, rootID string) ([]*models.File, error) {
	// Простая реализация - получаем все файлы пользователя
	files, _, err := r.ListFiles(ctx, "", ownerID, false, false, 10000, 0, "name", "ASC")
	return files, err
}

// File revision operations
func (r *dbRepository) CreateRevision(ctx context.Context, revision *models.FileRevision) (string, error) {
	query := `INSERT INTO homecloud.file_revisions (file_id, revision_id, md5_checksum, size, created_at, storage_path, mime_type, user_id)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		revision.FileID, revision.RevisionID, revision.MD5Checksum, revision.Size, revision.StoragePath, revision.MimeType, revision.UserID,
	).Scan(&id)
	return id, err
}

func (r *dbRepository) GetRevisions(ctx context.Context, fileID string) ([]*models.FileRevision, error) {
	query := `SELECT id, file_id, revision_id, md5_checksum, size, created_at, storage_path, mime_type, user_id FROM homecloud.file_revisions WHERE file_id=$1 ORDER BY revision_id DESC`
	rows, err := r.db.QueryContext(ctx, query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revisions []*models.FileRevision
	for rows.Next() {
		revision := &models.FileRevision{}
		err := rows.Scan(
			&revision.ID, &revision.FileID, &revision.RevisionID, &revision.MD5Checksum, &revision.Size, &revision.CreatedAt, &revision.StoragePath, &revision.MimeType, &revision.UserID,
		)
		if err != nil {
			return nil, err
		}
		revisions = append(revisions, revision)
	}

	return revisions, nil
}

func (r *dbRepository) GetRevision(ctx context.Context, fileID string, revisionID int64) (*models.FileRevision, error) {
	query := `SELECT id, file_id, revision_id, md5_checksum, size, created_at, storage_path, mime_type, user_id FROM homecloud.file_revisions WHERE file_id=$1 AND revision_id=$2`
	revision := &models.FileRevision{}
	err := r.db.QueryRowContext(ctx, query, fileID, revisionID).Scan(
		&revision.ID, &revision.FileID, &revision.RevisionID, &revision.MD5Checksum, &revision.Size, &revision.CreatedAt, &revision.StoragePath, &revision.MimeType, &revision.UserID,
	)
	if err != nil {
		return nil, err
	}
	return revision, nil
}

func (r *dbRepository) DeleteRevision(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM homecloud.file_revisions WHERE id=$1`, id)
	return err
}

// File permission operations
func (r *dbRepository) CreatePermission(ctx context.Context, permission *models.FilePermission) (string, error) {
	query := `INSERT INTO homecloud.file_permissions (file_id, grantee_id, grantee_type, role, allow_share, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		permission.FileID, permission.GranteeID, permission.GranteeType, permission.Role, permission.AllowShare,
	).Scan(&id)
	return id, err
}

func (r *dbRepository) GetPermissions(ctx context.Context, fileID string) ([]*models.FilePermission, error) {
	query := `SELECT id, file_id, grantee_id, grantee_type, role, allow_share, created_at FROM homecloud.file_permissions WHERE file_id=$1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*models.FilePermission
	for rows.Next() {
		permission := &models.FilePermission{}
		err := rows.Scan(
			&permission.ID, &permission.FileID, &permission.GranteeID, &permission.GranteeType, &permission.Role, &permission.AllowShare, &permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *dbRepository) UpdatePermission(ctx context.Context, permission *models.FilePermission) error {
	query := `UPDATE homecloud.file_permissions SET grantee_id=$1, grantee_type=$2, role=$3, allow_share=$4 WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query,
		permission.GranteeID, permission.GranteeType, permission.Role, permission.AllowShare, permission.ID,
	)
	return err
}

func (r *dbRepository) DeletePermission(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM homecloud.file_permissions WHERE id=$1`, id)
	return err
}

func (r *dbRepository) CheckPermission(ctx context.Context, fileID, userID, requiredRole string) (bool, error) {
	// Определяем иерархию ролей
	var roleHierarchy string
	switch requiredRole {
	case "READER":
		roleHierarchy = "'READER', 'COMMENTER', 'WRITER', 'FILE_OWNER', 'ORGANIZER', 'OWNER'"
	case "COMMENTER":
		roleHierarchy = "'COMMENTER', 'WRITER', 'FILE_OWNER', 'ORGANIZER', 'OWNER'"
	case "WRITER":
		roleHierarchy = "'WRITER', 'FILE_OWNER', 'ORGANIZER', 'OWNER'"
	case "FILE_OWNER":
		roleHierarchy = "'FILE_OWNER', 'ORGANIZER', 'OWNER'"
	case "ORGANIZER":
		roleHierarchy = "'ORGANIZER', 'OWNER'"
	case "OWNER":
		roleHierarchy = "'OWNER'"
	default:
		roleHierarchy = fmt.Sprintf("'%s'", requiredRole)
	}

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM homecloud.file_permissions WHERE file_id=$1 AND grantee_id=$2 AND role IN (%s))`, roleHierarchy)
	var exists bool
	err := r.db.QueryRowContext(ctx, query, fileID, userID).Scan(&exists)
	return exists, err
}

// File metadata operations
func (r *dbRepository) UpdateFileMetadata(ctx context.Context, fileID, metadata string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET indexable_text=$1, updated_at=NOW() WHERE id=$2`, metadata, fileID)
	return err
}

func (r *dbRepository) GetFileMetadata(ctx context.Context, fileID string) (string, error) {
	var metadata string
	err := r.db.QueryRowContext(ctx, `SELECT indexable_text FROM homecloud.files WHERE id=$1`, fileID).Scan(&metadata)
	if err != nil {
		return "", err
	}
	return metadata, nil
}

// File operations (star, move, copy, rename)
func (r *dbRepository) StarFile(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET starred=true, updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *dbRepository) UnstarFile(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET starred=false, updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *dbRepository) MoveFile(ctx context.Context, fileID, newParentID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET parent_id=$1, updated_at=NOW() WHERE id=$2`, newParentID, fileID)
	return err
}

func (r *dbRepository) CopyFile(ctx context.Context, fileID, newParentID, newName string) (*models.File, error) {
	// Get original file
	originalFile, err := r.GetFileByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Create new file
	newFile := *originalFile
	newFile.ID = ""
	newFile.ParentID = &newParentID
	newFile.Name = newName
	newFile.CreatedAt = time.Now()
	newFile.UpdatedAt = time.Now()
	newFile.Starred = false
	newFile.IsTrashed = false
	newFile.TrashedAt = nil
	newFile.LastViewedAt = nil
	newFile.ViewedByMe = false
	newFile.Version = 1

	newID, err := r.CreateFile(ctx, &newFile)
	if err != nil {
		return nil, err
	}

	return r.GetFileByID(ctx, newID)
}

func (r *dbRepository) RenameFile(ctx context.Context, fileID, newName string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE homecloud.files SET name=$1, updated_at=NOW() WHERE id=$2`, newName, fileID)
	return err
}

// File integrity operations
func (r *dbRepository) VerifyFileIntegrity(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM homecloud.files WHERE id=$1)`, id).Scan(&exists)
	return exists, err
}

func (r *dbRepository) CalculateFileChecksums(ctx context.Context, id string) (map[string]string, error) {
	file, err := r.GetFileByID(ctx, id)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string)
	if file.MD5Checksum != nil {
		checksums["md5"] = *file.MD5Checksum
	}
	if file.SHA256Checksum != nil {
		checksums["sha256"] = *file.SHA256Checksum
	}

	return checksums, nil
}
