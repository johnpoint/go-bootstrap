package berror

import (
	"fmt"
)

// 通用错误码
var (
	OK      = &BErr{Code: 0, Message: "OK"}
	Unknown = &BErr{Code: -1, Message: "未知错误"}
)

// BErr 定义错误
type BErr struct {
	Code      int    // 错误码
	Message   string // 展示给用户看的
	ErrorInfo error  // 保存内部错误信息
}

func (err *BErr) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.ErrorInfo)
}

// GetErrCode extracts the error code from an error.
// If the error is not a BErr type, it returns Unknown.Code.
func GetErrCode(err error) int {
	trueErr, ok := err.(*BErr)
	if !ok {
		return Unknown.Code
	}
	return trueErr.Code
}

// GetErrMessage extracts the error message from an error.
// If the error is not a BErr type, it returns Unknown.Message.
func GetErrMessage(err error) string {
	trueErr, ok := err.(*BErr)
	if !ok {
		return Unknown.Message
	}
	return trueErr.Message
}

// WrapErr wraps an error with additional error information.
// It preserves the code and message from the original error while attaching the error info.
// Returns Unknown error if err is nil to prevent panic.
func WrapErr(err *BErr, errInfo error) *BErr {
	if err == nil {
		err = Unknown
	}
	return &BErr{
		Code:      err.Code,
		Message:   err.Message,
		ErrorInfo: errInfo,
	}
}

// Deprecated: WarpErr is a misspelled alias for WrapErr. Use WrapErr instead.
func WarpErr(err *BErr, errInfo error) *BErr {
	return WrapErr(err, errInfo)
}

// DecodeErr decodes an error into its code and message components.
// If the error is nil, it returns OK code and message.
// If the error is not a BErr type, it returns Unknown code and message.
// If the error has ErrorInfo, it appends the error details to the message.
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
