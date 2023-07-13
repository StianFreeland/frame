package logHelper

import (
	"frame/protos"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateConfig(c *gin.Context, req *protos.UpdateLogConfigReq) {
	zlog.UpdateLevel(req.Level)
	c.JSON(http.StatusOK, protos.Success())
}
