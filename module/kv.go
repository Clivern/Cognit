// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/clivern/cognit/db"

	"github.com/samber/lo"
)

var (
	ErrKVNotFound     = errors.New("kv not found")
	ErrInvalidKVKey   = errors.New("invalid kv key")
	ErrFailedListKV   = errors.New("failed list kv")
	ErrFailedGetKV    = errors.New("failed get kv")
	ErrFailedPutKV    = errors.New("failed put kv")
	ErrFailedDeleteKV = errors.New("failed delete kv")
)

var kvKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_./-]{1,200}$`)

// KV is the module for workspace key/value config.
type KV struct {
	KVRepository        db.WorkspaceKVRepository
	WorkspaceRepository db.WorkspaceRepository
}

// NewKV creates a workspace KV module with the given repositories.
func NewKV(items db.WorkspaceKVRepository, workspaces db.WorkspaceRepository) *KV {
	return &KV{
		KVRepository:        items,
		WorkspaceRepository: workspaces,
	}
}

// PutKVRequest is the body for writing a workspace key.
type PutKVRequest struct {
	Value     string `json:"value" validate:"required" label:"Value"`
	ExpiresAt string `json:"expiresAt" validate:"omitempty,max=64" label:"Expires at"`
}

// KVResponse is a workspace key shaped for API responses.
type KVResponse struct {
	Id          db.Id   `json:"id"`
	WorkspaceId db.Id   `json:"workspaceId"`
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	ExpiresAt   *string `json:"expiresAt,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// ListKVResponse is returned when listing workspace keys.
type ListKVResponse struct {
	Items []*KVResponse
	Total int64
}

// ListKV returns non-expired keys that start with prefix.
func (k *KV) ListKV(workspaceId db.Id, prefix string) (*ListKVResponse, error) {
	prefixKey := strings.TrimSuffix(prefix, "/")
	valid := prefix == "" || kvKeyPattern.MatchString(prefixKey)
	if valid && prefix != "" {
		for _, part := range strings.Split(prefixKey, "/") {
			if part == "" || part == "." || part == ".." {
				valid = false
				break
			}
		}
	}
	if !valid {
		return nil, ErrInvalidKVKey
	}

	workspace, err := k.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	items, err := k.KVRepository.ListByPrefix(workspaceId, strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(prefix))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListKV, err)
	}

	list := make([]*KVResponse, 0, len(items))
	for _, item := range items {
		var expiresAt *string
		if item.ExpiresAt != nil {
			expiresAt = new(item.ExpiresAt.UTC().Format(time.RFC3339))
		}
		list = append(list, &KVResponse{
			Id:          item.Id,
			WorkspaceId: item.WorkspaceId,
			Key:         item.Key,
			Value:       item.Value,
			ExpiresAt:   expiresAt,
			CreatedAt:   item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	return &ListKVResponse{Items: list, Total: int64(len(list))}, nil
}

// GetKV returns one non-expired key.
func (k *KV) GetKV(workspaceId db.Id, key string) (*KVResponse, error) {
	valid := kvKeyPattern.MatchString(key)
	if valid {
		for _, part := range strings.Split(key, "/") {
			if part == "" || part == "." || part == ".." {
				valid = false
				break
			}
		}
	}
	if !valid {
		return nil, ErrInvalidKVKey
	}

	workspace, err := k.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	item, err := k.KVRepository.Get(workspaceId, key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetKV, err)
	}
	if item == nil {
		return nil, ErrKVNotFound
	}

	var expiresAt *string
	if item.ExpiresAt != nil {
		expiresAt = new(item.ExpiresAt.UTC().Format(time.RFC3339))
	}
	return &KVResponse{
		Id:          item.Id,
		WorkspaceId: item.WorkspaceId,
		Key:         item.Key,
		Value:       item.Value,
		ExpiresAt:   expiresAt,
		CreatedAt:   item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// PutKV writes a workspace key.
func (k *KV) PutKV(workspaceId db.Id, key string, req *PutKVRequest) (*KVResponse, error) {
	valid := kvKeyPattern.MatchString(key)
	if valid {
		for _, part := range strings.Split(key, "/") {
			if part == "" || part == "." || part == ".." {
				valid = false
				break
			}
		}
	}
	if !valid {
		return nil, ErrInvalidKVKey
	}

	workspace, err := k.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	var expiresAt *time.Time
	if lo.IsNotEmpty(req.ExpiresAt) {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return nil, ErrInvalidExpiresAt
		}
		expiresAt = new(t)
	}

	item := &db.WorkspaceKV{
		WorkspaceId: workspaceId,
		Key:         key,
		Value:       req.Value,
		ExpiresAt:   expiresAt,
	}
	err = k.KVRepository.Upsert(item)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedPutKV, err)
	}

	var expiresAtText *string
	if item.ExpiresAt != nil {
		expiresAtText = new(item.ExpiresAt.UTC().Format(time.RFC3339))
	}
	return &KVResponse{
		Id:          item.Id,
		WorkspaceId: item.WorkspaceId,
		Key:         item.Key,
		Value:       item.Value,
		ExpiresAt:   expiresAtText,
		CreatedAt:   item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeleteKV removes a workspace key.
func (k *KV) DeleteKV(workspaceId db.Id, key string) error {
	valid := kvKeyPattern.MatchString(key)
	if valid {
		for _, part := range strings.Split(key, "/") {
			if part == "" || part == "." || part == ".." {
				valid = false
				break
			}
		}
	}
	if !valid {
		return ErrInvalidKVKey
	}

	workspace, err := k.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	item, err := k.KVRepository.Get(workspaceId, key)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteKV, err)
	}
	if item == nil {
		return ErrKVNotFound
	}

	err = k.KVRepository.Delete(workspaceId, key)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteKV, err)
	}
	return nil
}
