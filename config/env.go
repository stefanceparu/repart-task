package config

import "github.com/caarlos0/env"

type LambdaConfig struct {
	Port int `env:"CUSTOM_PORT" envDefault:"8282"`
}

func ParseConfig() (LambdaConfig, error) {
	var lc LambdaConfig
	return lc, env.Parse(&lc)
}
