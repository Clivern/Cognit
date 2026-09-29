// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAgentAuthRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	agent := createTestAgent(t, database, workspace.Id, "auth-agent")
	repo := NewAgentAuthRepository(database)

	t.Run("create and get", func(t *testing.T) {
		auth := &AgentAuth{
			AgentId:  agent.Id,
			Name:     "gateway-key",
			Type:     AgentAuthTypeAPIKey,
			Config:   `{"header":"X-API-Key","value":"secret"}`,
			IsActive: true,
		}
		require.NoError(t, repo.Create(auth))

		got, err := repo.GetById(auth.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, agent.Id, got.AgentId)
		assert.Equal(t, "gateway-key", got.Name)
		assert.Equal(t, AgentAuthTypeAPIKey, got.Type)
		assert.True(t, got.IsActive)
		assert.Contains(t, got.Config, "X-API-Key")

		byName, err := repo.GetByAgentAndName(agent.Id, "gateway-key")
		require.NoError(t, err)
		require.NotNil(t, byName)
		assert.Equal(t, auth.Id, byName.Id)
	})

	t.Run("missing row", func(t *testing.T) {
		id, err := NewId()
		require.NoError(t, err)

		got, err := repo.GetById(id)
		require.NoError(t, err)
		assert.Nil(t, got)

		got, err = repo.GetByAgentAndName(agent.Id, "missing")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("update list and delete", func(t *testing.T) {
		basic := &AgentAuth{
			AgentId:  agent.Id,
			Name:     "basic",
			Type:     AgentAuthTypeBasicAuth,
			Config:   `{"username":"agent","password":"pass"}`,
			IsActive: true,
		}
		require.NoError(t, repo.Create(basic))

		oauth := &AgentAuth{
			AgentId:  agent.Id,
			Name:     "oauth",
			Type:     AgentAuthTypeClientCredentials,
			Config:   `{"token_url":"https://idp.test/token","client_id":"id","client_secret":"secret"}`,
			IsActive: true,
		}
		require.NoError(t, repo.Create(oauth))

		basic.Config = `{"username":"agent","password":"rotated"}`
		basic.IsActive = false
		require.NoError(t, repo.Update(basic))

		got, err := repo.GetById(basic.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.False(t, got.IsActive)
		assert.Contains(t, got.Config, "rotated")

		list, err := repo.ListByAgentId(agent.Id)
		require.NoError(t, err)
		assert.Len(t, list, 3)

		active, err := repo.ListActiveByAgentId(agent.Id)
		require.NoError(t, err)
		assert.Len(t, active, 2)
		for _, item := range active {
			assert.True(t, item.IsActive)
		}

		require.NoError(t, repo.Delete(basic.Id))
		got, err = repo.GetById(basic.Id)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("unique name per agent", func(t *testing.T) {
		duplicate := &AgentAuth{
			AgentId:  agent.Id,
			Name:     "gateway-key",
			Type:     AgentAuthTypeAPIKey,
			Config:   `{"header":"X-API-Key","value":"other"}`,
			IsActive: true,
		}
		assert.Error(t, repo.Create(duplicate))
	})

	t.Run("cascade delete with agent", func(t *testing.T) {
		other := createTestAgent(t, database, workspace.Id, "cascade-auth-agent")
		auth := &AgentAuth{
			AgentId:  other.Id,
			Name:     "gateway-key",
			Type:     AgentAuthTypeAPIKey,
			Config:   `{"header":"X-API-Key","value":"secret"}`,
			IsActive: true,
		}
		require.NoError(t, repo.Create(auth))
		require.NoError(t, NewAgentRepository(database).Delete(other.Id))

		got, err := repo.GetById(auth.Id)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}
