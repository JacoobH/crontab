package worker

import (
	"context"
	"fmt"
	"github.com/JacoobH/crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogSink struct {
	client     *mongo.Client
	collection *mongo.Collection
	logChan    chan *common.JobLog
}

var (
	G_logSink *LogSink
)

func (logSink *LogSink) writeLoop() {
	var (
		log *common.JobLog
	)
	for {
		select {
		case log = <-logSink.logChan:

		}
	}
}

func InitLogSink() (err error) {
	var (
		client     *mongo.Client
		clientOps  *options.ClientOptions
		collection *mongo.Collection
	)
	clientOps = options.Client().
		ApplyURI(G_config.MongodbUri).
		SetConnectTimeout(time.Duration(G_config.MongodbConnectTimeout) * time.Millisecond)
	if client, err = mongo.Connect(context.TODO(), clientOps); err != nil {
		fmt.Println(err)
		return
	}
	//select table my_collection
	collection = client.Database("cron").Collection("log")

	G_logSink = &LogSink{
		client:     client,
		collection: collection,
		logChan:    make(chan *common.JobLog, 1000),
	}
	G_logSink.writeLoop()
	return
}
