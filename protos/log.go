package protos

type UpdateLogConfigReq struct {
	Level int64 `json:"level" form:"level"`
}

type LogConfigData struct {
	Level int64 `json:"level"`
}
