package s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	s3Config "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/s3/config"
	"github.com/stretchr/testify/mock"
)

func NewSqsMock() *S3Mock {
	return &S3Mock{}
}

type S3Mock struct {
	mock.Mock
}

func (s *S3Mock) Name() string {
	args := s.Called()
	return args.Get(0).(string)
}

func (s *S3Mock) Start() error {
	args := s.Called()
	return args.Error(0)
}

func (s *S3Mock) Stop() error {
	args := s.Called()
	return args.Error(0)
}

func (s *S3Mock) ConfigFile() string {
	args := s.Called()
	return args.Get(0).(string)
}

func (s *S3Mock) Config() *s3Config.Config {
	args := s.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*s3Config.Config)
}

func (s *S3Mock) Client() *s3.Client {
	args := s.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*s3.Client)
}

// WithAdditionalConfigType sets an additional config type
func (s *S3Mock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := s.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (s *S3Mock) Started() bool {
	args := s.Called()
	return args.Get(0).(bool)
}
