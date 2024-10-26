package state_machine

import (
	"fmt"
	"path"
	"strings"

	stateMachineDomain "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/state_machine/instance"

	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/domain"
)

// StateMachineService service
type StateMachineService struct {
	// App
	app domain.IApp
	// Name
	name string
	// State Machines
	stateMachineMap map[string]stateMachineDomain.IStateMachine
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile state machine configuration file
	stateMachinesConfigPath = "state_machine"
)

// New creates a new state machine service
func New(app domain.IApp) *StateMachineService {
	s := &StateMachineService{
		name:            "StateMachine",
		app:             app,
		stateMachineMap: make(map[string]stateMachineDomain.IStateMachine),
	}

	return s
}

// Name gets the service name
func (s *StateMachineService) Name() string {
	var text []string

	for name := range s.stateMachineMap {
		text = append(text, fmt.Sprintf("%s [ %s ] ready", s.name, name))
	}

	text = append(text, s.name)

	return strings.Join(text, "\n")
}

// Start starts the state machine service
func (s *StateMachineService) Start() (err error) {
	for smName, sm := range s.stateMachineMap {
		err = sm.Load(path.Join("conf/", stateMachinesConfigPath, smName+".json"))
		if err != nil {
			return err
		}
	}

	s.started = true

	return nil
}

// Stop stops the state machine service
func (s *StateMachineService) Stop() error {
	if !s.started {
		return nil
	}
	s.started = false
	return nil
}

// Get get state machine by name
func (s *StateMachineService) Get(name string) stateMachineDomain.IStateMachine {
	sm, ok := s.stateMachineMap[name]
	if !ok {
		s.stateMachineMap[name] = stateMachineDomain.NewStateMachine()
		return s.stateMachineMap[name]
	}
	return sm
}

// WithAdditionalConfigType sets an additional config type
func (s *StateMachineService) WithAdditionalConfigType(obj interface{}) *StateMachineService {
	s.additionalConfigType = obj
	return s
}

// Started true if started
func (s *StateMachineService) Started() bool {
	return s.started
}
