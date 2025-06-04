package utils

import (
	"encoding/json"

	"net/http"
	"slacktui/config"
	"slacktui/structs"
	"time"
)

func GetChannelList() []structs.SidebarItem {
	var cfg, err = config.LoadConfig()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	token := cfg.SlackToken
	cookies := cfg.Cookies

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://slack.com/api/users.conversations", nil)
	if err != nil {
		panic("Error creating request: " + err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Cookie", cookies)

	resp, err := client.Do(req)
	if err != nil {
		panic("Error making request: " + err.Error())
	}
	defer resp.Body.Close()

	var raw struct {
		Channels []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			// Infer type based on flags
			IsIm      bool `json:"is_im"`
			IsMpim    bool `json:"is_mpim"`
			IsChannel bool `json:"is_channel"`
			IsPrivate bool `json:"is_private"`
		} `json:"channels"`
	}

	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		panic("Error decoding JSON: " + err.Error())
	}

	var items []structs.SidebarItem
	for _, ch := range raw.Channels {
		var t string
		switch {
		case ch.IsIm:
			t = "dm"
		case ch.IsMpim:
			t = "group_dm"
		case ch.IsChannel && ch.IsPrivate:
			t = "private_channel"
		case ch.IsChannel:
			t = "channel"
		default:
			t = "unknown"
		}
		items = append(items, structs.SidebarItem{
			Id:   ch.Id,
			Name: ch.Name,
			Type: t,
		})
	}

	return items
}
