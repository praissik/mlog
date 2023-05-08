package db

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

const (
	CollectionLogs = "logs"
)

func GetMongoClient() (*mongo.Client, func(), error) {
	clientOpts := options.Client().ApplyURI(viper.GetString("mongo.url"))
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}

	return client, func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Println(err.Error())
		}
	}, nil
}
