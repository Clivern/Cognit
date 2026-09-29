// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationIntegrationsRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	repo := NewIntegrationRepository(database)
	item := &Integration{
		WorkspaceId: workspace.Id,
		Type:        "webhook",
		Name:        "alerts",
		Config:      stringPtr(`{"url":"https://example.test"}`),
	}

	t.Run("create and get", func(t *testing.T) {
		require.NoError(t, repo.Create(item))
		got, err := repo.GetById(item.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "webhook", got.Type)
	})

	t.Run("update", func(t *testing.T) {
		item.Name = "alerts-prod"
		require.NoError(t, repo.Update(item))
		got, err := repo.GetById(item.Id)
		require.NoError(t, err)
		assert.Equal(t, "alerts-prod", got.Name)
	})

	t.Run("list and count", func(t *testing.T) {
		list, err := repo.ListByWorkspaceId(workspace.Id, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)

		count, err := repo.CountByWorkspaceId(workspace.Id)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("delete", func(t *testing.T) {
		require.NoError(t, repo.Delete(item.Id))
		got, err := repo.GetById(item.Id)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}
