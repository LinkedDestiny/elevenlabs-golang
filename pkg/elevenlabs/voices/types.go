package voices

// Voice represents a voice in the ElevenLabs system
type Voice struct {
	VoiceID                 string             `json:"voice_id"`
	Name                    string             `json:"name"`
	Samples                 []Sample           `json:"samples"`
	Category                string             `json:"category"`
	FinetuningState         string             `json:"fine_tuning_state"`
	Labels                  map[string]string  `json:"labels"`
	Description             string             `json:"description"`
	PreviewURL              string             `json:"preview_url"`
	AvailableForTiers       []string           `json:"available_for_tiers"`
	Settings                *VoiceSettings     `json:"settings,omitempty"`
	Sharing                 *VoiceSharing      `json:"sharing,omitempty"`
	HighQualityBaseModelIDs []string           `json:"high_quality_base_model_ids"`
	SafetyControl           *string            `json:"safety_control,omitempty"`
	VoiceVerification       *VoiceVerification `json:"voice_verification,omitempty"`
	PermissionOnResource    *string            `json:"permission_on_resource,omitempty"`
	IsLegacy                *bool              `json:"is_legacy,omitempty"`
	IsOwner                 *bool              `json:"is_owner,omitempty"`
}

// Sample represents a voice sample
type Sample struct {
	SampleID  string `json:"sample_id"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int    `json:"size_bytes"`
	Hash      string `json:"hash"`
}

// VoiceSettings represents voice configuration settings
type VoiceSettings struct {
	Stability       *float64 `json:"stability,omitempty"`
	SimilarityBoost *float64 `json:"similarity_boost,omitempty"`
	Style           *float64 `json:"style,omitempty"`
	UseSpeakerBoost *bool    `json:"use_speaker_boost,omitempty"`
}

// VoiceSharing represents voice sharing configuration
type VoiceSharing struct {
	Status              string            `json:"status"`
	HistoryItemSampleID string            `json:"history_item_sample_id"`
	OriginalVoiceID     string            `json:"original_voice_id"`
	PublicOwnerID       string            `json:"public_owner_id"`
	LikedByCount        int               `json:"liked_by_count"`
	ClonedByCount       int               `json:"cloned_by_count"`
	WhitelistedEmails   []string          `json:"whitelisted_emails"`
	Name                string            `json:"name"`
	Labels              map[string]string `json:"labels"`
	Description         string            `json:"description"`
	ReviewStatus        string            `json:"review_status"`
	ReviewMessage       string            `json:"review_message"`
	EnabledInLibrary    bool              `json:"enabled_in_library"`
	InstantCloning      bool              `json:"instant_cloning"`
	NoticePeroid        int               `json:"notice_period"`
	DisableLogs         bool              `json:"disable_logs"`
	VoiceMixingAllowed  bool              `json:"voice_mixing_allowed"`
	FeaturedDeeplyAi    bool              `json:"featured_deeply_ai"`
	CategoryOriginal    string            `json:"category_original"`
	CategoryHighQuality string            `json:"category_high_quality"`
}

// VoiceVerification represents voice verification status
type VoiceVerification struct {
	RequiresVerification bool     `json:"requires_verification"`
	IsVerified           bool     `json:"is_verified"`
	VerificationFailures []string `json:"verification_failures"`
	VerificationAttempts []string `json:"verification_attempts"`
	Language             *string  `json:"language,omitempty"`
	MinSpeakingDuration  *float64 `json:"min_speaking_duration,omitempty"`
}

// VoicesResponse represents the response from the voices list endpoint
type VoicesResponse struct {
	Voices []Voice `json:"voices"`
}

// GetAllOptions represents options for the GetAll method
type GetAllOptions struct {
	ShowLegacy *bool `json:"show_legacy,omitempty"`
}

// GetOptions represents options for the Get method
type GetOptions struct {
	WithSettings *bool `json:"with_settings,omitempty"`
}

// VoiceCloneRequest represents a voice cloning request
type VoiceCloneRequest struct {
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Files       []string          `json:"files"` // File paths or URLs
}

// AddVoiceRequest represents a request to add a new voice
type AddVoiceRequest struct {
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Files       []string          `json:"files"` // File paths
}

// EditVoiceRequest represents a request to edit a voice
type EditVoiceRequest struct {
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// VoiceGenerationSettings represents settings for voice generation
type VoiceGenerationSettings struct {
	Gender         *string  `json:"gender,omitempty"`          // "male" or "female"
	Age            *string  `json:"age,omitempty"`             // "young", "middle_aged", "old"
	Accent         *string  `json:"accent,omitempty"`          // "american", "british", etc.
	AccentStrength *float64 `json:"accent_strength,omitempty"` // 0.0 to 2.0
}
