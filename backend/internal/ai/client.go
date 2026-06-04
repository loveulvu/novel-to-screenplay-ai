package ai

type Config struct {
	Provider string
	APIKey   string
}

const ProviderMock = "mock"

// TODO: Replace MockClient with a real LLM-backed implementation.
// Keep the public workflow stable: chapter analysis -> story bible -> screenplay.
