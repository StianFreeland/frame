package userHelper

import (
	"errors"
	"frame/protos"
	"frame/services/cryptoService"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func AuthRoot(c *gin.Context, token string) {
	username, groupType, err := parseToken(token, cryptoService.TokenKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusOK, protos.TokenExpired)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, protos.Error(err))
		c.Abort()
		return
	}
	if username != protos.UserRoot || groupType != protos.GroupTypeRoot {
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		c.Abort()
		return
	}
	c.Set("username", username)
	c.Next()
}

func AuthAdmin(c *gin.Context, token string) {
	username, groupType, err := parseToken(token, cryptoService.TokenKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusOK, protos.TokenExpired)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, protos.Error(err))
		c.Abort()
		return
	}
	if groupType != protos.GroupTypeRoot && groupType != protos.GroupTypeAdmin {
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		c.Abort()
		return
	}
	c.Set("username", username)
	c.Next()
}

func AuthUser(c *gin.Context, token string) {
	username, groupType, err := parseToken(token, cryptoService.TokenKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusOK, protos.TokenExpired)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, protos.Error(err))
		c.Abort()
		return
	}
	if groupType != protos.GroupTypeRoot &&
		groupType != protos.GroupTypeAdmin &&
		groupType != protos.GroupTypeOperator {
		c.JSON(http.StatusOK, protos.InsufficientPermissions)
		c.Abort()
		return
	}
	c.Set("username", username)
	c.Next()
}
