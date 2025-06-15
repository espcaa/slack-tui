package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slacktui/config"
	"slacktui/structs"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	wsMutex       sync.Mutex
	wsInitialized bool
	wsConn        *websocket.Conn
)

type ChatHistoryUpdater interface {
	AppendMessages(newMessage structs.Message, threadbroadcast bool)
	GetSelectedChannelID() string
	DeleteMessage(messageId string)
	ModifyMessage(messageId string, newContent string)
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
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse the message
		var event map[string]interface{}
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		// Handle ping messages
		if event["type"] == "ping" {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"pong", "reply_to":`+fmt.Sprintf("%v", event["id"])+`}`))
			if err != nil {
				panic(fmt.Sprintf("Error sending pong: %v", err))
			}
			continue
		}

		// Handle channel messages
		if event["type"] == "message" {
			// check if it's not a deletion
			if event["subtype"] == "message_deleted" {
				if updater.GetSelectedChannelID() == event["channel"] {
					MessageID := fmt.Sprintf("%v", event["client_msg_id"])
					updater.DeleteMessage(MessageID)
				}
			} else if event["subtype"] == "message_changed" {
				if event["subtype"] == "message_changed" {
					if event["message"] != nil {
						message := event["message"].(map[string]interface{})
						if updater.GetSelectedChannelID() == event["channel"] {
							// Access `client_msg_id` and `text` from the `message` object
							clientMsgID := message["client_msg_id"]
							text := message["text"]

							// Modify the message in the updater
							updater.ModifyMessage(fmt.Sprintf("%v", clientMsgID), fmt.Sprintf("%v", text))
						}
					}
				}
			} else if event["subtype"] == nil || event["subtype"] == "" {
				// Normal message
				channel := event["channel"]
				text := event["text"]
				var timestamp int64
				if tsFloat, ok := event["ts"].(float64); ok {
					timestamp = int64(tsFloat)
				} else if tsString, ok := event["ts"].(string); ok {
					if parsedTs, err := strconv.ParseFloat(tsString, 64); err == nil {
						timestamp = int64(parsedTs)
					} else {
						continue
					}
				} else {
					continue
				}
				if updater.GetSelectedChannelID() == channel {
					updater.AppendMessages(structs.Message{
						MessageId: fmt.Sprintf("%v", event["client_msg_id"]),
						SenderId:  fmt.Sprintf("%v", event["user"]),
						Content:   fmt.Sprintf("%v", text),
						Timestamp: timestamp,
					}, false)
				}
			} else if event["subtype"] == "thread_broadcast" {
				// Handle thread broadcast messages
				channel := event["channel"]
				text := event["text"]
				var timestamp int64
				if tsFloat, ok := event["ts"].(float64); ok {
					timestamp = int64(tsFloat)
				} else if tsString, ok := event["ts"].(string); ok {
					if parsedTs, err := strconv.ParseFloat(tsString, 64); err == nil {
						timestamp = int64(parsedTs)
					} else {
						continue
					}
				} else {
					continue
				}
				if updater.GetSelectedChannelID() == channel {
					updater.AppendMessages(structs.Message{
						MessageId:       fmt.Sprintf("%v", event["client_msg_id"]),
						SenderId:        fmt.Sprintf("%v", event["user"]),
						Content:         fmt.Sprintf("%v", text),
						Timestamp:       timestamp,
						ThreadBroadcast: true,
					}, true)
				}
			}
		}
	}
}

func CheckWebSocketConnection() bool {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	if !wsInitialized || wsConn == nil {
		return false
	}

	if err := wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
		wsInitialized = false
		wsConn = nil
		return false
	}

	return true
}
