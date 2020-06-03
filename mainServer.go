package main

import (
	"InvestmentHelpver_V2/db"
	"InvestmentHelpver_V2/news"
	"InvestmentHelpver_V2/plot"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/configor"
	"net/http"
)

//Структура отражающая config.yml
type Config struct {
	DBConfig struct {
		Name           string `default:"dbName"`
		Collection     string `default:"dbCollection"`
		CollectionTest string `default:"dbCollectionTest"`
		Server         string `default:"dbServer"`
	}
	VentageKey string `default:"key"`
	LocalPort  string `default:"8888"`
}

//Метод считывающий config.yml и возвращающий его содержимое в экземпляре структуры Config
func loadConfig() Config {
	var config Config
	configor.Load(&config, "config.yml")
	return config
}

//Главная структура программы включающая в себя интерфейсы основных модулей(менеджеров)
type InvestmentServer struct {
	NewsManager news.NewsManager
	PlotManager plot.PlotManager
	DBManager   db.DBManager
}

func NewInvestmentServer(NewsManager news.NewsManager, PlotManager plot.PlotManager, DBManager db.DBManager) InvestmentServer {
	return InvestmentServer{NewsManager, PlotManager, DBManager}
}

//Метод обрабатывающий запросы на получение новостей, вызывает внутри себя метод GetNews и отправляет полученый список новостей в виде Json
func (server *InvestmentServer) NewsHandler(r *http.Request, w http.ResponseWriter) {
	symbol := r.URL.Query()["symbol"][0]
	news, err := server.NewsManager.GetNews(symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok == true {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	jsonData, err := json.Marshal(news)
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

//Метод обрабатывающий запросы на получение графика, вызывает внутри себя метод GetPlot и отправляет полученый список свечей в виде Json
func (server *InvestmentServer) PlotHandler(r *http.Request, w http.ResponseWriter) {
	symbol := r.URL.Query()["symbol"][0]
	plot, err := server.PlotManager.GetPlot(symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok == true {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	jsonData, err := json.Marshal(plot)
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

//Метод обрабатывающий запросы на работу с базой данных, вызывает внутри себя метод AddHistory и отправляет статус 200
//(GetHistory пока не используется т.к. пока нет реализации просмотра истории на сайте)
func (server *InvestmentServer) DBHandler(r *http.Request, w http.ResponseWriter) {
	user := r.URL.Query()["user"][0]
	symbol := r.URL.Query()["symbol"][0]
	err := server.DBManager.AddHistory(user, symbol)
	if err != nil {
		server.ErrorHandler(http.StatusInternalServerError, r, w)
		return
	}
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok == true {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(http.StatusOK)
}

//Метод срабатывающий в случае неправильного запроса со стороны сайта или возникновения ошибки во время обработки запроса,
//используется в остальных Handler-ах
func (server *InvestmentServer) ErrorHandler(httpStatus int, r *http.Request, w http.ResponseWriter) {
	pageServer := ""
	if origin, ok := r.Header["Origin"]; ok == true {
		pageServer = origin[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", pageServer)
	w.WriteHeader(httpStatus)
}

//Создаем экземпляры реализаций интерфейсов сервера, затем создаем экземпляр самого сервера с этими реализациями
var newsManager = news.NewNewsManagerYahoo()
var plotManager = plot.NewPlotManagerAlphaVantage(loadConfig().VentageKey)
var dbManager = db.NewDBManagerMongo(loadConfig().DBConfig.Name, loadConfig().DBConfig.Collection, loadConfig().DBConfig.Server)
var server = NewInvestmentServer(newsManager, plotManager, dbManager)

//Главный обработчик, вызывается при получении запроса на сервер, решает какой из Handler-ов должен этот запрос обработать
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
	err := http.ListenAndServe(localPort, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
