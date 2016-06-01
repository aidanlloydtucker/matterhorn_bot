package main

import (
	"net/http"

	"encoding/json"

	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
)

var p = proxy.New(&ace.Options{BaseDir: "views"})

func startWebsite() {

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(webNotFoundHandler)
	router.HandleFunc("/chat/{id}", webChatHandler).Methods("GET")
	router.HandleFunc("/chat/{id}", webChatChangeHandler).Methods("PUT")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./static/"))))

	allHandler := handlers.CompressHandler(router)

	http.ListenAndServe(":"+HttpPort, allHandler)
}

func webNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	tpl, err := p.Load("base", "error", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Status":  404,
		"Message": "Not Found: The page you requested could not be found.",
	}

	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func webChatHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exists, err := redis.Bool(redisConn.Do("EXISTS", formatRedisKey(chatId)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "ChatId is not in database", http.StatusNotFound)
		return
	}

	v, err := redis.Values(redisConn.Do("HGETALL", formatRedisKey(chatId)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err, chat := FromRedisChatInfo(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl, err := p.Load("base", "chat", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"ChatName": chat.Name,
		"ChatId":   chatId,
		"SettingsBool": map[string]interface{}{
			"NSFW": chat.NSFW,
		},
		"KeyWords":   chat.KeyWords,
		"AlertTimes": chat.AlertTimes,
	}

	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func webChatChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exists, err := redis.Bool(redisConn.Do("EXISTS", formatRedisKey(chatId)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "ChatId is not in database", http.StatusNotFound)
		return
	}

	var settings ChatSettings

	err = json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newKWs []KeyWord

	for _, kw := range settings.KeyWords {
		if kw.Key != "" && kw.Message != "" {
			newKWs = append(newKWs, kw)
		}
	}
	settings.KeyWords = newKWs

	var newATs []AlertTime
	// struct{}{} doesnt take up any space
	dupATs := make(map[string]struct{})
	var setAlerts int

	for _, at := range settings.AlertTimes {
		if at.Time != "" && at.Message != "" && timeRegex.MatchString(at.Time) {
			if _, ok := dupATs[at.Time]; !ok {
				dupATs[at.Time] = struct{}{}
				newATs = append(newATs, at)
				if setAlerts < MAX_ALERTS_ALLOWED {
					hour, min, err := parseTimes(at.Time)
					if err != nil {
						continue
					}
					go startReminder(hour, min, at.Message, chatId)
				}
				setAlerts++
			}
		}
	}
	settings.AlertTimes = newATs

	sJson, err := json.Marshal(settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = redisConn.Do("HSET", redis.Args{}.Add(formatRedisKey(chatId)).Add("settings").Add(string(sJson))...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
