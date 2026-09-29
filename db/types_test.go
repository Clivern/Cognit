// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitId(t *testing.T) {
	t.Run("NewId", func(t *testing.T) {
		id, err := NewId()
		require.NoError(t, err)
		assert.Len(t, id.String(), 36)
	})
}
