package elevenlabs

import (
	"net/http"
	"time"

	"github.com/elevenlabs/elevenlabs-golang/pkg/elevenlabs/core"
)

// Predefined environments
var (
	ProductionEnv = core.Environment{
		BaseURL:      "https://api.elevenlabs.io",
		WebSocketURL: "wss://api.elevenlabs.io",
	}
	ProductionUSEnv = core.Environment{
		BaseURL:      "https://api.us.elevenlabs.io",
		WebSocketURL: "wss://api.elevenlabs.io",
	}
	ProductionEUEnv = core.Environment{
		BaseURL:      "https://api.eu.residency.elevenlabs.io",
		WebSocketURL: "wss://api.elevenlabs.io",
	}
)

// Config represents the configuration for the ElevenLabs client
type Config struct {
	APIKey      string
	Environment core.Environment
	Timeout     time.Duration
	HTTPClient  *http.Client
	UserAgent   string
	RetryConfig core.RetryConfig
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Environment: ProductionEnv,
		Timeout:     240 * time.Second,
		UserAgent:   "elevenlabs-golang/v1.0.0",
		RetryConfig: core.DefaultRetryConfig(),
	}
}

// Option represents a configuration option
type Option func(*Config)

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithEnvironment sets the API environment
func WithEnvironment(env core.Environment) Option {
	return func(c *Config) {
		c.Environment = env
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithUserAgent sets a custom user agent
func WithUserAgent(userAgent string) Option {
	return func(c *Config) {
		c.UserAgent = userAgent
	}
}

// WithRetryConfig sets the retry configuration
func WithRetryConfig(retryConfig core.RetryConfig) Option {
	return func(c *Config) {
		c.RetryConfig = retryConfig
	}
}
