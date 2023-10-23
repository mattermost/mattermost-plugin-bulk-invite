package kvstore

import (
	"testing"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/require"
)

func TestCacheKeyStore(t *testing.T) {
	// Create a mock store for testing
	mockStore := NewMockStore()

	// Create a cache key store with a TTL of 1 second
	cacheStore := NewCacheKeyStore(mockStore, time.Second)

	// Test Store and Load methods
	t.Run("Store and Load", func(t *testing.T) {
		key := t.Name()
		value := []byte("value")

		// Store the value in the cache store
		err := cacheStore.Store(key, value)
		require.NoError(t, err)

		// Check that the value is in the cache store
		cachedValue, err := cacheStore.Load(key)
		require.NoError(t, err)
		require.Equal(t, value, cachedValue)

		// Check that the value is present in the cache of the cache store
		cachedValue, exists := cacheStore.(*cacheKeyStore).loadCache(key)
		require.True(t, exists)
		require.Equal(t, value, cachedValue)
	})

	// Test StoreTTL method
	t.Run("StoreTTL", func(t *testing.T) {
		t.Parallel()
		key := t.Name()
		value := []byte("value")
		ttlSeconds := int64(1)

		// Store the value in the cache store with a TTL
		err := cacheStore.StoreTTL(key, value, ttlSeconds)
		require.NoError(t, err)

		// Check that the value is in the cache store
		cachedValue, err := cacheStore.Load(key)
		require.NoError(t, err)
		require.Equal(t, value, cachedValue)

		// Check that the value is present in the cache of the cache store
		cachedValue, exists := cacheStore.(*cacheKeyStore).loadCache(key)
		require.True(t, exists)
		require.Equal(t, value, cachedValue)

		// Wait for the TTL to expire
		time.Sleep(time.Duration(ttlSeconds+1) * time.Second)

		// Avoiding checking the store because the mock store does not support TTL

		// Check that the value is not in the cache of the cache store
		cachedValue, exists = cacheStore.(*cacheKeyStore).loadCache(key)
		require.False(t, exists)
		require.Nil(t, cachedValue)
	})

	// Test StoreWithOptions method
	t.Run("StoreWithOptions", func(t *testing.T) {
		key := t.Name()
		value := []byte("value")
		opts := model.PluginKVSetOptions{
			Atomic: true,
		}

		// Store the value in the cache store with options
		ok, err := cacheStore.StoreWithOptions(key, value, opts)
		require.NoError(t, err)
		require.True(t, ok)

		// Check that the value is in the cache store
		cachedValue, err := cacheStore.Load(key)
		require.NoError(t, err)
		require.Equal(t, value, cachedValue)
	})

	// Test Delete method
	t.Run("Delete", func(t *testing.T) {
		key := t.Name()
		value := []byte("value")

		// Store the value in the cache store
		err := cacheStore.Store(key, value)
		require.NoError(t, err)

		// Check that the value is in the cache store
		cachedValue, exists := cacheStore.(*cacheKeyStore).loadCache(key)
		require.True(t, exists)
		require.Equal(t, value, cachedValue)

		// Delete the value from the cache store
		err = cacheStore.Delete(key)
		require.NoError(t, err)

		time.Sleep(time.Second * 2)

		// Avoiding checking the store because the mock store does not support TTL

		// Check that the value is not in the cache of the cache store
		cachedValue, exists = cacheStore.(*cacheKeyStore).loadCache(key)
		require.Nil(t, cachedValue)
		require.False(t, exists)
	})

	// Stop the cache store
	cacheStore.(*cacheKeyStore).cache.Stop()
}
