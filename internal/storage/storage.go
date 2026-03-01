package storage

import (
	"fmt"
	"time"

	"github.com/HiyawNT/envoke/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Storage defines the interface for secret storage operations
type Storage interface {
	// Environment operations
	CreateEnvironment(name string) (*models.Environment, error)
	GetEnvironment(name string) (*models.Environment, error)
	GetEnvironmentByID(id uint) (*models.Environment, error)
	ListEnvironments() ([]models.Environment, error)
	DeleteEnvironment(name string) error
	SetActiveEnvironment(name string) error
	GetActiveEnvironment() (*models.Environment, error)

	// Secret operations
	SetSecret(envName, key string, nonce, encryptedValue []byte) error
	GetSecret(envName, key string) (*models.Secret, error)
	ListSecrets(envName string) ([]models.Secret, error)
	DeleteSecret(envName, key string) error

	// Metadata operations
	SetMetadata(key, value string) error
	GetMetadata(key string) (string, error)

	// Utility
	Close() error
}

// SQLiteStorage implements Storage using SQLite
type SQLiteStorage struct {
	db *gorm.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate schemas
	if err := db.AutoMigrate(&models.Environment{}, &models.Secret{}, &models.Metadata{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

// CreateEnvironment creates a new environment
func (s *SQLiteStorage) CreateEnvironment(name string) (*models.Environment, error) {
	env := &models.Environment{
		Name:     name,
		IsActive: false,
	}

	if err := s.db.Create(env).Error; err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	return env, nil
}

// GetEnvironment retrieves an environment by name
func (s *SQLiteStorage) GetEnvironment(name string) (*models.Environment, error) {
	var env models.Environment
	if err := s.db.Where("name = ?", name).First(&env).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// GetEnvironmentByID retrieves an environment by ID
func (s *SQLiteStorage) GetEnvironmentByID(id uint) (*models.Environment, error) {
	var env models.Environment
	if err := s.db.First(&env, id).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// ListEnvironments lists all environments
func (s *SQLiteStorage) ListEnvironments() ([]models.Environment, error) {
	var envs []models.Environment
	if err := s.db.Order("created_at desc").Find(&envs).Error; err != nil {
		return nil, err
	}
	return envs, nil
}

// DeleteEnvironment deletes an environment and all its secrets
func (s *SQLiteStorage) DeleteEnvironment(name string) error {
	env, err := s.GetEnvironment(name)
	if err != nil {
		return err
	}

	// Delete all secrets in this environment first
	if err := s.db.Where("environment_id = ?", env.ID).Delete(&models.Secret{}).Error; err != nil {
		return fmt.Errorf("failed to delete secrets: %w", err)
	}

	// Delete the environment
	if err := s.db.Delete(env).Error; err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	return nil
}

// SetActiveEnvironment sets an environment as the active one
func (s *SQLiteStorage) SetActiveEnvironment(name string) error {
	// First, deactivate all environments
	if err := s.db.Model(&models.Environment{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to deactivate environments: %w", err)
	}

	// Then activate the specified environment
	env, err := s.GetEnvironment(name)
	if err != nil {
		return err
	}

	env.IsActive = true
	if err := s.db.Save(env).Error; err != nil {
		return fmt.Errorf("failed to activate environment: %w", err)
	}

	return nil
}

// GetActiveEnvironment retrieves the currently active environment
func (s *SQLiteStorage) GetActiveEnvironment() (*models.Environment, error) {
	var env models.Environment
	if err := s.db.Where("is_active = ?", true).First(&env).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// SetSecret creates or updates a secret in the specified environment
func (s *SQLiteStorage) SetSecret(envName, key string, nonce, encryptedValue []byte) error {
	env, err := s.GetEnvironment(envName)
	if err != nil {
		return fmt.Errorf("environment not found: %w", err)
	}

	// Check if secret already exists
	var existing models.Secret
	result := s.db.Where("environment_id = ? AND key = ?", env.ID, key).First(&existing)

	if result.Error == nil {
		// Update existing secret
		existing.Value = encryptedValue
		existing.Nonce = nonce
		existing.UpdatedAt = time.Now()
		if err := s.db.Save(&existing).Error; err != nil {
			return fmt.Errorf("failed to update secret: %w", err)
		}
	} else {
		// Create new secret
		secret := &models.Secret{
			EnvironmentID: env.ID,
			Key:           key,
			Value:         encryptedValue,
			Nonce:         nonce,
		}
		if err := s.db.Create(secret).Error; err != nil {
			return fmt.Errorf("failed to create secret: %w", err)
		}
	}

	return nil
}

// GetSecret retrieves a secret from the specified environment
func (s *SQLiteStorage) GetSecret(envName, key string) (*models.Secret, error) {
	env, err := s.GetEnvironment(envName)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	var secret models.Secret
	if err := s.db.Where("environment_id = ? AND key = ?", env.ID, key).First(&secret).Error; err != nil {
		return nil, err
	}

	return &secret, nil
}

// ListSecrets lists all secrets in the specified environment
func (s *SQLiteStorage) ListSecrets(envName string) ([]models.Secret, error) {
	env, err := s.GetEnvironment(envName)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	var secrets []models.Secret
	if err := s.db.Where("environment_id = ?", env.ID).Order("key asc").Find(&secrets).Error; err != nil {
		return nil, err
	}

	return secrets, nil
}

// DeleteSecret deletes a secret from the specified environment
func (s *SQLiteStorage) DeleteSecret(envName, key string) error {
	env, err := s.GetEnvironment(envName)
	if err != nil {
		return fmt.Errorf("environment not found: %w", err)
	}

	if err := s.db.Where("environment_id = ? AND key = ?", env.ID, key).Delete(&models.Secret{}).Error; err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

// SetMetadata stores a metadata key-value pair
func (s *SQLiteStorage) SetMetadata(key, value string) error {
	var metadata models.Metadata
	result := s.db.Where("key = ?", key).First(&metadata)

	if result.Error == nil {
		// Update existing
		metadata.Value = value
		metadata.UpdatedAt = time.Now()
		if err := s.db.Save(&metadata).Error; err != nil {
			return fmt.Errorf("failed to update metadata: %w", err)
		}
	} else {
		// Create new
		metadata = models.Metadata{
			Key:   key,
			Value: value,
		}
		if err := s.db.Create(&metadata).Error; err != nil {
			return fmt.Errorf("failed to create metadata: %w", err)
		}
	}

	return nil
}

// GetMetadata retrieves a metadata value by key
func (s *SQLiteStorage) GetMetadata(key string) (string, error) {
	var metadata models.Metadata
	if err := s.db.Where("key = ?", key).First(&metadata).Error; err != nil {
		return "", err
	}
	return metadata.Value, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
