package api

type ProjectStatus string

type Project struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Theme        string        `json:"theme"`
	StickerCount int           `json:"stickerCount"`
	Status       ProjectStatus `json:"status"`
	CharacterID  string        `json:"characterId"`
	AIProvider   string        `json:"aiProvider"`
	AIModel      string        `json:"aiModel"`
}

type ProjectCreateRequest struct {
	Title        string `json:"title"`
	StickerCount int    `json:"stickerCount"`
}

type ProjectUpdateRequest struct {
	Theme string `json:"theme"`
}

type AIConfigUpdateRequest struct {
	AIProvider string `json:"aiProvider"`
	AIModel    string `json:"aiModel"`
}

type CharacterCreateRequest struct {
	SourceType        string `json:"sourceType"`
	Prompt            string `json:"prompt"`
	ReferenceImageURL string `json:"referenceImageUrl"`
}

type Character struct {
	ID                string `json:"id"`
	SourceType        string `json:"sourceType"`
	ReferenceImageURL string `json:"referenceImageUrl"`
	Status            string `json:"status"`
}

type DraftUpdateRequest struct {
	Caption     string `json:"caption"`
	ImagePrompt string `json:"imagePrompt"`
}

type Draft struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectId"`
	Index       int    `json:"index"`
	Caption     string `json:"caption"`
	ImagePrompt string `json:"imagePrompt"`
	Status      string `json:"status"`
}

type Sticker struct {
	ID             string `json:"id"`
	ProjectID      string `json:"projectId"`
	DraftID        string `json:"draftId"`
	ImageURL       string `json:"imageUrl"`
	TransparentURL string `json:"transparentUrl"`
	Status         string `json:"status"`
}

type Job struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	Progress     int    `json:"progress"`
	ErrorMessage string `json:"errorMessage"`
}

type ThemeSuggestResponse struct {
	Suggestions []string `json:"suggestions"`
}

type ExportResponse struct {
	DownloadURL string `json:"downloadUrl"`
}
