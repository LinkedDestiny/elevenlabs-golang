# ElevenLabs Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/LinkedDestiny/elevenlabs-golang.svg)](https://pkg.go.dev/github.com/LinkedDestiny/elevenlabs-golang)
[![Go Report Card](https://goreportcard.com/badge/github.com/LinkedDestiny/elevenlabs-golang)](https://goreportcard.com/report/github.com/LinkedDestiny/elevenlabs-golang)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official Go SDK for [ElevenLabs](https://elevenlabs.io/). ElevenLabs brings the most compelling, rich and lifelike voices to creators and developers in just a few lines of code.

## 📖 API & Docs

Check out the [HTTP API documentation](https://elevenlabs.io/docs/api-reference).

## Install

```bash
go get github.com/LinkedDestiny/elevenlabs-golang
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "github.com/LinkedDestiny/elevenlabs-golang/pkg/elevenlabs"
    "github.com/LinkedDestiny/elevenlabs-golang/pkg/elevenlabs/text_to_speech"
)

func main() {
    client, err := elevenlabs.NewClient("YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    
    req := text_to_speech.ConvertRequest{
        Text:         "The first move is what sets everything in motion.",
        VoiceID:      "JBFqnCBsd6RMkjVDRZzb",
        ModelID:      elevenlabs.StringPtr("eleven_multilingual_v2"),
        OutputFormat: (*text_to_speech.OutputFormat)(elevenlabs.StringPtr("mp3_44100_128")),
    }
    
    audio, err := client.TextToSpeech.Convert(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    // Save or play the audio
    elevenlabs.SaveAudio(audio, "output.mp3")
}
```

## ✅ Implementation Status

This SDK provides a complete Go implementation of the ElevenLabs API with the following features:

### Core Infrastructure ✅
- ✅ HTTP client with retry logic and exponential backoff
- ✅ WebSocket client for real-time communication  
- ✅ Comprehensive error handling with typed errors
- ✅ Audio processing and playback utilities
- ✅ File upload and multipart form support
- ✅ Request/response streaming capabilities
- ✅ Environment configuration (Production, EU, US)
- ✅ Configurable timeouts and retry policies

### Text-to-Speech ✅
- ✅ Standard text-to-speech conversion
- ✅ Text-to-speech with timestamps
- ✅ Streaming audio generation
- ✅ Streaming with timestamps
- ✅ Real-time TTS via WebSocket
- ✅ Multiple output formats (MP3, PCM, WAV, μ-law)
- ✅ Voice settings (stability, similarity boost, style)
- ✅ Model selection (Multilingual v2, Flash v2.5, Turbo v2.5)

### Voice Management ✅ 
- ✅ List all voices with filtering options
- ✅ Get individual voice details
- ✅ Voice settings management
- ✅ Voice deletion
- ✅ Sample information retrieval

### Infrastructure Components ✅
- ✅ Pointer utility functions for optional fields
- ✅ Audio format detection and validation
- ✅ Cross-platform audio playback
- ✅ File I/O utilities
- ✅ Type-safe request/response models

## 🚧 Planned Features

The following modules are ready for implementation (infrastructure is complete):

- 🚧 Voice Cloning (IVC/PVC)
- 🚧 Conversational AI agents  
- 🚧 Audio isolation and processing
- 🚧 Speech-to-speech conversion
- 🚧 Speech-to-text transcription
- 🚧 Studio project management
- 🚧 Dubbing workflows
- 🚧 Pronunciation dictionaries
- 🚧 Usage analytics and history
- 🚧 Workspace and user management
- 🚧 Webhook integrations

## Features

- **Text-to-Speech**: Convert text to speech with multiple voice options
- **Streaming**: Real-time audio streaming with chunked responses
- **WebSocket Support**: Real-time bidirectional communication
- **Voice Management**: List, configure, and manage voices
- **Multiple Formats**: Support for MP3, PCM, WAV, and μ-law audio formats
- **Cross-Platform**: Works on macOS, Linux, and Windows
- **Comprehensive Types**: Fully typed API requests and responses
- **Error Handling**: Detailed error information with proper error types
- **Context Support**: Proper cancellation and timeout handling
- **Retry Logic**: Automatic retries with exponential backoff
- **Audio Utilities**: Built-in audio playback and file handling

## Configuration

### Environment Variables

```bash
export ELEVENLABS_API_KEY="your-api-key-here"
```

### Custom Configuration

```go
config := elevenlabs.Config{
    APIKey:      "YOUR_API_KEY",
    Environment: elevenlabs.ProductionEnv, // or ProductionEUEnv, ProductionUSEnv
    Timeout:     60 * time.Second,
    RetryConfig: core.RetryConfig{
        MaxAttempts:       3,
        InitialDelay:      500 * time.Millisecond,
        MaxDelay:          30 * time.Second,
        BackoffMultiplier: 2.0,
        JitterFactor:      0.25,
    },
}

client, err := elevenlabs.NewClientWithConfig(config)
```

## Audio Streaming

```go
req := text_to_speech.StreamRequest{
    Text:    "This is a streaming test with multiple sentences. Each sentence will be processed as it becomes available.",
    VoiceID: "JBFqnCBsd6RMkjVDRZzb",
    ModelID: elevenlabs.StringPtr("eleven_multilingual_v2"),
}

audioStream, err := client.TextToSpeech.Stream(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

for audioChunk := range audioStream {
    // Process each audio chunk as it arrives
    fmt.Printf("Received %d bytes of audio\n", len(audioChunk))
    // You can play chunks immediately or buffer them
}
```

## Voice Management

```go
// List all voices
voices, err := client.Voices.GetAll(context.Background(), voices.GetAllOptions{
    ShowLegacy: elevenlabs.BoolPtr(false),
})

// Get specific voice with settings
voice, err := client.Voices.Get(context.Background(), "voice_id", voices.GetOptions{
    WithSettings: elevenlabs.BoolPtr(true),
})

// Update voice settings
newSettings := voices.VoiceSettings{
    Stability:       elevenlabs.Float64Ptr(0.5),
    SimilarityBoost: elevenlabs.Float64Ptr(0.75),
    Style:           elevenlabs.Float64Ptr(0.3),
    UseSpeakerBoost: elevenlabs.BoolPtr(true),
}
updated, err := client.Voices.EditSettings(context.Background(), "voice_id", newSettings)
```

## Real-time TTS via WebSocket

```go
// Create a text stream
textStream := make(chan string)

req := text_to_speech.RealtimeRequest{
    VoiceID:    "JBFqnCBsd6RMkjVDRZzb", 
    ModelID:    elevenlabs.StringPtr("eleven_multilingual_v2"),
    TextStream: textStream,
}

audioStream, err := client.TextToSpeech.ConvertRealtime(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

// Send text chunks
go func() {
    textStream <- "Hello, "
    time.Sleep(1 * time.Second)
    textStream <- "this is "
    time.Sleep(1 * time.Second) 
    textStream <- "real-time speech!"
    close(textStream)
}()

// Receive audio chunks in real-time
for audioChunk := range audioStream {
    // Process audio immediately as it's generated
    fmt.Printf("Real-time audio chunk: %d bytes\n", len(audioChunk))
}
```

## Utility Functions

The SDK provides helpful utility functions for common tasks:

```go
// Pointer utilities for optional fields
modelID := elevenlabs.StringPtr("eleven_multilingual_v2")
stability := elevenlabs.Float64Ptr(0.5)
enabled := elevenlabs.BoolPtr(true)

// Safe value extraction
name := elevenlabs.StringValue(voice.Name)        // Returns "" if nil
boost := elevenlabs.Float64Value(settings.Boost) // Returns 0.0 if nil
flag := elevenlabs.BoolValue(voice.IsOwner)       // Returns false if nil

// Audio utilities
format := elevenlabs.DetectAudioFormat(audioData)
err := elevenlabs.SaveAudio(audioData, "output.mp3")
err = elevenlabs.PlayAudio(audioData) // Cross-platform playback
```

## Error Handling

The SDK provides comprehensive error handling:

```go
audio, err := client.TextToSpeech.Convert(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *elevenlabs.BadRequestError:
        fmt.Printf("Bad request: %v\n", e.Body())
    case *elevenlabs.ForbiddenError:
        fmt.Println("Insufficient permissions or quota exceeded")
    case *elevenlabs.NotFoundError:
        fmt.Println("Voice or resource not found")
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

## Examples

See the `cmd/examples/` directory for complete examples:

- `basic_tts.go` - Basic text-to-speech conversion with audio playback
- More examples coming as additional modules are implemented

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

## Architecture

The SDK is built with a modular architecture:

- **Core Layer**: HTTP/WebSocket clients, retry logic, streaming, error handling
- **Service Layer**: Individual API modules (TTS, Voices, etc.)
- **Type Layer**: Comprehensive type definitions for all API structures
- **Utility Layer**: Helper functions and audio processing utilities

Each service module is independent and can be used separately if needed.

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests to help improve this SDK.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 