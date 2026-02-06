package api

import (
	"errors"

	"example.com/app/internal/ai"
)

func (s *Store) VerifyAICredentials(projectID string) error {
	p, ok := s.GetProject(projectID)
	if !ok {
		return errors.New("project not found")
	}
	cred, ok := s.GetAICredentials(projectID)
	if !ok {
		return errors.New("missing credentials")
	}
	// verify default
	if err := aiPipelineFrom(p.AIProvider, p.AIModel, cred).Validate(); err != nil {
		return err
	}
	// verify per-task configs
	if err := aiPipelineFrom(p.TextProvider, p.TextModel, cred).Validate(); err != nil {
		return err
	}
	if err := aiPipelineFrom(p.ImageProvider, p.ImageModel, cred).Validate(); err != nil {
		return err
	}
	if err := aiPipelineFrom(p.BgProvider, p.BgModel, cred).Validate(); err != nil {
		return err
	}
	// record verified models
	s.addVerified(projectID, p.AIProvider, p.AIModel)
	s.addVerified(projectID, p.TextProvider, p.TextModel)
	s.addVerified(projectID, p.ImageProvider, p.ImageModel)
	s.addVerified(projectID, p.BgProvider, p.BgModel)
	return nil
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
