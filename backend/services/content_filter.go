package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type DeepSeekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ContentFilterService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewContentFilterService() *ContentFilterService {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	return &ContentFilterService{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/v1/chat/completions",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *ContentFilterService) CheckContent(title, content string) (bool, error) {
	if c.apiKey == "" {
		return false, fmt.Errorf("DEEPSEEK_API_KEY environment variable is not set")
	}

	prompt := fmt.Sprintf(`You are a content moderator. Analyze the following text for inappropriate content including profanity, hate speech, explicit content, or offensive language.

Title: %s
Content: %s

Respond with only "CLEAN" if the content is appropriate, or "INAPPROPRIATE" if it contains any offensive language, swear words, or inappropriate content. Do not provide explanations.`, title, content)

	request := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response DeepSeekResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		return false, fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return false, fmt.Errorf("no response from AI")
	}

	result := response.Choices[0].Message.Content
	isClean := result == "CLEAN"

	// Return true if content is clean, false if inappropriate
	return isClean, nil
}
