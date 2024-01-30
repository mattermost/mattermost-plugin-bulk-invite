package kvstore

import "github.com/mattermost/mattermost/server/public/model"

type mockStore struct {
	data map[string][]byte
}

func NewMockStore() KVStore {
	return &mockStore{
		data: make(map[string][]byte),
	}
}

func (s *mockStore) Load(key string) ([]byte, error) {
	return s.data[key], nil
}

func (s *mockStore) Store(key string, value []byte) error {
	s.data[key] = value
	return nil
}

func (s *mockStore) StoreTTL(key string, value []byte, _ int64) error {
	s.data[key] = value
	return nil
}

func (s *mockStore) StoreWithOptions(key string, value []byte, _ model.PluginKVSetOptions) (bool, error) {
	s.data[key] = value
	return true, nil
}

func (s *mockStore) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *mockStore) LoadAll() (map[string][]byte, error) {
	return s.data, nil
}

func (s *mockStore) DeleteAll() error {
	s.data = make(map[string][]byte)
	return nil
}

func (s *mockStore) Exists(key string) bool {
	_, ok := s.data[key]
	return ok
}

func (s *mockStore) Close() error {
	return nil
}
