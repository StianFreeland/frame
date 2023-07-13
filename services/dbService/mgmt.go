package dbService

import (
	"context"
	"errors"
	"frame/comm"
	"frame/models"
	"frame/protos"
	"frame/services/cryptoService"
	"frame/services/mdb"
	"frame/services/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

func initMgmt() {
	zlog.Info(moduleName, "init mgmt ...")
	createIndexesForMgmtLog()
	createIndexesForLoginLog()
	createIndexesForUserGroup()
	createIndexesForUserInfo()
	initMongoForUserGroup()
	initMongoForUserInfo()
}

func createIndexesForMgmtLog() {
	zlog.Debug(moduleName, "create indexes for mgmt log ...")
	doCreateIndexes(models.GetCollMgmtLog(), bson.D{{"event_type", -1}}, false)
	doCreateIndexes(models.GetCollMgmtLog(), bson.D{{"event_subtype", -1}}, false)
	doCreateIndexes(models.GetCollMgmtLog(), bson.D{{"operator_name", -1}}, false)
	doCreateIndexes(models.GetCollMgmtLog(), bson.D{{"target_name", -1}}, false)
	doCreateIndexes(models.GetCollMgmtLog(), bson.D{{"create_time", -1}}, false)
}

func createIndexesForLoginLog() {
	zlog.Debug(moduleName, "create indexes for login log ...")
	doCreateIndexes(models.GetCollLoginLog(), bson.D{{"username", -1}}, false)
	doCreateIndexes(models.GetCollLoginLog(), bson.D{{"create_time", -1}}, false)
}

func createIndexesForUserGroup() {
	zlog.Debug(moduleName, "create indexes for user group ...")
	doCreateIndexes(models.GetCollUserGroup(), bson.D{{"group_name", -1}}, true)
	doCreateIndexes(models.GetCollUserGroup(), bson.D{{"group_type", -1}}, false)
	doCreateIndexes(models.GetCollUserGroup(), bson.D{{"create_time", -1}}, false)
}

func createIndexesForUserInfo() {
	zlog.Debug(moduleName, "create indexes for user info ...")
	doCreateIndexes(models.GetCollUserInfo(), bson.D{{"username", -1}}, true)
	doCreateIndexes(models.GetCollUserInfo(), bson.D{{"account_state", -1}}, false)
	doCreateIndexes(models.GetCollUserInfo(), bson.D{{"group_name", -1}}, false)
	doCreateIndexes(models.GetCollUserInfo(), bson.D{{"create_time", -1}}, false)
}

func initMongoForUserGroup() {
	zlog.Debug(moduleName, "init mongo for user group ...")
	createUserGroup(protos.GroupRoot, "root", protos.GroupTypeRoot)
}

func initMongoForUserInfo() {
	zlog.Debug(moduleName, "init mongo for user info ...")
	createUserInfo(protos.UserRoot, "root", cryptoService.GetRootPwdSum(), protos.GroupRoot)
}

func createUserGroup(groupName string, groupDesc string, groupType protos.GroupType) {
	if !checkUserGroup(groupName) {
		return
	}
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"group_name", groupName},
		{"group_desc", groupDesc},
		{"group_type", groupType},
		{"group_menus", []int64{}},
		{"create_ip", "127.0.0.1"},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
		{"last_update_ip", ""},
		{"last_update_time", int64(0)},
		{"last_update_time_local", ""},
	}
	coll := mdb.Collection(models.GetCollUserGroup())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Fatal(moduleName, "create user group", zap.Error(err))
	}
}

func createUserInfo(username string, nickname string, password string, groupName string) {
	if !checkUserInfo(username) {
		return
	}
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"username", username},
		{"nickname", nickname},
		{"password", password},
		{"account_state", protos.CommStateEnabled},
		{"group_name", groupName},
		{"create_ip", "127.0.0.1"},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
		{"last_update_ip", ""},
		{"last_update_time", int64(0)},
		{"last_update_time_local", ""},
		{"last_login_ip", ""},
		{"last_login_time", int64(0)},
		{"last_login_time_local", ""},
	}
	coll := mdb.Collection(models.GetCollUserInfo())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Fatal(moduleName, "create user info", zap.Error(err))
	}
}

func checkUserGroup(groupName string) bool {
	filter := bson.D{{"group_name", groupName}}
	coll := mdb.Collection(models.GetCollUserGroup())
	if err := coll.FindOne(context.Background(), filter).Err(); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Fatal(moduleName, "check user group", zap.Error(err))
		}
		return true
	}
	return false
}

func checkUserInfo(username string) bool {
	filter := bson.D{{"username", username}}
	coll := mdb.Collection(models.GetCollUserInfo())
	if err := coll.FindOne(context.Background(), filter).Err(); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Fatal(moduleName, "check user info", zap.Error(err))
		}
		return true
	}
	return false
}
