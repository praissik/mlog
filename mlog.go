package mlog

import (
	"github.com/praissik/mlog/pkg/db"
	"github.com/spf13/viper"
	"log"
	"runtime/debug"
	"time"
)

type data struct {
	//ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CorrelationID string    `bson:"correlation_id,omitempty"`
	StackTrace    string    `bson:"stack_trace,omitempty"`
	Action        string    `bson:"action,omitempty"`
	Message       string    `bson:"message,omitempty"`
	Service       string    `bson:"service"`
	Datetime      time.Time `bson:"datetime"`
}

func Info(correlationID, action string) {
	d := data{
		CorrelationID: correlationID,
		Action:        action,
	}
	d.create()
}

func Error(correlationID string, err error) {
	d := data{
		CorrelationID: correlationID,
		StackTrace:    string(debug.Stack()),
		Message:       err.Error(),
	}
	d.create()
}

func (d *data) create() {
	d.Datetime = time.Now().UTC()
	d.Service = viper.GetString("server.name")

	mongoClient, deferF, err := db.GetMongoClient()
	defer deferF()
	if err != nil {
		return
	}
	_, err = mongoClient.
		Database(viper.GetString("mongo.db")).
		Collection(db.CollectionLogs).
		InsertOne(nil, d)
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}
