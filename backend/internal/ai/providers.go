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

type GenericAdapter struct{ Provider string }

func (g GenericAdapter) Validate(apiKey string, apiBase string, model string) error {
	if apiKey == "" {
		return errors.New("missing api key")
	}
	if model == "" {
		return errors.New("missing model")
	}
	return nil
}

func (g GenericAdapter) GenerateDrafts(apiKey, apiBase, model, theme string, count int, character CharacterInput) ([]DraftIdea, error) {
	return nil, errors.New("provider not implemented")
}

func (g GenericAdapter) GenerateImage(apiKey, apiBase, model, prompt string, character CharacterInput) (string, error) {
	return "", errors.New("provider not implemented")
}

func (g GenericAdapter) RemoveBackground(apiKey, apiBase, model, imageURL string) (string, error) {
	return imageURL, nil
}

func adapterFor(provider string) ProviderAdapter {
	switch provider {
	case "openai", "chatgpt":
		return OpenAIAdapter{}
	case "replicate":
		return ReplicateAdapter{}
	case "gemini", "copilot", "grok", "kimi", "deepseek", "other":
		return GenericAdapter{Provider: provider}
	default:
		return nil
	}
}
