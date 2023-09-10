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
