package ai

import (
	"context"

	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/novel"
	"novel-to-screenplay-ai/internal/screenplay"
	"novel-to-screenplay-ai/internal/story"
)

type Config struct {
	Provider string
	APIKey   string
	BaseURL  string
	Model    string
}

const ProviderMock = "mock"
const ProviderReal = "real"

type Client interface {
	AnalyzeChapter(ctx context.Context, chapter novel.Chapter) (analysis.ChapterAnalysis, error)
	MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (story.StoryBible, error)
	GenerateScreenplay(ctx context.Context, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (screenplay.Screenplay, error)
}
