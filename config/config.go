package config

import (
	"fmt"
	"strings"

	"gopkg.in/caarlos0/env.v2"
)

const (
	MIN_TOKEN_LENGTH = 16
)

type Config struct {
	ListenAddress string   `env:"PORT" envDefault:"80"`
	HealthFreqSec int      `env:"HELLOFOUR_HEALTH_FREQ_SEC" envDefault:"60"`
	EnvName       string   `env:"HELLOFOUR_ENV_NAME" envDefault:"dev"`
	Tokens        []string `env:"HELLOFOUR_TOKENS" envDefault:"aaaabbbb11112222"`
	ServiceName   string   `env:"HELLOFOUR_SERVICE_NAME" envDefault:"hellofour"`

	StatsDAddress string  `env:"HELLOFOUR_STATSD_ADDRESS" envDefault:"localhost:8125"`
	StatsDPrefix  string  `env:"HELLOFOUR_STATSD_PREFIX" envDefault:"statsd.hellofour.dev"`
	StatsDRate    float32 `env:"HELLOFOUR_STATSD_RATE" envDefault:"1.0"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadEnvVars() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("Unable to fetch env vars: %v", err.Error())
	}

	var errorList []string

	validations := []validation{
		nonEmptyStringSlice{s: c.Tokens, name: "HELLOFOUR_TOKENS"},
		tokenLength{s: c.Tokens, name: "HELLOFOUR_TOKENS"},
	}

	for _, v := range validations {
		if ok, e := v.validate(); !ok {
			errorList = append(errorList, e...)
		}
	}

	if len(errorList) != 0 {
		return fmt.Errorf(strings.Join(errorList, "; "))
	}

	return nil
}

type validation interface {
	validate() (bool, []string)
}

type nonEmptyString struct {
	s, name string
}

func (v nonEmptyString) validate() (bool, []string) {
	if v.s == "" {
		return false, []string{fmt.Sprintf("missing '%s' env var", v.name)}
	}

	return true, nil
}

type nonEmptyStringSlice struct {
	s    []string
	name string
}

func (v nonEmptyStringSlice) validate() (bool, []string) {
	if len(v.s) < 1 {
		return false, []string{fmt.Sprintf("missing '%s' env var", v.name)}
	}

	return true, nil
}

type tokenLength struct {
	s    []string
	name string
}

func (t tokenLength) validate() (bool, []string) {
	var errorList []string

	for _, token := range t.s {
		if len(token) < MIN_TOKEN_LENGTH {
			errorList = append(errorList, fmt.Sprintf("%v token must be at least %v chars long", token, MIN_TOKEN_LENGTH))
		}
	}

	if len(errorList) > 0 {
		return false, errorList
	}

	return true, nil
}
