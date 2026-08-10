package errx

import "github.com/lcylpzls/errx/internal/core"

type Code = core.Code
type CodeInfo = core.CodeInfo

const CodeUnknown = core.CodeUnknown

func RegisterCode(code Code, description string) { core.RegisterCode(code, description) }
func RegisterCodeKind(code Code, kind Kind)      { core.RegisterCodeKind(code, kind) }
func CodeKind(code Code) Kind                    { return core.CodeKind(code) }
func Describe(code Code) string                  { return core.Describe(code) }
func Codes() []CodeInfo                          { return core.Codes() }
func NewCode(code Code, msg string) *Error       { return core.NewCode(code, msg) }
func NewCodef(code Code, format string, args ...any) *Error {
	return core.NewCodef(code, format, args...)
}
func WrapCode(err error, code Code, msg string) *Error { return core.WrapCode(err, code, msg) }
func WrapCodef(err error, code Code, format string, args ...any) *Error {
	return core.WrapCodef(err, code, format, args...)
}
