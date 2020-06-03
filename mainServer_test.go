package main

import (
	"InvestmentHelpver_V2/db"
	"InvestmentHelpver_V2/news"
	"InvestmentHelpver_V2/plot"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testSymbolReal = "IBM"
var testSymbolUnreal = "unrealSymbol"
var testUser = "userTest"

func TestNewsHandler(t *testing.T) {
	newsManagerYahoo := news.NewNewsManagerYahoo()
	serverYahoo := NewInvestmentServer(newsManagerYahoo, nil, nil)

	t.Run(fmt.Sprintf("test response 200 newsManagerYahoo"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news?symbol=%s", testSymbolReal), nil)
		response := httptest.NewRecorder()
		serverYahoo.NewsHandler(request, response)
		wantCode := 200
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})

	t.Run(fmt.Sprintf("test response 500 newsManagerYahoo"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news?symbol=%s", testSymbolUnreal), nil)
		response := httptest.NewRecorder()
		serverYahoo.NewsHandler(request, response)
		wantCode := 500
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})
}

func TestPlotHandler(t *testing.T) {
	apiKey := loadConfig().VentageKey
	plotManagerAlphaVentage := plot.NewPlotManagerAlphaVantage(apiKey)
	serverAlphaVentage := NewInvestmentServer(nil, plotManagerAlphaVentage, nil)

	t.Run(fmt.Sprintf("test response 200 plotManagerAlphaVentage"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news?symbol=%s", testSymbolReal), nil)
		response := httptest.NewRecorder()
		serverAlphaVentage.PlotHandler(request, response)
		wantCode := 200
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})

	t.Run(fmt.Sprintf("test response 500 plotManagerAlphaVentage"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news?symbol=%s", testSymbolUnreal), nil)
		response := httptest.NewRecorder()
		serverAlphaVentage.PlotHandler(request, response)
		wantCode := 500
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})
}

func TestDBHandler(t *testing.T) {
	dbManagerMongo := db.NewDBManagerMongo(loadConfig().DBConfig.Name, loadConfig().DBConfig.CollectionTest, loadConfig().DBConfig.Server)
	serverDBManagerMongo := NewInvestmentServer(nil, nil, dbManagerMongo)

	t.Run(fmt.Sprintf("test response 200 dbManagerMongo"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/db?symbol=%s&user=%s", testSymbolReal, testUser), nil)
		response := httptest.NewRecorder()
		serverDBManagerMongo.DBHandler(request, response)
		wantCode := 200
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})
}
