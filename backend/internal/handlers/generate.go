package handlers

import (
	"net/http"

	"novel-to-screenplay-ai/internal/ai"
	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/fidelity"
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
	FidelityResult   fidelity.FidelityResult     `json:"fidelity_result"`
	Meta             generateMeta                `json:"meta"`
}

type generateMeta struct {
	AIProvider string `json:"ai_provider"`
	AIModel    string `json:"ai_model"`
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
	client, err := ai.NewClientFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	runtimeStatus := ai.RuntimeStatusFromEnv()
	analyzer := analysis.NewAnalyzer(client)
	merger := story.NewMerger(client)
	generator := screenplay.NewGenerator(client)
	fidelityChecker := fidelity.NewChecker(client)

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

	screenplayJSON, err := generator.Generate(ctx, storyBible, chapterAnalyses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	applySourceChaptersFromAnalyses(&screenplayJSON, chapterAnalyses)

	screenplayJSON, fidelityResult, err := fidelityChecker.CheckAndRepairOnce(ctx, screenplayJSON, storyBible, chapterAnalyses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	applySourceChaptersFromAnalyses(&screenplayJSON, chapterAnalyses)

	validation := screenplay.Validate(screenplayJSON)
	screenplayYAML := screenplay.ToYAML(screenplayJSON)

	c.JSON(http.StatusOK, generateResponse{
		ChapterCount:     len(chapters),
		ChapterAnalyses:  chapterAnalyses,
		StoryBible:       storyBible,
		ScreenplayJSON:   screenplayJSON,
		ScreenplayYAML:   screenplayYAML,
		ValidationResult: validation,
		FidelityResult:   fidelityResult,
		Meta: generateMeta{
			AIProvider: runtimeStatus.AIProvider,
			AIModel:    runtimeStatus.AIModel,
		},
	})
}

func applySourceChaptersFromAnalyses(target *screenplay.Screenplay, analyses []analysis.ChapterAnalysis) {
	sourceChapters := make([]screenplay.SourceChapter, 0, len(analyses))
	for _, chapter := range analyses {
		sourceChapters = append(sourceChapters, screenplay.SourceChapter{
			Number:  chapter.ChapterNumber,
			Title:   chapter.ChapterTitle,
			Summary: chapter.Summary,
		})
	}
	target.SourceChapters = sourceChapters
}
