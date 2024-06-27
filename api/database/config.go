package database

import (
	"bkoiki950/go-store/api/config"
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDatabase(DB_URI string) {
	if DB_URI == "" {
		DB_URI = config.GetEnv("DB_URI")
		if DB_URI == "" {
			log.Fatal("DB_URI not found")
		}
	}

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(DB_URI).SetServerAPIOptions(serverApi)
	client, err := mongo.Connect(context.TODO(), opts); if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
	mongoClient = client
}

func CloseDatabase(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}

	log.Println("Disconnected from MongoDB")
}

func GetActiveClient () (*mongo.Client, error) {
	if mongoClient == nil {
		return nil, mongo.ErrClientDisconnected
	}

	return mongoClient, nil
}

func GetDatabase(dbName string) (*mongo.Database, error) {
	if mongoClient == nil {
		return nil, mongo.ErrClientDisconnected
	}
	return mongoClient.Database(dbName), nil
}

func GetDefaultDatabase() (*mongo.Database, error) {
	// from db url get db name if db name is not provided
	databaseName := config.GetEnv("DB_NAME")
	if databaseName == "" {
		dbUri := config.GetEnv("DB_URI")
		if dbUri == "" {
			log.Fatal("DB_URI not found")
		}

		dbName := strings.Split(dbUri, "/")
		databaseName = dbName[len(dbName)-1]
	}

	db, err := GetDatabase(databaseName); if err != nil {
		return nil, err
	}
	return db, nil
}

func GetCollection(collectionName string) (*mongo.Collection, error) {
	db, err := GetDefaultDatabase(); if err != nil {
		return nil, err
	}
	return db.Collection(collectionName), nil
}