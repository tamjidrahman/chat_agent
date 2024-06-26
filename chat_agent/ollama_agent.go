package chat_agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Response represents the structure of the JSON response
type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

// generateRequest sends a POST request and returns the response as a Response struct
func generateRequest(url string, payload map[string]interface{}) (Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Response{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("error reading response body: %w", err)
	}

	var responseStruct Response
	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return Response{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return responseStruct, nil
}

type OllamaAgent struct {
	Url   string
	Model string
}

func (o *OllamaAgent) createPromptFromConversation(conversation *ChatConversation) string {
	prompt := "You're an assistant. Do your best to help Tamjid and Laura with their logistics!"
	for _, message := range conversation.ChatMessages {

		prompt_addend := fmt.Sprintf("%s: %s", message.Username, message.Message)
		prompt += prompt_addend + "\n"
	}
	return prompt
}
func (o *OllamaAgent) Query(chatConversation *ChatConversation) ChatMessage {
	payload := map[string]interface{}{
		"model":  o.Model,
		"prompt": o.createPromptFromConversation(chatConversation),
		"stream": false,
	}
	response, err := generateRequest(o.Url, payload)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ChatMessage{Username: "OllamaLarry", Message: "I'm sorry, I'm having trouble understanding you right now."}
	}
	return ChatMessage{Username: "OllamaLarry", Message: response.Response}

}

func (o *OllamaAgent) ShouldRespond(chatConversation *ChatConversation) bool {
	lastMessage := chatConversation.ChatMessages[len(chatConversation.ChatMessages)-1]
	query := strings.TrimSpace(lastMessage.Message)
	return strings.HasPrefix(query, "@larry")
}

func (o *OllamaAgent) AddToMessages(message ChatMessage) {
	// Do nothing
}
