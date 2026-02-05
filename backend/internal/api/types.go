package api

type ProjectStatus string

type Project struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Theme        string        `json:"theme"`
	StickerCount int           `json:"stickerCount"`
	Status       ProjectStatus `json:"status"`
	CharacterID  string        `json:"characterId"`
}

type Character struct {
	ID               string `json:"id"`
	SourceType       string `json:"sourceType"`
	ReferenceImageURL string `json:"referenceImageUrl"`
	Status           string `json:"status"`
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
