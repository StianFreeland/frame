package mgmtHelper

import (
	"context"
	"errors"
	"frame/models"
	"frame/protos"
	"frame/services/mdb"
	"frame/services/zlog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
)

func GetUserGroup(c *gin.Context, prefix string, groupName string) (*models.UserGroup, bool) {
	filter := bson.D{{"group_name", groupName}}
	data := &models.UserGroup{}
	coll := mdb.Collection(models.GetCollUserGroup())
	if err := coll.FindOne(context.Background(), filter).Decode(data); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Error(prefix, "get user group", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return nil, false
		}
		zlog.Error(
			prefix, "get user group",
			zap.String("error", "user group not exists"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("group_name", groupName),
		)
		c.JSON(http.StatusOK, protos.UserGroupNotExists)
		return nil, false
	}
	return data, true
}
