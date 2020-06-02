package db

import (
	"fmt"
	"testing"
)

func TestMongoDB(t *testing.T) {
	testSymbol := "IBM"
	testUser := "TestUser"
	dbName, collectionNameTest, mongoServer := "InvestmentHelper", "CollectionTest", "mongodb://127.0.0.1:27017"
	dbManagerTest := DBManagerMongo{dbName, collectionNameTest, mongoServer}
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
}
