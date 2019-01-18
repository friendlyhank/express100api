package express100api

//错误机制处理

type ErrCode int32

const (
	ErrCode_Success    ErrCode = 200
	ErrCode_SelfDefine ErrCode = 401
)

var ErrCode_name = map[int32]string{
	200: "Success",    //成功
	401: "SelfDefine", //自定义错误详情看message
}
