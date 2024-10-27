package elastic_search

import (
	"fmt"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/elastic_search/middlewares/tracer"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"

	"github.com/elastic/go-elasticsearch/v7/estransport"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	elasticSearchConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/elastic_search/config"
)

// ElasticSearch elastic search
type ElasticSearch struct {
	// Name
	name string
	// Configuration
	config *elasticSearchConfig.Config
	// Client
	client *elasticsearch.Client
	// App
	app domain.IApp
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile elastic search configuration file
	configFile = "elastic_search.yaml"
)

// New creates a new elastic search
func New(app domain.IApp, config *elasticSearchConfig.Config) *ElasticSearch {
	elasticSearch := &ElasticSearch{
		name: "Elastic Search",
		app:  app,
	}

	if config != nil {
		elasticSearch.config = config
	}

	return elasticSearch
}

// Name service name
func (es *ElasticSearch) Name() string {
	return es.name
}

// Start starts the elastic search service
func (es *ElasticSearch) Start() (err error) {
	if es.config == nil {
		es.config = &elasticSearchConfig.Config{}
		es.config.AdditionalConfig = es.additionalConfigType
		if err = config.Load(es.ConfigFile(), es.config); err != nil {
			err = errorCodes.ErrorLoadingConfigFile().Formats(es.ConfigFile(), err)
			message.ErrorMessage(es.Name(), err)
			return err
		}
	}

	address := es.config.Host

	if es.config.Port > 0 {
		address += fmt.Sprintf(":%d", es.config.Port)
	}

	config := elasticsearch.Config{
		Addresses: []string{address},
		Username:  es.config.User,
		Password:  es.config.Password,
	}

	if es.client, err = elasticsearch.NewClient(config); err != nil {
		return err
	}

	es.withMiddleware(tracer.NewTracerMiddleware(es.app, es.client))

	es.started = true

	return nil
}

// Stop stops the elastic search service
func (es *ElasticSearch) Stop() (err error) {
	if !es.started {
		return nil
	}
	es.started = false
	return nil
}

// Config gets the elastic search configuration
func (es *ElasticSearch) Config() *elasticSearchConfig.Config {
	return es.config
}

// ConfigFile gets the configuration file
func (es *ElasticSearch) ConfigFile() string {
	return configFile
}

// Client gets the client for elastic search
func (es *ElasticSearch) Client() *elasticsearch.Client {
	return es.client
}

func (es *ElasticSearch) withMiddleware(tracer estransport.Interface) {
	es.client.Transport = tracer
}

// WithAdditionalConfigType sets an additional config type
func (es *ElasticSearch) WithAdditionalConfigType(obj interface{}) domain.IElasticSearch {
	es.additionalConfigType = obj
	return es
}

// Started true if started
func (es *ElasticSearch) Started() bool {
	return es.started
}
