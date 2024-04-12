package storage

import (
	"slices"

	"github.com/Alvaroalonsobabbel/echo-go/internal/types"
)

type memoryStorage struct {
	types.EndpointsWrapper
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		types.EndpointsWrapper{
			Data: []types.Endpoint{},
		},
	}
}

func (m *memoryStorage) Read() *types.EndpointsWrapper {
	return &m.EndpointsWrapper
}

func (m *memoryStorage) Create(endpoint types.Endpoint) {
	m.Data = append(m.Data, endpoint)
}

func (m *memoryStorage) Update(id int, endpoint types.Endpoint) bool {
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

func (m *memoryStorage) Delete(id int) bool {
	for i, e := range m.Data {
		if e.ID == id {
			m.Data = slices.Delete(m.Data, i, i+1)
			return true
		}
	}
	return false
}

func (m *memoryStorage) Find(method string, path string) (*types.Endpoint, bool) {
	for _, e := range m.Data {
		if e.Attributes.Verb == method && e.Attributes.Path == path {
			return &e, true
		}
	}
	return &types.Endpoint{}, false
}
