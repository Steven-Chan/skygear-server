package service

import "github.com/skygeario/skygear-server/cmd/skynext/model"

type UserStore interface {
	Get(id string) *model.User
}

type MemoryUserStore struct {
	Store map[string]model.User
}

func (s MemoryUserStore) Get(id string) *model.User {
	user, ok := s.Store[id]
	if !ok {
		return nil
	}

	return &user
}
