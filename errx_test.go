package errx_test

import (
	"errors"
	"testing"

	"github.com/lcylpzls/errx"
)

// TestPublicAPI 黑盒冒烟测试：覆盖根包全部转发函数与类型别名，
// 保证 internal/core 重构后公开 API 行为一致。
func TestPublicAPI(t *testing.T) {
	if errx.Version != "v1.6.0" {
		t.Fatalf("Version = %s", errx.Version)
	}

	errx.SetStackCapture(true)
	errx.SetStackDepth(16)

	errx.RegisterCode(errx.Code("smoke_code"), "冒烟错误码")
	errx.RegisterCodeKind(errx.Code("smoke_code"), errx.KindBusiness)
	if errx.CodeKind(errx.Code("smoke_code")) != errx.KindBusiness {
		t.Fatal("CodeKind 不一致")
	}
	if errx.Describe(errx.Code("smoke_code")) != "冒烟错误码" {
		t.Fatal("Describe 不一致")
	}
	if len(errx.Codes()) == 0 {
		t.Fatal("Codes 为空")
	}
	if errx.KindHTTPStatus(errx.KindNotFound) != 404 {
		t.Fatal("KindHTTPStatus 不一致")
	}

	e := errx.New(errx.KindBusiness, errx.Code("smoke_code"), "冒烟")
	e2 := errx.Newf(errx.KindInvalid, errx.Code("smoke_fmt"), "%s", "格式化")
	e3 := errx.Wrap(errors.New("根因"), errx.KindTimeout, errx.Code("smoke_wrap"), "包装")
	e4 := errx.Wrapf(errors.New("根因"), errx.KindUnavailable, errx.Code("smoke_wrapf"), "%s", "包装")
	if e == nil || e2 == nil || e3 == nil || e4 == nil {
		t.Fatal("构造函数返回 nil")
	}

	if _, ok := errx.As(e); !ok {
		t.Fatal("As 失败")
	}
	if code, ok := errx.CodeOf(e); !ok || code != errx.Code("smoke_code") {
		t.Fatal("CodeOf 不一致")
	}
	if errx.KindOf(e) != errx.KindBusiness {
		t.Fatal("KindOf 不一致")
	}
	if !errx.Is(e, errx.Code("smoke_code")) {
		t.Fatal("Is 失败")
	}
	if errx.Retryable(e) {
		t.Fatal("Business 不应可重试")
	}
	if !errx.Retryable(e3) {
		t.Fatal("Timeout 应可重试")
	}
	we := errx.WithField(e, "key", "val")
	if we == nil {
		t.Fatal("WithField 返回 nil")
	}

	nc := errx.NewCode(errx.Code("smoke_nc"), "新码")
	nc2 := errx.NewCodef(errx.Code("smoke_ncf"), "%s", "新码格式化")
	wc := errx.WrapCode(errors.New("x"), errx.Code("smoke_wc"), "包装码")
	wc2 := errx.WrapCodef(errors.New("x"), errx.Code("smoke_wcf"), "%s", "包装码格式化")
	if nc == nil || nc2 == nil || wc == nil || wc2 == nil {
		t.Fatal("错误码构造返回 nil")
	}

	joined := errx.Join(errors.New("a"), errors.New("b"))
	if joined == nil {
		t.Fatal("Join 返回 nil")
	}

	hook := fakeHook{}
	errx.SetMetricsHook(hook)
	errx.ResetMetricsHook()

	var _ errx.KV = errx.KV{Key: "k", Value: "v"}
	var _ errx.Policy = errx.Policy{}
	var _ errx.Category = errx.CatBusiness
	var _ errx.CodeInfo = errx.CodeInfo{}
	var _ errx.StackFrame = errx.StackFrame{}
}

type fakeHook struct{}

func (fakeHook) IncCounter(string, ...string) {}
