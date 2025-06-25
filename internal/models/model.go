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

// File представляет файл в системе
type File struct {
	ID             string
	OwnerID        string
	ParentID       *string
	Name           string
	FileExtension  *string
	MimeType       string
	StoragePath    string
	Size           int64
	MD5Checksum    *string
	SHA256Checksum *string
	IsFolder       bool
	IsTrashed      bool
	TrashedAt      *time.Time
	Starred        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastViewedAt   *time.Time
	ViewedByMe     bool
	Version        int64
	RevisionID     *string
	IndexableText  *string
	ThumbnailLink  *string
	WebViewLink    *string
	WebContentLink *string
	IconLink       *string
}

// FileRevision представляет ревизию файла
type FileRevision struct {
	ID          string
	FileID      string
	RevisionID  int64
	MD5Checksum *string
	Size        int64
	CreatedAt   time.Time
	StoragePath string
	MimeType    *string
	UserID      *string
}

// FilePermission представляет права доступа к файлу
type FilePermission struct {
	ID          string
	FileID      string
	GranteeID   *string
	GranteeType string
	Role        string
	AllowShare  bool
	CreatedAt   time.Time
}

// FileMetadata представляет метаданные файла
type FileMetadata struct {
	FileID   string
	Metadata string
}
