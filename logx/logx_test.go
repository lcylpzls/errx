package logx

import (
	"errors"
	"testing"

	"github.com/lcylpzls/errx"
)

func TestFields_ErrX(t *testing.T) {
	e := errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").
		WithField("order_id", "10086")

	g := Fields(e)
	if g.Len() != 4 {
		t.Fatalf("字段数不符：%d", g.Len())
	}

	keys := map[string]bool{}
	var orderID any
	for i := 0; i < g.Len(); i++ {
		f := g.At(i)
		keys[f.Key] = true
		if f.Key == "order_id" {
			orderID = f.Value
		}
	}
	for _, want := range []string{"err.code", "err.kind", "err.message", "order_id"} {
		if !keys[want] {
			t.Errorf("缺少字段 %s", want)
		}
	}
	if orderID != "10086" {
		t.Errorf("order_id 值不符：%v", orderID)
	}
}

func TestFields_ErrXNoMessage(t *testing.T) {
	e := errx.New(errx.KindInvalid, "BAD_ARG", "")
	g := Fields(e)
	if g.Len() != 2 {
		t.Fatalf("无消息时应只有 code/kind：%d", g.Len())
	}
	for i := 0; i < g.Len(); i++ {
		if g.At(i).Key == "err.message" {
			t.Error("空消息不应输出 err.message")
		}
	}
}

func TestFields_PlainError(t *testing.T) {
	g := Fields(errors.New("普通错误"))
	if g.Len() != 2 {
		t.Fatalf("普通错误字段数不符：%d", g.Len())
	}
	if g.At(0).Key != "err.code" || g.At(1).Key != "err.kind" {
		t.Errorf("普通错误字段不符：%v %v", g.At(0).Key, g.At(1).Key)
	}
}

func TestFields_Nil(t *testing.T) {
	g := Fields(nil)
	if g.Len() != 2 {
		t.Fatalf("nil 错误字段数不符：%d", g.Len())
	}
}
