package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"slacktui/config"
	"slacktui/structs"
	"sort"
	"time"
)

type Channel struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	IsMpim    bool     `json:"is_mpim"`
	IsPrivate bool     `json:"is_private"`
	Members   []string `json:"members"`
}

type DM struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	IsIm    bool   `json:"is_im"`
	IsOpen  bool   `json:"is_open"`
	Updated int64  `json:"updated"`
}

type ApiResponse struct {
	Channels []Channel `json:"channels"`
	Dms      []DM      `json:"ims"`
}

func GetUserData() ([]structs.Channel, []structs.DMChannel, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://slack.com/api/client.userBoot", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+config.SlackToken)
	req.Header.Set("Cookie", config.Cookies)

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, nil, err
	}

	// Convert slices directly
	sort.Slice(apiResp.Channels, func(i, j int) bool {
		return apiResp.Channels[i].Name < apiResp.Channels[j].Name
	})

	// Convert slices directly
	structChannels := make([]structs.Channel, 0, len(apiResp.Channels))
	for _, ch := range apiResp.Channels {
		structChannels = append(structChannels, structs.Channel{
			ChannelId:     ch.ID,
			ChannelName:   ch.Name,
			Mention_count: 0,
		})
	}

	structDMs := make([]structs.DMChannel, 0, len(apiResp.Dms))
	for _, dm := range apiResp.Dms {
		// get the username from sqlite
		var username, err = GetNameFromID(dm.User, true)
		if err != nil {
			username = "error" // Fallback to User ID if username retrieval fails
		}

		structDMs = append(structDMs, structs.DMChannel{
			DmUserName:    username,
			DmID:          dm.ID,
			DmUserID:      dm.User,
			Latest:        dm.Updated, // Default value; update based on your logic
			Latest_text:   "",         // Default value; update based on your logic
			Mention_count: 0,          // Default value; update based on your logic
		})
	}

	// Sort DMs by Latest in descending order
	sort.Slice(structDMs, func(i, j int) bool {
		return i > j
	})

	// Interpret as structs.Channel and structs.DMChannel
	// This block is redundant and has been removed.

	return structChannels, structDMs, nil
}
