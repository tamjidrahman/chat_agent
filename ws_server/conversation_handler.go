package ws_server

import (
	"fmt"
	"sync"

	"github.com/tamjidrahman/chat_agent/chat_agent"
)

type ConversationHandler interface {
	AddMessage(chat_agent.ChatMessage)
	GetConversation() *chat_agent.ChatConversation
	GetBroadcast() *chan chat_agent.ChatMessage
}

type WSConversationHandler struct {
	conversation *chat_agent.ChatConversation
	broadcast    chan chat_agent.ChatMessage
	mutex        sync.Mutex
}

func newWSConversationHandler(conversation *chat_agent.ChatConversation) *WSConversationHandler {
	wsConvHandler := WSConversationHandler{conversation: conversation, broadcast: make(chan chat_agent.ChatMessage)}
	return &wsConvHandler
}

func (wsConvHandler *WSConversationHandler) AddMessage(message chat_agent.ChatMessage) {
	fmt.Println("Adding message to conversation")
	wsConvHandler.mutex.Lock()
	defer wsConvHandler.mutex.Unlock()

	(*wsConvHandler).broadcast <- message
	wsConvHandler.conversation.AddMessage(message)
}

func (wsConvHandler *WSConversationHandler) GetConversation() *chat_agent.ChatConversation {
	wsConvHandler.mutex.Lock()
	defer wsConvHandler.mutex.Unlock()
	return wsConvHandler.conversation
}

func (wsConvHandler *WSConversationHandler) GetBroadcast() *chan chat_agent.ChatMessage {
	return &wsConvHandler.broadcast
}
