package dbService

import (
	"context"
	"errors"
	"frame/comm"
	"frame/models"
	"frame/protos"
	"frame/services/mdb"
	"frame/services/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

func initMenu() {
	zlog.Info(moduleName, "init menu ...")
	createIndexesForMenuInfo()
	initMongoForMenuInfo()
}

func createIndexesForMenuInfo() {
	zlog.Debug(moduleName, "create indexes for menu info ...")
	doCreateIndexes(models.GetCollMenuInfo(), bson.D{{"menu_id", -1}}, true)
	doCreateIndexes(models.GetCollMenuInfo(), bson.D{{"menu_name", -1}}, true)
	doCreateIndexes(models.GetCollMenuInfo(), bson.D{{"menu_pid", -1}}, false)
	doCreateIndexes(models.GetCollMenuInfo(), bson.D{{"create_time", -1}}, false)
}

func initMongoForMenuInfo() {
	zlog.Debug(moduleName, "init mongo for menu info ...")
	createMenuForManagement()
}

func createMenuForManagement() {
	createMenu(10100, "Management", protos.RootMenu, 1)
	createMenu(10101, "User Group", 10100, 1)
	createMenu(10102, "User Info", 10100, 2)
	createMenu(10103, "Login Logs", 10100, 3)
	createMenu(10104, "Management Logs", 10100, 4)
}

func createMenu(menuID int64, menuName string, menuPID int64, menuPriority int64) {
	if !checkMenuID(menuID) {
		return
	}
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"menu_id", menuID},
		{"menu_name", menuName},
		{"menu_pid", menuPID},
		{"menu_priority", menuPriority},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMenuInfo())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Fatal(moduleName, "create menu", zap.Error(err))
	}
}

func checkMenuID(menuID int64) bool {
	filter := bson.D{{"menu_id", menuID}}
	coll := mdb.Collection(models.GetCollMenuInfo())
	if err := coll.FindOne(context.Background(), filter).Err(); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Fatal(moduleName, "check menu id", zap.Error(err))
		}
		return true
	}
	return false
}
