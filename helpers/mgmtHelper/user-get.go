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

func GetUsers(c *gin.Context, req *protos.GetUsersReq) {
	checkReqForGettingUsers(req)
	match := getMatchForGettingUsers(req)
	facet := getFacetForGettingUsers(req)
	project := getProjectForGettingUsers()
	pipe := mongo.Pipeline{match, facet, project}
	coll := mdb.Collection(models.GetCollUserInfo())
	cursor, err := coll.Aggregate(context.Background(), pipe)
	if err != nil {
		zlog.Error(moduleName, "get users", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	var data []protos.UserData
	if err := cursor.All(context.Background(), &data); err != nil {
		zlog.Error(moduleName, "get users", zap.Error(err))
		c.JSON(http.StatusOK, protos.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, protos.Success(data[0]))
}

func checkReqForGettingUsers(req *protos.GetUsersReq) {
	if req.Limit < 1 || req.Limit > protos.PageLimit {
		req.Limit = protos.PageLimit
	}
	if req.Page < 1 {
		req.Page = 1
	}
}

func getMatchForGettingUsers(req *protos.GetUsersReq) bson.D {
	conditions := bson.D{{"username", bson.D{{"$ne", protos.UserRoot}}}}
	if req.Username != "" {
		conditions = append(conditions, bson.E{Key: "username", Value: req.Username})
	}
	if req.AccountState != protos.CommStateNil {
		conditions = append(conditions, bson.E{Key: "account_state", Value: req.AccountState})
	}
	if req.GroupName != "" {
		conditions = append(conditions, bson.E{Key: "group_name", Value: req.GroupName})
	}
	if req.BeginTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$gte", req.BeginTime}}})
	}
	if req.EndTime != 0 {
		conditions = append(conditions, bson.E{Key: "create_time", Value: bson.D{{"$lte", req.EndTime}}})
	}
	return bson.D{{"$match", conditions}}
}

func getFacetForGettingUsers(req *protos.GetUsersReq) bson.D {
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
		getFacetLookupForGettingUsers(),
		getFacetProjectForGettingUsers(),
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

func getProjectForGettingUsers() bson.D {
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

func getFacetLookupForGettingUsers() bson.D {
	lookup := bson.D{
		{"$lookup",
			bson.D{
				{"from", models.GetCollUserGroup()},
				{"localField", "group_name"},
				{"foreignField", "group_name"},
				{"pipeline",
					bson.A{
						bson.D{
							{"$project",
								bson.D{
									{"_id", 0},
									{"group_type", 1},
									{"group_menus", 1},
								},
							},
						},
					},
				},
				{"as", "groups"},
			},
		},
	}
	return lookup
}

func getFacetProjectForGettingUsers() bson.D {
	project := bson.D{
		{"$project",
			bson.D{
				{"username", "$username"},
				{"nickname", "$nickname"},
				{"account_state", "$account_state"},
				{"group_name", "$group_name"},
				{"group_type", bson.D{{"$first", "$groups.group_type"}}},
				{"group_menus", bson.D{{"$first", "$groups.group_menus"}}},
				{"create_ip", "$create_ip"},
				{"create_time", "$create_time"},
				{"create_time_local", "$create_time_local"},
				{"last_update_ip", "$last_update_ip"},
				{"last_update_time", "$last_update_time"},
				{"last_update_time_local", "$last_update_time_local"},
				{"last_login_ip", "$last_login_ip"},
				{"last_login_time", "$last_login_time"},
				{"last_login_time_local", "$last_login_time_local"},
			},
		},
	}
	return project
}
