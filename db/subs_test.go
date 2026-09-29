// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSubscriptionRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	repo := NewSubscriptionRepository(database)
	customer := "cus_test"
	sub := &Subscription{
		WorkspaceId:        workspace.Id,
		ProviderCustomerId: &customer,
		AITokensBalance:    100,
	}

	t.Run("create and get", func(t *testing.T) {
		require.NoError(t, repo.Create(sub))

		got, err := repo.GetByWorkspaceId(workspace.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, int64(100), got.AITokensBalance)

		byId, err := repo.GetById(sub.Id)
		require.NoError(t, err)
		require.NotNil(t, byId)
	})

	t.Run("add and consume tokens", func(t *testing.T) {
		require.NoError(t, repo.AddTokens(workspace.Id, 50))
		got, err := repo.GetByWorkspaceId(workspace.Id)
		require.NoError(t, err)
		assert.Equal(t, int64(150), got.AITokensBalance)

		require.NoError(t, repo.ConsumeTokens(workspace.Id, 200))
		got, err = repo.GetByWorkspaceId(workspace.Id)
		require.NoError(t, err)
		assert.Equal(t, int64(0), got.AITokensBalance)
	})

	t.Run("update customer", func(t *testing.T) {
		customer2 := "cus_updated"
		sub.ProviderCustomerId = &customer2
		require.NoError(t, repo.Update(sub))
		got, err := repo.GetById(sub.Id)
		require.NoError(t, err)
		require.NotNil(t, got.ProviderCustomerId)
		assert.Equal(t, "cus_updated", *got.ProviderCustomerId)
	})
}
