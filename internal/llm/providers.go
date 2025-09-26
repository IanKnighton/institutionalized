package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Provider represents an LLM provider interface
type Provider interface {
	GenerateCommitMessage(ctx context.Context, diff string, useEmoji bool) (string, error)
	Name() string
}

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey string
}

// GeminiProvider implements the Provider interface for Google Gemini
type GeminiProvider struct {
	apiKey string
}

// NewOpenAIProvider creates a new OpenAI provider instance
func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: apiKey,
	}
}

// NewGeminiProvider creates a new Gemini provider instance
func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
	}
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "OpenAI"
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "Gemini"
}

// OpenAI API structures
type openAIRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []choice  `json:"choices"`
	Error   *apiError `json:"error,omitempty"`
}

type choice struct {
	Message message `json:"message"`
}

type apiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Gemini API structures
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      *geminiError      `json:"error,omitempty"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GenerateCommitMessage generates a commit message using OpenAI
func (p *OpenAIProvider) GenerateCommitMessage(ctx context.Context, diff string, useEmoji bool) (string, error) {
	emojiInstruction := ""
	if useEmoji {
		emojiInstruction = "\n- Add an appropriate emoji at the beginning of the commit type (✨ feat, 🐛 fix, 📚 docs, 💄 style, ♻️ refactor, ✅ test, 🔧 chore, ⚡ perf, 👷 ci, 🏗️ build, ⏪ revert)"
	}

	prompt := fmt.Sprintf(`Analyze the following git diff and generate a conventional commit message. 

The commit message should follow the Conventional Commits specification:
- Start with a type (feat, fix, docs, style, refactor, test, chore, etc.)
- Include a brief description in present tense
- Keep the first line under 50 characters if possible
- Add a body if the change is complex (separate with blank line)%s

Git diff:
%s

Return only the commit message, nothing else.`, emojiInstruction, diff)

	reqBody := openAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// GenerateCommitMessage generates a commit message using Gemini
func (p *GeminiProvider) GenerateCommitMessage(ctx context.Context, diff string, useEmoji bool) (string, error) {
	emojiInstruction := ""
	if useEmoji {
		emojiInstruction = "\n- Add an appropriate emoji at the beginning of the commit type (✨ feat, 🐛 fix, 📚 docs, 💄 style, ♻️ refactor, ✅ test, 🔧 chore, ⚡ perf, 👷 ci, 🏗️ build, ⏪ revert)"
	}

	prompt := fmt.Sprintf(`Analyze the following git diff and generate a conventional commit message. 

The commit message should follow the Conventional Commits specification:
- Start with a type (feat, fix, docs, style, refactor, test, chore, etc.)
- Include a brief description in present tense
- Keep the first line under 50 characters if possible
- Add a body if the change is complex (separate with blank line)%s

Git diff:
%s

Return only the commit message, nothing else.`, emojiInstruction, diff)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Gemini API endpoint
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", p.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("Gemini API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// ProviderManager manages multiple LLM providers with fallback capability
type ProviderManager struct {
	providers      []Provider
	delayThreshold time.Duration
}

// NewProviderManager creates a new provider manager
func NewProviderManager(providers []Provider, delayThreshold time.Duration) *ProviderManager {
	return &ProviderManager{
		providers:      providers,
		delayThreshold: delayThreshold,
	}
}

// GenerateCommitMessage tries providers in order with timeout and fallback
func (pm *ProviderManager) GenerateCommitMessage(diff string, useEmoji bool) (string, string, error) {
	for i, provider := range pm.providers {
		ctx, cancel := context.WithTimeout(context.Background(), pm.delayThreshold)
		defer cancel()

		result, err := provider.GenerateCommitMessage(ctx, diff, useEmoji)
		if err == nil {
			return result, provider.Name(), nil
		}

		// If this was the last provider or if it's not a timeout error, return the error
		if i == len(pm.providers)-1 {
			return "", provider.Name(), fmt.Errorf("all providers failed, last error from %s: %w", provider.Name(), err)
		}

		// Check if the error is due to context timeout
		if ctx.Err() == context.DeadlineExceeded {
			// Continue to next provider
			continue
		}

		// For non-timeout errors, still try the next provider but log this one
		continue
	}

	return "", "", fmt.Errorf("no providers available")
}
