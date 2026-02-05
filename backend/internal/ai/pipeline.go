package ai

type CharacterInput struct {
	Prompt            string
	ReferenceImageURL string
	SourceType        string
}

type DraftIdea struct {
	Caption     string
	ImagePrompt string
}

type Pipeline interface {
	GenerateDrafts(theme string, count int, character CharacterInput) ([]DraftIdea, error)
	GenerateImage(prompt string, character CharacterInput) (string, error)
	RemoveBackground(imageURL string) (string, error)
}

// MockPipeline provides deterministic placeholder output for MVP.
type MockPipeline struct{}

func (m MockPipeline) GenerateDrafts(theme string, count int, character CharacterInput) ([]DraftIdea, error) {
	ideas := make([]DraftIdea, 0, count)
	for i := 1; i <= count; i++ {
		ideas = append(ideas, DraftIdea{
			Caption:     theme + " " + "貼圖" + " " + itoa(i),
			ImagePrompt: character.Prompt + " / " + theme + " / action " + itoa(i),
		})
	}
	return ideas, nil
}

func (m MockPipeline) GenerateImage(prompt string, character CharacterInput) (string, error) {
	return "https://example.com/sticker.png", nil
}

func (m MockPipeline) RemoveBackground(imageURL string) (string, error) {
	return "https://example.com/sticker-transparent.png", nil
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}
