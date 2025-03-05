package config

import (
	"os"
	"time"
)

type TimeConfig struct {
	TimeAdd time.Duration
	TimeSub time.Duration
	TimeMul time.Duration
	TimeDiv time.Duration
}

type Config struct {
	Addr string
	TimeConf TimeConfig
}

func parseDurationFromEnv(strEnv string) time.Duration {
	duration_string := os.Getenv(strEnv)
	duration, err := time.ParseDuration(duration_string + "ms")
	if duration_string == "" || err != nil {
		return time.Second
	} else {
		return duration
	}
}

func ConfigFromEnv() *Config {
	config := new(Config)

	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}

	
	config.TimeConf.TimeAdd = parseDurationFromEnv("TIME_ADDITION_MS")
	config.TimeConf.TimeSub = parseDurationFromEnv("TIME_SUBTRACTION_MS")
	config.TimeConf.TimeMul = parseDurationFromEnv("TIME_MULTIPLICATIONS_MS")
	config.TimeConf.TimeDiv = parseDurationFromEnv("TIME_DIVISIONS_MS")

	return config
}