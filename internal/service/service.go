package service

import (
	"fmt"

	"github.com/HiyawNT/envoke/internal/crypto"
	"github.com/HiyawNT/envoke/internal/models"
	"github.com/HiyawNT/envoke/internal/storage"
)

// SecretService provides high-level operations for managing secrets
type SecretService struct {
	storage   storage.Storage
	encryptor *crypto.Encryptor
}

// NewSecretService creates a new secret service instance
func NewSecretService(storage storage.Storage, encryptor *crypto.Encryptor) *SecretService {
	return &SecretService{
		storage:   storage,
		encryptor: encryptor,
	}
}

// CreateEnvironment creates a new environment
func (s *SecretService) CreateEnvironment(name string) error {
	_, err := s.storage.CreateEnvironment(name)
	return err
}

// ListEnvironments returns all environments
func (s *SecretService) ListEnvironments() ([]models.Environment, error) {
	return s.storage.ListEnvironments()
}

// DeleteEnvironment deletes an environment
func (s *SecretService) DeleteEnvironment(name string) error {
	return s.storage.DeleteEnvironment(name)
}

// SetActiveEnvironment sets the active environment
func (s *SecretService) SetActiveEnvironment(name string) error {
	return s.storage.SetActiveEnvironment(name)
}

// GetActiveEnvironment returns the active environment
func (s *SecretService) GetActiveEnvironment() (*models.Environment, error) {
	return s.storage.GetActiveEnvironment()
}

// SetSecret encrypts and stores a secret
func (s *SecretService) SetSecret(envName, key, value string) error {
	nonce, ciphertext, err := s.encryptor.EncryptString(value)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	return s.storage.SetSecret(envName, key, nonce, ciphertext)
}

// GetSecret retrieves and decrypts a secret
func (s *SecretService) GetSecret(envName, key string) (string, error) {
	secret, err := s.storage.GetSecret(envName, key)
	if err != nil {
		return "", err
	}

	plaintext, err := s.encryptor.DecryptString(secret.Nonce, secret.Value)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return plaintext, nil
}

// SecretInfo contains metadata about a secret (without the value)
type SecretInfo struct {
	Key       string
	CreatedAt string
	UpdatedAt string
}

// ListSecrets returns information about all secrets in an environment (without decrypting values)
func (s *SecretService) ListSecrets(envName string) ([]SecretInfo, error) {
	secrets, err := s.storage.ListSecrets(envName)
	if err != nil {
		return nil, err
	}

	infos := make([]SecretInfo, len(secrets))
	for i, secret := range secrets {
		infos[i] = SecretInfo{
			Key:       secret.Key,
			CreatedAt: secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: secret.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return infos, nil
}

// ListSecretsWithValues returns all secrets with decrypted values
func (s *SecretService) ListSecretsWithValues(envName string) (map[string]string, error) {
	secrets, err := s.storage.ListSecrets(envName)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, secret := range secrets {
		plaintext, err := s.encryptor.DecryptString(secret.Nonce, secret.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt secret %s: %w", secret.Key, err)
		}
		result[secret.Key] = plaintext
	}

	return result, nil
}

// DeleteSecret deletes a secret
func (s *SecretService) DeleteSecret(envName, key string) error {
	return s.storage.DeleteSecret(envName, key)
}

// Close closes the underlying storage
func (s *SecretService) Close() error {
	return s.storage.Close()
}
