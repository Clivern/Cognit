// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSessionRepository(t *testing.T) {
	database := openTestDB(t)
	user := createTestUser(t, database)
	repo := NewSessionRepository(database)
	token := "sess-" + uuid.NewString()
	session := &Session{
		Token:     token,
		UserId:    user.Id,
		IPAddress: stringPtr("127.0.0.1"),
		UserAgent: stringPtr("test"),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}

	t.Run("create and get", func(t *testing.T) {
		require.NoError(t, repo.Create(session))

		got, err := repo.GetByToken(token)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, user.Id, got.UserId)

		byId, err := repo.GetById(session.Id)
		require.NoError(t, err)
		require.NotNil(t, byId)

		list, err := repo.GetByUserId(user.Id)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)
	})
}
