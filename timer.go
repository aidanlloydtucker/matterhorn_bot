package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

const INTERVAL_PERIOD time.Duration = 24 * time.Hour
const SECONDS_TO_TICK int = 10

func runningRoutine() {
	ticker := updateTicker()
	for {
		<-ticker.C
		fmt.Println(time.Now(), "- just ticked")
		ticker = updateTicker()
	}
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
		chat.AlertTimes

	}

}

func updateTicker(hour int, minutes int) *time.Ticker {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minutes, SECONDS_TO_TICK, 0, time.UTC)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	diff := nextTick.Sub(time.Now())
	return time.NewTicker(diff)
}
