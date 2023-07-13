package user

import (
	"frame/helpers/userHelper"
	"frame/protos"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func ChangePwd(c *gin.Context) {
	req := &protos.ChangePwdReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "change pwd", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	zlog.Warn(moduleName, "change pwd", zap.String("operator_name", c.GetString("username")))
	userHelper.ChangePwd(c, req)
}
