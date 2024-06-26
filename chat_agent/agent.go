package chat_agent

type ChatAgent interface {
	Query(string) string
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type MockChatAgent struct{}

func (m MockChatAgent) Query(_ string) string {
	return "I'm a Mock Chat Agent!"
}
