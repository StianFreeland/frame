package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(cors.Default())
	setupUserRouter(engine)
	setupLogRouter(engine)
	setupMgmtRouter(engine)
	setupMenuRouter(engine)
	return engine
}
