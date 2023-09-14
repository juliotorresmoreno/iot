package db

import (
	"context"
	"errors"

	"github.com/juliotorresmoreno/iot/etl/config"
	"github.com/juliotorresmoreno/iot/etl/db/faker"
	"github.com/juliotorresmoreno/iot/etl/entity"
	"github.com/juliotorresmoreno/iot/etl/sources"
)

var ErrOriginIsInvalid = errors.New("ErrOriginIsInvalid")

type Manager struct {
	limit        uint
	usersManager sources.Source
}

var DefaultManager *Manager

func init() {
	DefaultManager, _ = MakeManager(entity.Fake)
}

func MakeManager(origin entity.Source) (*Manager, error) {
	conf, _ := config.GetConfig()
	result := &Manager{
		limit: uint(conf.Limit),
	}
	if origin == entity.Fake {
		result.usersManager = faker.NewFakeUsers()
		return result, nil
	}

	return result, ErrOriginIsInvalid
}

func (m *Manager) GetUsers(callback func(rows any, err error)) (any, error) {
	limit := m.limit

	filter := &sources.Filter{
		Limit: limit,
		Skip:  0,
	}

	return m.usersManager.Get(filter, callback)
}

func (m *Manager) AddTask(ctx context.Context, callback func()) {

}
