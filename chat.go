package main

type ChatInfo struct {
	ChatSettings
	Name string `redis:"name"`
	Type string `redis:"type"`
}

func NewChatInfo() *ChatInfo {
	return &ChatInfo{
		ChatSettings: ChatSettings{
			NSFW: false,
		},
	}
}

type ChatSettings struct {
	NSFW       bool        `redis:"nsfw" web:"bool"`
	AlertTimes []AlertTime `redis:"alert_times" web:"list_string"`
	KeyWords   []KeyWord   `redis:"key_words" web:"list_string"`
}

type AlertTime struct {
	Time    int64
	Message string
}

type KeyWord struct {
	Key     string
	Message string
}
