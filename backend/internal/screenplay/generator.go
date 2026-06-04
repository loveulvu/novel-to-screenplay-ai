package screenplay

import (
	"context"

	"novel-to-screenplay-ai/internal/story"
)

type ScreenplayAIClient interface {
	GenerateScreenplay(ctx context.Context, bible story.StoryBible) (Screenplay, error)
}

type Generator struct {
	client ScreenplayAIClient
}

func NewGenerator(client ScreenplayAIClient) Generator {
	return Generator{client: client}
}

func (g Generator) Generate(ctx context.Context, bible story.StoryBible) (Screenplay, error) {
	return g.client.GenerateScreenplay(ctx, bible)
}
