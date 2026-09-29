// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// WorkspaceKV is a workspace-scoped key/value row with an optional expiry.
type WorkspaceKV struct {
	Id          Id
	WorkspaceId Id
	Key         string
	Value       string
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WorkspaceKVRepository is the interface for workspace KV persistence.
type WorkspaceKVRepository interface {
	Upsert(item *WorkspaceKV) error
	Get(workspaceId Id, key string) (*WorkspaceKV, error)
	Delete(workspaceId Id, key string) error
	ListByPrefix(workspaceId Id, prefix string) ([]*WorkspaceKV, error)
	DeleteExpired() (int64, error)
}

type WorkspaceKVRepositoryPostgres struct {
	db *sql.DB
}

const workspaceKVColumns = `id, workspace_id, key, value, expires_at, created_at, updated_at`

// NewWorkspaceKVRepository returns the repository for workspace_kv.
func NewWorkspaceKVRepository(db *sql.DB) WorkspaceKVRepository {
	return &WorkspaceKVRepositoryPostgres{db: db}
}

// Upsert inserts or replaces a workspace key/value row.
func (r *WorkspaceKVRepositoryPostgres) Upsert(item *WorkspaceKV) error {
	id, err := NewId()
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		`INSERT INTO workspace_kv (id, workspace_id, key, value, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (workspace_id, key) DO UPDATE SET
			value = EXCLUDED.value,
			expires_at = EXCLUDED.expires_at,
			updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
		RETURNING id, created_at, updated_at`,
		id.String(),
		item.WorkspaceId.String(),
		item.Key,
		item.Value,
		item.ExpiresAt,
	).Scan(
		&item.Id,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
}

// Get returns a non-expired value for a workspace key.
func (r *WorkspaceKVRepositoryPostgres) Get(workspaceId Id, key string) (*WorkspaceKV, error) {
	item := &WorkspaceKV{}
	err := r.db.QueryRow(
		`SELECT `+workspaceKVColumns+`
		FROM workspace_kv
		WHERE workspace_id = $1 AND key = $2
			AND (expires_at IS NULL OR expires_at > $3)`,
		workspaceId.String(),
		key,
		time.Now().UTC(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.Key,
		&item.Value,
		&item.ExpiresAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Delete removes a workspace key/value row.
func (r *WorkspaceKVRepositoryPostgres) Delete(workspaceId Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM workspace_kv WHERE workspace_id = $1 AND key = $2`,
		workspaceId.String(),
		key,
	)
	return err
}

// ListByPrefix lists non-expired keys in a workspace that start with prefix.
func (r *WorkspaceKVRepositoryPostgres) ListByPrefix(workspaceId Id, prefix string) ([]*WorkspaceKV, error) {
	rows, err := r.db.Query(
		`SELECT `+workspaceKVColumns+`
		FROM workspace_kv
		WHERE workspace_id = $1
			AND key LIKE $2
			AND (expires_at IS NULL OR expires_at > $3)
		ORDER BY key`,
		workspaceId.String(),
		prefix+"%",
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*WorkspaceKV
	for rows.Next() {
		item := &WorkspaceKV{}
		err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.Key,
			&item.Value,
			&item.ExpiresAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// DeleteExpired removes rows that have passed their expiry.
func (r *WorkspaceKVRepositoryPostgres) DeleteExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM workspace_kv WHERE expires_at IS NOT NULL AND expires_at <= $1`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
