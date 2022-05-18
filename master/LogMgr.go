package master

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
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

	G_logMgr = &LogMgr{
		client:        client,
		logCollection: collection,
	}
	return
}

func (logMgr *LogMgr) ListLog() {

}
