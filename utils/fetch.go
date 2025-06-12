package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"slacktui/config"
	"slacktui/structs"
	"slices"
	"strconv"
	"strings"
)

func FetchChannelData(channelid string, oldest int, useoldest bool) []structs.Message {
	// Fetch messages for the given channel
	var messages []structs.Message

	params := url.Values{}
	params.Set("channel", channelid)
	params.Set("limit", "100")
	if useoldest {
		params.Set("latest", strconv.FormatInt(int64(oldest), 10))
		params.Set("inclusive", "false")
	}

	var req, err = http.NewRequest("GET", "https://slack.com/api/conversations.history?"+params.Encode(), nil)
	if err != nil {
		return messages
	}

	var cfg, error = config.LoadConfig()
	if error != nil {
		return messages
	}

	req.Header.Set("Authorization", "Bearer "+cfg.SlackToken)
	req.Header.Set("Cookie", cfg.Cookies)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return messages
	}

	defer resp.Body.Close()

	var result struct {
		Ok       bool `json:"ok"`
		Messages []struct {
			Type    string `json:"type"`
			User    string `json:"user"`
			Text    string `json:"text"`
			Ts      string `json:"ts"`
			BotId   string `json:"bot_id,omitempty"`
			Subtype string `json:"subtype,omitempty"`
		} `json:"messages"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return messages
	}

	if !result.Ok {
		return messages
	}

	for _, m := range result.Messages {
		// Skip messages without a user (e.g., join/leave notices)
		if m.User == "" {
			continue
		}
		parts := strings.Split(m.Ts, ".")
		sec, _ := strconv.ParseInt(parts[0], 10, 64)

		var senderName, err = GetNameFromID(m.User, true)
		if err != nil {
			senderName = m.User // Fallback to user ID if name retrieval fails
		}

		messages = append(messages, structs.Message{
			MessageId:  m.Ts,
			SenderId:   m.User,
			SenderName: senderName,
			Content:    m.Text,
			Timestamp:  sec,
		})
	}

	// Shift the message in the inverse order
	slices.Reverse(messages)

	return messages
}
