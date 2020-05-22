package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPlotAlphaVantage(t *testing.T) {
	testSymbolReal := "IBM"
	testSymbolUnreal := "unrealSymbol"
	apiKey := loadConfig().VentageKeyTest
	plotManagerTest := PlotManagerAlphaVantage{apiKey}
	server := InvestmentServer{nil, plotManagerTest, nil}

	t.Run(fmt.Sprintf("test real stock symbol:%s", testSymbolReal), func(t *testing.T) {
		plot, err := plotManagerTest.GetPlot(testSymbolReal)
		if err != nil {
			t.Error(err)
		}
		if len(plot) == 0 {
			t.Error("empty plot")
		}
	})

	t.Run(fmt.Sprintf("test unreal stock symbol:%s", testSymbolUnreal), func(t *testing.T) {
		plot, err := plotManagerTest.GetPlot(testSymbolUnreal)
		if err == nil {
			t.Error(err)
		}
		if len(plot) != 0 {
			t.Error("not empty plot")
		}
	})

	t.Run(fmt.Sprintf("test response 200"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/plot?symbol=%s", testSymbolReal), nil)
		response := httptest.NewRecorder()
		server.PlotHandler(request, response)
		wantCode := 200
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})

	t.Run(fmt.Sprintf("test response 500"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/plot?symbol=%s", testSymbolUnreal), nil)
		response := httptest.NewRecorder()
		server.PlotHandler(request, response)
		wantCode := 500
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})

}
