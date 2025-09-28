package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Provider represents an LLM provider interface
type Provider interface {
	GenerateCommitMessage(ctx context.Context, diff string, useEmoji bool) (string, error)
	GeneratePRContent(ctx context.Context, commits string, currentBranch string, defaultBranch string, useEmoji bool, prTemplate string) (title string, body string, err error)
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

// ClaudeProvider implements the Provider interface for Anthropic Claude
type ClaudeProvider struct {
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

// NewClaudeProvider creates a new Claude provider instance
func NewClaudeProvider(apiKey string) *ClaudeProvider {
	return &ClaudeProvider{
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

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "Claude"
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

// Claude API structures
type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []claudeContent `json:"content"`
	Error   *claudeError    `json:"error,omitempty"`
}

type claudeContent struct {
	Text string `json:"text"`
}

type claudeError struct {
	Type    string `json:"type"`
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

// GeneratePRContent generates PR title and body using OpenAI
func (p *OpenAIProvider) GeneratePRContent(ctx context.Context, commits string, currentBranch string, defaultBranch string, useEmoji bool, prTemplate string) (string, string, error) {
	emojiInstruction := ""
	if useEmoji {
		emojiInstruction = "\n- You may add appropriate emojis to make the PR more engaging if it fits naturally"
	}

	templateInstruction := ""
	if prTemplate != "" {
		templateInstruction = fmt.Sprintf(`

IMPORTANT: This repository has a pull request template that you MUST follow. Please structure your response to match this template as closely as possible:

--- PR TEMPLATE START ---
%s
--- PR TEMPLATE END ---

When generating the PR body, use the template structure above but fill it with content based on the commit analysis. Maintain the same sections and format from the template.`, prTemplate)
	}

	prompt := fmt.Sprintf(`Analyze the following git commits and generate a comprehensive pull request title and body.

The pull request merges branch '%s' into '%s'.%s

Requirements:
- Generate a clear, concise PR title that summarizes the main purpose of the changes
- Create a detailed PR body with the following sections (unless overridden by template above):
  - ## Summary: Brief overview of what this PR accomplishes
  - ## Changes Made: Bullet points of key changes and improvements
  - ## Testing: Description of testing performed or needed
  - ## Additional Notes: Any important information for reviewers%s

Commits to analyze:
%s

Return the response in this exact format:
TITLE: [your generated title here]

BODY:
[your generated body here]`, currentBranch, defaultBranch, templateInstruction, emojiInstruction, commits)

	reqBody := openAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", "", fmt.Errorf("no response from OpenAI")
	}

	content := openAIResp.Choices[0].Message.Content
	return parsePRResponse(content)
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

// GeneratePRContent generates PR title and body using Gemini
func (p *GeminiProvider) GeneratePRContent(ctx context.Context, commits string, currentBranch string, defaultBranch string, useEmoji bool, prTemplate string) (string, string, error) {
	emojiInstruction := ""
	if useEmoji {
		emojiInstruction = "\n- You may add appropriate emojis to make the PR more engaging if it fits naturally"
	}

	templateInstruction := ""
	if prTemplate != "" {
		templateInstruction = fmt.Sprintf(`

IMPORTANT: This repository has a pull request template that you MUST follow. Please structure your response to match this template as closely as possible:

--- PR TEMPLATE START ---
%s
--- PR TEMPLATE END ---

When generating the PR body, use the template structure above but fill it with content based on the commit analysis. Maintain the same sections and format from the template.`, prTemplate)
	}

	prompt := fmt.Sprintf(`Analyze the following git commits and generate a comprehensive pull request title and body.

The pull request merges branch '%s' into '%s'.%s

Requirements:
- Generate a clear, concise PR title that summarizes the main purpose of the changes
- Create a detailed PR body with the following sections (unless overridden by template above):
  - ## Summary: Brief overview of what this PR accomplishes
  - ## Changes Made: Bullet points of key changes and improvements
  - ## Testing: Description of testing performed or needed
  - ## Additional Notes: Any important information for reviewers%s

Commits to analyze:
%s

Return the response in this exact format:
TITLE: [your generated title here]

BODY:
[your generated body here]`, currentBranch, defaultBranch, templateInstruction, emojiInstruction, commits)

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
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", p.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", "", fmt.Errorf("Gemini API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", "", fmt.Errorf("no response from Gemini")
	}

	content := geminiResp.Candidates[0].Content.Parts[0].Text
	return parsePRResponse(content)
}

// GenerateCommitMessage generates a commit message using Claude
func (p *ClaudeProvider) GenerateCommitMessage(ctx context.Context, diff string, useEmoji bool) (string, error) {
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

	reqBody := claudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 1024,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

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

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if claudeResp.Error != nil {
		return "", fmt.Errorf("Claude API error: %s", claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("no response from Claude")
	}

	return claudeResp.Content[0].Text, nil
}

// GeneratePRContent generates PR title and body using Claude
func (p *ClaudeProvider) GeneratePRContent(ctx context.Context, commits string, currentBranch string, defaultBranch string, useEmoji bool, prTemplate string) (string, string, error) {
	emojiInstruction := ""
	if useEmoji {
		emojiInstruction = "\n- You may add appropriate emojis to make the PR more engaging if it fits naturally"
	}

	templateInstruction := ""
	if prTemplate != "" {
		templateInstruction = fmt.Sprintf(`

IMPORTANT: This repository has a pull request template that you MUST follow. Please structure your response to match this template as closely as possible:

--- PR TEMPLATE START ---
%s
--- PR TEMPLATE END ---

When generating the PR body, use the template structure above but fill it with content based on the commit analysis. Maintain the same sections and format from the template.`, prTemplate)
	}

	prompt := fmt.Sprintf(`Analyze the following git commits and generate a comprehensive pull request title and body.

The pull request merges branch '%s' into '%s'.%s

Requirements:
- Generate a clear, concise PR title that summarizes the main purpose of the changes
- Create a detailed PR body with the following sections (unless overridden by template above):
  - ## Summary: Brief overview of what this PR accomplishes
  - ## Changes Made: Bullet points of key changes and improvements
  - ## Testing: Description of testing performed or needed
  - ## Additional Notes: Any important information for reviewers%s

Commits to analyze:
%s

Return the response in this exact format:
TITLE: [your generated title here]

BODY:
[your generated body here]`, currentBranch, defaultBranch, templateInstruction, emojiInstruction, commits)

	reqBody := claudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 2048,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if claudeResp.Error != nil {
		return "", "", fmt.Errorf("Claude API error: %s", claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return "", "", fmt.Errorf("no response from Claude")
	}

	content := claudeResp.Content[0].Text
	return parsePRResponse(content)
}

// parsePRResponse parses the LLM response to extract title and body
func parsePRResponse(content string) (string, string, error) {
	lines := strings.Split(content, "\n")
	var title, body string
	var inBody bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TITLE:") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "TITLE:"))
		} else if strings.HasPrefix(line, "BODY:") {
			inBody = true
		} else if inBody {
			if body != "" {
				body += "\n"
			}
			body += line
		}
	}

	if title == "" {
		return "", "", fmt.Errorf("no title found in LLM response")
	}

	if body == "" {
		return "", "", fmt.Errorf("no body found in LLM response")
	}

	// Clean up the body
	body = strings.TrimSpace(body)

	return title, body, nil
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

// GeneratePRContent tries providers in order to generate PR title and body
func (pm *ProviderManager) GeneratePRContent(commits string, currentBranch string, defaultBranch string, useEmoji bool, prTemplate string) (string, string, string, error) {
	for i, provider := range pm.providers {
		ctx, cancel := context.WithTimeout(context.Background(), pm.delayThreshold)
		defer cancel()

		title, body, err := provider.GeneratePRContent(ctx, commits, currentBranch, defaultBranch, useEmoji, prTemplate)
		if err == nil {
			return title, body, provider.Name(), nil
		}

		// If this was the last provider or if it's not a timeout error, return the error
		if i == len(pm.providers)-1 {
			return "", "", provider.Name(), fmt.Errorf("all providers failed, last error from %s: %w", provider.Name(), err)
		}

		// Check if the error is due to context timeout
		if ctx.Err() == context.DeadlineExceeded {
			// Continue to next provider
			continue
		}

		// For non-timeout errors, still try the next provider but log this one
		continue
	}

	return "", "", "", fmt.Errorf("no providers available")
}
