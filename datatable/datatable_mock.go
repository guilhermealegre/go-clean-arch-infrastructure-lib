package datatable

import (
	dtDatabase "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/datatable/database"
	dtElastic "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/datatable/elastic_search"
	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/domain"
	"github.com/stretchr/testify/mock"
)

func NewDatatableMock() *DatatableMock {
	return &DatatableMock{}
}

type DatatableMock struct {
	mock.Mock
}

func (d *DatatableMock) Name() string {
	args := d.Called()
	return args.Get(0).(string)
}

func (d *DatatableMock) Start() error {
	args := d.Called()
	return args.Error(0)
}

func (d *DatatableMock) Stop() error {
	args := d.Called()
	return args.Error(0)
}

// Database datatable
func (d *DatatableMock) Database() dtDatabase.IDatabase {
	args := d.Called()
	return args.Get(0).(dtDatabase.IDatabase)
}

// Elastic Elastic
func (d *DatatableMock) Elastic() dtElastic.IElastic {
	args := d.Called()
	return args.Get(0).(dtElastic.IElastic)
}

// Started true if started
func (d *DatatableMock) Started() bool {
	args := d.Called()
	return args.Get(0).(bool)
}

// WithAdditionalConfigType sets an additional config type
func (d *DatatableMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := d.Called(obj)
	return args.Get(0).(domain.IApp)
}
