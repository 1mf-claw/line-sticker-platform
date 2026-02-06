package api

// StoreAICredentials keeps credentials in-memory only (never persisted).
func (s *Store) StoreAICredentials(projectID string, req AICredentialsRequest) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.GetProject(projectID); !ok {
		return false
	}
	s.aiSecrets[projectID] = req
	return true
}

func (s *Store) addVerified(projectID, provider, model string) {
	if _, ok := s.aiVerified[projectID]; !ok {
		s.aiVerified[projectID] = map[string][]string{}
	}
	list := s.aiVerified[projectID][provider]
	for _, m := range list {
		if m == model {
			return
		}
	}
	s.aiVerified[projectID][provider] = append(list, model)
}

func (s *Store) ListVerifiedProviders(projectID string) []Provider {
	providers := []Provider{}
	m := s.aiVerified[projectID]
	for id, models := range m {
		providers = append(providers, Provider{ID: id, Name: id, Models: models})
	}
	return providers
}

func (s *Store) GetAICredentials(projectID string) (AICredentialsRequest, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.aiSecrets[projectID]
	return v, ok
}
