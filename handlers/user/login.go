package user

import (
	"frame/helpers/userHelper"
	"frame/protos"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Login(c *gin.Context) {
	req := &protos.UserLoginReq{}
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, protos.InvalidReqParams)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	userHelper.Login(c, req)
}
