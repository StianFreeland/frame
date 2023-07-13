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

func UpdateUser(c *gin.Context, req *protos.UpdateUserReq) {
	userMutex.Lock()
	defer userMutex.Unlock()
	userGroupMutex.Lock()
	defer userGroupMutex.Unlock()
	if !checkAccountState(c, req.AccountState) {
		return
	}
	if !checkRootPermsForUpdatingUser(c, req) {
		return
	}
	if !checkPermsForUpdatingUser(c, req) {
		return
	}
	doUpdateUser(c, req)
}

func checkAccountState(c *gin.Context, accountState protos.CommState) bool {
	if accountState < protos.CommStateEnabled || accountState > protos.CommStateDisabled {
		zlog.Error(
			moduleName, "check account state",
			zap.String("error", "invalid user account state"),
			zap.Int64("account", accountState),
		)
		c.JSON(http.StatusOK, protos.InvalidUserAccountState)
		return false
	}
	return true
}

func checkRootPermsForUpdatingUser(c *gin.Context, req *protos.UpdateUserReq) bool {
	if req.Username == protos.UserRoot || req.GroupName == protos.GroupRoot {
		zlog.Error(
			moduleName, "check root perms for updating user",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", req.Username),
			zap.String("group_name", req.GroupName),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkPermsForUpdatingUser(c *gin.Context, req *protos.UpdateUserReq) bool {
	userInfo, ok := GetUserInfo(c, moduleName, req.Username)
	if !ok {
		return false
	}
	if c.GetString("username") == protos.UserRoot {
		return true
	}
	if req.Username == c.GetString("username") {
		return checkSelfPermsForUpdatingUser(c, req, userInfo)
	}
	if !checkSourcePermsForUpdatingUser(c, userInfo) {
		return false
	}
	if req.GroupName != userInfo.GroupName {
		return checkTargetPermsForUpdatingUser(c, req)
	}
	return true
}

func doUpdateUser(c *gin.Context, req *protos.UpdateUserReq) {
	filter := bson.D{{"username", req.Username}}
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	update := bson.D{
		{"$set",
			bson.D{
				{"nickname", req.Nickname},
				{"account_state", req.AccountState},
				{"group_name", req.GroupName},
				{"last_update_ip", c.ClientIP()},
				{"last_update_time", now.Unix()},
				{"last_update_time_local", localTime},
			},
		},
	}
	coll := mdb.Collection(models.GetCollUserInfo())
	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		zlog.Error(moduleName, "do update user", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(result))
	go createLogForUpdatingUser(c, req)
}

func createLogForUpdatingUser(c *gin.Context, req *protos.UpdateUserReq) {
	extraInfo := fmt.Sprintf("account_state:%v, group_name:%v", req.AccountState, req.GroupName)
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"event_type", protos.EventTypeUser},
		{"event_subtype", protos.EventSubtypeUpdateUser},
		{"operator_name", c.GetString("username")},
		{"target_name", req.Username},
		{"extra_info", extraInfo},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMgmtLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for updating user", zap.Error(err))
		return
	}
}

func checkSelfPermsForUpdatingUser(c *gin.Context, req *protos.UpdateUserReq, userInfo *models.UserInfo) bool {
	if req.AccountState != userInfo.AccountState || req.GroupName != userInfo.GroupName {
		zlog.Error(
			moduleName, "check self perms for updating user",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", userInfo.Username),
			zap.Int64("account_state", userInfo.AccountState),
			zap.Int64("req account_state", req.AccountState),
			zap.String("group_name", userInfo.GroupName),
			zap.String("req group_name", req.GroupName),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkSourcePermsForUpdatingUser(c *gin.Context, userInfo *models.UserInfo) bool {
	userGroup, ok := GetUserGroup(c, moduleName, userInfo.GroupName)
	if !ok {
		return false
	}
	if userGroup.GroupType == protos.GroupTypeAdmin {
		zlog.Error(
			moduleName, "check source perms for updating user",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", userInfo.Username),
			zap.String("group_name", userInfo.GroupName),
			zap.Int64("group_type", userGroup.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkTargetPermsForUpdatingUser(c *gin.Context, req *protos.UpdateUserReq) bool {
	userGroup, ok := GetUserGroup(c, moduleName, req.GroupName)
	if !ok {
		return false
	}
	if userGroup.GroupType == protos.GroupTypeAdmin {
		zlog.Error(
			moduleName, "check target perms for updating user",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", req.Username),
			zap.String("group_name", req.GroupName),
			zap.Int64("group_type", userGroup.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}
