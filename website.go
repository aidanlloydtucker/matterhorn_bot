package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
)

var p = proxy.New(&ace.Options{BaseDir: "views"})

func startWebsite() {

	if HttpPort == "" {
		HttpPort = "8080"
	}

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(webNotFoundHandler)
	router.HandleFunc("/chat/{id:[0-9]+}", webChatHandler).Methods("GET")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./static/"))))

	allHandler := handlers.CompressHandler(router)

	http.ListenAndServe(":"+HttpPort, allHandler)
}

func webNotFoundHandler(w http.ResponseWriter, r *http.Request) {
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
	chatId := vars["id"]
	fmt.Fprintln(w, "chatid:", chatId)
}
