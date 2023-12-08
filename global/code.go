package global

// 响应码
// 请求成功，正常响应
var CodeOK = "0"

// 缺少必要参数
var CodeLackRequired = "1"

// 解码失败
var CodeUnmarshalFailed = "2"

// 查询失败
var CodeQueryFailed = "3"

// 格式错误
var CodeFormatError = "4"

// 获取参数失败
var CodeGetParamsFailed = "5"

// 数据为空，如"", [], {}, nil等
var CodeEmptyValue = "6"

// 创建文件失败
var CodeCreateFileError = "7"

// 打开文件失败
var CodeOpenFileError = "8"

// 复制数据失败
var CodeCopyFileError = "9"

// 存储失败
var CodeSaveFileError = "10"

// 内容已存在
var CodeExist = "11"

// 内容不存在
var CodeNotExist = "12"

// 内容未变更
var CodeNotModify = "13"

// 添加数据失败
var CodeCreateDataFailed = "14"

// 删除数据失败
var CodeDeleteFailed = "15"

// 参数无效
var CodeInvalid = "16"

// 未登录
var CodeNotLogin = "17"

// 数据错误
var CodeDataError = "18"

// 活动未开始
var CodeNotStart = "19"

// 活动已结束
var CodeFinished = "20"

// 参数有误
var CodeParamError = "21"

// 重复操作
var CodeRunAgain = "22"

// 无权限
var CodeAuthyLimited = "23"

// 奖池超限
var CodeAwardPoolOutLimited = "24"

// 更新失败
var CodeUpdateFailed = "25"

// token签名失败
var CodeSignTokenError = "26"

// token失效
var CodeTokenInvalid = "27"
