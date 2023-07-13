package mgmtHelper

import (
	"context"
	"errors"
	"fmt"
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

func CreateGroup(c *gin.Context, req *protos.CreateUserGroupReq) {
	groupMutex.Lock()
	defer groupMutex.Unlock()
	if !checkGroupType(c, req.GroupType) {
		return
	}
	if !checkRootPermsForCreatingGroup(c, req) {
		return
	}
	if !checkPermsForCreatingGroup(c, req) {
		return
	}
	if !checkGroupName(c, req.GroupName) {
		return
	}
	if !checkGroupMenus(c, req.GroupMenus) {
		return
	}
	if !doCreateGroup(c, req) {
		return
	}
}

func checkRootPermsForCreatingGroup(c *gin.Context, req *protos.CreateUserGroupReq) bool {
	if req.GroupType == protos.GroupTypeRoot {
		zlog.Error(
			moduleName, "check root perms for creating group",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.Int64("group_type", req.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkPermsForCreatingGroup(c *gin.Context, req *protos.CreateUserGroupReq) bool {
	if req.GroupType == protos.GroupTypeAdmin && c.GetString("username") != protos.UserRoot {
		zlog.Error(
			moduleName, "check perms for creating group",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.Int64("group_type", req.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkGroupName(c *gin.Context, groupName string) bool {
	filter := bson.D{{"group_name", groupName}}
	coll := mdb.Collection(models.GetCollUserGroup())
	if err := coll.FindOne(context.Background(), filter).Err(); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Error(moduleName, "check group name", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return false
		}
		return true
	}
	zlog.Error(
		moduleName, "check group name",
		zap.String("error", "user group name exists"),
		zap.String("operator_name", c.GetString("username")),
		zap.String("group_name", groupName),
	)
	c.JSON(http.StatusOK, protos.UserGroupNameExists)
	return false
}

func doCreateGroup(c *gin.Context, req *protos.CreateUserGroupReq) bool {
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"group_name", req.GroupName},
		{"group_desc", req.GroupDesc},
		{"group_type", req.GroupType},
		{"member_count", int64(0)},
		{"group_menus", req.GroupMenus},
		{"create_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
		{"last_update_ip", ""},
		{"last_update_time", int64(0)},
		{"last_update_time_local", ""},
	}
	coll := mdb.Collection(models.GetCollUserGroup())
	result, err := coll.InsertOne(context.Background(), doc)
	if err != nil {
		zlog.Error(moduleName, "do create group", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return false
	}
	c.JSON(http.StatusOK, protos.Success(result))
	go createLogForCreatingGroup(c, req)
	return true
}

func createLogForCreatingGroup(c *gin.Context, req *protos.CreateUserGroupReq) {
	extraInfo := fmt.Sprintf("group_type:%v", req.GroupType)
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"event_type", protos.EventTypeGroup},
		{"event_subtype", protos.EventSubtypeCreateGroup},
		{"operator_name", c.GetString("username")},
		{"target_name", req.GroupName},
		{"extra_info", extraInfo},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMgmtLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for creating group", zap.Error(err))
		return
	}
}
