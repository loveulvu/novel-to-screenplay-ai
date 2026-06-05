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
	if err := c.callJSON(ctx, "analyze chapter", BuildChapterAnalysisPrompt(chapter), chapterAnalysisSchemaDescription(), &result); err != nil {
		return result, err
	}
	result.ChapterNumber = chapter.Number
	result.ChapterTitle = chapter.Title
	return result, nil
}

func (c *RealClient) MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (story.StoryBible, error) {
	var result story.StoryBible
	if err := c.callJSON(ctx, "merge story bible", BuildStoryBiblePrompt(analyses), storyBibleSchemaDescription(), &result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *RealClient) GenerateScreenplay(ctx context.Context, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (screenplay.Screenplay, error) {
	var result screenplay.Screenplay
	prompt := BuildScreenplayPrompt(bible, analyses)
	if err := c.callJSON(ctx, "generate screenplay", prompt, screenplaySchemaDescription(), &result); err != nil {
		return result, err
	}

	validation := screenplay.Validate(result)
	if validation.Passed {
		return result, nil
	}

	repairPrompt := BuildRepairScreenplayPrompt(prompt, result, validation.Errors)
	var repaired screenplay.Screenplay
	if err := c.callJSONNoRepair(ctx, "generate screenplay validation repair", repairPrompt, &repaired); err != nil {
		return result, fmt.Errorf("generate screenplay: validation failed and repair parse failed: %w", err)
	}

	repairedValidation := screenplay.Validate(repaired)
	if !repairedValidation.Passed {
		return result, fmt.Errorf("generate screenplay: validation failed after repair: %s", strings.Join(repairedValidation.Errors, "; "))
	}

	return repaired, nil
}

func (c *RealClient) callJSON(ctx context.Context, step string, prompt string, schema string, target any) error {
	raw, err := c.callChatContent(ctx, step, prompt)
	if err != nil {
		return err
	}

	jsonText := extractJSON(raw)
	if err := json.Unmarshal([]byte(jsonText), target); err == nil {
		return nil
	} else {
		repairPrompt := BuildRepairJSONPrompt(prompt, raw, err, schema)
		repairRaw, repairErr := c.callChatContent(ctx, step+" JSON repair", repairPrompt)
		if repairErr != nil {
			return fmt.Errorf("%s: parse failed and repair request failed: %w", step, repairErr)
		}

		repairJSON := extractJSON(repairRaw)
		if repairUnmarshalErr := json.Unmarshal([]byte(repairJSON), target); repairUnmarshalErr != nil {
			return fmt.Errorf("%s: parse failed after repair: %w; output: %s", step, repairUnmarshalErr, summarize(repairJSON, 500))
		}

		return nil
	}
}

func (c *RealClient) callJSONNoRepair(ctx context.Context, step string, prompt string, target any) error {
	raw, err := c.callChatContent(ctx, step, prompt)
	if err != nil {
		return err
	}

	jsonText := extractJSON(raw)
	if err := json.Unmarshal([]byte(jsonText), target); err != nil {
		return fmt.Errorf("%s: parse failed: %w; output: %s", step, err, summarize(jsonText, 500))
	}

	return nil
}

func (c *RealClient) callChatContent(ctx context.Context, step string, prompt string) (string, error) {
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
		return "", fmt.Errorf("%s: marshal request: %w", step, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("%s: create request: %w", step, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: call chat completions: %w", step, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1_000_000))
	if err != nil {
		return "", fmt.Errorf("%s: read response: %w", step, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%s: chat completions returned status %d: %s", step, resp.StatusCode, summarize(string(body), 500))
	}

	var chatResp chatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("%s: parse chat completions response: %w; response: %s", step, err, summarize(string(body), 300))
	}
	if len(chatResp.Choices) == 0 || strings.TrimSpace(chatResp.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("%s: chat completions response has no message content", step)
	}

	return chatResp.Choices[0].Message.Content, nil
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
