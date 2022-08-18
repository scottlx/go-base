package main

import (
	"os"
	"sync"
	"testing"
)

type Config struct {
	GoRoot string
	GoPath string
}

var (
	once   sync.Once
	config *Config
)

func ReadWithOnce() *Config {
	once.Do(func() {
		config = &Config{
			GoRoot: os.Getenv("GOROOT"),
			GoPath: os.Getenv("GOPATH"),
		}
	})
	return config
}

func ReadConfig() *Config {

	config = &Config{
		GoRoot: os.Getenv("GOROOT"),
		GoPath: os.Getenv("GOPATH"),
	}
	return config
}

func BenchmarkReadWithOnce(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = ReadWithOnce()
	}
}

func BenchmarkReadConfig(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = ReadConfig()
	}
}
