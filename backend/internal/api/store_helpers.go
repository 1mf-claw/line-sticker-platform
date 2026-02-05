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
