package ai

import "errors"

// ProviderAdapter defines real provider integrations.
type ProviderAdapter interface {
	Validate(apiKey string, apiBase string, model string) error
	GenerateDrafts(apiKey, apiBase, model, theme string, count int, character CharacterInput) ([]DraftIdea, error)
	GenerateImage(apiKey, apiBase, model, prompt string, character CharacterInput) (string, error)
	RemoveBackground(apiKey, apiBase, model, imageURL string) (string, error)
}

type ReplicateAdapter struct{}
