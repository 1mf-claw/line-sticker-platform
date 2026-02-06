package ai

import "errors"

// ProviderAdapter is a placeholder for real provider integrations.
type ProviderAdapter interface {
	Validate(apiKey string, apiBase string, model string) error
}

type OpenAIAdapter struct{}

func (a OpenAIAdapter) Validate(apiKey string, apiBase string, model string) error {
	if apiKey == "" {
		return errors.New("missing api key")
	}
	return nil
}

type ReplicateAdapter struct{}

func (a ReplicateAdapter) Validate(apiKey string, apiBase string, model string) error {
	if apiKey == "" {
		return errors.New("missing api key")
	}
	return nil
}
