package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type OpenAIAdapter struct{}

type openAIChatRequest struct {
	Model    string         `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}

type openAIImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
}

type openAIImageResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func (a OpenAIAdapter) Validate(apiKey string, apiBase string, model string) error {
	if apiKey == "" {
		return errors.New("missing api key")
	}
	return nil
}

func (a OpenAIAdapter) GenerateDrafts(apiKey, apiBase, model, theme string, count int, character CharacterInput) ([]DraftIdea, error) {
	base := defaultBase(apiBase)
	payload := openAIChatRequest{
		Model: model,
		Messages: []openAIMessage{
			{Role: "system", Content: "You generate sticker drafts. Return JSON array with caption and imagePrompt."},
			{Role: "user", Content: "Theme: " + theme + ". Character: " + character.Prompt + ". Count: " + itoa(count)},
		},
	}

	body, _ := json.Marshal(payload)
	respBody, err := retry(3, 300*time.Millisecond, func() ([]byte, error) {
		return doJSON(base+"/v1/chat/completions", apiKey, body)
	})
	if err != nil {
		return nil, err
	}

	var chatResp openAIChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, err
	}
	if len(chatResp.Choices) == 0 {
		return nil, errors.New("no choices")
	}

	var ideas []DraftIdea
	if err := json.Unmarshal([]byte(chatResp.Choices[0].Message.Content), &ideas); err != nil {
		return nil, err
	}
	return ideas, nil
}

func (a OpenAIAdapter) GenerateImage(apiKey, apiBase, model, prompt string, character CharacterInput) (string, error) {
	base := defaultBase(apiBase)
	payload := openAIImageRequest{
		Model:  model,
		Prompt: prompt,
		Size:   "1024x1024",
	}
	body, _ := json.Marshal(payload)
	respBody, err := retry(3, 300*time.Millisecond, func() ([]byte, error) {
		return doJSON(base+"/v1/images/generations", apiKey, body)
	})
	if err != nil {
		return "", err
	}
	var imgResp openAIImageResponse
	if err := json.Unmarshal(respBody, &imgResp); err != nil {
		return "", err
	}
	if len(imgResp.Data) == 0 {
		return "", errors.New("no image")
	}
	return imgResp.Data[0].URL, nil
}

func (a OpenAIAdapter) RemoveBackground(apiKey, apiBase, model, imageURL string) (string, error) {
	// OpenAI does not provide background removal directly; return original for now.
	return imageURL, nil
}

func defaultBase(base string) string {
	if base == "" {
		return "https://api.openai.com"
	}
	return base
}

func doJSON(url, apiKey string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errors.New("provider error")
	}
	return ioReadAll(resp)
}

func ioReadAll(r *http.Response) ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r.Body)
	return buf.Bytes(), err
}
