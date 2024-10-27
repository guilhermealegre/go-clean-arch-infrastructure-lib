package sqs

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	sqsConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/sqs/config"
	"github.com/stretchr/testify/mock"
)

func NewSqsMock() *SqsMock {
	return &SqsMock{}
}

type SqsMock struct {
	mock.Mock
}

func (s *SqsMock) Name() string {
	args := s.Called()
	return args.Get(0).(string)
}

func (s *SqsMock) Start() error {
	args := s.Called()
	return args.Error(0)
}

func (s *SqsMock) Stop() error {
	args := s.Called()
	return args.Error(0)
}

func (s *SqsMock) ConfigFile() string {
	args := s.Called()
	return args.Get(0).(string)
}

func (s *SqsMock) Config() *sqsConfig.Config {
	args := s.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sqsConfig.Config)
}

func (s *SqsMock) WithConsumer(consumer domain.ISQSConsumer) domain.ISQS {
	args := s.Called(consumer)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ISQS)
}

func (s *SqsMock) Connection(name string) domain.ISQSConnection {
	args := s.Called(name)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ISQSConnection)
}

func (s *SqsMock) WithAdditionalConfigType(obj interface{}) domain.ISQS {
	args := s.Called(obj)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ISQS)
}

// Started true if started
func (s *SqsMock) Started() bool {
	args := s.Called()
	return args.Get(0).(bool)
}
