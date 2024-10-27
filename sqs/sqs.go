package sqs

import (
	"fmt"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
	sqsConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/sqs/config"
)

const (
	// configFile sqs configuration file
	configFile = "sqs.yaml"
)

// New creates a new sqs service
func New(app domain.IApp, config *sqsConfig.Config) *SQS {
	s := &SQS{
		app:         app,
		name:        "SQS",
		connections: map[string]*Connection{},
	}

	if config != nil {
		s.config = config
	}

	return s
}

// Name gets the service name
func (s *SQS) Name() string {
	return s.name
}

// Connection gets a connection
func (s *SQS) Connection(name string) domain.ISQSConnection {
	if conn, ok := s.connections[name]; ok {
		return conn
	}
	return nil
}

func (s *SQS) getConnection(name string, config *sqsConfig.Connection) *Connection {
	conn, ok := s.connections[name]
	if !ok {
		conn = &Connection{
			serviceName: fmt.Sprintf("%s :: %s", s.name, name),
			app:         s.app,
			name:        name,
		}
		s.connections[name] = conn
	}
	conn.config = config

	return conn
}

// WithConsumer adds a new consumer
func (s *SQS) WithConsumer(consumer domain.ISQSConsumer) domain.ISQS {
	conn := s.getConnection(consumer.GetConnection(), nil)
	conn.consumers = append(conn.consumers, consumer)

	return s
}

// Start starts the sqs service
func (s *SQS) Start() (err error) {
	if s.config == nil {
		s.config = &sqsConfig.Config{}
		s.config.AdditionalConfig = s.additionalConfigType
		if err = config.Load(s.ConfigFile(), s.config); err != nil {
			err = errors.ErrorLoadingConfigFile().Formats(s.ConfigFile(), err)
			message.ErrorMessage(s.Name(), err)
			return err
		}
	}

	for name, config := range s.config.Connections {
		conn := s.getConnection(name, config)
		if err = conn.Connect(); err != nil {
			return err
		}

		for _, consumer := range conn.consumers {
			maskedQueue := conn.maskQueue(s.app.Config().Env, consumer.GetQueue())

			if err = conn.migration.Run(maskedQueue); err != nil {
				return err
			}

			go conn.Consume(maskedQueue, consumer)
		}
	}

	s.started = true

	return nil
}

// Stop stops the sqs service
func (s *SQS) Stop() error {
	if !s.started {
		return nil
	}
	s.started = false
	return nil
}

// Config gets the service configuration
func (s *SQS) Config() *sqsConfig.Config {
	return s.config
}

// ConfigFile gets the configuration file
func (s *SQS) ConfigFile() string {
	return configFile
}

// WithAdditionalConfigType sets an additional config type
func (s *SQS) WithAdditionalConfigType(obj interface{}) domain.ISQS {
	s.additionalConfigType = obj
	return s
}

// Started true if started
func (s *SQS) Started() bool {
	return s.started
}
