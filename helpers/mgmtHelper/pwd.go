package mgmtHelper

import (
	"context"
	"frame/comm"
	"frame/models"
	"frame/protos"
	"frame/services/cryptoService"
	"frame/services/mdb"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func ResetPwd(c *gin.Context, req *protos.ResetPwdReq) {
	if !checkRootPermsForResettingPwd(c, req) {
		return
	}
	if !checkPermsForResettingPwd(c, req) {
		return
	}
	password, ok := rsaDecryptMsg(c, req.Password, cryptoService.PvtKey)
	if !ok {
		return
	}
	pwd := cryptoService.GetPwdSum(password)
	doResetPwd(c, req, pwd)
}

func checkRootPermsForResettingPwd(c *gin.Context, req *protos.ResetPwdReq) bool {
	if req.Username == protos.UserRoot {
		zlog.Error(
			moduleName, "check root perms for resetting pwd",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", req.Username),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkPermsForResettingPwd(c *gin.Context, req *protos.ResetPwdReq) bool {
	userInfo, ok := GetUserInfo(c, moduleName, req.Username)
	if !ok {
		return false
	}
	userGroup, ok := GetUserGroup(c, moduleName, userInfo.GroupName)
	if !ok {
		return false
	}
	if userGroup.GroupType == protos.GroupTypeAdmin && c.GetString("username") != protos.UserRoot {
		zlog.Error(
			moduleName, "check perms for resetting pwd",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", req.Username),
			zap.String("group_name", userInfo.GroupName),
			zap.Int64("group_type", userGroup.GroupType),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func doResetPwd(c *gin.Context, req *protos.ResetPwdReq, pwd string) {
	filter := bson.D{{"username", req.Username}}
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	update := bson.D{
		{"$set",
			bson.D{
				{"password", pwd},
				{"last_update_ip", c.ClientIP()},
				{"last_update_time", now.Unix()},
				{"last_update_time_local", localTime},
			},
		},
	}
	coll := mdb.Collection(models.GetCollUserInfo())
	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		zlog.Error(moduleName, "do reset pwd", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(result))
	go createLogForResettingPwd(c, req)
}

func createLogForResettingPwd(c *gin.Context, req *protos.ResetPwdReq) {
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"event_type", protos.EventTypeUser},
		{"event_subtype", protos.EventSubtypeResetPwd},
		{"operator_name", c.GetString("username")},
		{"target_name", req.Username},
		{"extra_info", ""},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMgmtLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for resetting pwd", zap.Error(err))
		return
	}
}
