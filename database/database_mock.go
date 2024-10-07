package database

import (
	"github.com/gocraft/dbr/v2"
	databaseConfig "github.com/guilhermealegre/be-clean-arch-infrastructure-lib/database/config"
	"github.com/guilhermealegre/be-clean-arch-infrastructure-lib/domain"
	"github.com/stretchr/testify/mock"
)

func NewDatabaseMock() *DatabaseMock {
	return &DatabaseMock{}
}

type DatabaseMock struct {
	mock.Mock
}

func (d *DatabaseMock) Name() string {
	args := d.Called()
	return args.Get(0).(string)
}

func (d *DatabaseMock) Start() error {
	args := d.Called()
	return args.Error(0)
}

func (d *DatabaseMock) Stop() error {
	args := d.Called()
	return args.Error(0)
}

func (d *DatabaseMock) ConfigFile() string {
	args := d.Called()
	return args.Get(0).(string)
}

func (d *DatabaseMock) Config() *databaseConfig.Config {
	args := d.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*databaseConfig.Config)
}

func (d *DatabaseMock) Read() *dbr.Session {
	args := d.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.Session)
}

func (d *DatabaseMock) Write() *dbr.Session {
	args := d.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.Session)
}

// WithAdditionalConfigType sets an additional config type
func (d *DatabaseMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := d.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (d *DatabaseMock) Started() bool {
	args := d.Called()
	return args.Get(0).(bool)
}
