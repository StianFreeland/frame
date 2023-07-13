package menu

import (
	"frame/helpers/menuHelper"
	"frame/protos"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func GetMenus(c *gin.Context) {
	req := &protos.GetMenusReq{}
	if err := c.ShouldBind(req); err != nil {
		zlog.Error(moduleName, "get menus", zap.Error(err))
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.MenuName = strings.TrimSpace(req.MenuName)
	menuHelper.GetMenus(c, req)
}
