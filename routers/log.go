package routers

import (
	"frame/handlers/log"
	"frame/handlers/user"
	"github.com/gin-gonic/gin"
)

func setupLogRouter(engine *gin.Engine) {
	engine.GET("/log/config", user.AuthRoot, log.GetConfig)
	engine.PUT("/log/config", user.AuthRoot, log.UpdateConfig)
}
