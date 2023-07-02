package dbclient

import (
	"errors"
	"fmt"

	"github.com/Progsilva/employee-service/config"
)

type Config struct {
	Uri string `envconfig:"DATABASE_URL" default:"server=localhost;user id=SA;password=Password.1;port=1433;database=master"`
}

func (c *Config) validate() error {
	if c.Uri == "" {
		return errors.New("uri must be set")
	}
	return nil
}

func defaultConfig() (*Config, error) {
	c := &Config{}
	if err := config.Load(c); err != nil {
		return nil, err
	}
	return c, nil
}

func getConfig() (*Config, error) {
	conf, err := defaultConfig()
	if err != nil {
		return nil, fmt.Errorf("couldn't get config from environment: %v", err)
	}
	if err = conf.validate(); err != nil {
		return nil, fmt.Errorf("config was invalid: %v", err)
	}
	return conf, nil
}
