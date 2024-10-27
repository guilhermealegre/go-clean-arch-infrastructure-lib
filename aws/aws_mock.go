package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/stretchr/testify/mock"
)

func NewAwsMock() *AwsMock {
	return &AwsMock{}
}

type AwsMock struct {
	mock.Mock
}

func (a *AwsMock) Name() string {
	args := a.Called()
	return args.Get(0).(string)
}

func (a *AwsMock) Start() error {
	args := a.Called()
	return args.Error(0)
}

func (a *AwsMock) Stop() error {
	args := a.Called()
	return args.Error(0)
}

func (a *AwsMock) Connection() *aws.Config {
	args := a.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*aws.Config)
}

// WithAdditionalConfigType sets an additional config type
func (a *AwsMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := a.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (a *AwsMock) Started() bool {
	args := a.Called()
	return args.Get(0).(bool)
}
