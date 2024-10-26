package datatable

import (
	databaseConfig "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/database/config"
	datatableConfig "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/datatable/config"
	dtDatabase "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/datatable/database"
	dtElastic "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/datatable/elastic_search"
	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/domain"
)

type Datatable struct {
	// Name
	name string
	// Configuration
	config *databaseConfig.Config //nolint:all
	// App
	app domain.IApp
	// Max results per pages
	maxPageResultLimit int
	// Log functions
	logFunction logFunction //nolint:all
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

type logFunction func(err error) error //nolint:all

// New creates a new database
func New(app domain.IApp, config *datatableConfig.Config) *Datatable {
	datatable := &Datatable{
		name:               "Datatable",
		maxPageResultLimit: dtDatabase.MaxPageResultLimit,
		app:                app,
	}

	if config != nil {
		datatable.maxPageResultLimit = config.MaxPageResultLimit
	}

	return datatable
}

// Name name of the service
func (d *Datatable) Name() string {
	return d.name
}

// Start starts the service
func (d *Datatable) Start() error {
	d.started = true

	return nil
}

// Stop stops the service
func (d *Datatable) Stop() error {
	if !d.started {
		return nil
	}
	d.started = false
	return nil
}

// Database datatable
func (d *Datatable) Database() dtDatabase.IDatabase {
	database := dtDatabase.New(
		dtDatabase.Client{
			Reader: d.app.Database().Read(),
			Writer: d.app.Database().Write(),
		},
		d.app.Logger().DBLog,
		d.maxPageResultLimit,
	)

	return database
}

// Elastic elastic
func (d *Datatable) Elastic() dtElastic.IElastic {
	// Not implemented yet due not any valid use case to use it
	return nil
}

// WithAdditionalConfigType sets an additional config type
func (d *Datatable) WithAdditionalConfigType(obj interface{}) domain.IDatatable {
	d.additionalConfigType = obj
	return d
}

// Started true if started
func (d *Datatable) Started() bool {
	return d.started
}
