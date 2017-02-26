package main

import (
	"net/http"

	"log"
	"strconv"

	"cloud.google.com/go/datastore"
	"fmt"
	chatpkg "github.com/billybobjoeaglt/matterhorn_bot/chat"
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	r.Handle("GET", "/quotes/:id", webQuotesHandler)
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

	chat, exists, err := DatastoreInst.GetChat(chatID)
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
		"QuotesDoc":  chat.Settings.QuotesDoc,
	}

	c.Header("Content-Type", "text/HTML")

	if err := tpl.Execute(c.Writer, data); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

type ChangedChatSettings struct {
	NSFW          bool                `json:"nsfw"`
	NewAlertTimes []chatpkg.AlertTime `json:"new_alert_times"`
	AlertTimes    map[int]bool        `json:"alert_times"`
	NewKeyWords   []chatpkg.KeyWord   `json:"new_key_words"`
	KeyWords      map[int]bool        `json:"key_words"`
	QuotesDoc     int                 `json:"quotes_doc"`
}

func webChatChangeHandler(c *gin.Context) {
	chatIdStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	_, exists, err := DatastoreInst.GetChat(chatID)
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

	updateFunc := func(oldChat chatpkg.Chat) chatpkg.Chat {
		newAT := []chatpkg.AlertTime{}
		for _, at := range oldChat.Settings.AlertTimes {
			if val, ok := changeSettings.AlertTimes[at.ID]; val || !ok {
				newAT = append(newAT, at)
			}
		}
		for _, at := range changeSettings.NewAlertTimes {
			if at.Time != "" && at.Message != "" {
				newAT = append(newAT, chatpkg.MakeAlertTime(at.Time, at.Message))
			}
		}

		newKW := []chatpkg.KeyWord{}
		for _, kw := range oldChat.Settings.KeyWords {
			if val, ok := changeSettings.KeyWords[kw.ID]; val || !ok {
				newKW = append(newKW, kw)
			}
		}
		for _, kw := range changeSettings.NewKeyWords {
			if len(kw.Key) > 2 && kw.Message != "" {
				newKW = append(newKW, chatpkg.MakeKeyWord(kw.Key, kw.Message))
			}
		}

		quotesDoc := oldChat.Settings.QuotesDoc
		if changeSettings.QuotesDoc != oldChat.Settings.QuotesDoc {
			if changeSettings.QuotesDoc == 0 {
				if oldChat.Settings.QuotesDoc != 0 {
					quotesDoc = oldChat.Settings.QuotesDoc
				} else {
					quotesDoc = int(rand.Int31())
				}
			} else {
				//TODO: Vet so you can only select a chat that exists
				quotesDoc = changeSettings.QuotesDoc
			}
		}

		newChat := chatpkg.Chat{
			Name: oldChat.Name,
			Type: oldChat.Type,
			Settings: chatpkg.ChatSettings{
				NSFW:       changeSettings.NSFW,
				AlertTimes: newAT,
				KeyWords:   newKW,
				QuotesDoc:  quotesDoc,
			},
		}
		return newChat
	}

	chat, err := DatastoreInst.UpdateChat(updateFunc, chatID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	insertTimersByChatID(chat.Settings.AlertTimes, chatID)
}

func webQuotesHandler(c *gin.Context) {
	quotesDocStr := c.Param("id")
	quotesDoc, err := strconv.ParseInt(quotesDocStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	quotes := []commands.Quote{}

	_, err = DatastoreInst.Client.GetAll(
		DatastoreInst.Ctx,
		datastore.NewQuery(commands.QuoteKeyKind).Namespace(chatpkg.KeyNamespace).Filter("Document =", quotesDoc),
		&quotes)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	quotesStr := make([]string, len(quotes))
	for i, quote := range quotes {
		quoteText := quote.Text
		if !quote.Manual {
			quoteText = fmt.Sprintf(`"%s" - %s %s`, quote.Text, quote.Author, quote.Date.Format("01/2/06"))
		}

		quotesStr[i] = quoteText
	}

	tpl, err := p.Load("base", "quotes", nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := map[string]interface{}{
		"QuotesDoc": quotesDoc,
		"Quotes":    quotesStr,
	}

	c.Header("Content-Type", "text/HTML")

	if err := tpl.Execute(c.Writer, data); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
