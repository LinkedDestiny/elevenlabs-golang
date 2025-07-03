package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/elevenlabs/elevenlabs-golang/pkg/elevenlabs"
	"github.com/elevenlabs/elevenlabs-golang/pkg/elevenlabs/text_to_speech"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		log.Fatal("ELEVENLABS_API_KEY environment variable is required")
	}

	// Create client
	client, err := elevenlabs.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("ElevenLabs Go SDK Example")
	fmt.Println("========================")

	// Example text
	text := "Hello! This is a demonstration of the ElevenLabs Go SDK. The SDK is working correctly!"

	// Voice ID for Rachel (you can use any valid voice ID)
	voiceID := "21m00Tcm4TlvDq8ikWAM"

	// Create TTS request using utility functions
	req := text_to_speech.ConvertRequest{
		Text:    text,
		VoiceID: voiceID,
		ModelID: elevenlabs.StringPtr("eleven_multilingual_v2"),
		VoiceSettings: &text_to_speech.VoiceSettings{
			Stability:       elevenlabs.Float64Ptr(0.5),
			SimilarityBoost: elevenlabs.Float64Ptr(0.75),
		},
		OutputFormat: (*text_to_speech.OutputFormat)(elevenlabs.StringPtr("mp3_44100_128")),
	}

	fmt.Printf("Converting text to speech:\n%s\n\n", text)

	// Convert text to speech
	audio, err := client.TextToSpeech.Convert(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to convert text to speech: %v", err)
	}

	fmt.Printf("Successfully generated %d bytes of audio\n", len(audio))

	// Save audio to file using utility function
	filename := "output.mp3"
	if err := elevenlabs.SaveAudio(audio, filename); err != nil {
		log.Fatalf("Failed to save audio: %v", err)
	}

	fmt.Printf("Audio saved to: %s\n", filename)

	// Detect audio format
	format := elevenlabs.DetectAudioFormat(audio)
	fmt.Printf("Detected audio format: %s\n", format)

	// Try to play the audio (optional)
	fmt.Println("Attempting to play audio...")
	if err := elevenlabs.PlayAudio(audio); err != nil {
		fmt.Printf("Could not play audio automatically: %v\n", err)
		fmt.Printf("You can manually play the file: %s\n", filename)
	} else {
		fmt.Println("Audio playback completed successfully!")
	}

	fmt.Println("\nExample completed successfully!")
}
