package routers

import (
	"frame/handlers/user"
	"github.com/gin-gonic/gin"
)

func setupUserRouter(engine *gin.Engine) {
	engine.POST("/user/login", user.Login)
	engine.POST("/user/pwd", user.AuthUser, user.ChangePwd)
}
