package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// AudioFormat represents supported audio formats
type AudioFormat string

const (
	AudioFormatMP3  AudioFormat = "mp3"
	AudioFormatWAV  AudioFormat = "wav"
	AudioFormatPCM  AudioFormat = "pcm"
	AudioFormatULAW AudioFormat = "ulaw"
)

// PlayAudio plays audio bytes using the system's default audio player
func PlayAudio(audio []byte) error {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "elevenlabs_audio_*.mp3")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write audio data to temporary file
	if _, err := tempFile.Write(audio); err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}

	// Sync to ensure data is written
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	// Play the audio using system player
	return playAudioFile(tempFile.Name())
}

// SaveAudio saves audio bytes to a file
func SaveAudio(audio []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.Write(audio)
	if err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}

	return nil
}

// CopyAudio copies audio from src to dst
func CopyAudio(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// DetectAudioFormat attempts to detect the audio format from the data
func DetectAudioFormat(data []byte) AudioFormat {
	if len(data) < 4 {
		return AudioFormatPCM // Default fallback
	}

	// Check for MP3 header
	if len(data) >= 3 && data[0] == 0xFF && (data[1]&0xE0) == 0xE0 {
		return AudioFormatMP3
	}

	// Check for WAV header (RIFF)
	if len(data) >= 4 && string(data[0:4]) == "RIFF" {
		return AudioFormatWAV
	}

	// Default to PCM for unknown formats
	return AudioFormatPCM
}

// ValidateAudioFormat checks if the audio format is supported
func ValidateAudioFormat(format AudioFormat) bool {
	switch format {
	case AudioFormatMP3, AudioFormatWAV, AudioFormatPCM, AudioFormatULAW:
		return true
	default:
		return false
	}
}

// playAudioFile plays an audio file using the system's default player
func playAudioFile(filename string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS - use afplay
		cmd = exec.Command("afplay", filename)
	case "linux":
		// Linux - try common audio players
		if _, err := exec.LookPath("paplay"); err == nil {
			cmd = exec.Command("paplay", filename)
		} else if _, err := exec.LookPath("aplay"); err == nil {
			cmd = exec.Command("aplay", filename)
		} else if _, err := exec.LookPath("mpg123"); err == nil {
			cmd = exec.Command("mpg123", filename)
		} else {
			return fmt.Errorf("no suitable audio player found on Linux")
		}
	case "windows":
		// Windows - use powershell to play
		cmd = exec.Command("powershell", "-c", fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", filename))
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if cmd == nil {
		return fmt.Errorf("no audio player command configured")
	}

	return cmd.Run()
}
