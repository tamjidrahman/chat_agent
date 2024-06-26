package main

import (
	// "os"

	"github.com/tamjidrahman/chat_agent/chat_agent"
	"github.com/tamjidrahman/ws_server"
)

func main() {
	// agent := chat_agent.MockChatAgent{}
	agent := chat_agent.OllamaAgent{Url: "http://localhost:11434/api/generate", Model: "mistral"}
	// OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	// agent := chat_agent.NewOpenAIClient(OPENAI_API_KEY)
	ws_server.Serve(&agent)
}
