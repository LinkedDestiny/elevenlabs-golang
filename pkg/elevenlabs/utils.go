package elevenlabs

import "github.com/LinkedDestiny/elevenlabs-golang/pkg/elevenlabs/core"

// PlayAudio is a convenience function that plays audio using the core audio utilities
func PlayAudio(audio []byte) error {
	return core.PlayAudio(audio)
}

// SaveAudio is a convenience function that saves audio to a file
func SaveAudio(audio []byte, filename string) error {
	return core.SaveAudio(audio, filename)
}

// DetectAudioFormat is a convenience function that detects audio format
func DetectAudioFormat(data []byte) core.AudioFormat {
	return core.DetectAudioFormat(data)
}

// ValidateAudioFormat is a convenience function that validates audio format
func ValidateAudioFormat(format core.AudioFormat) bool {
	return core.ValidateAudioFormat(format)
}

// Helper functions for working with pointers

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr returns a pointer to the given float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}

// StringValue returns the value of a string pointer or empty string if nil
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// IntValue returns the value of an int pointer or 0 if nil
func IntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// Float64Value returns the value of a float64 pointer or 0.0 if nil
func Float64Value(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

// BoolValue returns the value of a bool pointer or false if nil
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
