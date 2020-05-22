package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMongoDB(t *testing.T) {
	testSymbol := "IBM"
	testUser := "TestUser"
	dbConfig := loadConfig().DBConfig
	dbName, collectionNameTest, mongoServer := dbConfig.Name, dbConfig.CollectionTest, dbConfig.Server
	dbManagerTest := DBManagerMongo{dbName, collectionNameTest, mongoServer}
	server := InvestmentServer{nil, nil, dbManagerTest}
	t.Run(fmt.Sprintf("test add history"), func(t *testing.T) {
		err := dbManagerTest.AddHistory(testUser, testSymbol)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run(fmt.Sprintf("test read history"), func(t *testing.T) {
		history, err := dbManagerTest.GetHistory(testUser)
		if err != nil {
			t.Error(err)
		}
		if len(history) == 0 {
			t.Error("empty history")
		}
	})

	t.Run(fmt.Sprintf("test delete collection"), func(t *testing.T) {
		err := deleteMongoCollection(dbName, collectionNameTest, mongoServer)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run(fmt.Sprintf("test response 200"), func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/db?symbol=%s&user=%s", testSymbol, testUser), nil)
		response := httptest.NewRecorder()
		server.DBHandler(request, response)
		wantCode := 200
		if response.Code != wantCode {
			t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
		}
	})

	//t.Run(fmt.Sprintf("test response 500"), func(t *testing.T) {
	//	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s;_;news", testSymbolUnreal), nil)
	//	response := httptest.NewRecorder()
	//	server.NewsHandler(request, response)
	//	wantCode := 500
	//	if response.Code != wantCode {
	//		t.Error(fmt.Sprintf("wrong response code, want %d, get %d", wantCode, response.Code))
	//	}
	//})
}
