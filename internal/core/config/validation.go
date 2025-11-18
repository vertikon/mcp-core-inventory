// Package config provides configuration loading and validation
package config

import (
	"fmt"
)

// Validate validates the configuration
func Validate(cfg *Config) error {
	if err := validateServer(&cfg.Server); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if err := validateEngine(&cfg.Engine); err != nil {
		return fmt.Errorf("engine config: %w", err)
	}

	if err := validateCache(&cfg.Cache); err != nil {
		return fmt.Errorf("cache config: %w", err)
	}

	if err := validateNATS(&cfg.NATS); err != nil {
		return fmt.Errorf("nats config: %w", err)
	}

	if err := validateLogging(&cfg.Logging); err != nil {
		return fmt.Errorf("logging config: %w", err)
	}

	return nil
}

// validateServer validates server configuration
func validateServer(cfg *ServerConfig) error {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port: %d", cfg.Port)
	}

	if cfg.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}

	if cfg.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}

	return nil
}

// validateEngine validates engine configuration
func validateEngine(cfg *EngineConfig) error {
	workers := GetEngineWorkers(cfg)
	if workers < 0 {
		return fmt.Errorf("workers must be non-negative or 'auto'")
	}

	if cfg.QueueSize <= 0 {
		return fmt.Errorf("queue_size must be positive")
	}

	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	return nil
}

// validateCache validates cache configuration
func validateCache(cfg *CacheConfig) error {
	if cfg.L1Size < 0 {
		return fmt.Errorf("l1_size must be non-negative")
	}

	if cfg.L2TTL <= 0 {
		return fmt.Errorf("l2_ttl must be positive")
	}

	return nil
}

// validateNATS validates NATS configuration
func validateNATS(cfg *NATSConfig) error {
	if len(cfg.URLs) == 0 {
		return fmt.Errorf("nats urls cannot be empty")
	}

	return nil
}

// validateLogging validates logging configuration
func validateLogging(cfg *LoggingConfig) error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[cfg.Level] {
		return fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	validFormats := map[string]bool{
		"json":    true,
		"console": true,
	}

	if !validFormats[cfg.Format] {
		return fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	return nil
}
