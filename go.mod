module github.com/tamjidrahman/chat_agent

go 1.22.4

replace github.com/tamjidrahman/ws_server => ./ws_server

require (
	github.com/tamjidrahman/chat_agent/chat_agent v0.0.0-00010101000000-000000000000
	github.com/tamjidrahman/ws_server v0.0.0-00010101000000-000000000000
)

require github.com/gorilla/websocket v1.5.3 // indirect

replace github.com/tamjidrahman/chat_agent/chat_agent => ./chat_agent
