package userHelper

import (
	"context"
	"errors"
	"frame/comm"
	"frame/helpers/mgmtHelper"
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

func Login(c *gin.Context, req *protos.UserLoginReq) {
	userInfo, ok := getUserByName(c, req.Username)
	if !ok {
		return
	}
	if userInfo.AccountState != protos.CommStateEnabled {
		c.JSON(http.StatusOK, protos.UserAccountStateError)
		return
	}
	password, err := comm.RSADecryptMsg(req.Password, cryptoService.PvtKey)
	if err != nil {
		c.JSON(http.StatusOK, protos.Error(err))
		return
	}
	pwd := cryptoService.GetPwdSum(password)
	if pwd != userInfo.Password {
		c.JSON(http.StatusOK, protos.UsernameOrPasswordIncorrect)
		return
	}
	userGroup, ok := mgmtHelper.GetUserGroup(c, moduleName, userInfo.GroupName)
	if !ok {
		return
	}
	token, err := generateToken(userInfo.Username, userGroup.GroupType, cryptoService.TokenKey)
	if err != nil {
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	doLogin(c, userInfo, userGroup, token)
}

func getUserByName(c *gin.Context, username string) (*models.UserInfo, bool) {
	filter := bson.D{{"username", username}}
	data := &models.UserInfo{}
	coll := mdb.Collection(models.GetCollUserInfo())
	if err := coll.FindOne(context.Background(), filter).Decode(data); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Error(moduleName, "get user by name", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return nil, false
		}
		c.JSON(http.StatusOK, protos.UserNotExists)
		return nil, false
	}
	return data, true
}

func doLogin(c *gin.Context, userInfo *models.UserInfo, userGroup *models.UserGroup, token string) {
	now := time.Now()
	if !updateLogin(c, userInfo, now) {
		return
	}
	localTime := comm.GetLocalTime(now)
	data := &protos.UserLoginData{
		Token: token,
		UserInfo: protos.UserInfo{
			Username:            userInfo.Username,
			Nickname:            userInfo.Nickname,
			AccountState:        userInfo.AccountState,
			GroupName:           userGroup.GroupName,
			GroupType:           userGroup.GroupType,
			GroupMenus:          userGroup.GroupMenus,
			CreateIP:            userInfo.CreateIP,
			CreateTime:          userInfo.CreateTime,
			CreateTimeLocal:     userInfo.CreateTimeLocal,
			LastUpdateIP:        userInfo.LastUpdateIP,
			LastUpdateTime:      userInfo.LastUpdateTime,
			LastUpdateTimeLocal: userInfo.LastUpdateTimeLocal,
			LastLoginIP:         c.ClientIP(),
			LastLoginTime:       now.Unix(),
			LastLoginTimeLocal:  localTime,
		},
	}
	c.JSON(http.StatusOK, protos.Success(data))
	go createLogForLoggingIn(c, userInfo)
}

func updateLogin(c *gin.Context, userInfo *models.UserInfo, now time.Time) bool {
	filter := bson.D{{"username", userInfo.Username}}
	localTime := comm.GetLocalTime(now)
	update := bson.D{
		{"$set",
			bson.D{
				{"last_login_ip", c.ClientIP()},
				{"last_login_time", now.Unix()},
				{"last_login_time_local", localTime},
			},
		},
	}
	coll := mdb.Collection(models.GetCollUserInfo())
	if _, err := coll.UpdateOne(context.Background(), filter, update); err != nil {
		zlog.Error(moduleName, "update login", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return false
	}
	return true
}

func createLogForLoggingIn(c *gin.Context, userInfo *models.UserInfo) {
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	doc := bson.D{
		{"username", userInfo.Username},
		{"nickname", userInfo.Nickname},
		{"account_state", userInfo.AccountState},
		{"group_name", userInfo.GroupName},
		{"client_ip", c.ClientIP()},
		{"create_time", now.Unix()},
		{"create_time_local", localTime},
	}
	coll := mdb.Collection(models.GetCollLoginLog())
	if _, err := coll.InsertOne(context.Background(), doc); err != nil {
		zlog.Error(moduleName, "create log for logging in", zap.Error(err))
		return
	}
}
