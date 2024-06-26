package ws_server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/tamjidrahman/chat_agent/chat_agent"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func Serve(chatAgent chat_agent.ChatAgent) {
	conversation := chat_agent.ChatConversation{ChatMessages: []chat_agent.ChatMessage{}}
	conversationHandler := newWSConversationHandler(&conversation)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { serveHome(conversationHandler, w, r) })

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { handleConnections(conversationHandler, chatAgent, w, r) })

	go handleMessages(conversationHandler.GetBroadcast())

	fmt.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

// serveHome serves the HTML page
func serveHome(conversationHandler ConversationHandler, w http.ResponseWriter, r *http.Request) {
	tmpl, err := os.ReadFile("./template.html")
	if err != nil {
		http.Error(w, "Could not read template file"+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.New("home").Parse(string(tmpl))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not parse template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, struct {
		WebSocketURL string
		ChatMessages []chat_agent.ChatMessage
	}{
		WebSocketURL: "ws://" + r.Host + "/ws",
		ChatMessages: conversationHandler.GetConversation().ChatMessages,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}

func handleConnections(conversationHandler ConversationHandler, chatAgent chat_agent.ChatAgent, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg chat_agent.ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		conversationHandler.AddMessage(msg)
		if (chatAgent).ShouldRespond(conversationHandler.GetConversation()) {
			response := (chatAgent).Query(conversationHandler.GetConversation())
			conversationHandler.AddMessage(response)
		}
	}
}

func handleMessages(broadcast *chan chat_agent.ChatMessage) {
	for {
		msg := <-*broadcast
		fmt.Println("Broadcasting message: ", msg)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
