package api

type Provider struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Models []string `json:"models"`
}

func ListProviders() []Provider {
	return []Provider{
		{ID: "openai", Name: "OpenAI", Models: []string{"gpt-4o-mini", "gpt-4o"}},
		{ID: "replicate", Name: "Replicate", Models: []string{"sdxl", "flux-dev", "flux-schnell"}},
		{ID: "stability", Name: "Stability", Models: []string{"sd3", "sdxl"}},
	}
}
