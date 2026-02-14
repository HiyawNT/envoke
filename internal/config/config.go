package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	AppName        = "envoke"
	DBFileName     = "envoke.db"
	ConfigName     = "config"
	SaltKey        = "master_salt"
	InitializedKey = "initialized"
)

// Config holds the application configuration
type Config struct {
	v *viper.Viper
}

// New creates a new configuration instance
func New() (*Config, error) {
	v := viper.New()

	// Set config name and type
	v.SetConfigName(ConfigName)
	v.SetConfigType("yaml")

	// Add config search paths
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	v.AddConfigPath(configDir)
	v.AddConfigPath(".")

	// Set defaults
	v.SetDefault(InitializedKey, false)

	// Try to read config (don't error if it doesn't exist yet)
	_ = v.ReadInConfig()

	return &Config{v: v}, nil
}

// GetConfigDir returns the configuration directory for the application
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", AppName)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

// GetDBPath returns the full path to the database file
func GetDBPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, DBFileName), nil
}

// IsInitialized checks if the application has been initialized
func (c *Config) IsInitialized() bool {
	return c.v.GetBool(InitializedKey)
}

// SetInitialized marks the application as initialized
func (c *Config) SetInitialized(initialized bool) error {
	c.v.Set(InitializedKey, initialized)
	return c.Save()
}

// GetSalt retrieves the stored salt for key derivation
func (c *Config) GetSalt() string {
	return c.v.GetString(SaltKey)
}

// SetSalt stores the salt for key derivation
func (c *Config) SetSalt(salt string) error {
	c.v.Set(SaltKey, salt)
	return c.Save()
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, ConfigName+".yaml")
	return c.v.WriteConfigAs(configPath)
}

// Get retrieves a configuration value
func (c *Config) Get(key string) interface{} {
	return c.v.Get(key)
}

// GetString retrieves a string configuration value
func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) {
	c.v.Set(key, value)
}
