package errx

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

// maxChainDepth 是错误链遍历的最大深度，防御意外成环导致的死循环。
const maxChainDepth = 100

// Error 是结构化错误：携带错误码、分类、消息、结构化字段与可选调用栈。
// 通过 errors.Is / errors.As 与标准库错误链完全兼容。
type Error struct {
	code   Code
	kind   Kind
	msg    string
	fields []KV
	cause  error
	stack  []uintptr

	once   sync.Once // Error() 结果惰性缓存（并发安全）
	errStr string
}

// New 创建一个无原因的结构化错误。
func New(kind Kind, code Code, msg string) *Error {
	return &Error{
		kind:  kind,
		code:  code,
		msg:   msg,
		stack: captureStack(),
	}
}

// Newf 创建带格式化消息的结构化错误。
func Newf(kind Kind, code Code, format string, args ...any) *Error {
	return &Error{
		kind:  kind,
		code:  code,
		msg:   fmt.Sprintf(format, args...),
		stack: captureStack(),
	}
}

// Wrap 包装一个底层错误并附加分类与错误码。
// 当 err 为 nil 时返回 nil，便于直接 return errx.Wrap(err, ...) 的写法。
func Wrap(err error, kind Kind, code Code, msg string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		kind:  kind,
		code:  code,
		msg:   msg,
		cause: err,
		stack: captureStack(),
	}
}

// Wrapf 包装底层错误并附加格式化消息。
func Wrapf(err error, kind Kind, code Code, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		kind:  kind,
		code:  code,
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: captureStack(),
	}
}

// Error 返回格式为 "CODE: message: cause" 的文本；空字段自动省略。
// 结果惰性缓存，重复打印零额外开销。
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	e.once.Do(func() {
		var b strings.Builder
		b.WriteString(string(e.code))
		if e.msg != "" {
			b.WriteString(": ")
			b.WriteString(e.msg)
		}
		if e.cause != nil {
			b.WriteString(": ")
			b.WriteString(e.cause.Error())
		}
		e.errStr = b.String()
	})
	return e.errStr
}

// Format 支持 %s / %q / %v；%+v 额外输出创建时捕获的调用栈。
func (e *Error) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			io.WriteString(f, e.Error())
			for _, frame := range e.frames() {
				fmt.Fprintf(f, "\n\t%s:%d  %s", frame.file, frame.line, frame.fn)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(f, e.Error())
	case 'q':
		fmt.Fprintf(f, "%q", e.Error())
	default:
		fmt.Fprintf(f, "%%!%c(errx.Error=%s)", verb, e.Error())
	}
}

// Unwrap 返回被包装的底层错误，支持 errors.Is / errors.As 链路。
func (e *Error) Unwrap() error {
	return e.cause
}

// Cause 返回被包装的底层错误（兼容 pkg/errors 习惯）。
func (e *Error) Cause() error {
	return e.cause
}

// Code 返回错误码。
func (e *Error) Code() Code {
	return e.code
}

// Kind 返回错误分类。
func (e *Error) Kind() Kind {
	return e.kind
}

// Message 返回错误消息（不含错误码与原因）。
func (e *Error) Message() string {
	return e.msg
}

// Fields 返回错误携带的结构化字段快照。
func (e *Error) Fields() []KV {
	if len(e.fields) == 0 {
		return nil
	}
	out := make([]KV, len(e.fields))
	copy(out, e.fields)
	return out
}

// WithField 返回携带附加字段的新错误（不可变风格，原错误不受影响）。
func (e *Error) WithField(key string, val any) *Error {
	if e == nil {
		return nil
	}
	// 显式重建新实例：不复制 once/errStr（sync.Once 禁止复制），
	// Error() 文本不含字段，新实例重新惰性计算即可。
	ne := &Error{
		code:  e.code,
		kind:  e.kind,
		msg:   e.msg,
		cause: e.cause,
		stack: e.stack,
	}
	ne.fields = make([]KV, 0, len(e.fields)+1)
	ne.fields = append(ne.fields, e.fields...)
	ne.fields = append(ne.fields, KV{Key: key, Value: val})
	return ne
}

// frame 是格式化输出用的单帧信息。
type frame struct {
	file string
	line int
	fn   string
}

// frames 将捕获的 PC 转换为可读帧。
func (e *Error) frames() []frame {
	if len(e.stack) == 0 {
		return nil
	}
	callers := runtime.CallersFrames(e.stack)
	var out []frame
	for {
		fr, more := callers.Next()
		out = append(out, frame{file: fr.File, line: fr.Line, fn: fr.Function})
		if !more {
			break
		}
	}
	return out
}

// As 从错误链中取出第一个 *Error。
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// CodeOf 返回错误链中第一个结构化错误的错误码。
func CodeOf(err error) (Code, bool) {
	if e, ok := As(err); ok {
		return e.code, true
	}
	return "", false
}

// KindOf 返回错误链中第一个结构化错误的分类；无结构化错误时返回 KindUnknown。
func KindOf(err error) Kind {
	if e, ok := As(err); ok {
		return e.kind
	}
	return KindUnknown
}

// Is 判断错误链中是否存在指定错误码。
func Is(err error, code Code) bool {
	for depth := 0; err != nil && depth < maxChainDepth; depth++ {
		if e, ok := err.(*Error); ok && e.code == code {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Retryable 判断错误链中是否存在可重试分类（timeout/rate_limited/unavailable）。
func Retryable(err error) bool {
	for depth := 0; err != nil && depth < maxChainDepth; depth++ {
		if e, ok := err.(*Error); ok && e.kind.Retryable() {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// WithField 为任意错误附加结构化字段：
// 已是 *Error 时返回新实例；否则包装为 UNKNOWN 错误并保留原因链。
func WithField(err error, key string, val any) error {
	if err == nil {
		return nil
	}
	if e, ok := As(err); ok {
		return e.WithField(key, val)
	}
	return &Error{
		kind:  KindUnknown,
		code:  CodeUnknown,
		msg:   err.Error(),
		cause: err,
		fields: []KV{
			{Key: key, Value: val},
		},
		stack: captureStack(),
	}
}

// ---------------------------------------------------------------------------
// 调用栈捕获
// ---------------------------------------------------------------------------

var stackCapture atomic.Bool

func init() {
	stackCapture.Store(true)
}

// SetStackCapture 全局开关调用栈捕获。生产环境如对错误构造频率敏感可关闭。
func SetStackCapture(enabled bool) {
	stackCapture.Store(enabled)
}

// captureStack 捕获创建调用点的程序计数器（跳过本函数与构造函数）。
func captureStack() []uintptr {
	if !stackCapture.Load() {
		return nil
	}
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	return pcs[:n]
}
