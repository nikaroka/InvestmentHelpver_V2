package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRequest struct {
	UserID   string `bson:"userID,omitempty"`
	StockKey string `bson:"stockKey,omitempty"`
}

type DBManagerMongo struct {
	dbName         string
	collectionName string
	dbServer       string
}

func (dbManager DBManagerMongo) GetHistory(userID string) ([]UserRequest, error) {
	var userRequests []UserRequest
	dbName, collectionName, dbServer := dbManager.dbName, dbManager.collectionName, dbManager.dbServer
	collection, client, err := GetCollection(dbName, collectionName, dbServer)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	cur, err := collection.Find(context.TODO(), bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var userRequest UserRequest
		err := cur.Decode(&userRequest)
		if err != nil {
			return nil, err
		}
		userRequests = append(userRequests, userRequest)
	}
	return userRequests, nil
}

func (dbManager DBManagerMongo) AddHistory(userID string, stockKey string) error {
	dbName, collectionName, dbServer := dbManager.dbName, dbManager.collectionName, dbManager.dbServer
	collection, client, err := GetCollection(dbName, collectionName, dbServer)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	userReq := UserRequest{userID, stockKey}
	_, err = collection.InsertOne(context.TODO(), userReq)
	if err != nil {
		return err
	}
	return nil
}

func GetCollection(dbName string, collectionName string, mongoServer string) (*mongo.Collection, *mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoServer))

	if err != nil {
		return nil, nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, nil, err
	}
	collection := client.Database(dbName).Collection(collectionName)
	return collection, client, nil
}

func deleteMongoCollection(dbName string, collectionName string, mongoServer string) error {
	collection, client, err := GetCollection(dbName, collectionName, mongoServer)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	err = collection.Drop(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

//func requestHandlerReadDB(r *http.Request, dbChannel chan []byte, errChannel chan error) {
//	userId := strings.Split(r.URL.Path[1:], ";")[1]
//	dbConfig := loadConfig().DBConfig
//	dbName, collectionName, mongoServer := dbConfig.Name, dbConfig.Collection, dbConfig.Server
//	userRequests, err := GetHistory(dbName, collectionName, mongoServer, userId)
//	if err != nil {
//		errChannel <- err
//		return
//	}
//	jsonDB, err := json.Marshal(userRequests)
//	if err != nil {
//		errChannel <- err
//		return
//	}
//	dbChannel <- jsonDB
//}
//
//func requestHandlerWriteDB(r *http.Request, dbChannel chan []byte, errChannel chan error) {
//	symbol := strings.Split(r.URL.Path[1:], ";")[0]
//	userId := strings.Split(r.URL.Path[1:], ";")[1]
//	dbConfig := loadConfig().DBConfig
//	dbName, collectionName, mongoServer := dbConfig.Name, dbConfig.Collection, dbConfig.Server
//	err := AddHistory(userId, symbol, dbName, collectionName, mongoServer)
//	if err != nil {
//		errChannel <- err
//		return
//	}
//	jsonDB, err := json.Marshal("")
//	if err != nil {
//		errChannel <- err
//		return
//	}
//	dbChannel <- jsonDB
//}
