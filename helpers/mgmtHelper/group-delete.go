package mgmtHelper

import (
	"context"
	"errors"
	"frame/comm"
	"frame/models"
	"frame/protos"
	"frame/services/mdb"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func DeleteGroup(c *gin.Context, req *protos.DeleteUserGroupReq) {
	userGroupMutex.Lock()
	defer userGroupMutex.Unlock()
	if !checkRootPermsForDeletingGroup(c, req) {
		return
	}
	if !checkPermsForDeletingGroup(c, req) {
		return
	}
	if !checkGroupUsers(c, req.GroupName) {
		return
	}
	filter := bson.D{{"group_name", req.GroupName}}
	coll := mdb.Collection(models.GetCollUserGroup())
	result, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		zlog.Error(moduleName, "delete group", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(result))
	go createLogForDeletingGroup(c, req)
}

func checkRootPermsForDeletingGroup(c *gin.Context, req *protos.DeleteUserGroupReq) bool {
	if req.GroupName == protos.GroupRoot {
		zlog.Error(
			moduleName, "check root perms for deleting group",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("group_name", req.GroupName),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkPermsForDeletingGroup(c *gin.Context, req *protos.DeleteUserGroupReq) bool {
	userGroup, ok := GetUserGroup(c, moduleName, req.GroupName)
	if !ok {
		return false
	}
	if userGroup.GroupType == protos.GroupTypeAdmin && c.GetString("username") != protos.UserRoot {
		zlog.Error(
			moduleName, "check perms for deleting group",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("group_name", userGroup.GroupName),
			zap.Int64("group_type", userGroup.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkGroupUsers(c *gin.Context, groupName string) bool {
	filter := bson.D{{"group_name", groupName}}
	coll := mdb.Collection(models.GetCollUserInfo())
	if err := coll.FindOne(context.Background(), filter).Err(); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Error(moduleName, "check group users", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return false
		}
		return true
	}
	zlog.Error(
		moduleName, "check group users",
		zap.String("error", "user group not empty"),
		zap.String("operator_name", c.GetString("username")),
		zap.String("group_name", groupName),
	)
	c.JSON(http.StatusOK, protos.UserGroupNotEmpty)
	return false
}

func createLogForDeletingGroup(c *gin.Context, req *protos.DeleteUserGroupReq) {
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"event_type", protos.EventTypeGroup},
		{"event_subtype", protos.EventSubtypeDeleteGroup},
		{"operator_name", c.GetString("username")},
		{"target_name", req.GroupName},
		{"extra_info", ""},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMgmtLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for deleting group", zap.Error(err))
		return
	}
}
