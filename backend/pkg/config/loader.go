package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Timeout      time.Duration
	PollInterval time.Duration
	Logger       *slog.Logger
	ViperKey     string // Key prefix for storing secrets in viper
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:      30 * time.Second,
		PollInterval: 1 * time.Second,
		Logger:       slog.Default(),
		ViperKey:     "", // Empty means use root level
	}
}

// SecretData is the interface that all secret structures must implement
type SecretData interface {
	// Validate checks if the loaded secrets are valid
	Validate() error
	// GetType returns a string identifier for the secret type (useful for logging)
	GetType() string
}

type Loader[T SecretData] struct {
	config *Config
	viper  *viper.Viper
}

// NewLoader creates a new generic secrets loader
func NewLoader[T SecretData](config *Config) *Loader[T] {
	if config == nil {
		config = DefaultConfig()
	}
	var v *viper.Viper
	v = viper.New()
	v.SetConfigType("json")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &Loader[T]{
		config: config,
		viper:  v,
	}
}

// Load loads secrets with context support
func (l *Loader[T]) Load(ctx context.Context, path string) (T, error) {
	var zero T

	// Validate input path
	if path == "" {
		return zero, fmt.Errorf("secrets path cannot be empty")
	}

	// Clean the path
	path = filepath.Clean(path)

	// Create context with timeout if using background context
	if ctx == context.Background() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, l.config.Timeout)
		defer cancel()
	}

	l.config.Logger.Info("Waiting for secrets file: %s", path)

	ticker := time.NewTicker(l.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("context cancelled while waiting for secrets file: %s (%w)", path, ctx.Err())
		case <-ticker.C:
			if fileInfo, err := os.Stat(path); err == nil {
				if fileInfo.Size() == 0 {
					l.config.Logger.Warn("Secrets file exists but is empty, continuing to wait...")
					continue
				}

				l.config.Logger.Info("Secrets file found, size: %d bytes", fileInfo.Size())
				return l.readAndParseSecrets(path)
			}
		}
	}
}

// readAndParseSecrets handles the file reading and parsing
func (l *Loader[T]) readAndParseSecrets(path string) (T, error) {
	var zero T
	const maxRetries = 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		l.config.Logger.Info("Reading secrets file (attempt %d/%d)", attempt, maxRetries)

		var secrets T
		var err error

		secrets, err = l.loadWithViper(path, attempt, maxRetries)

		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
				continue
			}
			return zero, lastErr
		}

		// Validate the loaded secrets
		if err := secrets.Validate(); err != nil {
			lastErr = fmt.Errorf("secret validation failed: %w", err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
				continue
			}
			return zero, lastErr
		}

		l.config.Logger.Info("Successfully loaded and validated %s secrets", secrets.GetType())
		return secrets, nil
	}

	return zero, lastErr
}

// loadWithViper loads secrets using Viper
func (l *Loader[T]) loadWithViper(path string, attempt, maxRetries int) (T, error) {
	var zero T

	l.viper.SetConfigFile(path)

	if err := l.viper.ReadInConfig(); err != nil {
		return zero, fmt.Errorf("failed to read config file with viper: %w", err)
	}

	l.config.Logger.Info("Parsing secrets with Viper (attempt %d/%d)", attempt, maxRetries)

	var secrets T

	// Try to unmarshal from specific key first, then from root
	if l.config.ViperKey != "" {
		if err := l.viper.UnmarshalKey(l.config.ViperKey, &secrets); err != nil {
			return zero, fmt.Errorf("failed to unmarshal secrets from viper key %s: %w", l.config.ViperKey, err)
		}
	} else {
		if err := l.viper.Unmarshal(&secrets); err != nil {
			return zero, fmt.Errorf("failed to unmarshal secrets from viper: %w", err)
		}
	}

	return secrets, nil
}

// GetViper returns the underlying viper instance for advanced operations
func (l *Loader[T]) GetViper() *viper.Viper {
	return l.viper
}

// WatchSecrets sets up file watching for automatic secret reloading
func (l *Loader[T]) WatchSecrets(path string, onChange func(T, error)) error {
	if l.viper == nil {
		return fmt.Errorf("viper not initialized")
	}

	l.viper.SetConfigFile(path)
	l.viper.WatchConfig()

	if onChange != nil {
		l.viper.OnConfigChange(func(e fsnotify.Event) {
			l.config.Logger.Info("Secrets file changed: %s", e.Name)
			secrets, err := l.readAndParseSecrets(path)
			onChange(secrets, err)
		})
	}

	return nil
}
