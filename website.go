package main

import (
	"net/http"

	"log"
	"strconv"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
)

var p = proxy.New(&ace.Options{BaseDir: "views"})

func startWebsite(rel bool) {
	log.Println("Starting website")

	if rel {
		gin.SetMode(gin.ReleaseMode)
	}

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
	log.Println("Started Website")
}

func webChatHandler(c *gin.Context) {
	chatIdStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	chat, exists, err := getDatastoreChat(chatID)
	if err != nil {
		if !exists {
			c.AbortWithError(http.StatusNotFound, err)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	tpl, err := p.Load("base", "chat", nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := map[string]interface{}{
		"ChatName": chat.Name,
		"ChatId":   chatID,
		"SettingsBool": map[string]interface{}{
			"NSFW": chat.Settings.NSFW,
		},
		"KeyWords":   chat.Settings.KeyWords,
		"AlertTimes": chat.Settings.AlertTimes,
	}

	c.Header("Content-Type", "text/HTML")

	if err := tpl.Execute(c.Writer, data); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

type ChangedChatSettings struct {
	NSFW          bool         `json:"nsfw"`
	NewAlertTimes []AlertTime  `json:"new_alert_times"`
	AlertTimes    map[int]bool `json:"alert_times"`
	NewKeyWords   []KeyWord    `json:"new_key_words"`
	KeyWords      map[int]bool `json:"key_words"`
}

func webChatChangeHandler(c *gin.Context) {
	chatIdStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	_, exists, err := getDatastoreChat(chatID)
	if err != nil {
		if !exists {
			c.AbortWithError(http.StatusNotFound, err)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	var changeSettings ChangedChatSettings

	err = c.BindJSON(&changeSettings)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	updateFunc := func(oldChat Chat) Chat {
		newAT := []AlertTime{}
		for _, at := range oldChat.Settings.AlertTimes {
			if val, ok := changeSettings.AlertTimes[at.ID]; val || !ok {
				newAT = append(newAT, at)
			}
		}
		for _, at := range changeSettings.NewAlertTimes {
			if at.Time != "" && at.Message != "" {
				newAT = append(newAT, MakeAlertTime(at.Time, at.Message))
			}
		}

		newKW := []KeyWord{}
		for _, kw := range oldChat.Settings.KeyWords {
			if val, ok := changeSettings.KeyWords[kw.ID]; val || !ok {
				newKW = append(newKW, kw)
			}
		}
		for _, kw := range changeSettings.NewKeyWords {
			if len(kw.Key) > 2 && kw.Message != "" {
				newKW = append(newKW, MakeKeyWord(kw.Key, kw.Message))
			}
		}

		newChat := Chat{
			Name: oldChat.Name,
			Type: oldChat.Type,
			Settings: ChatSettings{
				NSFW:       changeSettings.NSFW,
				AlertTimes: newAT,
				KeyWords:   newKW,
			},
		}
		return newChat
	}

	chat, err := updateDatastoreChat(updateFunc, chatID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	insertTimersByChatID(chat.Settings.AlertTimes, chatID)
}
