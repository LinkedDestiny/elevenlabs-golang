package text_to_speech

// ConvertRequest represents a text-to-speech conversion request
type ConvertRequest struct {
	Text                            string                                   `json:"text"`
	VoiceID                         string                                   `json:"voice_id"`
	ModelID                         *string                                  `json:"model_id,omitempty"`
	LanguageCode                    *string                                  `json:"language_code,omitempty"`
	VoiceSettings                   *VoiceSettings                           `json:"voice_settings,omitempty"`
	PronunciationDictionaryLocators []*PronunciationDictionaryVersionLocator `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            *int                                     `json:"seed,omitempty"`
	PreviousText                    *string                                  `json:"previous_text,omitempty"`
	NextText                        *string                                  `json:"next_text,omitempty"`
	PreviousRequestIDs              []string                                 `json:"previous_request_ids,omitempty"`
	NextRequestIDs                  []string                                 `json:"next_request_ids,omitempty"`
	UsePVCAsIVC                     *bool                                    `json:"use_pvc_as_ivc,omitempty"`
	ApplyTextNormalization          *TextNormalization                       `json:"apply_text_normalization,omitempty"`
	ApplyLanguageTextNormalization  *bool                                    `json:"apply_language_text_normalization,omitempty"`
	EnableLogging                   *bool                                    `json:"enable_logging,omitempty"`
	OptimizeStreamingLatency        *int                                     `json:"optimize_streaming_latency,omitempty"`
	OutputFormat                    *OutputFormat                            `json:"output_format,omitempty"`
}

// StreamRequest represents a streaming text-to-speech request
type StreamRequest struct {
	Text                            string                                   `json:"text"`
	VoiceID                         string                                   `json:"voice_id"`
	ModelID                         *string                                  `json:"model_id,omitempty"`
	LanguageCode                    *string                                  `json:"language_code,omitempty"`
	VoiceSettings                   *VoiceSettings                           `json:"voice_settings,omitempty"`
	PronunciationDictionaryLocators []*PronunciationDictionaryVersionLocator `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            *int                                     `json:"seed,omitempty"`
	PreviousText                    *string                                  `json:"previous_text,omitempty"`
	NextText                        *string                                  `json:"next_text,omitempty"`
	PreviousRequestIDs              []string                                 `json:"previous_request_ids,omitempty"`
	NextRequestIDs                  []string                                 `json:"next_request_ids,omitempty"`
	UsePVCAsIVC                     *bool                                    `json:"use_pvc_as_ivc,omitempty"`
	ApplyTextNormalization          *TextNormalization                       `json:"apply_text_normalization,omitempty"`
	ApplyLanguageTextNormalization  *bool                                    `json:"apply_language_text_normalization,omitempty"`
	EnableLogging                   *bool                                    `json:"enable_logging,omitempty"`
	OptimizeStreamingLatency        *int                                     `json:"optimize_streaming_latency,omitempty"`
	OutputFormat                    *OutputFormat                            `json:"output_format,omitempty"`
}

// RealtimeRequest represents a real-time text-to-speech request
type RealtimeRequest struct {
	VoiceID       string         `json:"voice_id"`
	ModelID       *string        `json:"model_id,omitempty"`
	OutputFormat  *OutputFormat  `json:"output_format,omitempty"`
	VoiceSettings *VoiceSettings `json:"voice_settings,omitempty"`
	TextStream    <-chan string  `json:"-"` // Input text stream
}

// TimestampResponse represents a text-to-speech response with timing information
type TimestampResponse struct {
	AudioBase64 string     `json:"audio_base_64"`
	Alignment   *Alignment `json:"alignment,omitempty"`
}

// TimestampChunk represents a chunk of audio with timing information
type TimestampChunk struct {
	AudioBase64             string   `json:"audio_base_64"`
	CharacterStartTimestamp *float64 `json:"character_start_timestamp,omitempty"`
	CharacterEndTimestamp   *float64 `json:"character_end_timestamp,omitempty"`
	IsFinal                 bool     `json:"is_final"`
}

// Alignment represents timing alignment data
type Alignment struct {
	Characters                 []string  `json:"characters"`
	CharacterStartTimesSeconds []float64 `json:"character_start_times_seconds"`
	CharacterEndTimesSeconds   []float64 `json:"character_end_times_seconds"`
}

// VoiceSettings represents voice configuration settings
type VoiceSettings struct {
	Stability       *float64 `json:"stability,omitempty"`
	SimilarityBoost *float64 `json:"similarity_boost,omitempty"`
	Style           *float64 `json:"style,omitempty"`
	UseSpeakerBoost *bool    `json:"use_speaker_boost,omitempty"`
}

// PronunciationDictionaryVersionLocator represents a pronunciation dictionary reference
type PronunciationDictionaryVersionLocator struct {
	PronunciationDictionaryID string `json:"pronunciation_dictionary_id"`
	VersionID                 string `json:"version_id"`
}

// OutputFormat represents available audio output formats
type OutputFormat string

const (
	OutputFormatMP3_22050_32  OutputFormat = "mp3_22050_32"
	OutputFormatMP3_44100_32  OutputFormat = "mp3_44100_32"
	OutputFormatMP3_44100_64  OutputFormat = "mp3_44100_64"
	OutputFormatMP3_44100_96  OutputFormat = "mp3_44100_96"
	OutputFormatMP3_44100_128 OutputFormat = "mp3_44100_128"
	OutputFormatMP3_44100_192 OutputFormat = "mp3_44100_192"
	OutputFormatPCM_16000     OutputFormat = "pcm_16000"
	OutputFormatPCM_22050     OutputFormat = "pcm_22050"
	OutputFormatPCM_24000     OutputFormat = "pcm_24000"
	OutputFormatPCM_44100     OutputFormat = "pcm_44100"
	OutputFormatULAW_8000     OutputFormat = "ulaw_8000"
)

// TextNormalization represents text normalization options
type TextNormalization string

const (
	TextNormalizationAuto TextNormalization = "auto"
	TextNormalizationOn   TextNormalization = "on"
	TextNormalizationOff  TextNormalization = "off"
)
