package models

type MenuInfo struct {
	MenuID int64 `bson:"menu_id"`
}

func GetCollMenuInfo() string {
	return "menu_info"
}
