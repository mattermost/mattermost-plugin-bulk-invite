// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"time"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

var (
	lockedValue = []byte("")
	lockTTL     = 30 * time.Minute
	ErrIsLocked = errors.New("item is locked")
)

func getLockKey(key string) string {
	return "lock_" + key
}

type LockStore interface {
	Lock(key string) error
	Unlock(key string) error
	IsLocked(key string) bool
}

type lockStore struct {
	store KVStore
	ttl   time.Duration
}

func (s *lockStore) Lock(key string) error {
	if s.IsLocked(key) {
		return ErrIsLocked
	}
	return s.store.StoreTTL(getLockKey(key), lockedValue, int64(s.ttl))
}

func (s *lockStore) Unlock(key string) error {
	return s.store.Delete(getLockKey(key))
}

func (s *lockStore) IsLocked(key string) bool {
	return s.store.Exists(getLockKey(key))
}

func NewLockStore(api plugin.API) LockStore {
	return &lockStore{
		store: NewCacheKeyStore(NewPluginStore(api), 30*time.Second),
		ttl:   lockTTL,
	}
}
