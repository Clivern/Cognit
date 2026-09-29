// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationHealthCheckRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	agent := createTestAgent(t, database, workspace.Id, "invoice-extractor")
	instance := createTestInstance(t, database, agent.Id, "invoice-extractor-3")
	repo := NewHealthCheckRepository(database)

	t.Run("create starts critical", func(t *testing.T) {
		check := &HealthCheck{
			AgentInstanceId: instance.Id,
			CheckId:         HealthCheckIDTTL,
			Name:            "lease",
			Type:            HealthCheckTypeTTL,
			Definition:      stringPtr(`{"ttl":"30s"}`),
		}
		require.NoError(t, repo.Create(check))
		assert.Equal(t, HealthCheckStatusCritical, check.Status)

		got, err := repo.GetByInstanceAndCheckId(instance.Id, HealthCheckIDTTL)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, check.Id, got.Id)
		assert.Contains(t, *got.Definition, "30s")
	})

	t.Run("pass warn fail", func(t *testing.T) {
		expires := time.Now().UTC().Add(30 * time.Second)
		require.NoError(t, repo.Pass(
			mustCheck(t, repo, instance.Id, HealthCheckIDTTL).Id,
			"renewed",
			&expires,
		))

		got, err := repo.GetByInstanceAndCheckId(instance.Id, HealthCheckIDTTL)
		require.NoError(t, err)
		assert.Equal(t, HealthCheckStatusPassing, got.Status)
		require.NotNil(t, got.TTLExpiresAt)
		assert.WithinDuration(t, expires, *got.TTLExpiresAt, time.Second)
		require.NotNil(t, got.Output)
		assert.Equal(t, "renewed", *got.Output)

		httpCheck := &HealthCheck{
			AgentInstanceId: instance.Id,
			CheckId:         "http",
			Name:            "/health",
			Type:            HealthCheckTypeHTTP,
			Definition:      stringPtr(`{"http":"http://10.0.12.41:8080/health","interval":"10s"}`),
			Status:          HealthCheckStatusPassing,
		}
		require.NoError(t, repo.Create(httpCheck))

		require.NoError(t, repo.Warn(httpCheck.Id, "429"))
		got, err = repo.GetById(httpCheck.Id)
		require.NoError(t, err)
		assert.Equal(t, HealthCheckStatusWarning, got.Status)

		require.NoError(t, repo.Fail(httpCheck.Id, "connection refused"))
		got, err = repo.GetById(httpCheck.Id)
		require.NoError(t, err)
		assert.Equal(t, HealthCheckStatusCritical, got.Status)

		list, err := repo.ListByAgentInstanceId(instance.Id)
		require.NoError(t, err)
		assert.Len(t, list, 2)

		require.NoError(t, repo.Delete(httpCheck.Id))
		got, err = repo.GetById(httpCheck.Id)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

func mustCheck(t *testing.T, repo HealthCheckRepository, instanceId Id, checkId string) *HealthCheck {
	t.Helper()
	got, err := repo.GetByInstanceAndCheckId(instanceId, checkId)
	require.NoError(t, err)
	require.NotNil(t, got)
	return got
}
