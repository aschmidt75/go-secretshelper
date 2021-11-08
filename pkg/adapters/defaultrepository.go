package adapters

import (
	"go-secretshelper/pkg/core"
	"sync"
)

// DefaultRepository stores all item in a map
type DefaultRepository struct {
	items map[string]interface{}
	m *sync.Mutex
}

// NewDefaultRepository creates a new DefaultRepository
func NewDefaultRepository() *DefaultRepository {
	return &DefaultRepository{
		items: make(map[string]interface{}),
		m: &sync.Mutex{},
	}
}

// Put places varName with content in repository
func (r *DefaultRepository) Put(varName string, content interface{}) {
	r.m.Lock()
	defer r.m.Unlock()

	r.items[varName] = content
}

// Get returns varName or an error
func (r *DefaultRepository) Get(varName string) (interface{}, error) {
	r.m.Lock()
	defer r.m.Unlock()

	res, ex := r.items[varName]
	if !ex {
		return nil, core.RepositoryError{Reason: "No such variable", Info: varName}
	}
	return res, nil
}
