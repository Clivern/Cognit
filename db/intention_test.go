// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationIntentionRepository(t *testing.T) {
	database := openTestDB(t)
	workspace := createTestWorkspace(t, database)
	repo := NewIntentionRepository(database)

	t.Run("create get update delete", func(t *testing.T) {
		skill := "invoice.extract"
		intention := &Intention{
			WorkspaceId:      workspace.Id,
			SourceAgent:      "support-assistant",
			DestinationAgent: "invoice-extractor",
			Skill:            &skill,
			Action:           IntentionActionAllow,
		}
		require.NoError(t, repo.Create(intention))

		got, err := repo.GetById(intention.Id)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, IntentionActionAllow, got.Action)
		require.NotNil(t, got.Skill)
		assert.Equal(t, skill, *got.Skill)

		got.Action = IntentionActionDeny
		require.NoError(t, repo.Update(got))

		got, err = repo.GetById(intention.Id)
		require.NoError(t, err)
		assert.Equal(t, IntentionActionDeny, got.Action)

		list, err := repo.ListByWorkspaceId(workspace.Id, 50, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)

		require.NoError(t, repo.Delete(intention.Id))
		got, err = repo.GetById(intention.Id)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("match prefers named skill", func(t *testing.T) {
		require.NoError(t, repo.Create(&Intention{
			WorkspaceId:      workspace.Id,
			SourceAgent:      "support-assistant",
			DestinationAgent: "user-profile",
			Action:           IntentionActionDeny,
		}))
		skill := "user.profile.get"
		require.NoError(t, repo.Create(&Intention{
			WorkspaceId:      workspace.Id,
			SourceAgent:      "support-assistant",
			DestinationAgent: "user-profile",
			Skill:            &skill,
			Action:           IntentionActionAllow,
		}))

		matched, err := repo.Match(workspace.Id, "support-assistant", "user-profile", "user.profile.get")
		require.NoError(t, err)
		require.NotNil(t, matched)
		assert.Equal(t, IntentionActionAllow, matched.Action)
		require.NotNil(t, matched.Skill)
		assert.Equal(t, skill, *matched.Skill)

		matched, err = repo.Match(workspace.Id, "support-assistant", "user-profile", "other.skill")
		require.NoError(t, err)
		require.NotNil(t, matched)
		assert.Equal(t, IntentionActionDeny, matched.Action)
		assert.Nil(t, matched.Skill)

		matched, err = repo.Match(workspace.Id, "missing", "user-profile", "user.profile.get")
		require.NoError(t, err)
		assert.Nil(t, matched)
	})
}
