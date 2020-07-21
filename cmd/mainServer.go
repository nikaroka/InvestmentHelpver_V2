package main

import (
	"InvestmentHelpver_V2/internal/db"
	"InvestmentHelpver_V2/internal/news"
	"InvestmentHelpver_V2/internal/plot"
	"fmt"
	"os"

	"github.com/jinzhu/configor"

	"encoding/json"
	"log"
	"net/http"
)

// Структура отражающая config.yml
type Config struct {
	DBConfig struct {
		Name           string `default:"dbName"`
		Collection     string `default:"dbCollection"`
		CollectionTest string `default:"dbCollectionTest"`
		DBserver       string `default:"dbServer"`
	}
	VentageKey string `default:"key"`
	LocalPort  string `default:"8888"`
}

// Метод считывающий config.yml и возвращающий его содержимое в экземпляре структуры Config
func loadConfig() Config {
	config := Config{}
	err := configor.Load(&config, "config.yml")
	if err != nil {
		panic(err)
	}
	dbServer, ok := os.LookupEnv("dbserver")
	if ok {
		config.DBConfig.DBserver = dbServer
	}
	return config
}

// Главная структура программы включающая в себя интерфейсы основных модулей(менеджеров)
type InvestmentServer struct {
	NewsManager news.NewsManager
	PlotManager plot.PlotManager
	DBManager   db.DBManager
}

func NewInvestmentServer(newsManager news.NewsManager, plotManager plot.PlotManager, dbManager db.DBManager) InvestmentServer {
	return InvestmentServer{newsManager, plotManager, dbManager}
}

// Метод обрабатывающий запросы на получение новостей, вызывает внутри себя метод GetNews и отправляет полученый список новостей в виде Json
func (server *InvestmentServer) NewsHandler(r *http.Request, w http.ResponseWriter) {
	symbol := r.URL.Query()["symbol"][0]
	newsSLice, err := server.NewsManager.GetNews(symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	jsonData, err := json.Marshal(newsSLice)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
}

// Метод обрабатывающий запросы на получение графика, вызывает внутри себя метод GetPlot и отправляет полученый список свечей в виде Json
func (server *InvestmentServer) PlotHandler(r *http.Request, w http.ResponseWriter) {
	symbol := r.URL.Query()["symbol"][0]
	plotSlice, err := server.PlotManager.GetPlot(symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	jsonData, err := json.Marshal(plotSlice)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
}

// Метод обрабатывающий запросы на работу с базой данных, вызывает внутри себя метод AddHistory и отправляет статус 200
// (GetHistory пока не используется т.к. пока нет реализации просмотра истории на сайте)
func (server *InvestmentServer) DBHandler(r *http.Request, w http.ResponseWriter) {
	user := r.URL.Query()["user"][0]
	symbol := r.URL.Query()["symbol"][0]
	err := server.DBManager.AddHistory(user, symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(http.StatusOK)
}

// Метод срабатывающий в случае неправильного запроса со стороны сайта или возникновения ошибки во время обработки запроса,
// используется в остальных Handler-ах
func (server *InvestmentServer) ErrorHandler(httpStatus int, r *http.Request, w http.ResponseWriter) {
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(httpStatus)
}

// Создаем экземпляры реализаций интерфейсов сервера, затем создаем экземпляр самого сервера с этими реализациями
var newsManager = news.NewNewsManagerYahoo()
var plotManager = plot.NewPlotManagerAlphaVantage(loadConfig().VentageKey)
var dbManager = db.NewDBManagerMongo(loadConfig().DBConfig.Name, loadConfig().DBConfig.Collection, loadConfig().DBConfig.DBserver)
var server = NewInvestmentServer(newsManager, plotManager, dbManager)

// Главный обработчик, вызывается при получении запроса на сервер, решает какой из Handler-ов должен этот запрос обработать
func mainHandler(w http.ResponseWriter, r *http.Request) {
	command := r.URL.EscapedPath()
	switch {
	case command == "/db":
		log.Printf("%s\n", "db")
		server.DBHandler(r, w)
	case command == "/news":
		log.Printf("%s\n", "news")
		server.NewsHandler(r, w)
	case command == "/plot":
		log.Printf("%s\n", "plot")
		server.PlotHandler(r, w)
	default:
		log.Printf("%s\n", "wrong command")
		server.ErrorHandler(http.StatusBadRequest, r, w)
	}
}

func main() {
	localPort := ":" + loadConfig().LocalPort
	fmt.Println("jok")
	http.HandleFunc("/", mainHandler)
	log.Printf("%s\n", "Server is Up")
	err := http.ListenAndServe(localPort, nil)
	if err != nil {
		log.Print(err)
		return
	}
}
