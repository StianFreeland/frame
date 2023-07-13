package mgmtHelper

import (
	"context"
	"frame/comm"
	"frame/models"
	"frame/protos"
	"frame/services/mdb"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var groupMutex sync.Mutex

func checkGroupType(c *gin.Context, groupType protos.GroupType) bool {
	if groupType < protos.GroupTypeRoot || groupType > protos.GroupTypeOperator {
		zlog.Error(
			moduleName, "check group type",
			zap.String("error", "invalid user group type"),
			zap.String("operator_name", c.GetString("username")),
			zap.Int64("group_type", groupType),
		)
		c.JSON(http.StatusOK, protos.InvalidUserGroupType)
		return false
	}
	return true
}

func checkGroupMenus(c *gin.Context, groupMenus []int64) bool {
	if len(groupMenus) == 0 {
		return true
	}
	menus, ok := getMenus(c)
	if !ok {
		return false
	}
	for _, menu := range groupMenus {
		if _, ok := menus[menu]; !ok {
			zlog.Error(
				moduleName, "check group menus",
				zap.String("error", "menu not exists"),
				zap.String("operator_name", c.GetString("username")),
				zap.Int64("menu", menu),
			)
			c.JSON(http.StatusOK, protos.MenuNotExists)
			return false
		}
	}
	return true
}

func updateGroupMemberCount(c *gin.Context, groupName string, delta int64) error {
	filter := bson.D{{"group_name", groupName}}
	now := time.Now()
	localTime := comm.GetLocalTime(now)
	update := bson.D{
		{"$set",
			bson.D{
				{"last_update_ip", c.ClientIP()},
				{"last_update_time", now.Unix()},
				{"last_update_time_local", localTime},
			},
		},
		{"$inc",
			bson.D{{"member_count", delta}},
		},
	}
	coll := mdb.Collection(models.GetCollUserGroup())
	if _, err := coll.UpdateOne(context.Background(), filter, update); err != nil {
		zlog.Error(moduleName, "update group member count", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return err
	}
	return nil
}

func getMenus(c *gin.Context) (map[int64]bool, bool) {
	opt := options.Find().SetLimit(comm.FindLimit)
	coll := mdb.Collection(models.GetCollMenuInfo())
	cursor, err := coll.Find(context.Background(), bson.D{}, opt)
	if err != nil {
		zlog.Error(moduleName, "get menus", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return nil, false
	}
	var data []models.MenuInfo
	if err := cursor.All(context.Background(), &data); err != nil {
		zlog.Error(moduleName, "get menus", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return nil, false
	}
	menus := make(map[int64]bool)
	for _, d := range data {
		menus[d.MenuID] = true
	}
	return menus, true
}
