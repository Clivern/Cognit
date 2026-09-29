// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationUserRepository(t *testing.T) {
	database := openTestDB(t)
	repo := NewUserRepository(database)
	user := createTestUser(t, database)

	t.Run("get by id and email", func(t *testing.T) {
		got, err := repo.GetById(user.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, user.Email, got.Email)

		byEmail, err := repo.GetByEmail(user.Email)
		require.NoError(t, err)
		require.NotNil(t, byEmail)
		assert.Equal(t, user.Id, byEmail.Id)
	})

	t.Run("update and get by provider", func(t *testing.T) {
		providerId := "gh-" + user.Id.String()[:8]
		user.Provider = UserProviderGithub
		user.ProviderUserId = &providerId
		user.Name = "updated"
		require.NoError(t, repo.Update(user))

		byProvider, err := repo.GetByProvider(UserProviderGithub, providerId)
		require.NoError(t, err)
		require.NotNil(t, byProvider)
		assert.Equal(t, "updated", byProvider.Name)
	})

	t.Run("update last login", func(t *testing.T) {
		require.NoError(t, repo.UpdateLastLogin(user.Id))
		afterLogin, err := repo.GetById(user.Id)
		require.NoError(t, err)
		assert.True(t, afterLogin.LastLoginAt.After(time.Now().UTC().Add(-time.Minute)))
	})

	t.Run("list and count", func(t *testing.T) {
		list, err := repo.List(50, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)

		count, err := repo.Count()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))
	})

	t.Run("missing email", func(t *testing.T) {
		missing, err := repo.GetByEmail("missing@t.test")
		require.NoError(t, err)
		assert.Nil(t, missing)
	})
}

func TestIntegrationUserGetByAPIKey(t *testing.T) {
	database := openTestDB(t)
	user := createTestUser(t, database)
	key := &APIKey{UserId: user.Id, Name: "cli", Key: "tok-" + user.Id.String()[:8]}
	require.NoError(t, NewAPIKeyRepository(database).Create(key))
	repo := NewUserRepository(database)

	t.Run("matching key", func(t *testing.T) {
		got, err := repo.GetByAPIKey(key.Key)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, user.Id, got.Id)
	})

	t.Run("empty key", func(t *testing.T) {
		got, err := repo.GetByAPIKey("")
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}
