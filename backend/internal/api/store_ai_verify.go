package api

import (
	"errors"

	"example.com/app/internal/ai"
)

func (s *Store) VerifyAICredentials(projectID string) error {
	provider, model := s.getProjectAI(projectID)
	cred, ok := s.GetAICredentials(projectID)
	if !ok {
		return errors.New("missing credentials")
	}
	p := aiPipelineFrom(provider, model, cred)
	return p.Validate()
}

func aiPipelineFrom(provider string, model string, cred AICredentialsRequest) aiPipeline {
	return aiPipeline{
		Provider: provider,
		Model:    model,
		APIKey:   cred.APIKey,
		APIBase:  cred.APIBase,
	}
}

type aiPipeline struct {
	Provider string
	Model    string
	APIKey   string
	APIBase  string
}

func (p aiPipeline) Validate() error {
	return ai.BYOKPipeline{
		Provider: p.Provider,
		Model:    p.Model,
		APIKey:   p.APIKey,
		APIBase:  p.APIBase,
	}.Validate()
}
