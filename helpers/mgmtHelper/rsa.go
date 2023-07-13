package mgmtHelper

import (
	"frame/comm"
	"frame/protos"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func rsaDecryptMsg(c *gin.Context, msg string, pvtKey []byte) (string, bool) {
	decrypted, err := comm.RSADecryptMsg(msg, pvtKey)
	if err != nil {
		zlog.Error(moduleName, "decrypt msg", zap.Error(err))
		c.JSON(http.StatusOK, protos.Error(err))
		return "", false
	}
	return decrypted, true
}
