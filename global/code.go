package global

// 响应码
var (
	CodeOK                  = "0"  // 请求成功，正常响应
	CodeLackRequired        = "1"  // 缺少必要参数
	CodeUnmarshalFailed     = "2"  // 解码失败
	CodeQueryFailed         = "3"  // 查询失败
	CodeFormatError         = "4"  // 格式错误
	CodeGetParamsFailed     = "5"  // 获取参数失败
	CodeEmptyValue          = "6"  // 数据为空，如"", [], {}, nil等
	CodeCreateFileError     = "7"  // 创建文件失败
	CodeOpenFileError       = "8"  // 打开文件失败
	CodeCopyFileError       = "9"  // 复制数据失败
	CodeSaveFileError       = "10" // 存储失败
	CodeExist               = "11" // 内容已存在
	CodeNotExist            = "12" // 内容不存在
	CodeNotModify           = "13" // 内容未变更
	CodeCreateDataFailed    = "14" // 添加数据失败
	CodeDeleteFailed        = "15" // 删除数据失败
	CodeInvalid             = "16" // 参数无效
	CodeNotLogin            = "17" // 未登录
	CodeDataError           = "18" // 数据错误
	CodeNotStart            = "19" // 活动未开始
	CodeFinished            = "20" // 活动已结束
	CodeParamError          = "21" // 参数有误
	CodeRunAgain            = "22" // 重复操作
	CodeAuthyLimited        = "23" // 无权限
	CodeAwardPoolOutLimited = "24" // 奖池超限
	CodeUpdateFailed        = "25" // 更新失败
	CodeSignTokenError      = "26" // token签名失败
	CodeTokenInvalid        = "27" // token失效
	CodeNoActivityInfo      = "28" // 没有活动信息
	CodeActivityInfoError   = "29" // 活动信息错误
)
