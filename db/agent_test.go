// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAgentRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	repo := NewAgentRepository(database)

	t.Run("create and get", func(t *testing.T) {
		agent := createTestAgent(t, database, workspace.Id, "invoice-extractor")

		got, err := repo.GetById(agent.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, agent.Name, got.Name)
		assert.Equal(t, workspace.Id, got.WorkspaceId)
		assert.Equal(t, agent.CardChecksum, got.CardChecksum)
		assert.Equal(t, "1.0.0", got.Version)
		assert.Contains(t, got.Card, "invoice-extractor")

		byName, err := repo.GetByWorkspaceAndName(workspace.Id, "invoice-extractor")
		require.NoError(t, err)
		require.NotNil(t, byName)
		assert.Equal(t, agent.Id, byName.Id)
	})

	t.Run("missing row", func(t *testing.T) {
		id, err := NewId()
		require.NoError(t, err)

		got, err := repo.GetById(id)
		require.NoError(t, err)
		assert.Nil(t, got)

		got, err = repo.GetByWorkspaceAndName(workspace.Id, "missing")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("update list count and delete", func(t *testing.T) {
		agent := createTestAgent(t, database, workspace.Id, "support-assistant")
		createTestAgent(t, database, workspace.Id, "other-agent")

		agent.Version = "1.1.0"
		agent.Card = `{"name":"support-assistant","version":"1.1.0"}`
		agent.CardChecksum = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		require.NoError(t, repo.Update(agent))

		got, err := repo.GetById(agent.Id)
		require.NoError(t, err)
		assert.Equal(t, "1.1.0", got.Version)
		assert.Equal(t, agent.CardChecksum, got.CardChecksum)

		list, err := repo.ListByWorkspaceId(workspace.Id, 50, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 2)

		count, err := repo.CountByWorkspaceId(workspace.Id)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2))

		require.NoError(t, repo.Delete(agent.Id))
		got, err = repo.GetById(agent.Id)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("unique name per workspace", func(t *testing.T) {
		createTestAgent(t, database, workspace.Id, "unique-agent")
		duplicate := &Agent{
			WorkspaceId:  workspace.Id,
			Name:         "unique-agent",
			Card:         `{"name":"unique-agent"}`,
			CardChecksum: "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
			Version:      "1.0.0",
		}
		assert.Error(t, repo.Create(duplicate))
	})
}

func TestIntegrationAgentMetaRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	agent := createTestAgent(t, database, workspace.Id, "meta-agent")
	repo := NewAgentMetaRepository(database)

	t.Run("create and get", func(t *testing.T) {
		require.NoError(t, repo.Create(agent.Id, "region", "eu"))
		got, err := repo.Get(agent.Id, "region")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "eu", got.Value)
	})

	t.Run("upsert", func(t *testing.T) {
		require.NoError(t, repo.Upsert(agent.Id, "region", "us"))
		require.NoError(t, repo.Upsert(agent.Id, "tier", "gold"))
		got, err := repo.Get(agent.Id, "region")
		require.NoError(t, err)
		assert.Equal(t, "us", got.Value)
	})

	t.Run("list and delete", func(t *testing.T) {
		list, err := repo.ListByAgentId(agent.Id)
		require.NoError(t, err)
		assert.Len(t, list, 2)

		require.NoError(t, repo.Delete(agent.Id, "tier"))
		got, err := repo.Get(agent.Id, "tier")
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}
