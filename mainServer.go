package main

import (
	"github.com/jinzhu/configor"
	"net/http"
	"strings"
)

type Config struct {
	DBConfig struct {
		Name     string `default:"dbName"`
		Collection     string `default:"dbCollection"`
		Server string `default:"dbServer"`
	}
	VentageKey string `default:"dbName"`
	PageServer string `default:"localhost"`
	LocalPort string `default:"8888"`
}

func loadConfig() Config{
	var config Config
	configor.Load(&config, "config.yml")
	return config
}

func sendToSite (w http.ResponseWriter, jsn []byte){
	pageServer := loadConfig().PageServer
	if jsn != nil{
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", pageServer)
		w.WriteHeader(http.StatusOK)
		w.Write(jsn)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", pageServer)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	command := strings.Split(r.URL.Path[1:], ";")[2]
	resultChannel := make(chan []byte)
	errChannel := make(chan error)

	switch {
		case command =="writeDB":
			go requestHandlerWriteDB(r, resultChannel, errChannel)
		case command == "news":
			go requestHandlerNews(r, resultChannel, errChannel)
		case command == "plot":
			go requestHandlerPlot(r, resultChannel, errChannel)
		case command == "readDB":
			go requestHandlerReadDB(r, resultChannel, errChannel)
		default:
			sendToSite(w, nil)
			return
	}

	select {
		case result:= <- resultChannel:
			sendToSite(w, result)
		case <- errChannel:
			sendToSite(w, nil)
	}
}

func main() {
	localPort := ":" + loadConfig().LocalPort
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(localPort, nil)
	}

