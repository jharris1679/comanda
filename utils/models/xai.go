package models

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// XAIProvider handles X.AI family of models
type XAIProvider struct {
	apiKey  string
	config  ModelConfig
	verbose bool
}

// NewXAIProvider creates a new X.AI provider instance
func NewXAIProvider() *XAIProvider {
	return &XAIProvider{
		config: ModelConfig{
			Temperature: 0.7,
			MaxTokens:   2000,
			TopP:        1.0,
		},
	}
}

// Name returns the provider name
func (x *XAIProvider) Name() string {
	return "xai"
}

// debugf prints debug information if verbose mode is enabled
func (x *XAIProvider) debugf(format string, args ...interface{}) {
	if x.verbose {
		fmt.Printf("[DEBUG][XAI] "+format+"\n", args...)
	}
}

// SupportsModel checks if the given model name is supported by X.AI
func (x *XAIProvider) SupportsModel(modelName string) bool {
	x.debugf("Checking if model is supported: %s", modelName)
	modelName = strings.ToLower(modelName)

	// Accept any model name that starts with our known prefixes
	validPrefixes := []string{
		"grok-",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(modelName, prefix) {
			x.debugf("Model %s is supported (matches prefix %s)", modelName, prefix)
			return true
		}
	}

	x.debugf("Model %s is not supported (no matching prefix)", modelName)
	return false
}

// Configure sets up the provider with necessary credentials
func (x *XAIProvider) Configure(apiKey string) error {
	x.debugf("Configuring X.AI provider")
	if apiKey == "" {
		return fmt.Errorf("API key is required for X.AI provider")
	}
	x.apiKey = apiKey
	x.debugf("API key configured successfully")
	return nil
}

// SendPrompt sends a prompt to the specified model and returns the response
func (x *XAIProvider) SendPrompt(modelName string, prompt string) (string, error) {
	x.debugf("Preparing to send prompt to model: %s", modelName)
	x.debugf("Prompt length: %d characters", len(prompt))

	if x.apiKey == "" {
		return "", fmt.Errorf("X.AI provider not configured: missing API key")
	}

	if !x.SupportsModel(modelName) {
		return "", fmt.Errorf("invalid X.AI model: %s", modelName)
	}

	x.debugf("Model validation passed, preparing API call")
	x.debugf("Using configuration: Temperature=%.2f, MaxTokens=%d, TopP=%.2f",
		x.config.Temperature, x.config.MaxTokens, x.config.TopP)

	config := openai.DefaultConfig(x.apiKey)
	config.BaseURL = "https://api.x.ai/v1"
	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: modelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: float32(x.config.Temperature),
			MaxTokens:   x.config.MaxTokens,
			TopP:        float32(x.config.TopP),
		},
	)

	if err != nil {
		return "", fmt.Errorf("X.AI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from X.AI")
	}

	response := resp.Choices[0].Message.Content
	x.debugf("API call completed, response length: %d characters", len(response))

	return response, nil
}

// ValidateModel checks if the specific X.AI model variant is valid
func (x *XAIProvider) ValidateModel(modelName string) bool {
	return x.SupportsModel(modelName)
}

// SetConfig updates the provider configuration
func (x *XAIProvider) SetConfig(config ModelConfig) {
	x.debugf("Updating provider configuration")
	x.debugf("Old config: Temperature=%.2f, MaxTokens=%d, TopP=%.2f",
		x.config.Temperature, x.config.MaxTokens, x.config.TopP)
	x.config = config
	x.debugf("New config: Temperature=%.2f, MaxTokens=%d, TopP=%.2f",
		x.config.Temperature, x.config.MaxTokens, x.config.TopP)
}

// GetConfig returns the current provider configuration
func (x *XAIProvider) GetConfig() ModelConfig {
	return x.config
}

// SetVerbose enables or disables verbose mode
func (x *XAIProvider) SetVerbose(verbose bool) {
	x.verbose = verbose
}
