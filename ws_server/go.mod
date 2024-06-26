module github.com/tamjidrahman/ws_server

go 1.22.4

require (
	github.com/gorilla/websocket v1.5.3
	github.com/tamjidrahman/chat_agent/chat_agent v0.0.0-00010101000000-000000000000
)

replace github.com/tamjidrahman/chat_agent/chat_agent => ../chat_agent
