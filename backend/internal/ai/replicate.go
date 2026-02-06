package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type replicateRequest struct {
	Version string                 `json:"version"`
	Input   map[string]interface{} `json:"input"`
}

type replicateResponse struct {
	URLs struct {
		Get string `json:"get"`
	} `json:"urls"`
}

type replicateGetResponse struct {
	Status string      `json:"status"`
	Output interface{} `json:"output"`
}

// Very minimal Replicate adapter (async polling)
func (a ReplicateAdapter) Validate(apiKey string, apiBase string, model string) error {
	if apiKey == "" {
		return errors.New("missing api key")
	}
	if model == "" {
		return errors.New("missing model version")
	}
	return nil
}

func (a ReplicateAdapter) GenerateDrafts(apiKey, apiBase, model, theme string, count int, character CharacterInput) ([]DraftIdea, error) {
	base := apiBase
	if base == "" {
		base = "https://api.replicate.com/v1"
	}
	payload := replicateRequest{
		Version: model,
		Input: map[string]interface{}{
			"prompt": "Theme: " + theme + ", Character: " + character.Prompt + ", Count: " + itoa(count),
		},
	}
	body, _ := json.Marshal(payload)
	respBody, err := retry(3, 300*time.Millisecond, func() ([]byte, error) {
		return doReplicateJSON(base+"/predictions", apiKey, body)
	})
	if err != nil {
		return nil, err
	}
	var r replicateResponse
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.URLs.Get == "" {
		return nil, errors.New("no prediction url")
	}
	out, err := retry(3, 300*time.Millisecond, func() (interface{}, error) {
		return replicatePoll(r.URLs.Get, apiKey)
	})
	if err != nil {
		return nil, err
	}
	if arr, ok := out.([]interface{}); ok {
		ideas := []DraftIdea{}
		for _, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				caption, _ := m["caption"].(string)
				prompt, _ := m["imagePrompt"].(string)
				ideas = append(ideas, DraftIdea{Caption: caption, ImagePrompt: prompt})
				continue
			}
			if s, ok := item.(string); ok {
				ideas = append(ideas, DraftIdea{Caption: s, ImagePrompt: s})
			}
		}
		return ideas, nil
	}
	return nil, errors.New("unexpected output")
}

func (a ReplicateAdapter) GenerateImage(apiKey, apiBase, model, prompt string, character CharacterInput) (string, error) {
	base := apiBase
	if base == "" {
		base = "https://api.replicate.com/v1"
	}
	payload := replicateRequest{
		Version: model,
		Input: map[string]interface{}{
			"prompt": prompt,
		},
	}
	body, _ := json.Marshal(payload)
	respBody, err := retry(3, 300*time.Millisecond, func() ([]byte, error) {
		return doReplicateJSON(base+"/predictions", apiKey, body)
	})
	if err != nil {
		return "", err
	}
	var r replicateResponse
	if err := json.Unmarshal(respBody, &r); err != nil {
		return "", err
	}
	// poll once
	if r.URLs.Get == "" {
		return "", errors.New("no prediction url")
	}
	out, err := retry(3, 300*time.Millisecond, func() (interface{}, error) {
		return replicatePoll(r.URLs.Get, apiKey)
	})
	if err != nil {
		return "", err
	}
	if url, ok := out.(string); ok {
		return url, nil
	}
	return "", errors.New("unexpected output")
}

func (a ReplicateAdapter) RemoveBackground(apiKey, apiBase, model, imageURL string) (string, error) {
	return imageURL, nil
}

func doReplicateJSON(url, apiKey string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Token "+apiKey)
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

func replicatePoll(url, apiKey string) (interface{}, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Token "+apiKey)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errors.New("provider error")
	}
	body, err := ioReadAll(resp)
	if err != nil {
		return nil, err
	}
	var r replicateGetResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	return r.Output, nil
}
