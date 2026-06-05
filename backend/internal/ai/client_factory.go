package ai

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func NewClientFromEnv() (Client, error) {
	loadDotEnvIfPresent("../.env")
	loadDotEnvIfPresent(".env")

	provider := strings.ToLower(strings.TrimSpace(os.Getenv("AI_PROVIDER")))
	if provider == "" || provider == ProviderMock {
		return NewMockClient(), nil
	}

	if provider != ProviderReal {
		return nil, fmt.Errorf("unknown AI_PROVIDER %q; expected mock or real", provider)
	}

	cfg := Config{
		Provider: provider,
		APIKey:   strings.TrimSpace(os.Getenv("AI_API_KEY")),
		BaseURL:  strings.TrimSpace(os.Getenv("AI_BASE_URL")),
		Model:    strings.TrimSpace(os.Getenv("AI_MODEL")),
	}

	var missing []string
	if cfg.APIKey == "" {
		missing = append(missing, "AI_API_KEY")
	}
	if cfg.BaseURL == "" {
		missing = append(missing, "AI_BASE_URL")
	}
	if cfg.Model == "" {
		missing = append(missing, "AI_MODEL")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("AI_PROVIDER=real requires %s", strings.Join(missing, ", "))
	}

	return NewRealClient(cfg), nil
}

func loadDotEnvIfPresent(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		_ = os.Setenv(key, value)
	}
}
