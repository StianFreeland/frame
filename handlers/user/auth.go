package user

import (
	"frame/helpers/userHelper"
	"frame/protos"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthRoot(c *gin.Context) {
	token, _ := c.GetQuery("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, protos.TokenEmpty)
		c.Abort()
		return
	}
	userHelper.AuthRoot(c, token)
}

func AuthAdmin(c *gin.Context) {
	token, _ := c.GetQuery("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, protos.TokenEmpty)
		c.Abort()
		return
	}
	userHelper.AuthAdmin(c, token)
}

func AuthUser(c *gin.Context) {
	token, _ := c.GetQuery("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, protos.TokenEmpty)
		c.Abort()
		return
	}
	userHelper.AuthUser(c, token)
}
