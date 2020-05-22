package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNewsYahoo(t *testing.T) {
	testSymbolReal := "IBM"
	testSymbolUnreal := "unrealSymbol"
	newsManagerTest := NewsManagerYahoo{}
	server := InvestmentServer{newsManagerTest, nil, nil}
	t.Run(fmt.Sprintf("test real stock symbol:%s", testSymbolReal), func(t *testing.T) {
		news, err := newsManagerTest.GetNews(testSymbolReal)
		if err != nil {
			t.Error(err)
		}
		if len(news) == 0 {
			t.Error("empty news")
		}
	})

	t.Run(fmt.Sprintf("test unreal stock symbol:%s", testSymbolUnreal), func(t *testing.T) {
		news, err := newsManagerTest.GetNews(testSymbolUnreal)
		if err == nil {
			t.Error(err)
		}
		if len(news) != 0 {
			t.Error("not empty news")
		}
	})

	t.Run(fmt.Sprintf("test response 200"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news?symbol=%s", testSymbolReal), nil)
		response := httptest.NewRecorder()
		server.NewsHandler(request, response)
		wantCode := 200
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})

	t.Run(fmt.Sprintf("test response 500"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news?symbol=%s", testSymbolUnreal), nil)
		response := httptest.NewRecorder()
		server.NewsHandler(request, response)
		wantCode := 500
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})
}
