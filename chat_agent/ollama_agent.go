package chat_agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	body, err := ioutil.ReadAll(resp.Body)
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
	Url string
}

func (o OllamaAgent) Query(query string) string {
	payload := map[string]interface{}{
		"model":  "llama3",
		"prompt": query,
		"stream": false,
	}
	response, err := generateRequest(o.Url, payload)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return response.Response

}
