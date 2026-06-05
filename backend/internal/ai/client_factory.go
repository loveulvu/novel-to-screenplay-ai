package ai

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var loadEnvOnce sync.Once

type RuntimeStatus struct {
	AIProvider          string `json:"ai_provider"`
	AIModel             string `json:"ai_model"`
	AIBaseURLConfigured bool   `json:"ai_base_url_configured"`
	AIAPIKeyConfigured  bool   `json:"ai_api_key_configured"`
}

func LoadEnv() {
	loadEnvOnce.Do(func() {
		for _, path := range dotenvCandidates() {
			loadDotEnvIfPresent(path)
		}
	})
}

func NewClientFromEnv() (Client, error) {
	LoadEnv()

	provider := normalizedProvider(os.Getenv("AI_PROVIDER"))
	if provider == ProviderMock {
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

func RuntimeStatusFromEnv() RuntimeStatus {
	LoadEnv()

	return RuntimeStatus{
		AIProvider:          normalizedProvider(os.Getenv("AI_PROVIDER")),
		AIModel:             strings.TrimSpace(os.Getenv("AI_MODEL")),
		AIBaseURLConfigured: strings.TrimSpace(os.Getenv("AI_BASE_URL")) != "",
		AIAPIKeyConfigured:  strings.TrimSpace(os.Getenv("AI_API_KEY")) != "",
	}
}

func normalizedProvider(value string) string {
	provider := strings.ToLower(strings.TrimSpace(value))
	if provider == "" {
		return ProviderMock
	}
	return provider
}

func dotenvCandidates() []string {
	paths := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("..", "..", ".env"),
		filepath.Join("backend", ".env"),
		filepath.Join("..", "backend", ".env"),
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		paths = append(paths,
			filepath.Join(exeDir, ".env"),
			filepath.Join(exeDir, "..", ".env"),
			filepath.Join(exeDir, "..", "..", ".env"),
		)
	}

	return uniquePaths(paths)
}

func uniquePaths(paths []string) []string {
	seen := make(map[string]bool, len(paths))
	unique := make([]string, 0, len(paths))
	for _, path := range paths {
		cleaned := filepath.Clean(path)
		if seen[cleaned] {
			continue
		}
		seen[cleaned] = true
		unique = append(unique, cleaned)
	}
	return unique
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
