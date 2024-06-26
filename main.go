package main

import (
	"github.com/tamjidrahman/chat_agent/chat_agent"
	"github.com/tamjidrahman/ws_server"
)

func main() {
	// mockAgent := chat_agent.MockChatAgent{}
	ollama_agent := chat_agent.OllamaAgent{Url: "http://localhost:11434/api/generate"}
	ws_server.Serve(ollama_agent)
}
