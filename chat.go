package main

type ChatInfo struct {
	Name string `redis:"name"`
	NSFW bool   `redis:"nsfw"`
	Type string `redis:"type"`
}

func NewChatInfo() *ChatInfo {
	return &ChatInfo{
		NSFW: false,
	}
}

type ChatSettings struct {
	NSFW bool `redis:"nsfw"`
}
