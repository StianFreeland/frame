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

func GetUserInfo(c *gin.Context, prefix string, username string) (*models.UserInfo, bool) {
	filter := bson.D{{"username", username}}
	data := &models.UserInfo{}
	coll := mdb.Collection(models.GetCollUserInfo())
	if err := coll.FindOne(context.Background(), filter).Decode(data); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			zlog.Error(prefix, "get user info", zap.Error(err))
			c.JSON(http.StatusOK, protos.InternalServerError)
			return nil, false
		}
		zlog.Error(
			prefix, "get user info",
			zap.String("error", "user not exists"),
			zap.String("operator_name", c.GetString("username")),
			zap.String("username", username),
		)
		c.JSON(http.StatusOK, protos.UserNotExists)
		return nil, false
	}
	return data, true
}
