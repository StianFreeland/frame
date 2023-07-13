package menuHelper

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

func GetMenus(c *gin.Context, req *protos.GetMenusReq) {
	checkReqForGettingMenus(req)
	match := getMatchForGettingMenus(req)
	facet := getFacetForGettingMenus(req)
	project := getProjectForGettingMenus()
	pipe := mongo.Pipeline{match, facet, project}
	coll := mdb.Collection(models.GetCollMenuInfo())
	cursor, err := coll.Aggregate(context.Background(), pipe)
	if err != nil {
		zlog.Error(moduleName, "get menus", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	var data []protos.MenuData
	if err := cursor.All(context.Background(), &data); err != nil {
		zlog.Error(moduleName, "get menus", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(data[0]))
}

func checkReqForGettingMenus(req *protos.GetMenusReq) {
	if req.Limit < 1 || req.Limit > protos.PageLimit {
		req.Limit = protos.PageLimit
	}
	if req.Page < 1 {
		req.Page = 1
	}
}

func getMatchForGettingMenus(req *protos.GetMenusReq) bson.D {
	conditions := bson.D{}
	if req.MenuID != 0 {
		conditions = append(conditions, bson.E{Key: "menu_id", Value: req.MenuID})
	}
	if req.MenuName != "" {
		conditions = append(conditions, bson.E{Key: "menu_name", Value: req.MenuName})
	}
	if req.MenuPID != 0 {
		conditions = append(conditions, bson.E{Key: "menu_pid", Value: req.MenuPID})
	}
	if req.BeginTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$gte", req.BeginTime}}})
	}
	if req.EndTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$lte", req.EndTime}}})
	}
	return bson.D{{"$match", conditions}}
}

func getFacetForGettingMenus(req *protos.GetMenusReq) bson.D {
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

func getProjectForGettingMenus() bson.D {
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
