package kerrors

import (
	"fmt"
	"log"
	"runtime"
)

// KError is struct for custom error response
type KError struct {
	Prefix           string                 `json:"prefix"`
	Code             ErrorCode              `json:"code"`
	Message          ErrorMessage           `json:"message"`
	LogMessage       string                 `json:"log_message"`
	OriginalError    error                  `json:"-"`
	OriginalErrorStr string                 `json:"original_error_str"`
	Stack            string                 `json:"stack"`
	Data             map[string]interface{} `json:"data"`
}

const (
	// stackline | [prefix]: log_message
	stackLineFormat                      = "%s | "
	printErrorFormatWithOutOriginalError = stackLineFormat + "[%s]: %v"
	// stackline | [prefix]: log_message \n recursive original error
	printErrorFormatWithOriginalError = printErrorFormatWithOutOriginalError + "\n%v"

	logFormatWithMoreInfo = "%s | [MORE INFO] %s"
)

func showCodeLine(filePath string, line int) string {
	return fmt.Sprintf("%v:%v", filePath, line)
}

func (e KError) Error() string {
	if e.OriginalError == nil { // last error message
		return fmt.Sprintf(printErrorFormatWithOutOriginalError, e.Stack, e.Prefix, e.LogMessage)
	}
	return e.Message.String()
}

func (e KError) Extract() map[string]interface{} {
	m := make(map[string]interface{}, len(e.Data)+4)
	m["prefix"] = e.Prefix
	m["code"] = e.Code
	m["stack"] = e.Stack
	m["original_error"] = e.OriginalErrorStr

	for key, value := range e.Data {
		m[key] = value
	}

	return m
}

// Wrap returns an error, and adds context to original error
func (e KError) Wrap(err error, code ErrorCode, data map[string]interface{}) error {
	msg := ErrDictionary.Get(code)
	_, filePath, line := caller()
	return KError{
		Prefix:           e.Prefix,
		Message:          msg,
		Code:             code,
		LogMessage:       msg.String(),
		Stack:            showCodeLine(filePath, line),
		OriginalError:    err,
		Data:             data,
		OriginalErrorStr: err.Error(),
	}
}
func (e KError) WrapManual(err error, code ErrorCode, msg string, data map[string]interface{}) error {
	_, filePath, line := caller()
	return KError{
		Prefix:           e.Prefix,
		Message:          ErrorMessage(msg),
		Code:             code,
		LogMessage:       msg,
		Stack:            showCodeLine(filePath, line),
		OriginalError:    err,
		Data:             data,
		OriginalErrorStr: err.Error(),
	}
}

// Wrapf returns an error, adds context to original error, and can create a new format more information
func (e KError) Wrapf(err error, code ErrorCode, format string, args ...interface{}) error {
	msg := ErrDictionary.Get(code)
	_, filePath, line := caller()
	kError := KError{
		Prefix:           e.Prefix,
		Message:          msg,
		Code:             code,
		LogMessage:       msg.String(),
		Stack:            showCodeLine(filePath, line),
		OriginalError:    err,
		OriginalErrorStr: err.Error(),
	}
	if len(format) != 0 { // if format, and args are not empty, add more info with new formats.
		kError.LogMessage = fmt.Sprintf(logFormatWithMoreInfo, msg, fmt.Sprintf(format, args...))
		return kError
	}
	return kError
}

// WithPrefix adding prefix to error message [prefix]: error_msg
func WithPrefix(prefix string) KError {
	return KError{
		Prefix: prefix,
	}
}

func caller() (string, string, int) {
	fpcs := make([]uintptr, 32)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		log.Println("MSG: NO CALLER")
	}
	c := runtime.FuncForPC(fpcs[0] - 1)
	if c == nil {
		log.Println("MSG CALLER WAS NIL")
		return "", "", 0
	}
	filepath, line := c.FileLine(fpcs[0] - 1)
	return c.Name(), filepath, line
}
