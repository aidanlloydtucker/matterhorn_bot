package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"reflect"
)

type ChatInfo struct {
	ChatSettings
	Name string `redis:"name"`
	Type string `redis:"type"`
}

type RedisChatInfo struct {
	Settings string `redis:"settings"`
	Name     string `redis:"name"`
	Type     string `redis:"type"`
}

func NewRedisChatInfo() *RedisChatInfo {
	csJson, _ := json.Marshal(ChatSettings{
		NSFW: false,
	})

	return &RedisChatInfo{
		Settings: string(csJson),
	}
}

func ToRedisChatInfo(chat *ChatInfo) (error, *RedisChatInfo) {
	rChat := RedisChatInfo{}

	sJson, err := json.Marshal(chat.ChatSettings)
	if err != nil {
		return err, &RedisChatInfo{}
	}

	rChat.Settings = string(sJson)

	mutable := reflect.ValueOf(&rChat).Elem()

	v := reflect.ValueOf(chat).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.Struct {
			switch v.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				mutable.FieldByName(v.Type().Field(i).Name).SetInt(v.Field(i).Int())
			case reflect.Bool:
				mutable.FieldByName(v.Type().Field(i).Name).SetBool(v.Field(i).Bool())
			case reflect.String:
				mutable.FieldByName(v.Type().Field(i).Name).SetString(v.Field(i).String())
			}
		}
	}
	return nil, &rChat

}

func FromRedisChatInfo(value []interface{}) (error, *ChatInfo) {
	var chat RedisChatInfo

	if err := redis.ScanStruct(value, &chat); err != nil {
		return err, &ChatInfo{}
	}

	tChat := ChatInfo{}

	mutable := reflect.ValueOf(&tChat).Elem()

	v := reflect.ValueOf(chat)

	var csJson ChatSettings
	err := json.Unmarshal([]byte(chat.Settings), &csJson)
	if err != nil {
		return err, &ChatInfo{}
	}
	tChat.ChatSettings = csJson

	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name != "Settings" {
			switch v.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				mutable.FieldByName(v.Type().Field(i).Name).SetInt(v.Field(i).Int())
			case reflect.Bool:
				mutable.FieldByName(v.Type().Field(i).Name).SetBool(v.Field(i).Bool())
			case reflect.String:
				mutable.FieldByName(v.Type().Field(i).Name).SetString(v.Field(i).String())
			}
		}
	}

	return nil, &tChat
}

// TODO: Fix the sketchy way that we store Structs
// Have to do it because of redis
type ChatSettings struct {
	NSFW       bool
	AlertTimes []AlertTime `type:"json"`
	KeyWords   []KeyWord   `type:"json"`
}

type AlertTime struct {
	Time    int64  `json:"time"`
	Message string `json:"msg"`
}

type KeyWord struct {
	Key     string `json:"key"`
	Message string `json:"msg"`
}
