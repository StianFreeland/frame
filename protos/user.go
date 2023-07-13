package protos

type UserLoginReq struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserLoginData struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"user_info"`
}

type ChangePwdReq struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
}
