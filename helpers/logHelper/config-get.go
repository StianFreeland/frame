package logHelper

import (
	"frame/protos"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfig(c *gin.Context) {
	data := &protos.LogConfigData{
		Level: zlog.GetLevel(),
	}
	c.JSON(http.StatusOK, protos.Success(data))
}
