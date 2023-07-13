package routers

import (
	"frame/handlers/menu"
	"frame/handlers/user"
	"github.com/gin-gonic/gin"
)

func setupMenuRouter(engine *gin.Engine) {
	engine.GET("/menu/menus", user.AuthUser, menu.GetMenus)
}
