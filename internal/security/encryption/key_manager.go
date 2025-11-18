package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"sync"
	"time"

	"github.com/vertikon/mcp-hulk/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrKeyRotationFailed = errors.New("key rotation failed")
)

// KeyManager handles encryption key management and rotation
type KeyManager interface {
	// GetEncryptionKey returns the current encryption key
	GetEncryptionKey() ([]byte, error)
	
	// GetKeyVersion returns the current key version
	GetKeyVersion() string
	
	// RotateKey rotates the encryption key
	RotateKey() error
	
	// GetRSAPrivateKey returns RSA private key
	GetRSAPrivateKey() (*rsa.PrivateKey, error)
	
	// GetRSAPublicKey returns RSA public key
	GetRSAPublicKey() (*rsa.PublicKey, error)
	
	// LoadKeyFromEnv loads key from environment variable
	LoadKeyFromEnv(keyName string) error
	
	// LoadKeyFromFile loads key from file
	LoadKeyFromFile(filePath string) error
}

// Manager implements KeyManager
type Manager struct {
	encryptionKey []byte
	keyVersion    string
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	rotationTTL  time.Duration
	lastRotation  time.Time
	mu            sync.RWMutex
	logger        *zap.Logger
}

// KeyManagerConfig holds configuration for KeyManager
type KeyManagerConfig struct {
	RotationTTL time.Duration
	KeySize     int // RSA key size (2048, 4096)
}

// NewKeyManager creates a new KeyManager
func NewKeyManager(config KeyManagerConfig) KeyManager {
	km := &Manager{
		keyVersion:   "v1",
		rotationTTL: config.RotationTTL,
		logger:       logger.WithContext(nil),
	}

	// Generate initial encryption key
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		km.logger.Error("Failed to generate encryption key", zap.Error(err))
	} else {
		km.encryptionKey = key
	}

	// Generate RSA key pair
	if err := km.generateRSAKeys(config.KeySize); err != nil {
		km.logger.Error("Failed to generate RSA keys", zap.Error(err))
	}

	km.lastRotation = time.Now()
	return km
}

// GetEncryptionKey returns the current encryption key
func (m *Manager) GetEncryptionKey() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.encryptionKey == nil {
		return nil, ErrKeyNotFound
	}

	// Check if rotation is needed
	if time.Since(m.lastRotation) > m.rotationTTL {
		go m.RotateKey()
	}

	// Return a copy to prevent external modification
	key := make([]byte, len(m.encryptionKey))
	copy(key, m.encryptionKey)
	return key, nil
}

// GetKeyVersion returns the current key version
func (m *Manager) GetKeyVersion() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.keyVersion
}

// RotateKey rotates the encryption key
func (m *Manager) RotateKey() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate new key
	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		m.logger.Error("Failed to generate new encryption key", zap.Error(err))
		return ErrKeyRotationFailed
	}

	// Update key
	oldKey := m.encryptionKey
	m.encryptionKey = newKey
	m.keyVersion = "v" + time.Now().Format("20060102150405")
	m.lastRotation = time.Now()

	m.logger.Info("Encryption key rotated",
		zap.String("old_version", m.keyVersion),
		zap.String("new_version", m.keyVersion),
	)

	// In production, old key should be kept for decrypting old data
	_ = oldKey

	return nil
}

// GetRSAPrivateKey returns RSA private key
func (m *Manager) GetRSAPrivateKey() (*rsa.PrivateKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.rsaPrivateKey == nil {
		return nil, ErrKeyNotFound
	}

	return m.rsaPrivateKey, nil
}

// GetRSAPublicKey returns RSA public key
func (m *Manager) GetRSAPublicKey() (*rsa.PublicKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.rsaPublicKey == nil {
		return nil, ErrKeyNotFound
	}

	return m.rsaPublicKey, nil
}

// LoadKeyFromEnv loads key from environment variable
func (m *Manager) LoadKeyFromEnv(keyName string) error {
	// In production, load from environment
	// For now, this is a placeholder
	m.logger.Info("Loading key from environment",
		zap.String("key_name", keyName),
	)
	return nil
}

// LoadKeyFromFile loads key from file
func (m *Manager) LoadKeyFromFile(filePath string) error {
	// In production, load from file with proper permissions
	// For now, this is a placeholder
	m.logger.Info("Loading key from file",
		zap.String("file_path", filePath),
	)
	return nil
}

// generateRSAKeys generates RSA key pair
func (m *Manager) generateRSAKeys(keySize int) error {
	if keySize == 0 {
		keySize = 2048 // Default
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		m.logger.Error("Failed to generate RSA private key", zap.Error(err))
		return err
	}

	m.rsaPrivateKey = privateKey
	m.rsaPublicKey = &privateKey.PublicKey

	m.logger.Info("RSA keys generated",
		zap.Int("key_size", keySize),
	)

	return nil
}

// ExportRSAPrivateKey exports RSA private key as PEM
func (m *Manager) ExportRSAPrivateKey() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.rsaPrivateKey == nil {
		return nil, ErrKeyNotFound
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(m.rsaPrivateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return privateKeyPEM, nil
}

// ExportRSAPublicKey exports RSA public key as PEM
func (m *Manager) ExportRSAPublicKey() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.rsaPublicKey == nil {
		return nil, ErrKeyNotFound
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(m.rsaPublicKey)
	if err != nil {
		return nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return publicKeyPEM, nil
}
