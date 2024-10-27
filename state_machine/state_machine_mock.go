package state_machine

import (
	domain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	stateMachineDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/state_machine/instance"
	"github.com/stretchr/testify/mock"
)

func NewStateMachineServiceMock() *StateMachineServiceMock {
	return &StateMachineServiceMock{}
}

type StateMachineServiceMock struct {
	mock.Mock
}

func (s *StateMachineServiceMock) Name() string {
	args := s.Called()
	return args.Get(0).(string)
}

func (s *StateMachineServiceMock) Get(name string) (stateMachineDomain.IStateMachine, error) {
	args := s.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(stateMachineDomain.IStateMachine), args.Error(1)
}

func (s *StateMachineServiceMock) Start() (err error) {
	args := s.Called()
	return args.Error(0)
}

func (s *StateMachineServiceMock) Stop() (err error) {
	args := s.Called()
	return args.Error(0)
}

// WithAdditionalConfigType sets an additional config type
func (s *StateMachineServiceMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := s.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (s *StateMachineServiceMock) Started() bool {
	args := s.Called()
	return args.Get(0).(bool)
}
