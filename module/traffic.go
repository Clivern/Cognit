// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"

	"github.com/clivern/cognit/db"
)

var (
	ErrTrafficNotFound = errors.New("traffic not found")
)

// Traffic is the module for recorded gateway calls.
type Traffic struct {
	WorkspaceRepository db.WorkspaceRepository
}

// NewTraffic creates a traffic module with the given repositories.
func NewTraffic(workspaces db.WorkspaceRepository) *Traffic {
	return &Traffic{WorkspaceRepository: workspaces}
}

// TrafficCallResponse is one gateway call shaped for API responses.
type TrafficCallResponse struct {
	Id          db.Id  `json:"id"`
	WorkspaceId db.Id  `json:"workspaceId"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Skill       string `json:"skill"`
	Verb        string `json:"verb"`
	Decision    string `json:"decision"`
	Result      string `json:"result"`
	Instance    string `json:"instance,omitempty"`
	LatencyMs   int64  `json:"latencyMs,omitempty"`
	TaskId      string `json:"taskId,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

// ListTrafficResponse is returned when listing gateway calls.
type ListTrafficResponse struct {
	Calls []*TrafficCallResponse
	Total int64
}

// ListTraffic returns recorded gateway calls for a workspace.
// Call rows are not stored yet, so the page is empty.
func (t *Traffic) ListTraffic(workspaceId db.Id, _, _ int) (*ListTrafficResponse, error) {
	workspace, err := t.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	return &ListTrafficResponse{Calls: []*TrafficCallResponse{}, Total: 0}, nil
}

// GetTraffic returns one gateway call by id.
func (t *Traffic) GetTraffic(workspaceId, _ db.Id) (*TrafficCallResponse, error) {
	workspace, err := t.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	return nil, ErrTrafficNotFound
}
