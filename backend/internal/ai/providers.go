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

func (a ReplicateAdapter) Validate(apiKey string, apiBase string, model string) error {
	if apiKey == "" {
		return errors.New("missing api key")
	}
	return nil
}

func (a ReplicateAdapter) GenerateDrafts(apiKey, apiBase, model, theme string, count int, character CharacterInput) ([]DraftIdea, error) {
	return nil, errors.New("replicate adapter not implemented")
}

func (a ReplicateAdapter) GenerateImage(apiKey, apiBase, model, prompt string, character CharacterInput) (string, error) {
	return "", errors.New("replicate adapter not implemented")
}

func (a ReplicateAdapter) RemoveBackground(apiKey, apiBase, model, imageURL string) (string, error) {
	return imageURL, nil
}
