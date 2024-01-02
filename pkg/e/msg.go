package e

var MsgFlags = map[int]string{
	Success:                 "ok",
	Error:                   "fail",
	InvalidParams:           "参数错误",
	ErrorExistsUser:         "用户名已存在",
	ErrorFailEncryption:     "密码加密失败",
	ErrorExistsUserNotFound: "用户不存在",
	ErrorNotCompare:         "密码不匹配",
	ErrorAuthToken:          "token 认证失败",
	ErrorAuthCheckToken:     "token 过期",
	ErrorUploadFail:         "图片上传失败",
}

// GetMsg 获取状态码对应的信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
