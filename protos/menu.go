package protos

const RootMenu int64 = 1

type GetMenusReq struct {
	Limit     int64  `form:"limit"`
	Page      int64  `form:"page"`
	MenuID    int64  `form:"menu_id"`
	MenuName  string `form:"menu_name"`
	MenuPID   int64  `form:"menu_pid"`
	BeginTime int64  `form:"begin_time"`
	EndTime   int64  `form:"end_time"`
}

type MenuInfo struct {
	MenuID          int64  `json:"menu_id" bson:"menu_id"`
	MenuName        string `json:"menu_name" bson:"menu_name"`
	MenuPID         int64  `json:"menu_pid" bson:"menu_pid"`
	MenuPriority    int64  `json:"menu_priority" bson:"menu_priority"`
	CreateTime      int64  `json:"create_time" bson:"create_time"`
	CreateTimeLocal string `json:"create_time_local" bson:"create_time_local"`
}

type MenuData struct {
	Count int64      `json:"count"`
	Data  []MenuInfo `json:"data"`
}
