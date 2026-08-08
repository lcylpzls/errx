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
	countConstruct(kind)
	return &Error{
		kind:  kind,
		code:  normalizeCode(code),
		msg:   msg,
		stack: captureStack(),
	}
}

// Newf 创建带格式化消息的结构化错误。
func Newf(kind Kind, code Code, format string, args ...any) *Error {
	countConstruct(kind)
	return &Error{
		kind:  kind,
		code:  normalizeCode(code),
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
	countConstruct(kind)
	return &Error{
		kind:  kind,
		code:  normalizeCode(code),
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
	countConstruct(kind)
	return &Error{
		kind:  kind,
		code:  normalizeCode(code),
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: captureStack(),
	}
}

// normalizeCode 将空错误码归一为 CodeUnknown，保证错误文本格式稳定。
func normalizeCode(code Code) Code {
	if code == "" {
		return CodeUnknown
	}
	return code
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
		// message 与 cause 文本相同时只输出一次（如 WithField 包装普通错误），
		// 避免 "UNKNOWN: 普通错误: 普通错误" 式重复。
		if e.msg != "" && (e.cause == nil || e.msg != e.cause.Error()) {
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
				fmt.Fprintf(f, "\n\t%s:%d  %s", frame.File, frame.Line, frame.Function)
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

// StackFrame 是调用栈中的单帧信息，供日志与监控程序化读取。
type StackFrame struct {
	File     string
	Line     int
	Function string
}

// frames 将捕获的 PC 转换为可读帧。
func (e *Error) frames() []StackFrame {
	if len(e.stack) == 0 {
		return nil
	}
	callers := runtime.CallersFrames(e.stack)
	var out []StackFrame
	for {
		fr, more := callers.Next()
		out = append(out, StackFrame{File: fr.File, Line: fr.Line, Function: fr.Function})
		if !more {
			break
		}
	}
	return out
}

// StackTrace 返回创建时捕获的调用栈；栈捕获关闭或无栈时返回 nil。
func (e *Error) StackTrace() []StackFrame {
	return e.frames()
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
	countQuery()
	if e, ok := As(err); ok {
		return e.code, true
	}
	return "", false
}

// KindOf 返回错误链中第一个结构化错误的分类；无结构化错误时返回 KindUnknown。
func KindOf(err error) Kind {
	countQuery()
	if e, ok := As(err); ok {
		return e.kind
	}
	return KindUnknown
}

// codeSentinel 是错误码匹配哨兵，作为 errors.Is 的目标。
type codeSentinel Code

// Error 实现 error 接口。
func (c codeSentinel) Error() string {
	return string(c)
}

// retryableSentinel 是可重试匹配哨兵，作为 errors.Is 的目标。
type retryableSentinel struct{}

// Error 实现 error 接口。
func (retryableSentinel) Error() string {
	return "retryable"
}

// Is 支持 errors.Is 按错误码或可重试分类匹配（沿错误链与聚合子错误展开）。
func (e *Error) Is(target error) bool {
	switch t := target.(type) {
	case codeSentinel:
		return e.code == Code(t)
	case retryableSentinel:
		return e.kind.Retryable()
	default:
		return false
	}
}

// Is 判断错误链中是否存在指定错误码（支持单链与 Aggregate 多错误展开）。
func Is(err error, code Code) bool {
	countQuery()
	return errors.Is(err, codeSentinel(code))
}

// Retryable 判断错误链中是否存在可重试分类（支持单链与 Aggregate 多错误展开）。
func Retryable(err error) bool {
	countQuery()
	return errors.Is(err, retryableSentinel{})
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
var stackDepth atomic.Int32

func init() {
	stackCapture.Store(true)
	stackDepth.Store(32)
}

// SetStackCapture 全局开关调用栈捕获。生产环境如对错误构造频率敏感可关闭。
func SetStackCapture(enabled bool) {
	stackCapture.Store(enabled)
}

// SetStackDepth 设置栈捕获的最大帧数；depth <= 0 时恢复默认 32。
func SetStackDepth(depth int) {
	if depth <= 0 {
		depth = 32
	}
	stackDepth.Store(int32(depth))
}

// captureStack 捕获创建调用点的程序计数器（跳过本函数与构造函数）。
func captureStack() []uintptr {
	if !stackCapture.Load() {
		return nil
	}
	d := int(stackDepth.Load())
	pcs := make([]uintptr, d)
	n := runtime.Callers(2, pcs)
	return pcs[:n]
}
