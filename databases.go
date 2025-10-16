package sevalla

import (
	"context"
	"fmt"
)

// DatabasesService handles communication with the database-related
// methods of the Sevalla API.
type DatabasesService struct {
	client *Client
}

// CreateDatabaseRequest represents a request to create a new database
type CreateDatabaseRequest struct {
	Name       string `json:"name"`
	Type       Engine `json:"type"`
	Version    string `json:"version,omitempty"`
	Region     Region `json:"location,omitempty"`
	Size       string `json:"size,omitempty"`
	Storage    int    `json:"storage_gb,omitempty"`
	Backups    bool   `json:"backups_enabled,omitempty"`
	SSLEnabled bool   `json:"ssl_enabled,omitempty"`
}

// UpdateDatabaseRequest represents a request to update a database
type UpdateDatabaseRequest struct {
	Name       *string `json:"name,omitempty"`
	Size       *string `json:"size,omitempty"`
	Storage    *int    `json:"storage_gb,omitempty"`
	Backups    *bool   `json:"backups_enabled,omitempty"`
	SSLEnabled *bool   `json:"ssl_enabled,omitempty"`
}

// CreateBackupRequest represents a request to create a database backup
type CreateBackupRequest struct {
	Type        string `json:"type,omitempty"` // "manual" or "automatic"
	Description string `json:"description,omitempty"`
}

// RestoreBackupRequest represents a request to restore a database from backup
type RestoreBackupRequest struct {
	BackupID string `json:"backup_id"`
}

// List returns all databases
func (s *DatabasesService) List(ctx context.Context, opts *ListOptions) ([]*Database, *Response, error) {
	u := "databases"
	req, err := s.client.NewRequestWithQuery(ctx, "GET", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var databases []*Database
	resp, err := s.client.Do(req, &databases)
	if err != nil {
		return nil, resp, err
	}

	return databases, resp, nil
}

// Get returns a single database by ID
func (s *DatabasesService) Get(ctx context.Context, id string) (*Database, *Response, error) {
	u := fmt.Sprintf("databases/%s", id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	database := new(Database)
	resp, err := s.client.Do(req, &database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, nil
}

// Create creates a new database
func (s *DatabasesService) Create(ctx context.Context, createReq *CreateDatabaseRequest) (*Database, *Response, error) {
	u := "databases"
	req, err := s.client.NewRequest(ctx, "POST", u, createReq)
	if err != nil {
		return nil, nil, err
	}

	database := new(Database)
	resp, err := s.client.Do(req, &database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, nil
}

// Update updates an existing database
func (s *DatabasesService) Update(ctx context.Context, id string, updateReq *UpdateDatabaseRequest) (*Database, *Response, error) {
	u := fmt.Sprintf("databases/%s", id)
	req, err := s.client.NewRequest(ctx, "PATCH", u, updateReq)
	if err != nil {
		return nil, nil, err
	}

	database := new(Database)
	resp, err := s.client.Do(req, &database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, nil
}

// Delete deletes a database
func (s *DatabasesService) Delete(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("databases/%s", id)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetCredentials retrieves database connection credentials
func (s *DatabasesService) GetCredentials(ctx context.Context, id string) (*Database, *Response, error) {
	u := fmt.Sprintf("databases/%s/credentials", id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	database := new(Database)
	resp, err := s.client.Do(req, &database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, nil
}

// ResetPassword resets the database password
func (s *DatabasesService) ResetPassword(ctx context.Context, id string) (*Database, *Response, error) {
	u := fmt.Sprintf("databases/%s/reset-password", id)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	database := new(Database)
	resp, err := s.client.Do(req, &database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, nil
}

// ListBackups lists all backups for a database
func (s *DatabasesService) ListBackups(ctx context.Context, id string, opts *ListOptions) ([]*Backup, *Response, error) {
	u := fmt.Sprintf("databases/%s/backups", id)
	req, err := s.client.NewRequestWithQuery(ctx, "GET", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var backups []*Backup
	resp, err := s.client.Do(req, &backups)
	if err != nil {
		return nil, resp, err
	}

	return backups, resp, nil
}

// CreateBackup creates a new backup for a database
func (s *DatabasesService) CreateBackup(ctx context.Context, id string, backupReq *CreateBackupRequest) (*Backup, *Response, error) {
	u := fmt.Sprintf("databases/%s/backups", id)
	req, err := s.client.NewRequest(ctx, "POST", u, backupReq)
	if err != nil {
		return nil, nil, err
	}

	backup := new(Backup)
	resp, err := s.client.Do(req, &backup)
	if err != nil {
		return nil, resp, err
	}

	return backup, resp, nil
}

// GetBackup gets a specific backup
func (s *DatabasesService) GetBackup(ctx context.Context, dbID, backupID string) (*Backup, *Response, error) {
	u := fmt.Sprintf("databases/%s/backups/%s", dbID, backupID)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	backup := new(Backup)
	resp, err := s.client.Do(req, &backup)
	if err != nil {
		return nil, resp, err
	}

	return backup, resp, nil
}

// DeleteBackup deletes a backup
func (s *DatabasesService) DeleteBackup(ctx context.Context, dbID, backupID string) (*Response, error) {
	u := fmt.Sprintf("databases/%s/backups/%s", dbID, backupID)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RestoreFromBackup restores a database from a backup
func (s *DatabasesService) RestoreFromBackup(ctx context.Context, id string, restoreReq *RestoreBackupRequest) (*Response, error) {
	u := fmt.Sprintf("databases/%s/restore", id)
	req, err := s.client.NewRequest(ctx, "POST", u, restoreReq)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetUsage retrieves usage metrics for a database
func (s *DatabasesService) GetUsage(ctx context.Context, id string, period string) (*Usage, *Response, error) {
	u := fmt.Sprintf("databases/%s/usage", id)
	if period != "" {
		u = fmt.Sprintf("%s?period=%s", u, period)
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	usage := new(Usage)
	resp, err := s.client.Do(req, &usage)
	if err != nil {
		return nil, resp, err
	}

	return usage, resp, nil
}

// EnablePublicAccess enables public access to a database
func (s *DatabasesService) EnablePublicAccess(ctx context.Context, id string) (*Database, *Response, error) {
	u := fmt.Sprintf("databases/%s/public-access", id)
	req, err := s.client.NewRequest(ctx, "PUT", u, map[string]bool{"enabled": true})
	if err != nil {
		return nil, nil, err
	}

	database := new(Database)
	resp, err := s.client.Do(req, &database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, nil
}

// DisablePublicAccess disables public access to a database
func (s *DatabasesService) DisablePublicAccess(ctx context.Context, id string) (*Database, *Response, error) {
	u := fmt.Sprintf("databases/%s/public-access", id)
	req, err := s.client.NewRequest(ctx, "PUT", u, map[string]bool{"enabled": false})
	if err != nil {
		return nil, nil, err
	}

	database := new(Database)
	resp, err := s.client.Do(req, &database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, nil
}
