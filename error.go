package errx

import "github.com/lcylpzls/errx/internal/core"

// Error 是结构化错误（实现见 internal/core）。
type Error = core.Error

// StackFrame 是调用栈帧（实现见 internal/core）。
type StackFrame = core.StackFrame

func New(kind Kind, code Code, msg string) *Error { return core.New(kind, code, msg) }
func Newf(kind Kind, code Code, format string, args ...any) *Error {
	return core.Newf(kind, code, format, args...)
}
func Wrap(err error, kind Kind, code Code, msg string) *Error { return core.Wrap(err, kind, code, msg) }
func Wrapf(err error, kind Kind, code Code, format string, args ...any) *Error {
	return core.Wrapf(err, kind, code, format, args...)
}
func As(err error) (*Error, bool)   { return core.As(err) }
func CodeOf(err error) (Code, bool) { return core.CodeOf(err) }
func KindOf(err error) Kind         { return core.KindOf(err) }
func Is(err error, code Code) bool  { return core.Is(err, code) }
func Retryable(err error) bool      { return core.Retryable(err) }
func WithField(err error, key string, val any) error {
	return core.WithField(err, key, val)
}
func SetStackCapture(enabled bool) { core.SetStackCapture(enabled) }
func SetStackDepth(depth int)      { core.SetStackDepth(depth) }
