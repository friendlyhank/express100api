package express100api

//错误机制处理

type ErrCode int32

const (
	ErrCode_Success       ErrCode = 200
	ErrCode_SelfDefine    ErrCode = 401
	ErrCide_AppidIllegal  ErrCode = 40001
	ErrCide_AppUidIllegal ErrCode = 40002
)

var ErrCode_name = map[int32]string{
	200:  "Success",    //成功
	401:  "SelfDefine", //自定义错误详情看message
	4001: "参数appid非法或参数client_id非法 ",
	4002: "参数 appuid 非法",
}
