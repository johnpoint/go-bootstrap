package errorhelper

import (
	"fmt"
)

var (
	// 通用错误码
	OK      = &BErr{Code: 0, Message: "OK"}
	Unknown = &BErr{Code: -1, Message: "未知错误"}
)

// Err 定义错误
type BErr struct {
	Code      int    // 错误码
	Message   string // 展示给用户看的
	ErrorInfo error  // 保存内部错误信息
}

func (err *BErr) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.ErrorInfo)
}

func GetErrCode(err error) int {
	trueErr, ok := err.(*BErr)
	if !ok {
		return Unknown.Code
	}
	return trueErr.Code
}

func GetErrMessage(err error) string {
	trueErr, ok := err.(*BErr)
	if !ok {
		return Unknown.Message
	}
	return trueErr.Message
}

func WarpErr(err *BErr, errInfo error) *BErr {
	return &BErr{
		Code:      err.Code,
		Message:   err.Message,
		ErrorInfo: errInfo,
	}
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	trueErr, ok := err.(*BErr)
	if !ok {
		return Unknown.Code, Unknown.Message
	}
	if trueErr.ErrorInfo != nil {
		trueErr.Message = fmt.Sprintf("%s: %+v", trueErr.Message, trueErr.ErrorInfo.Error())
	}
	return trueErr.Code, trueErr.Message
}
