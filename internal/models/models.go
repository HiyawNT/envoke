package models

import (
	"time"

	"gorm.io/gorm"
)

// Environment represents an isolated namespace for secrets (dev, staging, prod, etc.)
type Environment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`
	IsActive  bool           `gorm:"default:false" json:"is_active"`
}

// Secret represents an encrypted key-value pair within an environment
// The Value field stores the encrypted blob, never plaintext
type Secret struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	EnvironmentID uint           `gorm:"not null;index" json:"environment_id"`
	Environment   Environment    `gorm:"foreignKey:EnvironmentID" json:"-"`
	Key           string         `gorm:"not null;index" json:"key"`
	// Value stores the encrypted secret (NaCl secretbox output: nonce + ciphertext)
	Value []byte `gorm:"type:blob;not null" json:"-"`
	// Nonce stores the 24-byte nonce used for encryption (extracted for easier querying)
	Nonce []byte `gorm:"type:blob;size:24;not null" json:"-"`
}

// Metadata stores configuration and key derivation parameters
type Metadata struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
}

// TableName overrides for cleaner DB schema
func (Environment) TableName() string {
	return "environments"
}

func (Secret) TableName() string {
	return "secrets"
}

func (Metadata) TableName() string {
	return "metadata"
}
