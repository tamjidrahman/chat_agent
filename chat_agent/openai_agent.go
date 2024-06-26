package chat_agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Message represents a message in the chat.
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request body for the OpenAI API.
type ChatRequest struct {
	Model          string          `json:"model"`
	OpenAIMessages []OpenAIMessage `json:"messages"`
}

// ChatResponse represents the response from the OpenAI API.
type ChatResponse struct {
	Choices []struct {
		OpenAIMessage OpenAIMessage `json:"message"`
	} `json:"choices"`
}

// OpenAIClient is a client for interacting with the OpenAI API.
type OpenAIClient struct {
	APIKey         string
	OpenAIMessages []OpenAIMessage
}

// NewOpenAIClient creates a new OpenAI client.
func NewOpenAIClient(apiKey string) OpenAIClient {
	prompt_string := "You are an AI assistant named Larry that helps Tamjid and Laura with their household logistics"
	messages := []OpenAIMessage{{Role: "user", Content: prompt_string}}

	return OpenAIClient{
		APIKey:         apiKey,
		OpenAIMessages: messages,
	}
}

func (client *OpenAIClient) ShouldRespond(query string) bool {
	query = strings.TrimSpace(query)
	return strings.HasPrefix(query, "@larry")
}

func (client *OpenAIClient) AddToMessages(message ChatMessage) {
	client.OpenAIMessages = append(client.OpenAIMessages, OpenAIMessage{Role: "user", Content: message.Username + ": " + message.Message})
}

func (client *OpenAIClient) CreateMessage(message string) ChatMessage {
	client.OpenAIMessages = append(client.OpenAIMessages, OpenAIMessage{Role: "assistant", Content: message})
	return ChatMessage{Username: "Larry", Message: message}
}

// Query sends a chat completion request to the OpenAI API and returns the response.
func (client *OpenAIClient) Query(message ChatMessage) ChatMessage {
	url := "https://api.openai.com/v1/chat/completions"

	client.OpenAIMessages = append(client.OpenAIMessages, OpenAIMessage{Role: "user", Content: message.Username + ": " + message.Message})

	chatRequest := ChatRequest{
		Model:          "gpt-4o",
		OpenAIMessages: client.OpenAIMessages,
	}

	requestBody, err := json.Marshal(chatRequest)
	fmt.Println("Request Body: ", requestBody)
	fmt.Println(chatRequest)
	if err != nil {
		fmt.Println("Error marshalling JSON: %w", err)
		return client.CreateMessage("")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request: %w", err)
		return client.CreateMessage("")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.APIKey)

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		fmt.Println("Error sending request: %w", err)
		return client.CreateMessage("")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-200 response code: " + string(resp.StatusCode))
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		return client.CreateMessage("")
	}

	var chatResponse ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		fmt.Println("Error reading response body: %w", err)
		return client.CreateMessage("")
	}

	if len(chatResponse.Choices) == 0 {
		fmt.Println("No choices in response")
		return client.CreateMessage("")
	}

	return client.CreateMessage(chatResponse.Choices[0].OpenAIMessage.Content)
}
