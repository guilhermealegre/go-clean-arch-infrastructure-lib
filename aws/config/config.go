package config

// Config aws configuration
type Config struct {
	// Region
	Region string `yaml:"region"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
