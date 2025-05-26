package config

import (
	"github.com/hengadev/leviosa/pkg/envmode"
)

// Manager provides a convenient way to manage multiple secret loaders
type Manager struct {
	loaders map[string]interface{}
	config  *Config
}

// NewManager creates a new secrets manager
func NewManager(config *Config) *Manager {
	if config == nil {
		config = DefaultConfig()
	}

	return &Manager{
		loaders: make(map[string]interface{}),
		config:  config,
	}
}

// RegisterLoader registers a typed loader with a name
func RegisterLoader[T SecretData](m *Manager, name string, customConfig *Config, env envmode.Mode) *Loader[T] {
	config := m.config
	if customConfig != nil {
		config = customConfig
	}

	loader := NewLoader[T](config)
	m.loaders[name] = loader
	return loader
}

// GetLoader retrieves a registered loader by name
func GetLoader[T SecretData](m *Manager, name string) (*Loader[T], bool) {
	if loader, exists := m.loaders[name]; exists {
		if typedLoader, ok := loader.(*Loader[T]); ok {
			return typedLoader, true
		}
	}
	return nil, false
}
