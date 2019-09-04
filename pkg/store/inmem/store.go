package inmem

import (
	"sync"

	"github.com/pkg/errors"
)

// Store.
type Store struct {
	hashmap map[string][]byte
	mu      sync.RWMutex
}

// NewStore initializes a Store.
func NewStore() *Store {
	return &Store{hashmap: make(map[string][]byte), mu: sync.RWMutex{}}
}

func (m *Store) Insert(key string, value []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hashmap[key] = value
}

func (m *Store) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.hashmap, key)
}

func (m *Store) Read(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.hashmap[key]
	if !ok {
		return nil, errors.Errorf("unable to get secret key %s", key)
	}

	return v, nil
}
