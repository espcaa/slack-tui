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
	Latest        int64  // Timestamp of the latest message
	Latest_text   string // Text of the latest message
	Mention_count int    // Unread mention count
}

type Channel struct {
	ChannelId     string
	ChannelName   string
	Mention_count int  // Unread mention count
	IsPrivate     bool // Indicates if the channel is private
}

type Message struct {
	MessageId       string
	SenderId        string
	SenderName      string
	Content         string
	Timestamp       int64 // Unix timestamp
	ThreadBroadcast bool  // Indicates if the message is a thread broadcast
}
