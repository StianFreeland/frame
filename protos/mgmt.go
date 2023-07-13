package protos

type GroupType = int64
type EventType = string
type EventSubtype = string

const (
	GroupRoot string = "root"
	UserRoot  string = "root"
)

const (
	GroupTypeNil GroupType = iota
	GroupTypeRoot
	GroupTypeAdmin
	GroupTypeOperator
)

const (
	EventTypeGroup EventType = "group"
	EventTypeUser  EventType = "user"
)

const (
	EventSubtypeCreateGroup EventSubtype = "create_group"
	EventSubtypeUpdateGroup EventSubtype = "update_group"
	EventSubtypeDeleteGroup EventSubtype = "delete_group"
	EventSubtypeCreateUser  EventSubtype = "create_user"
	EventSubtypeUpdateUser  EventSubtype = "update_user"
	EventSubtypeDeleteUser  EventSubtype = "delete_user"
	EventSubtypeResetPwd    EventSubtype = "reset_pwd"
)

type CreateUserGroupReq struct {
	GroupName  string    `json:"group_name" form:"group_name" binding:"required"`
	GroupDesc  string    `json:"group_desc" form:"group_desc" binding:"required"`
	GroupType  GroupType `json:"group_type" form:"group_type" binding:"required"`
	GroupMenus []int64   `json:"group_menus" form:"group_menus" binding:"required"`
}

type GetUserGroupsReq struct {
	Limit     int64     `form:"limit"`
	Page      int64     `form:"page"`
	GroupName string    `form:"group_name"`
	GroupType GroupType `form:"group_type"`
	BeginTime int64     `form:"begin_time"`
	EndTime   int64     `form:"end_time"`
}

type UpdateUserGroupReq struct {
	GroupName  string    `json:"group_name" form:"group_name" binding:"required"`
	GroupDesc  string    `json:"group_desc" form:"group_desc"`
	GroupType  GroupType `json:"group_type" form:"group_type" binding:"required"`
	GroupMenus []int64   `json:"group_menus" form:"group_menus" binding:"required"`
}

type DeleteUserGroupReq struct {
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
}

type UserGroup struct {
	GroupName           string    `json:"group_name" bson:"group_name"`
	GroupDesc           string    `json:"group_desc" bson:"group_desc"`
	GroupType           GroupType `json:"group_type" bson:"group_type"`
	MemberCount         int64     `json:"member_count" bson:"member_count"`
	GroupMenus          []int64   `json:"group_menus" bson:"group_menus"`
	CreateIP            string    `json:"create_ip" bson:"create_ip"`
	CreateTime          int64     `json:"create_time" bson:"create_time"`
	CreateTimeLocal     string    `json:"create_time_local" bson:"create_time_local"`
	LastUpdateIP        string    `json:"last_update_ip" bson:"last_update_ip"`
	LastUpdateTime      int64     `json:"last_update_time" bson:"last_update_time"`
	LastUpdateTimeLocal string    `json:"last_update_time_local" bson:"last_update_time_local"`
}

type UserGroupData struct {
	Count int64       `json:"count"`
	Data  []UserGroup `json:"data"`
}

type CreateUserReq struct {
	Username  string `json:"username" form:"username" binding:"required"`
	Nickname  string `json:"nickname" form:"nickname" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
}

type GetUsersReq struct {
	Limit        int64     `form:"limit"`
	Page         int64     `form:"page"`
	Username     string    `form:"username"`
	AccountState CommState `form:"account_state"`
	GroupName    string    `form:"group_name"`
	BeginTime    int64     `form:"begin_time"`
	EndTime      int64     `form:"end_time"`
}

type UpdateUserReq struct {
	Username     string    `json:"username" form:"username" binding:"required"`
	Nickname     string    `json:"nickname" form:"nickname" binding:"required"`
	AccountState CommState `json:"account_state" form:"account_state" binding:"required"`
	GroupName    string    `json:"group_name" form:"group_name" binding:"required"`
}

type DeleteUserReq struct {
	Username string `json:"username" form:"username" binding:"required"`
}

type UserInfo struct {
	Username            string    `json:"username" bson:"username"`
	Nickname            string    `json:"nickname" bson:"nickname"`
	AccountState        CommState `json:"account_state" bson:"account_state"`
	GroupName           string    `json:"group_name" bson:"group_name"`
	GroupType           GroupType `json:"group_type" bson:"group_type"`
	GroupMenus          []int64   `json:"group_menus" bson:"group_menus"`
	CreateIP            string    `json:"create_ip" bson:"create_ip"`
	CreateTime          int64     `json:"create_time" bson:"create_time"`
	CreateTimeLocal     string    `json:"create_time_local" bson:"create_time_local"`
	LastUpdateIP        string    `json:"last_update_ip" bson:"last_update_ip"`
	LastUpdateTime      int64     `json:"last_update_time" bson:"last_update_time"`
	LastUpdateTimeLocal string    `json:"last_update_time_local" bson:"last_update_time_local"`
	LastLoginIP         string    `json:"last_login_ip" bson:"last_login_ip"`
	LastLoginTime       int64     `json:"last_login_time" bson:"last_login_time"`
	LastLoginTimeLocal  string    `json:"last_login_time_local" bson:"last_login_time_local"`
}

type UserData struct {
	Count int64      `json:"count"`
	Data  []UserInfo `json:"data"`
}

type ResetPwdReq struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type GetLoginLogsReq struct {
	Limit     int64  `form:"limit"`
	Page      int64  `form:"page"`
	Username  string `form:"username"`
	BeginTime int64  `form:"begin_time"`
	EndTime   int64  `form:"end_time"`
}

type LoginLog struct {
	Username        string    `json:"username" bson:"username"`
	Nickname        string    `json:"nickname" bson:"nickname"`
	AccountState    CommState `json:"account_state" bson:"account_state"`
	GroupName       string    `json:"group_name" bson:"group_name"`
	ClientIP        string    `json:"client_ip" bson:"client_ip"`
	CreateTime      int64     `json:"create_time" bson:"create_time"`
	CreateTimeLocal string    `json:"create_time_local" bson:"create_time_local"`
}

type LoginLogData struct {
	Count int64      `json:"count"`
	Data  []LoginLog `json:"data"`
}

type GetMgmtLogsReq struct {
	Limit        int64        `form:"limit"`
	Page         int64        `form:"page"`
	EventType    EventType    `form:"event_type"`
	EventSubtype EventSubtype `form:"event_subtype"`
	OperatorName string       `form:"operator_name"`
	TargetName   string       `form:"target_name"`
	BeginTime    int64        `form:"begin_time"`
	EndTime      int64        `form:"end_time"`
}

type MgmtLog struct {
	EventType       EventType    `json:"event_type" bson:"event_type"`
	EventSubtype    EventSubtype `json:"event_subtype" bson:"event_subtype"`
	OperatorName    string       `json:"operator_name" bson:"operator_name"`
	TargetName      string       `json:"target_name" bson:"target_name"`
	ExtraInfo       string       `json:"extra_info" bson:"extra_info"`
	ClientIP        string       `json:"client_ip" bson:"client_ip"`
	CreateTime      int64        `json:"create_time" bson:"create_time"`
	CreateTimeLocal string       `json:"create_time_local" bson:"create_time_local"`
}

type MgmtLogData struct {
	Count int64     `json:"count"`
	Data  []MgmtLog `json:"data"`
}
