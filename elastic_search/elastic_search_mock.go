package elastic_search

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	elasticSearchConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/elastic_search/config"
	"github.com/stretchr/testify/mock"
)

func NewElasticSearchMock() *ElasticSearchMock {
	return &ElasticSearchMock{}
}

type ElasticSearchMock struct {
	mock.Mock
}

func (e *ElasticSearchMock) Name() string {
	args := e.Called()
	return args.Get(0).(string)
}

func (e *ElasticSearchMock) Start() error {
	args := e.Called()
	return args.Error(0)
}

func (e *ElasticSearchMock) Stop() error {
	args := e.Called()
	return args.Error(0)
}

func (e *ElasticSearchMock) ConfigFile() string {
	args := e.Called()
	return args.Get(0).(string)
}

func (e *ElasticSearchMock) Config() *elasticSearchConfig.Config {
	args := e.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*elasticSearchConfig.Config)
}

func (e *ElasticSearchMock) Client() *elasticsearch.Client {
	args := e.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*elasticsearch.Client)
}

// WithAdditionalConfigType sets an additional config type
func (e *ElasticSearchMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := e.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (e *ElasticSearchMock) Started() bool {
	args := e.Called()
	return args.Get(0).(bool)
}
