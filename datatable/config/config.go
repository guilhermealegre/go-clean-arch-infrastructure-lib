package config

// Config configurations for the database
type Config struct {
	MaxPageResultLimit int
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
