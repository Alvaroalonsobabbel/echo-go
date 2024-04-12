package storage

import (
	"slices"

	"github.com/Alvaroalonsobabbel/echo-go/internal/types"
)

type MemoryStorage struct {
	types.EndpointsWrapper
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		types.EndpointsWrapper{
			Data: []types.Endpoint{},
		},
	}
}

func (m *MemoryStorage) Read() *types.EndpointsWrapper {
	return &m.EndpointsWrapper
}

func (m *MemoryStorage) Create(endpoint types.Endpoint) {
	m.Data = append(m.Data, endpoint)
}

func (m *MemoryStorage) Update(id int, endpoint types.Endpoint) bool {
	for i, e := range m.Data {
		if e.ID == id {
			m.Data = append(m.Data[:i], m.Data[i+1:]...)
			endpoint.ID = id
			m.Data = append(m.Data, endpoint)
			return true
		}
	}
	return false
}

func (m *MemoryStorage) Delete(id int) bool {
	for i, e := range m.Data {
		if e.ID == id {
			m.Data = slices.Delete(m.Data, i, i+1)
			return true
		}
	}
	return false
}

func (m *MemoryStorage) Find(method string, path string) (*types.Endpoint, bool) {
	for _, e := range m.Data {
		if e.Attributes.Verb == method && e.Attributes.Path == path {
			return &e, true
		}
	}
	return &types.Endpoint{}, false
}
