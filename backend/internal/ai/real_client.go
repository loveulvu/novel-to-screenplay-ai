package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/novel"
	"novel-to-screenplay-ai/internal/screenplay"
	"novel-to-screenplay-ai/internal/story"
)

type RealClient struct {
	apiKey     string
	model      string
	endpoint   string
	httpClient *http.Client
}

func NewRealClient(cfg Config) *RealClient {
	return &RealClient{
		apiKey:   cfg.APIKey,
		model:    cfg.Model,
		endpoint: chatCompletionsEndpoint(cfg.BaseURL),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *RealClient) AnalyzeChapter(ctx context.Context, chapter novel.Chapter) (analysis.ChapterAnalysis, error) {
	var result analysis.ChapterAnalysis
	if err := c.completeJSON(ctx, "analyze chapter", BuildChapterAnalysisPrompt(chapter), &result); err != nil {
		return result, err
	}
	if result.ChapterNumber == 0 {
		result.ChapterNumber = chapter.Number
	}
	if result.ChapterTitle == "" {
		result.ChapterTitle = chapter.Title
	}
	return result, nil
}

func (c *RealClient) MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (story.StoryBible, error) {
	var result story.StoryBible
	if err := c.completeJSON(ctx, "merge story bible", BuildStoryBiblePrompt(analyses), &result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *RealClient) GenerateScreenplay(ctx context.Context, bible story.StoryBible) (screenplay.Screenplay, error) {
	var result screenplay.Screenplay
	if err := c.completeJSON(ctx, "generate screenplay", BuildScreenplayPrompt(bible), &result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *RealClient) completeJSON(ctx context.Context, step string, prompt string, target any) error {
	reqBody := chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("%s: marshal request: %w", step, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("%s: create request: %w", step, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: call chat completions: %w", step, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1_000_000))
	if err != nil {
		return fmt.Errorf("%s: read response: %w", step, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s: chat completions returned status %d: %s", step, resp.StatusCode, summarize(string(body), 500))
	}

	var chatResp chatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return fmt.Errorf("%s: parse chat completions response: %w; response: %s", step, err, summarize(string(body), 300))
	}
	if len(chatResp.Choices) == 0 || strings.TrimSpace(chatResp.Choices[0].Message.Content) == "" {
		return fmt.Errorf("%s: chat completions response has no message content", step)
	}

	jsonText := extractJSON(chatResp.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(jsonText), target); err != nil {
		return fmt.Errorf("%s: unmarshal model JSON: %w; output: %s", step, err, summarize(jsonText, 500))
	}

	return nil
}

func chatCompletionsEndpoint(baseURL string) string {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasSuffix(base, "/chat/completions") {
		return base
	}
	return base + "/chat/completions"
}

func summarize(value string, limit int) string {
	text := strings.TrimSpace(value)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}
