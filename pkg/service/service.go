package service

import (
	"github.com/richardstrnad/gotmx/pkg/infra/inmemory"
	"github.com/richardstrnad/gotmx/pkg/store"
)

type Configuration func(s *Service) error

type Service struct {
	Store store.Store
}

func New(cfgs ...Configuration) (*Service, error) {
	s := &Service{}
	for _, cfg := range cfgs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func WithInMemoryStore(store *inmemory.InMemoryDataStore) Configuration {
	return func(s *Service) error {
		s.Store = store
		return nil
	}
}
