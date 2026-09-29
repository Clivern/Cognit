// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAuditEventRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	user := createTestUser(t, database)
	repo := NewAuditEventRepository(database)
	event := &AuditEvent{
		WorkspaceId:  workspace.Id,
		UserId:       &user.Id,
		Action:       "agent.register",
		ResourceType: stringPtr("agent"),
		ResourceId:   &user.Id,
		IPAddress:    stringPtr("10.0.0.1"),
		Meta:         stringPtr(`{"agent":"invoice-extractor"}`),
	}

	t.Run("create and get", func(t *testing.T) {
		require.NoError(t, repo.Create(event))

		got, err := repo.GetById(event.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "agent.register", got.Action)
		require.NotNil(t, got.UserId)
		assert.Equal(t, user.Id, *got.UserId)
	})

	t.Run("list and count", func(t *testing.T) {
		list, err := repo.ListByWorkspaceId(workspace.Id, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)

		count, err := repo.CountByWorkspaceId(workspace.Id)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))
	})

	t.Run("missing row", func(t *testing.T) {
		missing, err := repo.GetById(user.Id)
		require.NoError(t, err)
		assert.Nil(t, missing)
	})
}
