package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/telegram-bot-api.v4"
)

const INTERVAL_PERIOD time.Duration = 24 * time.Hour
const SECONDS_TO_TICK int = 10

const MAX_ALERTS_ALLOWED int = 3

func startReminder(hour int, minutes int, message string, chatid int64) {
	ticker := updateTicker(hour, minutes)
	for {
		<-ticker.C
		msg := tgbotapi.NewMessage(chatid, message)
		mainBot.Send(msg)
		ticker = updateTicker(hour, minutes)
	}
}

func parseTimes(timeStr string) (int, int, error) {
	tm, err := time.ParseInLocation(time.Kitchen, timeStr, time.UTC)
	if err != nil {
		return 0, 0, err
	}

	return tm.Hour(), tm.Minute(), nil
}

func loadTimeReminders() {
	rc := redisPool.Get()
	defer rc.Close()

	chats, err := redis.Strings(rc.Do("KEYS", REDIS_KEY_PREFIX+"*"))
	if err != nil {
		return
	}
	for _, chatKey := range chats {
		v, err := redis.String(rc.Do("GET", chatKey))
		if err != nil {
			continue
		}

		chat, err := DecodeRedisChatInfo(v)
		if err != nil {
			continue
		}

		chatIdStr := strings.TrimPrefix(chatKey, REDIS_KEY_PREFIX)
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			continue
		}

		for i, at := range chat.Settings.AlertTimes {
			if i >= MAX_ALERTS_ALLOWED {
				break
			}
			hour, min, err := parseTimes(at.Time)
			if err != nil {
				continue
			}
			go startReminder(hour, min, at.Message, chatId)
		}
	}

}

func updateTicker(hour int, minutes int) *time.Ticker {
	utcTimeNow := time.Now().UTC()
	nextTick := time.Date(utcTimeNow.Year(), utcTimeNow.Month(), utcTimeNow.Day(), hour, minutes, SECONDS_TO_TICK, 0, time.UTC)
	if nextTick.Before(utcTimeNow) || nextTick.Equal(utcTimeNow) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	diff := nextTick.Sub(utcTimeNow)
	return time.NewTicker(diff)
}
