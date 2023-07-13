package mdb

import (
	"context"
	"fmt"
	"frame/models"
	"frame/services/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

const moduleName = "mdb"

var mongoDatabase string
var mongoClient *mongo.Client

func Start() {
	zlog.Warn(moduleName, "start ...")
	loadConfig()
	uri := fmt.Sprintf("%v://%v", mongoScheme, mongoAddress)
	if mongoUsername != "" && mongoPassword != "" {
		uri = fmt.Sprintf("%v://%v:%v@%v", mongoScheme, mongoUsername, mongoPassword, mongoAddress)
	}
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		zlog.Fatal(moduleName, "connect", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		zlog.Fatal(moduleName, "ping", zap.Error(err))
	}
	mongoClient = client
	mongoDatabase = models.GetDatabase()
}

func Stop() {
	zlog.Warn(moduleName, "stop ...")
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		zlog.Error(moduleName, "disconnect", zap.Error(err))
		return
	}
}

func Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return mongoClient.Database(mongoDatabase).Collection(name, opts...)
}

func UseSession(fn func(mongo.SessionContext) error) error {
	return mongoClient.UseSession(context.Background(), fn)
}
