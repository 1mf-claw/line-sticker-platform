package ai

import "errors"

// BYOKPipeline routes calls based on provider/model selected by user.
// This MVP uses MockPipeline for all providers, but keeps the interface for real integration.
type BYOKPipeline struct {
	Provider string
	Model    string
	APIKey   string
	APIBase  string
	Fallback MockPipeline
	Adapter  ProviderAdapter
}

func (p BYOKPipeline) Validate() error {
	if p.Provider == "" || p.Model == "" {
		return errors.New("provider/model required")
	}
	if p.APIKey == "" {
		return errors.New("api key required")
	}
	adapter := p.Adapter
	if adapter == nil {
		adapter = adapterFor(p.Provider)
	}
	if adapter == nil {
		return errors.New("unsupported provider")
	}
	return adapter.Validate(p.APIKey, p.APIBase, p.Model)
}

func (p BYOKPipeline) GenerateDrafts(theme string, count int, character CharacterInput) ([]DraftIdea, error) {
	adapter := p.Adapter
	if adapter == nil {
		adapter = adapterFor(p.Provider)
	}
	if adapter == nil {
		return p.Fallback.GenerateDrafts(theme, count, character)
	}
	return adapter.GenerateDrafts(p.APIKey, p.APIBase, p.Model, theme, count, character)
}

func (p BYOKPipeline) GenerateImage(prompt string, character CharacterInput) (string, error) {
	adapter := p.Adapter
	if adapter == nil {
		adapter = adapterFor(p.Provider)
	}
	if adapter == nil {
		return p.Fallback.GenerateImage(prompt, character)
	}
	return adapter.GenerateImage(p.APIKey, p.APIBase, p.Model, prompt, character)
}

func (p BYOKPipeline) RemoveBackground(imageURL string) (string, error) {
	adapter := p.Adapter
	if adapter == nil {
		adapter = adapterFor(p.Provider)
	}
	if adapter == nil {
		return p.Fallback.RemoveBackground(imageURL)
	}
	return adapter.RemoveBackground(p.APIKey, p.APIBase, p.Model, imageURL)
}

func adapterFor(provider string) ProviderAdapter {
	switch provider {
	case "openai":
		return OpenAIAdapter{}
	case "replicate":
		return ReplicateAdapter{}
	default:
		return nil
	}
}

func (p BYOKPipeline) GenerateDrafts(theme string, count int, character CharacterInput) ([]DraftIdea, error) {
	return p.Fallback.GenerateDrafts(theme, count, character)
}

func (p BYOKPipeline) GenerateImage(prompt string, character CharacterInput) (string, error) {
	return p.Fallback.GenerateImage(prompt, character)
}

func (p BYOKPipeline) RemoveBackground(imageURL string) (string, error) {
	return p.Fallback.RemoveBackground(imageURL)
}
