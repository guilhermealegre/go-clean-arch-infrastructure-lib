package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/aws/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"
	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

// Aws service
type Aws struct {
	// Name
	name string
	// Configuration
	config *awsConfig.Config
	// App
	app domain.IApp
	// Connection
	connection *aws.Config
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile aws configuration file
	configFile = "aws.yaml"
)

// New creates a new aws service
func New(app domain.IApp, config *awsConfig.Config) *Aws {
	aws := &Aws{
		name: "Aws",
		app:  app,
	}

	if config != nil {
		aws.config = config
	}

	return aws
}

// Name gets the service name
func (a *Aws) Name() string {
	return a.name
}

// Start starts the aws service
func (a *Aws) Start() (err error) {
	if a.config == nil {
		a.config = &awsConfig.Config{}
		a.config.AdditionalConfig = a.additionalConfigType
		if err = config.Load(a.ConfigFile(), a.config); err != nil {
			err = errorCodes.ErrorLoadingConfigFile().Formats(a.ConfigFile(), err)
			message.ErrorMessage(a.Name(), err)
			return err
		}
	}

	a.connection = &aws.Config{
		Region: a.config.Region,
	}

	a.started = true

	return nil
}

// Stop stops the aws service
func (a *Aws) Stop() error {
	if !a.started {
		return nil
	}

	a.started = false
	return nil
}

// Config gets the configurations
func (a *Aws) Config() *awsConfig.Config {
	return a.config
}

// Connection gets the aws connection config
func (a *Aws) Connection() *aws.Config {
	return a.connection
}

// ConfigFile gets the configuration file
func (a *Aws) ConfigFile() string {
	return configFile
}

// WithAdditionalConfigType sets an additional config type
func (a *Aws) WithAdditionalConfigType(obj interface{}) domain.IAws {
	a.additionalConfigType = obj
	return a
}

// Started true if started
func (a *Aws) Started() bool {
	return a.started
}
