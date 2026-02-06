package api

func (s *Store) UpdateProjectPipeline(projectID string, req AIPipelineConfigRequest) (*Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, _ := s.db.Exec(`UPDATE projects SET text_provider=?, text_model=?, image_provider=?, image_model=?, bg_provider=?, bg_model=? WHERE id=?`,
		req.TextProvider, req.TextModel, req.ImageProvider, req.ImageModel, req.BgProvider, req.BgModel, projectID,
	)
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return nil, false
	}
	return s.GetProject(projectID)
}
