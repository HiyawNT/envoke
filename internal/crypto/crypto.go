package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
)

const (
	// KeySize is the required size for NaCl secretbox keys (32 bytes)
	KeySize = 32
	// NonceSize is the required size for NaCl secretbox nonces (24 bytes)
	NonceSize = 24
	// SaltSize is the size of the salt used for key derivation (16 bytes)
	SaltSize = 16
)

var (
	ErrInvalidKeySize    = errors.New("invalid key size: must be 32 bytes")
	ErrInvalidNonceSize  = errors.New("invalid nonce size: must be 24 bytes")
	ErrDecryptionFailed  = errors.New("decryption failed: message authentication failed")
	ErrInvalidCiphertext = errors.New("invalid ciphertext: too short")
)

// Encryptor handles encryption and decryption of secrets using NaCl secretbox
type Encryptor struct {
	key [KeySize]byte
}

// NewEncryptor creates a new encryptor with the given 32-byte key
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKeySize
	}

	var keyArray [KeySize]byte
	copy(keyArray[:], key)

	return &Encryptor{key: keyArray}, nil
}

// DeriveKey derives a 32-byte encryption key from a passphrase using Argon2id
// Parameters chosen for ~100ms on modern hardware (2 iterations, 64MB memory, 4 threads)
func DeriveKey(passphrase string, salt []byte) []byte {
	if len(salt) != SaltSize {
		// Generate new salt if not provided or invalid
		salt = make([]byte, SaltSize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			panic(fmt.Sprintf("failed to generate salt: %v", err))
		}
	}

	// Argon2id parameters:
	// - time: 2 iterations (adjust based on security requirements)
	// - memory: 64MB (64 * 1024 KB)
	// - threads: 4 (parallelism factor)
	// - keyLen: 32 bytes (256 bits)
	return argon2.IDKey([]byte(passphrase), salt, 2, 64*1024, 4, KeySize)
}

// GenerateSalt creates a cryptographically secure random salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// Encrypt encrypts plaintext using NaCl secretbox (XSalsa20-Poly1305)
func (e *Encryptor) Encrypt(plaintext []byte) (nonce []byte, ciphertext []byte, err error) {
	// Generate a random nonce
	var nonceArray [NonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonceArray[:]); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt using secretbox
	// Output format: Poly1305 tag (16 bytes) + ciphertext
	encrypted := secretbox.Seal(nil, plaintext, &nonceArray, &e.key)

	return nonceArray[:], encrypted, nil
}

// Decrypt decrypts ciphertext using NaCl secretbox
func (e *Encryptor) Decrypt(nonce, ciphertext []byte) ([]byte, error) {
	if len(nonce) != NonceSize {
		return nil, ErrInvalidNonceSize
	}

	var nonceArray [NonceSize]byte
	copy(nonceArray[:], nonce)

	// Decrypt and verify using secretbox
	// This verifies the Poly1305 MAC before decrypting
	plaintext, ok := secretbox.Open(nil, ciphertext, &nonceArray, &e.key)
	if !ok {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// EncryptString is a convenience method for encrypting strings
func (e *Encryptor) EncryptString(plaintext string) (nonce []byte, ciphertext []byte, err error) {
	return e.Encrypt([]byte(plaintext))
}

// DecryptString is a convenience method for decrypting to strings
func (e *Encryptor) DecryptString(nonce, ciphertext []byte) (string, error) {
	plaintext, err := e.Decrypt(nonce, ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncodeSalt encodes a salt to base64 for storage
func EncodeSalt(salt []byte) string {
	return base64.StdEncoding.EncodeToString(salt)
}

// DecodeSalt decodes a base64-encoded salt
func DecodeSalt(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// This stores a hash of the derived key, not the passphrase itself
func CreateVerificationHash(passphrase string, salt []byte) string {

	key := DeriveKey(passphrase, salt)
	hash := make([]byte, 32)
	copy(hash, key)

	return base64.StdEncoding.EncodeToString(hash)
}

func VerifyPassphrase(passphrase string, salt []byte, storedHash string) bool {
	key := DeriveKey(passphrase, salt)
	currentHash := base64.StdEncoding.EncodeToString(key)

	// compare the hashes
	return currentHash == storedHash
}
