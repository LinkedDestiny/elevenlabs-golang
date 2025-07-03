package elevenlabs

import (
	"github.com/elevenlabs/elevenlabs-golang/pkg/elevenlabs/core"
	"github.com/elevenlabs/elevenlabs-golang/pkg/elevenlabs/text_to_speech"
	"github.com/elevenlabs/elevenlabs-golang/pkg/elevenlabs/voices"
)

// Client is the main ElevenLabs API client
type Client struct {
	httpClient *core.HTTPClient
	config     Config

	// Service clients (others will be added as modules are implemented)
	TextToSpeech *text_to_speech.Client
	Voices       *voices.Client
}

// NewClient creates a new ElevenLabs client with the specified API key and options
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	config := DefaultConfig()
	config.APIKey = apiKey

	for _, opt := range opts {
		opt(&config)
	}

	return NewClientWithConfig(config)
}

// NewClientWithConfig creates a new ElevenLabs client with the specified configuration
func NewClientWithConfig(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, &APIError{
			statusCode: 0,
			message:    "API key is required",
		}
	}

	// Create core HTTP client config
	coreConfig := core.Config{
		APIKey:      config.APIKey,
		Environment: config.Environment,
		Timeout:     config.Timeout,
		HTTPClient:  config.HTTPClient,
		UserAgent:   config.UserAgent,
		RetryConfig: config.RetryConfig,
	}

	httpClient := core.NewHTTPClient(coreConfig)

	client := &Client{
		httpClient: httpClient,
		config:     config,
	}

	// Initialize service clients
	client.TextToSpeech = text_to_speech.NewClient(httpClient)
	client.Voices = voices.NewClient(httpClient)

	return client, nil
}
