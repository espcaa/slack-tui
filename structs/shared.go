package structs

type SidebarItem struct {
	Type string
	Name string
	Id   string
}

type DMChannel struct {
	DmID          string
	DmUserID      string
	DmUserName    string
	Latest        int    // Timestamp of the latest message
	Latest_text   string // Text of the latest message
	Mention_count int    // Unread mention count
}

type Channel struct {
	ChannelId     string
	ChannelName   string
	Mention_count int // Unread mention count
}
