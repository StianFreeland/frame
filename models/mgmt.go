package models

import "frame/protos"

type UserGroup struct {
	GroupName  string           `bson:"group_name"`
	GroupType  protos.GroupType `bson:"group_type"`
	GroupMenus []int64          `bson:"group_menus"`
}

type UserInfo struct {
	Username            string           `bson:"username"`
	Nickname            string           `bson:"nickname"`
	Password            string           `bson:"password"`
	AccountState        protos.CommState `bson:"account_state"`
	GroupName           string           `bson:"group_name"`
	CreateIP            string           `bson:"create_ip"`
	CreateTime          int64            `bson:"create_time"`
	CreateTimeLocal     string           `bson:"create_time_local"`
	LastUpdateIP        string           `bson:"last_update_ip"`
	LastUpdateTime      int64            `bson:"last_update_time"`
	LastUpdateTimeLocal string           `bson:"last_update_time_local"`
	LastLoginIP         string           `bson:"last_login_ip"`
	LastLoginTime       int64            `bson:"last_login_time"`
	LastLoginTimeLocal  string           `bson:"last_login_time_local"`
}

func GetCollUserGroup() string {
	return "user_group"
}

func GetCollUserInfo() string {
	return "user_info"
}

func GetCollLoginLog() string {
	return "login_log"
}

func GetCollMgmtLog() string {
	return "mgmt_log"
}
