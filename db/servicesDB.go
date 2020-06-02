package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Структура UserRequest содержит ID пользователя и символ финансового актива информацию по которому он запрашивал
type UserRequest struct {
	UserID      string `bson:"userID,omitempty"`   // индентификатор пользователя в системе
	StockSymbol string `bson:"stockKey,omitempty"` // символ акции
}

//интерфейс менеджера графиков, реализующие его струтуры должны иметь метод GetHistory принимающий ID пользователя и возвращающий историю его запросов в виде списка экземпляров UserRequest
//и метод AddHistory принимающий ID пользователя и символ финансового актива, и записыющий эту информацию в базу данных
type DBManager interface {
	GetHistory(string) ([]UserRequest, error) // принимает ID пользователя, возвращает историю его запросов в виде списка экземпляров UserRequest
	AddHistory(string, string) error          // принимает ID пользователя и символ финансового актива, записывает эту информацию в базу данных
}

//Реализация интерфейса DBManager, отвечает за работу с MonboDB, имеет параметры
//DBName - имя базы данных, CollectionName - имя колекции, DBServer - сервер базы данных
type DBManagerMongo struct {
	DBName         string
	CollectionName string
	DBServer       string
}

//Конструктор для структуры DBManagerMongo
func NewDBManagerMongo(dbName string, collectionName string, dbServer string) DBManager {
	dbManager := DBManagerMongo{dbName, collectionName, dbServer}
	return dbManager
}

//Метод структуры DBManagerMongo, принимает ID пользователя, возвращает список экземпляров структуры UserRequest
func (dbManager DBManagerMongo) GetHistory(userID string) ([]UserRequest, error) {
	var userRequests []UserRequest
	dbName, collectionName, dbServer := dbManager.DBName, dbManager.CollectionName, dbManager.DBServer
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

//Метод структуры DBManagerMongo, принимает ID пользователя и символ финансового актива, возвращает ошибку если она есть
func (dbManager DBManagerMongo) AddHistory(userID string, symbol string) error {
	dbName, collectionName, dbServer := dbManager.DBName, dbManager.CollectionName, dbManager.DBServer
	collection, client, err := GetCollection(dbName, collectionName, dbServer)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	userReq := UserRequest{userID, symbol}
	_, err = collection.InsertOne(context.TODO(), userReq)
	if err != nil {
		return err
	}
	return nil
}

//Вспомогательный метод возвращаюий указатели на mongo.Collection и mongo.Сlient
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

//Метод удаляющий коллекцию, сейчас используется для удаление тестовой коллекции после прохождения тестов
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
