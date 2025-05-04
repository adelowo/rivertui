package config

import (
	"time"
)

// ENUM(prod,dev)
type LogMode string

type Config struct {
	Logging struct {
		Mode LogMode `yaml:"mode" mapstructure:"mode"`
	} `yaml:"logging" mapstructure:"logging"`

	Database struct {
		DSN        string `yaml:"dsn" mapstructure:"dsn"`
		LogQueries bool   `yaml:"log_queries" mapstructure:"log_queries"`
		// How much timeout should be used to run queries agains the db or the
		// context dies
		QueryTimeout time.Duration `yaml:"query_timeout" mapstructure:"query_timeout"`
	} `yaml:"database" mapstructure:"database"`
}

func (c *Config) Validate() error {

	return nil
}
