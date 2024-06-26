package chat_agent

type ChatAgent interface {
	Query(ChatMessage) ChatMessage
	ShouldRespond(string) bool
	AddToMessages(ChatMessage)
}

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type MockChatAgent struct{}

func (m *MockChatAgent) Query(_ ChatMessage) ChatMessage {
	return ChatMessage{Username: "MockChatAgent", Message: "I'm a Mock Chat Agent!"}
}

func (m *MockChatAgent) ShouldRespond(_ string) bool {
	return true
}

func (m *MockChatAgent) AddToMessages(_ ChatMessage) {
	// Do nothing
}
