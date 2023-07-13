package protos

type CommRsp struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func Success(data ...interface{}) *CommRsp {
	if len(data) == 0 {
		return &CommRsp{}
	}
	return &CommRsp{Data: data[0]}
}

func Error(err error) *CommRsp {
	return &CommRsp{Code: UnknownError.Code, Msg: err.Error()}
}
