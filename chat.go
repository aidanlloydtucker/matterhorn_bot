package main

import (
	"bytes"
	"encoding/gob"
	"regexp"
	"strings"
)

type ChatInfo struct {
	Settings ChatSettings `json:"settings"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
}

func NewChatInfo() *ChatInfo {
	return &ChatInfo{}
}

func DecodeRedisChatInfo(chatInfoGob string) (ChatInfo, error) {
	decoder := gob.NewDecoder(strings.NewReader(chatInfoGob))
	ci := ChatInfo{}
	err := decoder.Decode(&ci)
	return ci, err
}

func EncodeRedisChatInfo(ci ChatInfo) (string, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(&ci)
	return buf.String(), err
}

type ChatSettings struct {
	NSFW       bool        `json:"nsfw"`
	AlertTimes []AlertTime `json:"alert_times"`
	KeyWords   []KeyWord   `json:"key_words"`
}

type AlertTime struct {
	Time    string `json:"time"`
	Message string `json:"msg"`
}

type KeyWord struct {
	Key     string `json:"key"`
	Message string `json:"msg"`
}

var timeRegex = regexp.MustCompile(`^([01]\d|2[0-3]):?([0-5]\d)$`)
