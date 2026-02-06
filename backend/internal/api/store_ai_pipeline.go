package api

import "example.com/app/internal/ai"

func (s *Store) getPipeline(projectID string) (ai.Pipeline, error) {
	provider, model := s.getProjectAI(projectID)
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

type aiErr string

func (e aiErr) Error() string { return string(e) }
