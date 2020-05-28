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
	symbol := r.URL.Query()["symbol"][0]
	news, err := server.newsManager.GetNews(symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	server.SendResponseToSite(news, r, w)
}

func (server *InvestmentServer) PlotHandler(r *http.Request, w http.ResponseWriter) {
	symbol := r.URL.Query()["symbol"][0]
	plot, err := server.plotManager.GetPlot(symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	server.SendResponseToSite(plot, r, w)
}

func (server *InvestmentServer) DBHandler(r *http.Request, w http.ResponseWriter) {
	user := r.URL.Query()["user"][0]
	symbol := r.URL.Query()["symbol"][0]
	err := server.dbManager.AddHistory(user, symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	server.SendResponseToSite(nil, r, w)
}

func (server *InvestmentServer) ErrorHandler(httpStatus int, r *http.Request, w http.ResponseWriter) {
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok == true {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(httpStatus)
}

func (server *InvestmentServer) SendResponseToSite(responseData interface{}, r *http.Request, w http.ResponseWriter) {
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok == true {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(http.StatusOK)
	if responseData != nil {
		jsonData, err := json.Marshal(responseData)
		if err != nil {
			server.ErrorHandler(http.StatusInternalServerError, r, w)
			return
		}
		w.Write(jsonData)
		w.Header().Set("Content-Type", "application/json")
	}
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
	default:
		fmt.Println("wrong command")
		server.ErrorHandler(http.StatusBadRequest, r, w)
	}
}

func main() {
	fmt.Println("GO")
	localPort := ":" + loadConfig().LocalPort
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(localPort, nil)
}
