// Package encryption currently exposes lightweight stubs that satisfy the
// configuration layer while the full crypto pipeline is implemented.
package encryption

import "time"

// EncryptionManager is a minimal interface for encrypt/decrypt operations.
type EncryptionManager interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

// KeyManager represents a minimal key manager contract.
type KeyManager interface {
	Rotate() error
	ActiveKey() []byte
}

// KeyManagerConfig carries the basic settings for the stub manager.
type KeyManagerConfig struct {
	RotationTTL time.Duration
	KeySize     int
}

type stubKeyManager struct {
	config KeyManagerConfig
	key    []byte
}

// NewKeyManager returns a no-op key manager.
func NewKeyManager(cfg KeyManagerConfig) KeyManager {
	return &stubKeyManager{config: cfg, key: make([]byte, cfg.KeySize)}
}

func (km *stubKeyManager) Rotate() error { return nil }
func (km *stubKeyManager) ActiveKey() []byte {
	return km.key
}

type stubEncryptionManager struct {
	keyManager KeyManager
}

// NewEncryptionManager returns a no-op encryption manager.
func NewEncryptionManager(km KeyManager) EncryptionManager {
	return &stubEncryptionManager{keyManager: km}
}

func (m *stubEncryptionManager) Encrypt(data []byte) ([]byte, error) {
	// Return a copy so callers don't mutate the original slice.
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func (m *stubEncryptionManager) Decrypt(data []byte) ([]byte, error) {
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}
