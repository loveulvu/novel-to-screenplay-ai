package handlers

import (
	"net/http"

	"novel-to-screenplay-ai/internal/ai"
	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/novel"
	"novel-to-screenplay-ai/internal/screenplay"
	"novel-to-screenplay-ai/internal/story"

	"github.com/gin-gonic/gin"
)

type generateRequest struct {
	NovelText string `json:"novel_text"`
}

type generateResponse struct {
	ChapterCount     int                         `json:"chapter_count"`
	ChapterAnalyses  []analysis.ChapterAnalysis  `json:"chapter_analyses"`
	StoryBible       story.StoryBible            `json:"story_bible"`
	ScreenplayJSON   screenplay.Screenplay       `json:"screenplay_json"`
	ScreenplayYAML   string                      `json:"screenplay_yaml"`
	ValidationResult screenplay.ValidationResult `json:"validation"`
}

func Generate(c *gin.Context) {

	var req generateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json body",
		})
		return
	}

	chapters := novel.ParseChapters(req.NovelText)
	if len(chapters) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "at least 3 chapters are required",
		})
		return
	}
	mockClient := ai.NewMockClient()
	analyzer := analysis.NewAnalyzer(mockClient)
	merger := story.NewMerger(mockClient)
	generator := screenplay.NewGenerator(mockClient)

	ctx := c.Request.Context()
	chapterAnalyses, err := analyzer.AnalyzeChapters(ctx, chapters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	storyBible, err := merger.Merge(ctx, chapterAnalyses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	screenplayJSON, err := generator.Generate(ctx, storyBible)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	validation := screenplay.Validate(screenplayJSON)
	screenplayYAML := screenplay.ToYAML(screenplayJSON)

	c.JSON(http.StatusOK, generateResponse{
		ChapterCount:     len(chapters),
		ChapterAnalyses:  chapterAnalyses,
		StoryBible:       storyBible,
		ScreenplayJSON:   screenplayJSON,
		ScreenplayYAML:   screenplayYAML,
		ValidationResult: validation,
	})
}
