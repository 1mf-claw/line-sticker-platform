package api

func (s *Store) UpdateProjectAI(projectID string, req AIConfigUpdateRequest) (*Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, _ := s.db.Exec(`UPDATE projects SET ai_provider=?, ai_model=? WHERE id=?`, req.AIProvider, req.AIModel, projectID)
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return nil, false
	}
	return s.GetProject(projectID)
}
