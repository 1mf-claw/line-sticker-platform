package api

type Provider struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Models []string `json:"models"`
}

func ListProviders() []Provider {
	return []Provider{
		{ID: "openai", Name: "OpenAI (ChatGPT)", Models: []string{"gpt-4o-mini", "gpt-4o"}},
		{ID: "replicate", Name: "Replicate", Models: []string{"sdxl", "flux-dev", "flux-schnell"}},
		{ID: "gemini", Name: "Gemini", Models: []string{"gemini-1.5-pro", "gemini-1.5-flash"}},
		{ID: "copilot", Name: "Copilot", Models: []string{"gpt-4o"}},
		{ID: "grok", Name: "Grok", Models: []string{"grok-beta"}},
		{ID: "kimi", Name: "Kimi", Models: []string{"kimi-k2"}},
		{ID: "deepseek", Name: "DeepSeek", Models: []string{"deepseek-chat", "deepseek-reasoner"}},
		{ID: "other", Name: "Other (Custom)", Models: []string{}},
	}
}
