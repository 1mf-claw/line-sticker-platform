package api

import "example.com/app/internal/ai"

func (s *Store) getCharacterInput(projectID string) ai.CharacterInput {
	row := s.db.QueryRow(`SELECT source_type, reference_image_url FROM characters WHERE project_id=? ORDER BY rowid DESC LIMIT 1`, projectID)
	var sourceType, refURL string
	_ = row.Scan(&sourceType, &refURL)
	return ai.CharacterInput{
		Prompt:            "main character",
		ReferenceImageURL: refURL,
		SourceType:        sourceType,
	}
}

func (s *Store) getProjectAI(projectID string) (string, string) {
	row := s.db.QueryRow(`SELECT ai_provider, ai_model FROM projects WHERE id=?`, projectID)
	var provider, model string
	_ = row.Scan(&provider, &model)
	return provider, model
}

func resolveProviderModel(primaryProvider, primaryModel, fallbackProvider, fallbackModel string) (string, string) {
	provider := primaryProvider
	model := primaryModel
	if provider == "" {
		provider = fallbackProvider
	}
	if model == "" {
		model = fallbackModel
	}
	return provider, model
}
