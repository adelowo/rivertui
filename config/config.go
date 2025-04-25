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

	Frontend struct {
		AppURL string `mapstructure:"app_url" yaml:"app_url"`
	} `mapstructure:"frontend" yaml:"frontend"`

	Database struct {
		Postgres struct {
			DSN        string `yaml:"dsn" mapstructure:"dsn"`
			LogQueries bool   `yaml:"log_queries" mapstructure:"log_queries"`
			// How much timeout should be used to run queries agains the db or the
			// context dies
			QueryTimeout time.Duration `yaml:"query_timeout" mapstructure:"query_timeout"`
		} `yaml:"postgres" mapstructure:"postgres"`
	} `yaml:"database" mapstructure:"database"`
}

func (c *Config) Validate() error {

	return nil
}
