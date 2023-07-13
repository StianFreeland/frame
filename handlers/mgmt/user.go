package mgmt

import (
	"frame/helpers/mgmtHelper"
	"frame/protos"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func CreateUser(c *gin.Context) {
	req := &protos.CreateUserReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "create user", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Password = strings.TrimSpace(req.Password)
	req.GroupName = strings.TrimSpace(req.GroupName)
	zlog.Warn(
		moduleName, "create user",
		zap.String("operator_name", c.GetString("username")),
		zap.String("username", req.Username),
		zap.String("nickname", req.Nickname),
		zap.String("group_name", req.GroupName),
	)
	mgmtHelper.CreateUser(c, req)
}

func GetUsers(c *gin.Context) {
	req := &protos.GetUsersReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "get users", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.GroupName = strings.TrimSpace(req.GroupName)
	mgmtHelper.GetUsers(c, req)
}

func UpdateUser(c *gin.Context) {
	req := &protos.UpdateUserReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "update user", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.GroupName = strings.TrimSpace(req.GroupName)
	zlog.Warn(
		moduleName, "update user",
		zap.String("operator_name", c.GetString("username")),
		zap.String("username", req.Username),
		zap.String("nickname", req.Nickname),
		zap.Int64("account_state", req.AccountState),
		zap.String("group_name", req.GroupName),
	)
	mgmtHelper.UpdateUser(c, req)
}

func DeleteUser(c *gin.Context) {
	req := &protos.DeleteUserReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "delete user", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	zlog.Warn(
		moduleName, "delete user",
		zap.String("operator_name", c.GetString("username")),
		zap.String("username", req.Username),
	)
	mgmtHelper.DeleteUser(c, req)
}
