// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

const (
	AgentAuthTypeAPIKey            = "api_key"
	AgentAuthTypeBasicAuth         = "basic_auth"
	AgentAuthTypeClientCredentials = "client_credentials"
)

// AgentAuth is a single row in the agent_auths table.
type AgentAuth struct {
	Id        Id
	AgentId   Id
	Name      string
	Type      string
	Config    string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AgentAuthRepository is the interface for agent auth CRUD.
type AgentAuthRepository interface {
	Create(auth *AgentAuth) error
	GetById(id Id) (*AgentAuth, error)
	GetByAgentAndName(agentId Id, name string) (*AgentAuth, error)
	Update(auth *AgentAuth) error
	Delete(id Id) error
	ListByAgentId(agentId Id) ([]*AgentAuth, error)
	ListActiveByAgentId(agentId Id) ([]*AgentAuth, error)
}

type AgentAuthRepositoryPostgres struct {
	db *sql.DB
}

// NewAgentAuthRepository returns the repository for the agent_auths table.
func NewAgentAuthRepository(db *sql.DB) AgentAuthRepository {
	return &AgentAuthRepositoryPostgres{db: db}
}

// Create inserts an agent auth row.
func (r *AgentAuthRepositoryPostgres) Create(auth *AgentAuth) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	auth.Id = id

	return r.db.QueryRow(
		`INSERT INTO agent_auths (id, agent_id, name, type, config, is_active)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6)
		RETURNING created_at, updated_at`,
		auth.Id.String(),
		auth.AgentId.String(),
		auth.Name,
		auth.Type,
		auth.Config,
		auth.IsActive,
	).Scan(&auth.CreatedAt, &auth.UpdatedAt)
}

// GetById returns an agent auth by id.
func (r *AgentAuthRepositoryPostgres) GetById(id Id) (*AgentAuth, error) {
	auth := &AgentAuth{}
	err := r.db.QueryRow(
		`SELECT id, agent_id, name, type, config, is_active, created_at, updated_at
		FROM agent_auths
		WHERE id = $1`,
		id.String(),
	).Scan(
		&auth.Id,
		&auth.AgentId,
		&auth.Name,
		&auth.Type,
		&auth.Config,
		&auth.IsActive,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return auth, err
}

// GetByAgentAndName returns an agent auth by agent and name.
func (r *AgentAuthRepositoryPostgres) GetByAgentAndName(agentId Id, name string) (*AgentAuth, error) {
	auth := &AgentAuth{}
	err := r.db.QueryRow(
		`SELECT id, agent_id, name, type, config, is_active, created_at, updated_at
		FROM agent_auths
		WHERE agent_id = $1 AND name = $2`,
		agentId.String(),
		name,
	).Scan(
		&auth.Id,
		&auth.AgentId,
		&auth.Name,
		&auth.Type,
		&auth.Config,
		&auth.IsActive,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return auth, err
}

// Update updates an agent auth row.
func (r *AgentAuthRepositoryPostgres) Update(auth *AgentAuth) error {
	_, err := r.db.Exec(
		`UPDATE agent_auths
		SET name = $1, type = $2, config = $3::jsonb, is_active = $4, updated_at = $5
		WHERE id = $6`,
		auth.Name,
		auth.Type,
		auth.Config,
		auth.IsActive,
		time.Now().UTC(),
		auth.Id.String(),
	)
	return err
}

// Delete removes an agent auth row.
func (r *AgentAuthRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(`DELETE FROM agent_auths WHERE id = $1`, id.String())
	return err
}

// ListByAgentId lists the auth rows of an agent.
func (r *AgentAuthRepositoryPostgres) ListByAgentId(agentId Id) ([]*AgentAuth, error) {
	rows, err := r.db.Query(
		`SELECT id, agent_id, name, type, config, is_active, created_at, updated_at
		FROM agent_auths
		WHERE agent_id = $1
		ORDER BY name`,
		agentId.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*AgentAuth
	for rows.Next() {
		auth := &AgentAuth{}
		err := rows.Scan(
			&auth.Id,
			&auth.AgentId,
			&auth.Name,
			&auth.Type,
			&auth.Config,
			&auth.IsActive,
			&auth.CreatedAt,
			&auth.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, auth)
	}
	return list, rows.Err()
}

// ListActiveByAgentId lists the active auth rows of an agent.
func (r *AgentAuthRepositoryPostgres) ListActiveByAgentId(agentId Id) ([]*AgentAuth, error) {
	rows, err := r.db.Query(
		`SELECT id, agent_id, name, type, config, is_active, created_at, updated_at
		FROM agent_auths
		WHERE agent_id = $1 AND is_active = TRUE
		ORDER BY name`,
		agentId.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*AgentAuth
	for rows.Next() {
		auth := &AgentAuth{}
		err := rows.Scan(
			&auth.Id,
			&auth.AgentId,
			&auth.Name,
			&auth.Type,
			&auth.Config,
			&auth.IsActive,
			&auth.CreatedAt,
			&auth.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, auth)
	}
	return list, rows.Err()
}
