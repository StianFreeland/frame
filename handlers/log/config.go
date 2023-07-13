package log

import (
	"frame/helpers/logHelper"
	"frame/protos"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func GetConfig(c *gin.Context) {
	logHelper.GetConfig(c)
}

func UpdateConfig(c *gin.Context) {
	req := &protos.UpdateLogConfigReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "update config", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	zlog.Warn(
		moduleName, "update config",
		zap.String("operator_name", c.GetString("username")),
		zap.Int64("level", req.Level),
	)
	logHelper.UpdateConfig(c, req)
}
