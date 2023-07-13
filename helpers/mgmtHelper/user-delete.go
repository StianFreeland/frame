package mgmtHelper

import (
	"context"
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

func DeleteUser(c *gin.Context, req *protos.DeleteUserReq) {
	if !checkRootPermsForDeletingUser(c, req) {
		return
	}
	userInfo, ok := GetUserInfo(c, moduleName, req.Username)
	if !ok {
		return
	}
	if !checkPermsForDeletingUser(c, userInfo) {
		return
	}
	if !doDeleteUserTransaction(c, req, userInfo) {
		return
	}
}

func checkRootPermsForDeletingUser(c *gin.Context, req *protos.DeleteUserReq) bool {
	if req.Username == protos.UserRoot {
		zlog.Error(
			moduleName, "check root perms for deleting user",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", req.Username),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkPermsForDeletingUser(c *gin.Context, userInfo *models.UserInfo) bool {
	userGroup, ok := GetUserGroup(c, moduleName, userInfo.GroupName)
	if !ok {
		return false
	}
	if userGroup.GroupType == protos.GroupTypeAdmin && c.GetString("username") != protos.UserRoot {
		zlog.Error(
			moduleName, "check perms for deleting user",
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

func doDeleteUserTransaction(c *gin.Context, req *protos.DeleteUserReq, userInfo *models.UserInfo) bool {
	fn := func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			zlog.Error(moduleName, "do delete user transaction", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return err
		}
		result, err := doDeleteUser(c, req)
		if err != nil {
			if err := sc.AbortTransaction(sc); err != nil {
				zlog.Error(moduleName, "do delete user transaction", zap.Error(err))
				return err
			}
			return err
		}
		if err := updateGroupMemberCount(c, userInfo.GroupName, -1); err != nil {
			if err := sc.AbortTransaction(sc); err != nil {
				zlog.Error(moduleName, "do delete user transaction", zap.Error(err))
				return err
			}
			return err
		}
		if err := sc.CommitTransaction(sc); err != nil {
			zlog.Error(moduleName, "do delete user transaction", zap.Error(err))
			return err
		}
		c.JSON(http.StatusOK, protos.Success(result))
		return nil
	}
	if err := mdb.UseSession(fn); err != nil {
		return false
	}
	return true
}

func doDeleteUser(c *gin.Context, req *protos.DeleteUserReq) (*mongo.DeleteResult, error) {
	filter := bson.D{{"username", req.Username}}
	coll := mdb.Collection(models.GetCollUserInfo())
	result, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		zlog.Error(moduleName, "do delete user", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return nil, err
	}
	go createLogForDeletingUser(c, req)
	return result, nil
}

func createLogForDeletingUser(c *gin.Context, req *protos.DeleteUserReq) {
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"event_type", protos.EventTypeUser},
		{"event_subtype", protos.EventSubtypeDeleteUser},
		{"operator_name", c.GetString("username")},
		{"target_name", req.Username},
		{"extra_info", ""},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMgmtLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for deleting user", zap.Error(err))
		return
	}
}
