package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashTable_Get(t *testing.T) {
	t.Parallel()
	t.Run("return existing value for key", func(t *testing.T) {
		table := NewHashTable()
		table.Set("key1", "value1")
		value, found := table.Get("key1")
		require.True(t, found)
		require.Equal(t, "value1", value)
	})

	t.Run("return empty string and false if key does not exist", func(t *testing.T) {
		table := NewHashTable()
		value, found := table.Get("key1")
		require.False(t, found)
		require.Empty(t, value)
	})
}

func TestHashTable_Del(t *testing.T) {
	t.Parallel()

	t.Run("Deletion of key-value pair", func(t *testing.T) {
		table := NewHashTable()
		table.Set("key1", "value1")
		table.Del("key1")
		value, found := table.Get("key1")
		require.False(t, found)
		require.Empty(t, value)
	})

	t.Run("Deletion of not existing key (no error)", func(t *testing.T) {
		table := NewHashTable()
		table.Del("key1")
		value, found := table.Get("key1")
		require.False(t, found)
		require.Empty(t, value)
	})
}

func TestMemoryTable_Set(t *testing.T) {
	t.Parallel()
	t.Run("correct setting key and value", func(t *testing.T) {
		table := NewHashTable()
		table.Set("key1", "value1")
		value, found := table.Get("key1")
		require.True(t, found)
		require.Equal(t, "value1", value)
	})

	t.Run("correct overwriting existing key-value", func(t *testing.T) {
		table := NewHashTable()
		table.Set("key1", "value1")
		table.Set("key1", "value2")
		value, found := table.Get("key1")
		require.True(t, found)
		require.Equal(t, "value2", value)
	})
}
