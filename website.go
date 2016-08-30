package main

import (
	"net/http"

	"strconv"

	"errors"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
)

var p = proxy.New(&ace.Options{BaseDir: "views"})

func startWebsite() {
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	/*tpl, err := p.Load("base", "chat", nil)
	if err != nil {
		panic(err)
	}

	r.SetHTMLTemplate(tpl)*/

	r.Handle("GET", "/chat/:id", webChatHandler)
	r.Handle("PUT", "/chat/:id", webChatChangeHandler)
	r.StaticFS("/public/", http.Dir("./static/"))

	r.Run(":" + HttpPort)
}

func webChatHandler(c *gin.Context) {
	chatIdStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rc := redisPool.Get()
	defer rc.Close()

	ciPtr, httpCode, err := getRedisChatInfo(rc, chatID)
	if err != nil || httpCode != http.StatusOK || ciPtr == nil {
		c.AbortWithError(httpCode, err)
		return
	}

	ci := *ciPtr

	tpl, err := p.Load("base", "chat", nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := map[string]interface{}{
		"ChatName": ci.Name,
		"ChatId":   chatID,
		"SettingsBool": map[string]interface{}{
			"NSFW": ci.Settings.NSFW,
		},
		"KeyWords":   ci.Settings.KeyWords,
		"AlertTimes": ci.Settings.AlertTimes,
	}

	c.Header("Content-Type", "text/HTML")

	if err := tpl.Execute(c.Writer, data); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func webChatChangeHandler(c *gin.Context) {
	chatIdStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	rc := redisPool.Get()
	defer rc.Close()

	ciPtr, httpCode, err := getRedisChatInfo(rc, chatID)
	if err != nil || httpCode != http.StatusOK || ciPtr == nil {
		c.AbortWithError(httpCode, err)
		return
	}

	ci := *ciPtr

	var settings ChatSettings

	err = c.BindJSON(&settings)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var newKWs []KeyWord

	for _, kw := range settings.KeyWords {
		if kw.Key == "" || kw.Message == "" {
			continue
		}
		newKWs = append(newKWs, kw)
	}

	var newATs []AlertTime

	for _, at := range settings.AlertTimes {
		if at.Time == "" || at.Message == "" {
			continue
		}
		newATs = append(newATs, at)
	}

	insertTimersByChatID(newATs, chatID)

	ci.Settings.KeyWords = newKWs
	ci.Settings.AlertTimes = newATs
	ci.Settings.NSFW = settings.NSFW

	err = setRedisChatInfo(rc, ci, chatID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func getRedisChatInfo(rc redis.Conn, chatID int64) (*ChatInfo, int, error) {
	exists, err := redis.Bool(rc.Do("EXISTS", formatRedisKey(chatID)))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !exists {
		return nil, http.StatusNotFound, errors.New("ChatId is not in database")
	}

	v, err := redis.String(rc.Do("GET", formatRedisKey(chatID)))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	ci, err := DecodeRedisChatInfo(v)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &ci, http.StatusOK, nil
}

func setRedisChatInfo(rc redis.Conn, ci ChatInfo, chatID int64) error {
	ciStr, err := EncodeRedisChatInfo(ci)
	if err != nil {
		return err
	}
	_, err = rc.Do("SET", redis.Args{}.Add(formatRedisKey(chatID)).Add(ciStr)...)
	return err
}
