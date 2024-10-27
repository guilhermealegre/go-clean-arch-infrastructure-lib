package s3

import (
	"context"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/s3/middlewares/tracer"

	awsMiddleware "github.com/aws/smithy-go/middleware"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	s3Config "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/s3/config"
)

// S3 service
type S3 struct {
	// Name
	name string
	// App
	app domain.IApp
	// Configuration
	config *s3Config.Config
	// Client
	s3Client *s3.Client
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile sqs configuration file
	configFile = "s3.yaml"
)

// New creates a new sqs service
func New(app domain.IApp, config *s3Config.Config) *S3 {
	s := &S3{
		app:  app,
		name: "s3",
	}

	if config != nil {
		s.config = config
	}

	return s
}

// Name gets the service name
func (s *S3) Name() string {
	return s.name
}

// Start starts the sqs service
func (s *S3) Start() (err error) {
	// Load config File
	if s.config == nil {
		s.config = &s3Config.Config{}
		s.config.AdditionalConfig = s.additionalConfigType
		if err = config.Load(s.ConfigFile(), s.config); err != nil {
			err = errors.ErrorLoadingConfigFile().Formats(s.ConfigFile(), err)
			message.ErrorMessage(s.Name(), err)
			return err
		}
	}

	// Init s3 config default
	cfg, _ := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithSharedConfigProfile(
			s.Config().Role,
		),
	)

	// middlewares
	cfg.APIOptions = s.withMiddlewares(cfg.APIOptions, tracer.NewTracerMiddleware(s.app))

	// Set Region
	cfg.Region = s.config.Region
	// Init S3 client
	s.s3Client = s3.NewFromConfig(cfg)

	s.started = true

	return nil
}

// Client stops the s3 service
func (s *S3) Client() *s3.Client {
	return s.s3Client
}

// Stop stops the s3 service
func (s *S3) Stop() error {
	if !s.started {
		return nil
	}
	s.started = false
	return nil
}

// Config gets the service configuration
func (s *S3) Config() *s3Config.Config {
	return s.config
}

// ConfigFile gets the configuration file
func (s *S3) ConfigFile() string {
	return configFile
}

// withMiddlewares adds middlewares to s3
func (s *S3) withMiddlewares(apiOptions []func(*awsMiddleware.Stack) error,
	middlewares ...awsMiddleware.InitializeMiddleware) []func(*awsMiddleware.Stack) error {
	for _, middleware := range middlewares {
		apiOptions = append(apiOptions, func(stack *awsMiddleware.Stack) error {
			return stack.Initialize.Add(middleware, awsMiddleware.Before)
		})
	}
	return apiOptions
}

// Started true if started
func (s *S3) Started() bool {
	return s.started
}
