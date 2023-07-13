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

func CreateGroup(c *gin.Context) {
	req := &protos.CreateUserGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "create group", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.GroupName = strings.TrimSpace(req.GroupName)
	req.GroupDesc = strings.TrimSpace(req.GroupDesc)
	zlog.Warn(
		moduleName, "create group",
		zap.String("operator_name", c.GetString("username")),
		zap.String("group_name", req.GroupName),
		zap.String("group_desc", req.GroupDesc),
		zap.Int64("group_type", req.GroupType),
		zap.Int64s("group_menus", req.GroupMenus),
	)
	mgmtHelper.CreateGroup(c, req)
}

func GetGroups(c *gin.Context) {
	req := &protos.GetUserGroupsReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "get groups", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.GroupName = strings.TrimSpace(req.GroupName)
	mgmtHelper.GetGroups(c, req)
}

func UpdateGroup(c *gin.Context) {
	req := &protos.UpdateUserGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "update group", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.GroupName = strings.TrimSpace(req.GroupName)
	req.GroupDesc = strings.TrimSpace(req.GroupDesc)
	zlog.Warn(
		moduleName, "update group",
		zap.String("operator_name", c.GetString("username")),
		zap.String("group_name", req.GroupName),
		zap.String("group_desc", req.GroupDesc),
		zap.Int64("group_type", req.GroupType),
		zap.Int64s("group_menus", req.GroupMenus),
	)
	if req.GroupName == protos.GroupRoot {
		req.GroupMenus = nil
	}
	mgmtHelper.UpdateGroup(c, req)
}

func DeleteGroup(c *gin.Context) {
	req := &protos.DeleteUserGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "delete group", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.GroupName = strings.TrimSpace(req.GroupName)
	zlog.Warn(
		moduleName, "delete group",
		zap.String("operator_name", c.GetString("username")),
		zap.String("group_name", req.GroupName),
	)
	mgmtHelper.DeleteGroup(c, req)
}
