package dbService

import (
	"context"
	"frame/services/mdb"
	"frame/services/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func doCreateIndexes(collection string, keys bson.D, unique bool) {
	coll := mdb.Collection(collection)
	im := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(unique),
	}
	if _, err := coll.Indexes().CreateOne(context.Background(), im); err != nil {
		zlog.Fatal(moduleName, "do create indexes", zap.Error(err))
	}
}
