package main

import (
	"time"

	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/telegram-bot-api.v4"
)

const INTERVAL_PERIOD time.Duration = 24 * time.Hour
const SECONDS_TO_TICK int = 10

var TimersByChatID = make(map[int64][]*time.Timer)

func startTimer(when time.Time, message string, chatID int64) *time.Timer {
	timeNow := time.Now().In(when.Location())
	nextTick := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), when.Hour(), when.Minute(), SECONDS_TO_TICK, 0, when.Location())
	if nextTick.UnixNano() < timeNow.UnixNano() || nextTick.Equal(timeNow) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	diff := nextTick.Sub(timeNow)

	timer := time.NewTimer(diff)

	go func() {
		for {
			<-timer.C
			msg := tgbotapi.NewMessage(chatID, message)
			mainBot.Send(msg)

			timer = time.NewTimer(INTERVAL_PERIOD)
		}
	}()
	return timer
}

func insertTimersByChatID(ats []AlertTime, chatID int64) {
	timers, ok := TimersByChatID[chatID]
	if ok {
		for _, timer := range timers {
			timer.Stop()
		}
	}
	newTimers := make([]*time.Timer, 0)
	for _, at := range ats {
		timerTime, err := parseTimes(at.Time)
		if err != nil {
			continue
		}
		newTimers = append(newTimers, startTimer(timerTime, at.Message, chatID))
	}
	TimersByChatID[chatID] = newTimers
}

func parseTimes(timeStr string) (time.Time, error) {
	return time.Parse(`3:04PM -07`, timeStr)
}

func initTimers() {
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

		insertTimersByChatID(chat.Settings.AlertTimes, chatId)
	}
}

/*type AlertTimer struct {
	Message string
	ChatID  int64
}

var TimerMap = make(map[string][]AlertTimer)

func buildTimerKey(hour, minute int) string {
	return fmt.Sprintf("%d:%d", hour, minute)
}

func decodeTimerKey(str string) (int, int, error) {
	strArr := strings.Split(str, ":")
	if len(strArr) != 2 {
		return 0, 0, fmt.Errorf("Length of timer is %d", len(strArr))
	}

	hour, err := strconv.Atoi(strArr[0])
	if err != nil {
		return 0, 0, err
	}

	min, err := strconv.Atoi(strArr[0])
	return hour, min, err
}

func startReminder(hour int, minutes int) {
	ticker := updateTicker(hour, minutes)
	for {
		<-ticker.C
		for _, alert := range TimerMap[buildTimerKey(hour, minutes)] {
			msg := tgbotapi.NewMessage(alert.ChatID, alert.Message)
			mainBot.Send(msg)
		}
		if len(TimerMap[buildTimerKey(hour, minutes)]) == 0 {
			ticker.Stop()
		}
		delete(TimerMap, buildTimerKey(hour, minutes))

		ticker = updateTicker(hour, minutes)
	}
}

func insertReminder(hour, minute int, message string, chatID int64) {
	ats, ok := TimerMap[buildTimerKey(hour, minute)]
	if !ok {
		ats = make([]AlertTimer, 0)
	}
	ats = append(ats, AlertTimer{
		Message: message,
		ChatID:  chatID,
	})
	TimerMap[buildTimerKey(hour, minute)] = ats
	if !ok {
		go startReminder(hour, minute)
	}
}

func parseTimes(timeStr string) (int, int, error) {
	tm, err := time.Parse(`3:04PM MST`, timeStr)
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
			ats, ok := TimerMap[buildTimerKey(hour, min)]
			if !ok {
				ats = make([]AlertTimer, 0)
			}
			ats = append(ats, AlertTimer{
				Message: at.Message,
				ChatID:  chatId,
			})
			TimerMap[buildTimerKey(hour, min)] = ats
		}
	}
	for key := range TimerMap {
		hour, min, err := decodeTimerKey(key)
		if err != nil {
			continue
		}
		go startReminder(hour, min)
	}
}

func updateTicker(hour int, minutes int) *time.Ticker {
	timeNow := time.Now()
	nextTick := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), hour, minutes, SECONDS_TO_TICK, 0, time.UTC)
	if nextTick.Before(timeNow) || nextTick.Equal(timeNow) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	diff := nextTick.Sub(timeNow)
	log.Println("Starting a reminder for", diff)
	return time.NewTicker(diff)
}
*/
