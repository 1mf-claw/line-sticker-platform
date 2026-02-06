package api

import "example.com/app/internal/ai"

type TaskPipeline struct {
	Provider string
	Model    string
	Cred     AICredentialsRequest
}

func (s *Store) getTaskPipeline(projectID string, provider string, model string) (ai.BYOKPipeline, error) {
	cred, ok := s.GetAICredentials(projectID)
	if !ok {
		return ai.BYOKPipeline{}, aiErr("missing credentials")
	}
	p := ai.BYOKPipeline{
		Provider: provider,
		Model:    model,
		APIKey:   cred.APIKey,
		APIBase:  cred.APIBase,
		Fallback: ai.MockPipeline{},
	}
	if err := p.Validate(); err != nil {
		return p, err
	}
	return p, nil
}
