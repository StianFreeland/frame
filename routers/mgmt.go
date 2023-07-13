package routers

import (
	"frame/handlers/mgmt"
	"frame/handlers/user"
	"github.com/gin-gonic/gin"
)

func setupMgmtRouter(engine *gin.Engine) {
	engine.POST("/mgmt/group", user.AuthAdmin, mgmt.CreateGroup)
	engine.GET("/mgmt/groups", user.AuthAdmin, mgmt.GetGroups)
	engine.PUT("/mgmt/group", user.AuthAdmin, mgmt.UpdateGroup)
	engine.DELETE("/mgmt/group", user.AuthAdmin, mgmt.DeleteGroup)
	engine.POST("/mgmt/user", user.AuthAdmin, mgmt.CreateUser)
	engine.GET("/mgmt/users", user.AuthAdmin, mgmt.GetUsers)
	engine.PUT("/mgmt/user", user.AuthAdmin, mgmt.UpdateUser)
	engine.DELETE("/mgmt/user", user.AuthAdmin, mgmt.DeleteUser)
	engine.POST("/mgmt/pwd", user.AuthAdmin, mgmt.ResetPwd)
	engine.GET("/mgmt/login-logs", user.AuthAdmin, mgmt.GetLoginLogs)
	engine.GET("/mgmt/mgmt-logs", user.AuthAdmin, mgmt.GetMgmtLogs)
}
