package faker

import "github.com/juliotorresmoreno/iot/etl/sources"

type FakeUsers struct{}

func NewFakeUsers() *FakeUsers {
	return &FakeUsers{}
}

func (m *FakeUsers) Get(filter *sources.Filter, callback func(rows any, err error)) (any, error) {
	return nil, nil
}

func (m *FakeUsers) Count() uint {
	return 0
}
