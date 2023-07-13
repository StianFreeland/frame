package mgmtHelper

import (
	"context"
	"fmt"
	"frame/comm"
	"frame/models"
	"frame/protos"
	"frame/services/mdb"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func UpdateGroup(c *gin.Context, req *protos.UpdateUserGroupReq) {
	if !checkGroupType(c, req.GroupType) {
		return
	}
	if !checkRootPermsForUpdatingGroup(c, req) {
		return
	}
	if !checkPermsForUpdatingGroup(c, req) {
		return
	}
	if !checkGroupMenus(c, req.GroupMenus) {
		return
	}
	doUpdateGroup(c, req)
}

func checkRootPermsForUpdatingGroup(c *gin.Context, req *protos.UpdateUserGroupReq) bool {
	if req.GroupName == protos.GroupRoot || req.GroupType == protos.GroupTypeRoot {
		zlog.Error(
			moduleName, "check root perms for updating group",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("group_name", req.GroupName),
			zap.Int64("group_type", req.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkPermsForUpdatingGroup(c *gin.Context, req *protos.UpdateUserGroupReq) bool {
	userGroup, ok := GetUserGroup(c, moduleName, req.GroupName)
	if !ok {
		return false
	}
	if (userGroup.GroupType == protos.GroupTypeAdmin || req.GroupType == protos.GroupTypeAdmin) &&
		c.GetString("username") != protos.UserRoot {
		zlog.Error(
			moduleName, "check perms for updating group",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("group_name", req.GroupName),
			zap.Int64("group_type", userGroup.GroupType),
			zap.Int64("req group_type", req.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func doUpdateGroup(c *gin.Context, req *protos.UpdateUserGroupReq) {
	filter := bson.D{{"group_name", req.GroupName}}
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	update := bson.D{
		{"$set",
			bson.D{
				{"group_desc", req.GroupDesc},
				{"group_type", req.GroupType},
				{"group_menus", req.GroupMenus},
				{"last_update_ip", c.ClientIP()},
				{"last_update_time", now.Unix()},
				{"last_update_time_local", localTime},
			},
		},
	}
	coll := mdb.Collection(models.GetCollUserGroup())
	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		zlog.Error(moduleName, "do update group", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(result))
	go createLogForUpdatingGroup(c, req)
}

func createLogForUpdatingGroup(c *gin.Context, req *protos.UpdateUserGroupReq) {
	extraInfo := fmt.Sprintf("group_type:%v", req.GroupType)
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"event_type", protos.EventTypeGroup},
		{"event_subtype", protos.EventSubtypeUpdateGroup},
		{"operator_name", c.GetString("username")},
		{"target_name", req.GroupName},
		{"extra_info", extraInfo},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMgmtLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for updating group", zap.Error(err))
		return
	}
}
