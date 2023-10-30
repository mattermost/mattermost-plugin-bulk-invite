package kvstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLockStore(t *testing.T) {
	store := NewMockStore()
	lockStore := &lockStore{
		store: store,
		ttl:   1 * time.Second,
	}

	// Test Lock
	err := lockStore.Lock("test")
	assert.NoError(t, err)

	// Test IsLocked
	assert.True(t, lockStore.IsLocked("test"))

	// Test Lock with existing lock
	err = lockStore.Lock("test")
	assert.Error(t, err)
	assert.ErrorIs(t, ErrIsLocked, err)

	// Test Unlock
	err = lockStore.Unlock("test")
	assert.NoError(t, err)

	// Test Unlock with unlocked key
	err = lockStore.Unlock("test")
	assert.NoError(t, err)

	// Test IsLocked with unlocked key
	assert.False(t, lockStore.IsLocked("test"))
}
