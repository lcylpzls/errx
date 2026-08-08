package errx

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	inner := New(KindBusiness, "INNER", "内层")
	e := Wrap(inner, KindTimeout, "OUTER", "外层").WithField("req", "abc")

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("MarshalJSON 失败：%v", err)
	}
	s := string(data)
	for _, want := range []string{
		`"code":"OUTER"`, `"kind":"timeout"`, `"message":"外层"`,
		`"fields":[{"key":"req","value":"abc"}]`, `"cause"`, `"code":"INNER"`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("JSON 缺少 %s：%s", want, s)
		}
	}
}

func TestMarshalJSONNil(t *testing.T) {
	var e *Error
	data, err := json.Marshal(e)
	if err != nil || string(data) != "null" {
		t.Errorf("nil Error 应序列化为 null：%s %v", data, err)
	}
	// 直接调用 MarshalJSON 的 nil 分支
	data, err = e.MarshalJSON()
	if err != nil || string(data) != "null" {
		t.Errorf("直接调用 nil MarshalJSON 不符：%s %v", data, err)
	}
}

func TestMarshalJSONPlainCause(t *testing.T) {
	e := Wrap(errors.New("plain cause"), KindInternal, "WRAP", "包装")
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("MarshalJSON 失败：%v", err)
	}
	if !strings.Contains(string(data), "plain cause") {
		t.Errorf("非 errx 原因文本应保留：%s", data)
	}
}

func TestUnmarshalJSONRoundTrip(t *testing.T) {
	inner := New(KindBusiness, "INNER", "内层").WithField("k", "v")
	original := Wrap(inner, KindTimeout, "OUTER", "外层").WithField("req", "abc")

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal 失败：%v", err)
	}

	var got Error
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal 失败：%v", err)
	}
	if got.Code() != "OUTER" || got.Kind() != KindTimeout || got.Message() != "外层" {
		t.Errorf("外层恢复不符：%+v", &got)
	}
	cause := got.Unwrap()
	if e, ok := cause.(*Error); !ok || e.Code() != "INNER" || e.Kind() != KindBusiness {
		t.Errorf("内层恢复不符：%v", cause)
	}
	fields := got.Fields()
	if len(fields) != 1 || fields[0].Key != "req" {
		t.Errorf("外层字段恢复不符：%+v", fields)
	}
}

func TestUnmarshalJSONUnknownKind(t *testing.T) {
	data := []byte(`{"code":"X","kind":"not_a_kind","message":"m"}`)
	var e Error
	if err := json.Unmarshal(data, &e); err != nil {
		t.Fatalf("Unmarshal 失败：%v", err)
	}
	if e.Kind() != KindUnknown || e.Code() != "X" {
		t.Errorf("未知 kind 应解析为 unknown：%+v", &e)
	}
}

func TestUnmarshalJSONInvalid(t *testing.T) {
	var e Error
	if err := e.UnmarshalJSON([]byte("{invalid")); err == nil {
		t.Error("非法 JSON 应返回错误")
	}
}

func TestFromWireNil(t *testing.T) {
	if got := fromWire(nil); got != nil {
		t.Errorf("fromWire(nil) 应返回 nil：%v", got)
	}
}

func TestKindHTTPStatus(t *testing.T) {
	cases := map[Kind]int{
		KindInvalid:          400,
		KindUnauthorized:     401,
		KindForbidden:        403,
		KindNotFound:         404,
		KindAlreadyExists:    409,
		KindConflict:         409,
		KindRateLimited:      429,
		KindQuotaExceeded:    429,
		KindBusiness:         422,
		KindCancelled:        499,
		KindNotImplemented:   501,
		KindUnavailable:      503,
		KindDeadlineExceeded: 504,
		KindTimeout:          504,
		KindInternal:         500,
		KindDataLoss:         500,
		KindUnknown:          500,
		Kind(255):            500,
	}
	for k, want := range cases {
		if got := KindHTTPStatus(k); got != want {
			t.Errorf("Kind(%d) HTTPStatus 不符：got %d, want %d", k, got, want)
		}
	}
}

func TestErrorHTTPStatus(t *testing.T) {
	if got := New(KindNotFound, "N", "n").HTTPStatus(); got != 404 {
		t.Errorf("HTTPStatus 不符：%d", got)
	}
	var e *Error
	if got := e.HTTPStatus(); got != 500 {
		t.Errorf("nil Error HTTPStatus 应为 500：%d", got)
	}
}
