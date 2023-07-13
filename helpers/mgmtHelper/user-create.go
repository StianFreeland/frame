package mgmtHelper

import (
	"context"
	"errors"
	"fmt"
	"frame/comm"
	"frame/models"
	"frame/protos"
	"frame/services/cryptoService"
	"frame/services/mdb"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func CreateUser(c *gin.Context, req *protos.CreateUserReq) {
	userMutex.Lock()
	defer userMutex.Unlock()
	userGroupMutex.Lock()
	defer userGroupMutex.Unlock()
	if !checkRootPermsForCreatingUser(c, req) {
		return
	}
	if !checkPermsForCreatingUser(c, req) {
		return
	}
	if !checkUsername(c, req.Username) {
		return
	}
	password, ok := rsaDecryptMsg(c, req.Password, cryptoService.PvtKey)
	if !ok {
		return
	}
	pwd := cryptoService.GetPwdSum(password)
	if !ok {
		return
	}
	if !doCreateUserTransaction(c, req, pwd) {
		return
	}
}

func checkRootPermsForCreatingUser(c *gin.Context, req *protos.CreateUserReq) bool {
	if req.GroupName == protos.GroupRoot {
		zlog.Error(
			moduleName, "check root perms for creating user",
			zap.String("error", "insufficient permissions"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("group_name", req.GroupName),
		)
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		return false
	}
	return true
}

func checkPermsForCreatingUser(c *gin.Context, req *protos.CreateUserReq) bool {
	userGroup, ok := GetUserGroup(c, moduleName, req.GroupName)
	if !ok {
		return false
	}
	if userGroup.GroupType == protos.GroupTypeAdmin && c.GetString("username") != protos.UserRoot {
		zlog.Error(
			moduleName, "check perms for creating user",
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

func checkUsername(c *gin.Context, username string) bool {
	filter := bson.D{{"username", username}}
	coll := mdb.Collection(models.GetCollUserInfo())
	if err := coll.FindOne(context.Background(), filter).Err(); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Error(moduleName, "check username", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return false
		}
		return true
	}
	zlog.Error(
		moduleName, "check username",
		zap.String("error", "username exists"),
		zap.String("operator_name", c.GetString("username")),
		zap.String("username", username),
	)
	c.JSON(http.StatusOK, protos.UsernameExists)
	return false
}

func doCreateUserTransaction(c *gin.Context, req *protos.CreateUserReq, pwd string) bool {
	fn := func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			zlog.Error(moduleName, "do create user transaction", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return err
		}
		result, err := doCreateUser(c, req, pwd)
		if err != nil {
			if err := sc.AbortTransaction(sc); err != nil {
				zlog.Error(moduleName, "do create user transaction", zap.Error(err))
				return err
			}
			return err
		}
		if err := updateGroupMemberCount(c, req.GroupName, 1); err != nil {
			if err := sc.AbortTransaction(sc); err != nil {
				zlog.Error(moduleName, "do create user transaction", zap.Error(err))
				return err
			}
			return err
		}
		if err := sc.CommitTransaction(sc); err != nil {
			zlog.Error(moduleName, "do create user transaction", zap.Error(err))
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

func doCreateUser(c *gin.Context, req *protos.CreateUserReq, pwd string) (*mongo.InsertOneResult, error) {
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"username", req.Username},
		{"nickname", req.Nickname},
		{"password", pwd},
		{"account_state", protos.CommStateEnabled},
		{"group_name", req.GroupName},
		{"create_ip", c.ClientIP()},
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
	result, err := coll.InsertOne(context.Background(), doc)
	if err != nil {
		zlog.Error(moduleName, "do create user", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return nil, err
	}
	go createLogForCreatingUser(c, req)
	return result, nil
}

func createLogForCreatingUser(c *gin.Context, req *protos.CreateUserReq) {
	extraInfo := fmt.Sprintf("group_name:%v", req.GroupName)
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"event_type", protos.EventTypeUser},
		{"event_subtype", protos.EventSubtypeCreateUser},
		{"operator_name", c.GetString("username")},
		{"target_name", req.Username},
		{"extra_info", extraInfo},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollMgmtLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for creating user", zap.Error(err))
		return
	}
}
