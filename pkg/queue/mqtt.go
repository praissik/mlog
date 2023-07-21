package queue

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type LogInfo struct {
	//ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CorrelationID string    `bson:"correlation_id,omitempty"`
	Level         string    `bson:"level,omitempty"`
	Message       string    `bson:"message,omitempty"`
	Service       string    `bson:"service"`
	Datetime      time.Time `bson:"datetime"`
	Duration      int64    `bson:"duration"`
}

type LogError struct {
	//ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CorrelationID string    `bson:"correlation_id,omitempty"`
	StackTrace    string    `bson:"stack_trace,omitempty"`
	Level         string    `bson:"level,omitempty"`
	Message       string    `bson:"message,omitempty"`
	Service       string    `bson:"service"`
	Datetime      time.Time `bson:"datetime"`
}

func newMQTTClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("pub")
	opts.SetCleanSession(false)

	return mqtt.NewClient(opts)
}

func newMQTTConnection() mqtt.Client {
	client := newMQTTClient()
	connect(client)
	return client
}

func connect(client mqtt.Client) {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func (l *LogError) PrepareLog() string {
	var buffer bytes.Buffer

	buffer.WriteString(`{"correlation_id":"`)
	buffer.WriteString(l.CorrelationID)

	buffer.WriteString(`","level":"`)
	buffer.WriteString(l.Level)

	buffer.WriteString(`","message":"`)
	buffer.WriteString(l.Message)

	buffer.WriteString(`","service":"`)
	buffer.WriteString(l.Service)

	buffer.WriteString(`","datetime":"`)
	buffer.WriteString(l.Datetime.Format(time.RFC3339Nano))

	buffer.WriteString(`","stack_trace":"`)
	st := strings.ReplaceAll(l.StackTrace, "\t", "")
	buffer.WriteString(strings.ReplaceAll(st, "\n", "\\n"))
	buffer.WriteString(`"}`)

	return buffer.String()
}

func (l *LogInfo) PrepareLog() string {
	var buffer bytes.Buffer

	buffer.WriteString(`{"correlation_id":"`)
	buffer.WriteString(l.CorrelationID)

	buffer.WriteString(`","level":"`)
	buffer.WriteString(l.Level)

	buffer.WriteString(`","message":"`)
	buffer.WriteString(l.Message)

	buffer.WriteString(`","service":"`)
	buffer.WriteString(l.Service)

	buffer.WriteString(`","duration":`)
	buffer.WriteString(strconv.FormatInt(l.Duration, 10))

	buffer.WriteString(`,"datetime":"`)
	buffer.WriteString(l.Datetime.Format(time.RFC3339Nano))

	buffer.WriteString(`","stack_trace":""}`)

	return buffer.String()
}

func Publish(payload string) {
	client := newMQTTConnection()
	token := client.Publish("logs", 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Fatal(token.Error())
	} else {
		log.Printf(".")
	}
	client.Disconnect(250)
}
