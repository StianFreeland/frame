package mgmtHelper

import (
	"context"
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

func GetGroups(c *gin.Context, req *protos.GetUserGroupsReq) {
	checkReqForGettingGroups(req)
	match := getMatchForGettingGroups(req)
	facet := getFacetForGettingGroups(req)
	project := getProjectForGettingGroups()
	pipe := mongo.Pipeline{match, facet, project}
	coll := mdb.Collection(models.GetCollUserGroup())
	cursor, err := coll.Aggregate(context.Background(), pipe)
	if err != nil {
		zlog.Error(moduleName, "get groups", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	var data []protos.UserGroupData
	if err := cursor.All(context.Background(), &data); err != nil {
		zlog.Error(moduleName, "get groups", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(data[0]))
}

func checkReqForGettingGroups(req *protos.GetUserGroupsReq) {
	if req.Limit < 1 || req.Limit > protos.PageLimit {
		req.Limit = protos.PageLimit
	}
	if req.Page < 1 {
		req.Page = 1
	}
}

func getMatchForGettingGroups(req *protos.GetUserGroupsReq) bson.D {
	conditions := bson.D{{"group_name", bson.D{{"$ne", protos.GroupRoot}}}}
	if req.GroupName != "" {
		conditions = append(conditions, bson.E{Key: "group_name", Value: req.GroupName})
	}
	if req.GroupType != protos.GroupTypeNil {
		conditions = append(conditions, bson.E{Key: "group_type", Value: req.GroupType})
	}
	if req.BeginTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$gte", req.BeginTime}}})
	}
	if req.EndTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$lte", req.EndTime}}})
	}
	return bson.D{{"$match", conditions}}
}

func getFacetForGettingGroups(req *protos.GetUserGroupsReq) bson.D {
	meta := bson.A{
		bson.D{
			{"$group",
				bson.D{
					{"_id", nil},
					{"count", bson.D{{"$sum", 1}}},
				},
			},
		},
	}
	data := bson.A{
		bson.D{{"$sort", bson.D{{"_id", -1}}}},
		bson.D{{"$skip", (req.Page - 1) * req.Limit}},
		bson.D{{"$limit", req.Limit}},
	}
	facet := bson.D{
		{"$facet",
			bson.D{
				{"meta", meta},
				{"data", data},
			},
		},
	}
	return facet
}

func getProjectForGettingGroups() bson.D {
	project := bson.D{
		{"$project",
			bson.D{
				{"count", bson.D{{"$first", "$meta.count"}}},
				{"data", "$data"},
			},
		},
	}
	return project
}
