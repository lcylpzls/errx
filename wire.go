package errx

import (
	"encoding/json"
)

// wireError 是 Error 的跨服务传输形态。
type wireError struct {
	Code    Code       `json:"code"`
	Kind    string     `json:"kind"`
	Message string     `json:"message,omitempty"`
	Fields  []KV       `json:"fields,omitempty"`
	Cause   *wireError `json:"cause,omitempty"`
}

var kindByNames = func() map[string]Kind {
	m := make(map[string]Kind, len(kindNames))
	for k, name := range kindNames {
		m[name] = k
	}
	return m
}()

// MarshalJSON 将 Error 序列化为跨服务可传输的 JSON（含原因链与字段）。
func (e *Error) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	return json.Marshal(toWire(e))
}

// UnmarshalJSON 从 JSON 恢复 Error。调用栈不跨服务传输。
func (e *Error) UnmarshalJSON(data []byte) error {
	var w wireError
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	*e = *fromWire(&w)
	return nil
}

// toWire 将错误链转换为传输形态；非 errx 原因保留其文本。
func toWire(err error) *wireError {
	if e, ok := err.(*Error); ok {
		w := &wireError{
			Code:    e.code,
			Kind:    e.kind.String(),
			Message: e.msg,
			Fields:  e.fields,
		}
		if e.cause != nil {
			w.Cause = toWire(e.cause)
		}
		return w
	}
	return &wireError{
		Code:    CodeUnknown,
		Kind:    KindUnknown.String(),
		Message: err.Error(),
	}
}

// fromWire 从传输形态恢复错误链。
func fromWire(w *wireError) *Error {
	if w == nil {
		return nil
	}
	e := &Error{
		code:   w.Code,
		kind:   parseKind(w.Kind),
		msg:    w.Message,
		fields: w.Fields,
	}
	if w.Cause != nil {
		e.cause = fromWire(w.Cause)
	}
	return e
}

// parseKind 按名称解析 Kind，未知返回 KindUnknown。
func parseKind(s string) Kind {
	if k, ok := kindByNames[s]; ok {
		return k
	}
	return KindUnknown
}

// KindHTTPStatus 返回 Kind 对应的 HTTP 状态码。
func KindHTTPStatus(kind Kind) int {
	switch kind {
	case KindInvalid:
		return 400
	case KindUnauthorized:
		return 401
	case KindForbidden:
		return 403
	case KindNotFound:
		return 404
	case KindAlreadyExists, KindConflict:
		return 409
	case KindRateLimited, KindQuotaExceeded:
		return 429
	case KindBusiness:
		return 422
	case KindCancelled:
		return 499
	case KindNotImplemented:
		return 501
	case KindUnavailable:
		return 503
	case KindDeadlineExceeded, KindTimeout:
		return 504
	default:
		return 500
	}
}

// HTTPStatus 返回该错误对应的 HTTP 状态码。
func (e *Error) HTTPStatus() int {
	if e == nil {
		return 500
	}
	return KindHTTPStatus(e.kind)
}
