// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package kvstore

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"
)

type KVStore interface {
	Load(key string) ([]byte, error)
	Store(key string, data []byte) error
	StoreTTL(key string, data []byte, ttlSeconds int64) error
	StoreWithOptions(key string, value []byte, opts model.PluginKVSetOptions) (bool, error)
	Delete(key string) error
	Exists(key string) bool
}

var ErrNotFound = errors.New("not found")
