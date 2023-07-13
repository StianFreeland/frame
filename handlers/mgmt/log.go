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

func GetLoginLogs(c *gin.Context) {
	req := &protos.GetLoginLogsReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "get login logs", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	mgmtHelper.GetLoginLogs(c, req)
}

func GetMgmtLogs(c *gin.Context) {
	req := &protos.GetMgmtLogsReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "get mgmt logs", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.EventType = strings.TrimSpace(req.EventType)
	req.EventSubtype = strings.TrimSpace(req.EventSubtype)
	req.OperatorName = strings.TrimSpace(req.OperatorName)
	req.TargetName = strings.TrimSpace(req.TargetName)
	mgmtHelper.GetMgmtLogs(c, req)
}
