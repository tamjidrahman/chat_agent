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
	APIKey string
	Prompt string
}

// NewOpenAIClient creates a new OpenAI client.
func NewOpenAIClient(apiKey string) OpenAIClient {
	prompt_string := "You are an AI assistant named Larry that helps Tamjid and Laura with their household logistics"

	return OpenAIClient{
		APIKey: apiKey,
		Prompt: prompt_string,
	}
}

func (client *OpenAIClient) ShouldRespond(chatConversation *ChatConversation) bool {

	last_message := chatConversation.ChatMessages[len(chatConversation.ChatMessages)-1].Message
	query := strings.TrimSpace(last_message)
	return strings.HasPrefix(query, "@larry")
}

func (client *OpenAIClient) CreateOpenAIMessages(chatConversation *ChatConversation) []OpenAIMessage {

	var openAIMessages []OpenAIMessage
	openAIMessages = append(openAIMessages, OpenAIMessage{Role: "system", Content: client.Prompt})
	for _, msg := range chatConversation.ChatMessages {
		openAIMessages = append(openAIMessages, OpenAIMessage{Role: "user", Content: msg.Username + ": " + msg.Message})
	}
	return openAIMessages
}

func (client *OpenAIClient) CreateMessage(content string) ChatMessage {
	return ChatMessage{Username: "Larry", Message: content}
}

// Query sends a chat completion request to the OpenAI API and returns the response.
func (client *OpenAIClient) Query(chatConversation *ChatConversation) ChatMessage {
	url := "https://api.openai.com/v1/chat/completions"

	openAIMessages := client.CreateOpenAIMessages(chatConversation)

	chatRequest := ChatRequest{
		Model:          "gpt-4o",
		OpenAIMessages: openAIMessages,
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
