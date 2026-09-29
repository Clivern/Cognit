// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package lb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitPicker(t *testing.T) {
	a := Instance{ID: "a", Datacenter: "eu-west-1", Outstanding: 2}
	b := Instance{ID: "b", Datacenter: "us-east-1", Outstanding: 0}
	c := Instance{ID: "c", Datacenter: "eu-west-1", Outstanding: 0}

	t.Run("empty pool", func(t *testing.T) {
		_, err := New().Pick(nil, Options{})
		assert.ErrorIs(t, err, ErrEmpty)
	})

	t.Run("single instance", func(t *testing.T) {
		got, err := New().Pick([]Instance{a}, Options{})
		require.NoError(t, err)
		assert.Equal(t, "a", got.ID)
	})

	t.Run("unknown strategy", func(t *testing.T) {
		_, err := New().Pick([]Instance{a}, Options{Strategy: "weighted"})
		assert.ErrorIs(t, err, ErrStrategy)
	})

	t.Run("random uses injected index", func(t *testing.T) {
		p := New()
		p.intn = func(n int) int {
			assert.Equal(t, 3, n)
			return 1
		}
		got, err := p.Pick([]Instance{a, b, c}, Options{Strategy: Random})
		require.NoError(t, err)
		assert.Equal(t, "b", got.ID)
	})

	t.Run("round robin cycles in id order", func(t *testing.T) {
		p := New()
		pool := []Instance{c, a, b}
		var ids []string
		for range 4 {
			got, err := p.Pick(pool, Options{Strategy: RoundRobin, Pool: "invoice-extractor"})
			require.NoError(t, err)
			ids = append(ids, got.ID)
		}
		assert.Equal(t, []string{"a", "b", "c", "a"}, ids)
	})

	t.Run("least outstanding picks the lightest", func(t *testing.T) {
		p := New()
		p.intn = func(int) int { return 0 }
		got, err := p.Pick([]Instance{a, b}, Options{Strategy: LeastOutstandingTasks})
		require.NoError(t, err)
		assert.Equal(t, "b", got.ID)
	})

	t.Run("least outstanding breaks ties at random", func(t *testing.T) {
		p := New()
		p.intn = func(n int) int {
			assert.Equal(t, 2, n)
			return 1
		}
		got, err := p.Pick([]Instance{a, b, c}, Options{Strategy: LeastOutstandingTasks})
		require.NoError(t, err)
		assert.Equal(t, "c", got.ID)
	})

	t.Run("prefer datacenter", func(t *testing.T) {
		p := New()
		p.intn = func(n int) int {
			assert.Equal(t, 2, n)
			return 1
		}
		got, err := p.Pick([]Instance{a, b, c}, Options{
			Strategy:         Random,
			PreferDatacenter: "eu-west-1",
		})
		require.NoError(t, err)
		assert.Equal(t, "c", got.ID)
	})

	t.Run("prefer datacenter falls back", func(t *testing.T) {
		got, err := New().Pick([]Instance{a, b}, Options{
			PreferDatacenter: "ap-south-1",
		})
		require.NoError(t, err)
		assert.Contains(t, []string{"a", "b"}, got.ID)
	})

	t.Run("sticky key is stable", func(t *testing.T) {
		p := New()
		pool := []Instance{a, b, c}
		first, err := p.Pick(pool, Options{StickyKey: "task-1"})
		require.NoError(t, err)
		second, err := p.Pick(pool, Options{Strategy: Random, StickyKey: "task-1"})
		require.NoError(t, err)
		assert.Equal(t, first.ID, second.ID)
	})

	t.Run("sticky follows remaining members", func(t *testing.T) {
		p := New()
		full, err := p.Pick([]Instance{a, b, c}, Options{StickyKey: "task-9"})
		require.NoError(t, err)

		remaining := make([]Instance, 0, 2)
		for _, inst := range []Instance{a, b, c} {
			if inst.ID != full.ID {
				remaining = append(remaining, inst)
			}
		}
		got, err := p.Pick(remaining, Options{StickyKey: "task-9"})
		require.NoError(t, err)
		assert.NotEqual(t, full.ID, got.ID)
		assert.Contains(t, []string{remaining[0].ID, remaining[1].ID}, got.ID)
	})
}

func TestUnitShuffle(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, New().Shuffle(nil))
	})

	t.Run("does not mutate input", func(t *testing.T) {
		in := []Instance{{ID: "a"}, {ID: "b"}, {ID: "c"}}
		p := New()
		p.intn = func(n int) int { return n - 1 }
		out := p.Shuffle(in)
		assert.Equal(t, "a", in[0].ID)
		assert.Len(t, out, 3)
	})
}
