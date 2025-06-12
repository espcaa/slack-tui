package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slacktui/config"
	"slacktui/structs"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	wsMutex       sync.Mutex
	wsInitialized bool
	wsConn        *websocket.Conn
)

type ChatHistoryUpdater interface {
	AppendMessages(newMessage structs.Message)
	GetSelectedChannelID() string
}

// Connect to the Slack WebSocket API

func ConnectWebSocket() (*websocket.Conn, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	headers := http.Header{}
	headers.Add("Cookie", cfg.Cookies)

	url := "wss://wss-primary.slack.com/?token=" + cfg.SlackToken
	conn, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	return conn, nil
}

func InitializeWebSocket(updater ChatHistoryUpdater) error {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	if wsInitialized {
		return nil
	}

	conn, err := ConnectWebSocket()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to WebSocket: %v", err))
	}

	wsConn = conn
	wsInitialized = true
	go ListenWebSocket(wsConn, updater) // Pass the appropriate ChatHistoryUpdater instance
	return nil
}

func ListenWebSocket(conn *websocket.Conn, updater ChatHistoryUpdater) {
	defer conn.Close()

	for {
		// Read message from WebSocket
		_, message, err := conn.ReadMessage()
		if err != nil {
			panic(fmt.Sprintf("Error reading message: %v", err))
			break
		}

		// Parse the message
		var event map[string]interface{}
		if err := json.Unmarshal(message, &event); err != nil {
			panic(fmt.Sprintf("Error parsing message: %v", err))
			continue
		}

		// Handle ping messages
		if event["type"] == "ping" {
			panic("Received ping, sending pong...")
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"pong", "reply_to":`+fmt.Sprintf("%v", event["id"])+`}`))
			if err != nil {
				panic(fmt.Sprintf("Error sending pong: %v", err))
			}
			continue
		}

		// Handle channel messages
		if event["type"] == "message" {
			channel := event["channel"]
			text := event["text"]
			var timestamp int64
			if tsFloat, ok := event["ts"].(float64); ok {
				timestamp = int64(tsFloat)
			} else if tsString, ok := event["ts"].(string); ok {
				if parsedTs, err := strconv.ParseFloat(tsString, 64); err == nil {
					timestamp = int64(parsedTs)
				} else {
					panic(fmt.Sprintf("Unexpected timestamp format: %v", event["ts"]))
					continue
				}
			} else {
				panic(fmt.Sprintf("Unexpected timestamp format: %v", event["ts"]))
				continue
			}
			if updater.GetSelectedChannelID() == channel {
				updater.AppendMessages(structs.Message{
					MessageId: fmt.Sprintf("%v", event["ts"]),
					SenderId:  fmt.Sprintf("%v", event["user"]),
					Content:   fmt.Sprintf("%v", text),
					Timestamp: timestamp,
				})
			}
		}
	}
}
