package chat_agent

type ChatAgent interface {
	Query(*ChatConversation) ChatMessage
	ShouldRespond(*ChatConversation) bool
}

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type ChatConversation struct {
	ChatMessages []ChatMessage `json:"messages"`
}

func (c *ChatConversation) AddMessage(message ChatMessage) {
	c.ChatMessages = append(c.ChatMessages, message)
}

type MockChatAgent struct{}

func (m *MockChatAgent) Query(_ *ChatConversation) ChatMessage {
	return ChatMessage{Username: "MockChatAgent", Message: "I'm a Mock Chat Agent!"}
}

func (m *MockChatAgent) ShouldRespond(_ *ChatConversation) bool {
	return true
}
