package meter

import (
	meterConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/meter/config"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric"
)

func NewMeterMock() *MeterMock {
	return &MeterMock{}
}

type MeterMock struct {
	mock.Mock
}

func (m *MeterMock) Name() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MeterMock) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MeterMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MeterMock) Config() *meterConfig.Config {
	args := m.Called()
	return args.Get(0).(*meterConfig.Config)
}

func (m *MeterMock) ConfigFile() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MeterMock) InitObservers() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MeterMock) Prometheus() metric.Meter {
	args := m.Called()
	return args.Get(0).(metric.Meter)
}

// Started true if started
func (m *MeterMock) Started() bool {
	args := m.Called()
	return args.Get(0).(bool)
}
