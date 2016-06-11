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
		v, err := redis.Values(redisConn.Do("HGETALL", REDIS_KEY_PREFIX+chatid))
		if err != nil {
			return
		}
		err, chat := FromRedisChatInfo(v)
		if err != nil {
			return
		}
		for _, at := range chat.AlertTimes {
			atHour, atMin, err := parseTimes(at.Time)
			if err != nil {
				continue
			}
			if atHour == hour && atMin == minutes {
				return
			}
		}

		msg := tgbotapi.NewMessage(chatid, message)
		mainBot.Send(msg)
		ticker = updateTicker(hour, minutes)
	}
}

func parseTimes(timeStr string) (int, int, error) {
	splitNums := strings.Split(timeStr, ":")
	hour, err := strconv.Atoi(splitNums[0])
	if err != nil {
		return 0, 0, err
	}
	min, err := strconv.Atoi(splitNums[1])
	if err != nil {
		return 0, 0, err
	}
	return hour, min, nil
}

func loadTimeReminders() {
	vals, err := redis.Strings(redisConn.Do("KEYS", REDIS_KEY_PREFIX+"*"))
	if err != nil {
		return
	}
	for _, key := range vals {
		v, err := redis.Values(redisConn.Do("HGETALL", key))
		if err != nil {
			continue
		}

		err, chat := FromRedisChatInfo(v)
		if err != nil {
			continue
		}

		chatIdStr := strings.TrimPrefix(key, REDIS_KEY_PREFIX)
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			continue
		}

		for i, at := range chat.AlertTimes {
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
