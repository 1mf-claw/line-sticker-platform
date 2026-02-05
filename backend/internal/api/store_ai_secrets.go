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

func (s *Store) GetAICredentials(projectID string) (AICredentialsRequest, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.aiSecrets[projectID]
	return v, ok
}
