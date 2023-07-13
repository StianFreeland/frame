package protos

var (
	InternalServerError = &CommRsp{Code: 10000, Msg: "Internal Server Error"}
)

var (
	InvalidReqParams        = &CommRsp{Code: 20000, Msg: "Invalid Request Parameters"}
	TokenEmpty              = &CommRsp{Code: 20001, Msg: "token empty"}
	InvalidUserGroupType    = &CommRsp{Code: 20200, Msg: "Invalid User Group Type"}
	InvalidUserAccountState = &CommRsp{Code: 20201, Msg: "Invalid User Account State"}
)

var (
	InsufficientPermissions     = &CommRsp{Code: 30000, Msg: "Insufficient Permissions"}
	UsernameOrPasswordIncorrect = &CommRsp{Code: 30100, Msg: "Username Or Password Incorrect"}
	UserAccountStateError       = &CommRsp{Code: 30101, Msg: "User Account State Error"}
	PasswordIncorrect           = &CommRsp{Code: 30102, Msg: "Password Incorrect"}
	TokenExpired                = &CommRsp{Code: 30110, Msg: "Please Login Again"}
	UserGroupNameExists         = &CommRsp{Code: 30200, Msg: "User Group Name Exists"}
	UserGroupNotExists          = &CommRsp{Code: 30201, Msg: "User Group Not Exists"}
	UserGroupNotEmpty           = &CommRsp{Code: 30202, Msg: "User Group Not Empty"}
	UsernameExists              = &CommRsp{Code: 30210, Msg: "Username Exists"}
	UserNotExists               = &CommRsp{Code: 30211, Msg: "User Not Exists"}
	MenuNotExists               = &CommRsp{Code: 30300, Msg: "Menu Not Exists"}
)

var (
	UnknownError = &CommRsp{Code: 90000, Msg: "Unknown Error"}
)
