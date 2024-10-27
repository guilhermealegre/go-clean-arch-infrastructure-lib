package meter

import (
	"context"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"
	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
	meterConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/meter/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Meter service
type Meter struct {
	// Name
	name string
	// Configuration
	config *meterConfig.Config
	// App
	app domain.IApp
	// exporter
	*prometheus.Exporter
	// meter
	metric.Meter
	// Started
	started bool
}

const (
	// configFile tracer configuration file
	configFile = "meter.yaml"
)

// New creates a new tracer service
func New(app domain.IApp, config *meterConfig.Config) *Meter {
	meter := &Meter{
		name: "Meter",
		app:  app,
	}

	if config != nil {
		meter.config = config
	}

	return meter
}

// Name gets the service name
func (m *Meter) Name() string {
	return m.name
}

// Start starts the meter service
func (m *Meter) Start() (err error) {
	if m.config == nil {
		m.config = &meterConfig.Config{}
		if err = config.Load(m.ConfigFile(), m.config); err != nil {
			err = errorCodes.ErrorLoadingConfigFile().Formats(m.ConfigFile(), err)
			message.ErrorMessage(m.Name(), err)
			return err
		}
	}

	// prometheus
	if m.config.Enabled {
		exporter, errPrometheus := prometheus.New()
		if errPrometheus != nil {
			return errPrometheus
		}

		res, errResource := resource.New(
			context.Background(),
			resource.WithAttributes(semconv.ServiceName(m.app.Name())))
		if errResource != nil {
			return errResource
		}

		// Creating the MeterProvider and registering it as the global meter provider
		otel.SetMeterProvider(sdkMetric.NewMeterProvider(
			sdkMetric.WithReader(exporter),
			sdkMetric.WithResource(res),
		))

		m.Exporter = exporter
		m.Meter = otel.Meter(m.app.Name())

		err = m.InitObservers()
		if err != nil {
			return err
		}
	} else {
		// setting the meter as a no-op to avoid panics
		m.Meter = noop.Meter{}
	}

	m.started = true

	return nil
}

// ConfigFile gets the configuration file
func (m *Meter) ConfigFile() string {
	return configFile
}

// Config gets the configurations
func (m *Meter) Config() *meterConfig.Config {
	return m.config
}

// Started true if started
func (m *Meter) Started() bool {
	return m.started
}

// Stop stops the meter service
func (m *Meter) Stop() error {
	if !m.started {
		return nil
	}

	if m.config.Enabled {
		return m.Exporter.Shutdown(context.Background())
	}

	m.started = false
	return nil
}

// InitObservers initializes the async observers
func (m *Meter) InitObservers() error {
	for _, initializer := range observers {
		err := initializer(m.Meter)
		if err != nil {
			return err
		}
	}
	return nil
}

// Prometheus returns the Prometheus meter
func (m *Meter) Prometheus() metric.Meter {
	return m.Meter
}
