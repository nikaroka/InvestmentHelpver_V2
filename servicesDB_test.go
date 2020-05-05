package main

import (
	"testing"
)

func TestAddToMonga(t *testing.T)  {
	dbConfig := loadConfig().DBConfig
	mongoServer := dbConfig.Server
	err := addToMongo("testUser", "testKey", "testDB", "testCollection", mongoServer)
	if err != nil{
		t.Error(err)
	}
}
func TestDeleteMongoCollection(t *testing.T)  {
	dbConfig := loadConfig().DBConfig
	mongoServer := dbConfig.Server
	err := deleteMongoCollection("testDB", "testCollection", mongoServer)
	if err != nil{
		t.Error(err)
	}
}
