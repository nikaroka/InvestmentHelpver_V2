package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/configor"
	"net/http"
)

type Config struct {
	DBConfig struct {
		Name           string `default:"dbName"`
		Collection     string `default:"dbCollection"`
		CollectionTest string `default:"dbCollectionTest"`
		Server         string `default:"dbServer"`
	}
	VentageKey     string `default:"key"`
	VentageKeyTest string `default:"keytest"`
	PageServer     string `default:"localhost"`
	LocalPort      string `default:"8888"`
}

func loadConfig() Config {
	var config Config
	configor.Load(&config, "config.yml")
	return config
}

type InvestmentServer struct {
	newsManager NewsManager
	plotManager PlotManager
	dbManager   DBManager
}

type NewsManager interface {
	GetNews(string) ([]News, error)
}
type PlotManager interface {
	GetPlot(string) ([]Day, error)
}
type DBManager interface {
	GetHistory(string) ([]UserRequest, error)
	AddHistory(string, string) error
}

func (server *InvestmentServer) NewsHandler(r *http.Request, w http.ResponseWriter) {
	pageServer := loadConfig().PageServer
	//symbol := strings.Split(r.URL.Path[1:], ";")[0]
	symbol := r.URL.Query()["symbol"][0]
	news, err := server.newsManager.GetNews(symbol)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", pageServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonNews, err := json.Marshal(news)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", pageServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonNews)
}

func (server *InvestmentServer) PlotHandler(r *http.Request, w http.ResponseWriter) {
	pageServer := loadConfig().PageServer
	symbol := r.URL.Query()["symbol"][0]
	plot, err := server.plotManager.GetPlot(symbol)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", pageServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonPlot, err := json.Marshal(plot)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", pageServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonPlot)
}

func (server *InvestmentServer) DBHandler(r *http.Request, w http.ResponseWriter) {
	pageServer := loadConfig().PageServer
	user := r.URL.Query()["user"][0]
	symbol := r.URL.Query()["symbol"][0]
	err := server.dbManager.AddHistory(user, symbol)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", pageServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(http.StatusOK)
}

var newsManager = NewsManagerYahoo{}
var plotManager = PlotManagerAlphaVantage{apiKey: loadConfig().VentageKey}
var dbManager = DBManagerMongo{dbName: loadConfig().DBConfig.Name,
	collectionName: loadConfig().DBConfig.Collection,
	dbServer:       loadConfig().DBConfig.Server}
var server = InvestmentServer{newsManager, plotManager, dbManager}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	command := r.URL.EscapedPath()
	switch {
	case command == "/db":
		fmt.Println("db")
		server.DBHandler(r, w)
	case command == "/news":
		fmt.Println("news")
		server.NewsHandler(r, w)
	case command == "/plot":
		fmt.Println("plot")
		server.PlotHandler(r, w)
		//default:
		//	sendToSite(w, nil)
		//	return
	}
}

func main() {
	localPort := ":" + loadConfig().LocalPort
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(localPort, nil)
}
