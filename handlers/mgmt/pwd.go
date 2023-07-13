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

func ResetPwd(c *gin.Context) {
	req := &protos.ResetPwdReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "reset pwd", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	zlog.Warn(
		moduleName, "reset pwd",
		zap.String("operator_name", c.GetString("username")),
		zap.String("username", req.Username),
	)
	mgmtHelper.ResetPwd(c, req)
}
