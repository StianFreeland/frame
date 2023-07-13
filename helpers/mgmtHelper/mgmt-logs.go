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

func GetMgmtLogs(c *gin.Context, req *protos.GetMgmtLogsReq) {
	checkReqForMgmtLogs(req)
	match := getMatchForMgmtLogs(req)
	facet := getFacetForMgmtLogs(req)
	project := getProjectForMgmtLogs()
	pipe := mongo.Pipeline{match, facet, project}
	coll := mdb.Collection(models.GetCollMgmtLog())
	cursor, err := coll.Aggregate(context.Background(), pipe)
	if err != nil {
		zlog.Error(moduleName, "get mgmt logs", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	var data []protos.MgmtLogData
	if err := cursor.All(context.Background(), &data); err != nil {
		zlog.Error(moduleName, "get mgmt logs", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(data[0]))
}

func checkReqForMgmtLogs(req *protos.GetMgmtLogsReq) {
	if req.Limit < 1 || req.Limit > protos.PageLimit {
		req.Limit = protos.PageLimit
	}
	if req.Page < 1 {
		req.Page = 1
	}
}

func getMatchForMgmtLogs(req *protos.GetMgmtLogsReq) bson.D {
	conditions := bson.D{}
	if req.EventType != "" {
		conditions = append(conditions, bson.E{Key: "event_type", Value: req.EventType})
	}
	if req.EventSubtype != "" {
		conditions = append(conditions, bson.E{Key: "event_subtype", Value: req.EventSubtype})
	}
	if req.OperatorName != "" {
		conditions = append(conditions, bson.E{Key: "operator_name", Value: req.OperatorName})
	}
	if req.TargetName != "" {
		conditions = append(conditions, bson.E{Key: "target_name", Value: req.TargetName})
	}
	if req.BeginTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$gte", req.BeginTime}}})
	}
	if req.EndTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$lte", req.EndTime}}})
	}
	return bson.D{{"$match", conditions}}
}

func getFacetForMgmtLogs(req *protos.GetMgmtLogsReq) bson.D {
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

func getProjectForMgmtLogs() bson.D {
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
