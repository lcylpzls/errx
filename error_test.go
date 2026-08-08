package errx

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

const testCode Code = "TEST_ERROR"

func TestNew(t *testing.T) {
	e := New(KindBusiness, testCode, "业务失败")
	if e.Code() != testCode || e.Kind() != KindBusiness || e.Message() != "业务失败" {
		t.Errorf("New 字段不符：%+v", e)
	}
	if e.Cause() != nil || e.Unwrap() != nil {
		t.Error("New 不应携带原因")
	}
	if e.Error() != "TEST_ERROR: 业务失败" {
		t.Errorf("Error() 不符：%s", e.Error())
	}
}

func TestNewf(t *testing.T) {
	e := Newf(KindInvalid, "INVALID", "参数 %s 无效", "id")
	if e.Error() != "INVALID: 参数 id 无效" {
		t.Errorf("Newf Error() 不符：%s", e.Error())
	}
}

func TestEmptyCodeNormalized(t *testing.T) {
	cases := []*Error{
		New(KindBusiness, "", "空码"),
		Newf(KindBusiness, "", "空码 %s", "x"),
		Wrap(errors.New("cause"), KindTimeout, "", "包装"),
		Wrapf(errors.New("cause"), KindTimeout, "", "包装 %s", "x"),
	}
	for _, e := range cases {
		if e.Code() != CodeUnknown {
			t.Errorf("空错误码应归一为 UNKNOWN：%v", e)
		}
	}
	if got := New(KindBusiness, "", "空码").Error(); got != "UNKNOWN: 空码" {
		t.Errorf("归一后 Error() 不符：%s", got)
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("底层失败")
	e := Wrap(cause, KindTimeout, "TIMEOUT", "调用超时")
	if e.Unwrap() != cause || e.Cause() != cause {
		t.Error("Wrap 原因链不符")
	}
	if e.Error() != "TIMEOUT: 调用超时: 底层失败" {
		t.Errorf("Wrap Error() 不符：%s", e.Error())
	}
}

func TestWrapNil(t *testing.T) {
	if e := Wrap(nil, KindInternal, "X", "y"); e != nil {
		t.Errorf("Wrap(nil) 应返回 nil：%v", e)
	}
	if e := Wrapf(nil, KindInternal, "X", "y"); e != nil {
		t.Errorf("Wrapf(nil) 应返回 nil：%v", e)
	}
}

func TestWrapf(t *testing.T) {
	e := Wrapf(errors.New("db"), KindUnavailable, "DB_DOWN", "连接 %s 失败", "mysql")
	if e.Error() != "DB_DOWN: 连接 mysql 失败: db" {
		t.Errorf("Wrapf Error() 不符：%s", e.Error())
	}
}

func TestErrorStringVariants(t *testing.T) {
	cases := []struct {
		name string
		e    *Error
		want string
	}{
		{"code only", New(KindUnknown, "CODE", ""), "CODE"},
		{"code+cause", &Error{code: "CODE", cause: errors.New("c")}, "CODE: c"},
		{"code+msg+cause", Wrap(errors.New("c"), KindUnknown, "CODE", "m"), "CODE: m: c"},
		{"nil", nil, ""},
	}
	for _, tc := range cases {
		if got := tc.e.Error(); got != tc.want {
			t.Errorf("%s：got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestFormat(t *testing.T) {
	e := New(KindInternal, "FMT", "格式化")

	if got := fmt.Sprintf("%s", e); got != "FMT: 格式化" {
		t.Errorf("%%s 不符：%s", got)
	}
	if got := fmt.Sprintf("%q", e); got != `"FMT: 格式化"` {
		t.Errorf("%%q 不符：%s", got)
	}
	if got := fmt.Sprintf("%v", e); got != "FMT: 格式化" {
		t.Errorf("%%v 不符：%s", got)
	}
	if got := fmt.Sprintf("%d", e); !strings.Contains(got, "FMT: 格式化") {
		t.Errorf("%%d 回退不符：%s", got)
	}
}

func TestFormatWithStack(t *testing.T) {
	e := New(KindInternal, "STACK", "带栈")
	got := fmt.Sprintf("%+v", e)
	if !strings.Contains(got, "STACK: 带栈") {
		t.Errorf("%%+v 缺少错误文本：%s", got)
	}
	if !strings.Contains(got, "error_test.go") {
		t.Errorf("%%+v 缺少调用栈：%s", got)
	}
}

func TestFormatNoStack(t *testing.T) {
	e := &Error{code: "NOSTACK", msg: "无栈"}
	if got := fmt.Sprintf("%+v", e); got != "NOSTACK: 无栈" {
		t.Errorf("无栈 %%+v 不符：%s", got)
	}
}

func TestAs(t *testing.T) {
	e := New(KindBusiness, "AS", "查找")
	if got, ok := As(e); !ok || got != e {
		t.Error("As 未命中 *Error")
	}
	if _, ok := As(errors.New("plain")); ok {
		t.Error("As 不应命中普通错误")
	}
}

func TestCodeOf(t *testing.T) {
	e := Wrap(errors.New("x"), KindBusiness, "CODE_OF", "包裹")
	if code, ok := CodeOf(e); !ok || code != "CODE_OF" {
		t.Errorf("CodeOf 不符：%v %v", code, ok)
	}
	if _, ok := CodeOf(errors.New("plain")); ok {
		t.Error("CodeOf 不应命中普通错误")
	}
}

func TestKindOf(t *testing.T) {
	if got := KindOf(New(KindRateLimited, "K", "k")); got != KindRateLimited {
		t.Errorf("KindOf 不符：%v", got)
	}
	if got := KindOf(errors.New("plain")); got != KindUnknown {
		t.Errorf("普通错误 KindOf 应为 Unknown：%v", got)
	}
}

func TestIs(t *testing.T) {
	inner := New(KindBusiness, "INNER", "内部")
	outer := Wrap(inner, KindInternal, "OUTER", "外层")
	if !Is(outer, "INNER") || !Is(outer, "OUTER") {
		t.Error("Is 应命中链上错误码")
	}
	if Is(outer, "MISSING") {
		t.Error("Is 不应命中未存在错误码")
	}
	if Is(errors.New("plain"), "INNER") {
		t.Error("普通错误链 Is 不应命中")
	}
}

func TestRetryable(t *testing.T) {
	if !Retryable(New(KindTimeout, "T", "超时")) {
		t.Error("timeout 应可重试")
	}
	if !Retryable(New(KindRateLimited, "R", "限流")) {
		t.Error("rate_limited 应可重试")
	}
	if !Retryable(New(KindUnavailable, "U", "不可用")) {
		t.Error("unavailable 应可重试")
	}
	if Retryable(New(KindBusiness, "B", "业务")) {
		t.Error("business 不应可重试")
	}
	if Retryable(errors.New("plain")) {
		t.Error("普通错误不应可重试")
	}
	// 链上全部不可重试时不可重试
	if Retryable(Wrap(New(KindBusiness, "B", "业务"), KindInternal, "I", "内层")) {
		t.Error("全部不可重试的链不应可重试")
	}
	// 链上任一节点可重试即视为可重试
	if !Retryable(Wrap(New(KindTimeout, "T", "超时"), KindInternal, "I", "内层")) {
		t.Error("链上 timeout 应可重试")
	}
}

func TestWithFieldMethod(t *testing.T) {
	e := New(KindBusiness, "WF", "带字段").WithField("order_id", "10086")
	if e == nil || e.Code() != "WF" {
		t.Fatal("WithField 返回错误不符")
	}
	fields := e.Fields()
	if len(fields) != 1 || fields[0].Key != "order_id" || fields[0].Value != "10086" {
		t.Errorf("字段不符：%+v", fields)
	}

	// 不可变：继续追加不影响原错误
	e2 := e.WithField("user_id", "42")
	if len(e.Fields()) != 1 || len(e2.Fields()) != 2 {
		t.Error("WithField 应保持不可变")
	}
	// Fields 返回快照，修改不影响内部
	snap := e.Fields()
	snap[0].Key = "hacked"
	if e.Fields()[0].Key != "order_id" {
		t.Error("Fields 应返回拷贝")
	}
}

func TestWithFieldNilReceiver(t *testing.T) {
	var e *Error
	if got := e.WithField("k", "v"); got != nil {
		t.Errorf("nil 接收者应返回 nil：%v", got)
	}
}

func TestFieldsEmpty(t *testing.T) {
	e := New(KindBusiness, "E", "无字段")
	if e.Fields() != nil {
		t.Errorf("无字段时应返回 nil：%v", e.Fields())
	}
}

func TestWithFieldTopLevel(t *testing.T) {
	if got := WithField(nil, "k", "v"); got != nil {
		t.Errorf("WithField(nil) 应返回 nil：%v", got)
	}

	ee := New(KindBusiness, "TF", "顶层")
	got := WithField(ee, "k", "v")
	if e, ok := got.(*Error); !ok || len(e.Fields()) != 1 {
		t.Errorf("errx 错误 WithField 不符：%v", got)
	}

	plain := errors.New("普通错误")
	wrapped := WithField(plain, "k", "v")
	if e, ok := wrapped.(*Error); !ok || e.Code() != CodeUnknown || e.Unwrap() != plain {
		t.Errorf("普通错误 WithField 包装不符：%v", wrapped)
	}
	if e, ok := wrapped.(*Error); ok {
		if e.Error() != "UNKNOWN: 普通错误" {
			t.Errorf("普通错误 WithField 文本不应重复：%s", e.Error())
		}
		if e.Message() != "普通错误" {
			t.Errorf("普通错误 Message 应保留原文：%s", e.Message())
		}
	}
}

func TestStackTrace(t *testing.T) {
	e := New(KindInternal, "STACK_TRACE", "带栈")
	frames := e.StackTrace()
	if len(frames) == 0 {
		t.Fatal("StackTrace 应返回调用栈")
	}
	found := false
	for _, f := range frames {
		if strings.Contains(f.File, "error_test.go") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("栈帧应包含测试文件：%+v", frames)
	}

	SetStackCapture(false)
	if got := New(KindInternal, "NO_STACK", "无栈").StackTrace(); got != nil {
		t.Errorf("关闭栈捕获后 StackTrace 应为 nil：%v", got)
	}
	SetStackCapture(true)
}

func TestSetStackCapture(t *testing.T) {
	SetStackCapture(false)
	e := New(KindInternal, "NOSTACK", "x")
	if got := fmt.Sprintf("%+v", e); got != "NOSTACK: x" {
		t.Errorf("关闭栈捕获后 %%+v 不符：%s", got)
	}

	SetStackCapture(true)
	e2 := New(KindInternal, "STACK2", "y")
	if !strings.Contains(fmt.Sprintf("%+v", e2), "error_test.go") {
		t.Error("恢复后应捕获调用栈")
	}
}

func TestErrorsCompat(t *testing.T) {
	inner := New(KindBusiness, "COMPAT", "内层")
	outer := Wrap(inner, KindInternal, "OUTER", "外层")

	var target *Error
	if !errors.As(outer, &target) || target.Code() != "OUTER" {
		t.Error("errors.As 应命中外层")
	}
	if !errors.Is(outer, inner) {
		t.Error("errors.Is 应沿 Unwrap 链命中内层")
	}
}

func TestSentinelError(t *testing.T) {
	if got := codeSentinel("SENT").Error(); got != "SENT" {
		t.Errorf("codeSentinel Error 不符：%s", got)
	}
	if got := (retryableSentinel{}).Error(); got != "retryable" {
		t.Errorf("retryableSentinel Error 不符：%s", got)
	}
}
