package ws_server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tamjidrahman/chat_agent/chat_agent"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan chat_agent.Message)

func Serve(chatAgent chat_agent.ChatAgent) {
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { handleConnections(chatAgent, w, r) })

	go handleMessages()

	fmt.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

// serveHome serves the HTML page
func serveHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := ioutil.ReadFile("./template.html")
	if err != nil {
		http.Error(w, "Could not read template file"+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.New("home").Parse(string(tmpl))
	if err != nil {
		http.Error(w, "Could not parse template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, struct {
		WebSocketURL string
	}{
		WebSocketURL: "ws://" + r.Host + "/ws",
	})
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}

func handleConnections(chatAgent chat_agent.ChatAgent, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg chat_agent.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		broadcast <- msg
		response := (chatAgent).Query(msg.Message)
		fmt.Println("Response: ", response)
		broadcast <- chat_agent.Message{Username: "Chat Agent", Message: response}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
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
