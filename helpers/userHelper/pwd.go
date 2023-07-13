package userHelper

import (
	"context"
	"frame/comm"
	"frame/helpers/mgmtHelper"
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

func ChangePwd(c *gin.Context, req *protos.ChangePwdReq) {
	userInfo, ok := mgmtHelper.GetUserInfo(c, moduleName, c.GetString("username"))
	if !ok {
		return
	}
	password, ok := rsaDecryptMsg(c, req.OldPassword, cryptoService.PvtKey)
	if !ok {
		return
	}
	pwd := cryptoService.GetPwdSum(password)
	if !ok {
		return
	}
	if pwd != userInfo.Password {
		zlog.Error(
			moduleName, "change pwd",
			zap.String("error", "password incorrect"),
			zap.String("operator_name", c.GetString("username")),
		)
		c.JSON(http.StatusOK, protos.PasswordIncorrect)
		return
	}
	password, ok = rsaDecryptMsg(c, req.NewPassword, cryptoService.PvtKey)
	if !ok {
		return
	}
	pwd = cryptoService.GetPwdSum(password)
	doChangePwd(c, pwd)
}

func doChangePwd(c *gin.Context, pwd string) {
	filter := bson.D{{"username", c.GetString("username")}}
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
		zlog.Error(moduleName, "change pwd", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(result))
}
