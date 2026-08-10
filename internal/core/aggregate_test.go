package core

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestAggregateFilterAndPassthrough(t *testing.T) {
	if got := Join(); got != nil {
		t.Errorf("空聚合应返回 nil：%v", got)
	}
	if got := Join(nil, nil); got != nil {
		t.Errorf("全 nil 应返回 nil：%v", got)
	}
	single := New(KindBusiness, "SINGLE", "单个")
	if got := Join(single, nil); got != single {
		t.Error("单个错误应直接透传")
	}
}

func TestAggregateError(t *testing.T) {
	a := Join(
		New(KindBusiness, "A1", "错误一"),
		New(KindTimeout, "A2", "错误二"),
	)
	if a == nil {
		t.Fatal("聚合不应为 nil")
	}
	msg := a.Error()
	if !strings.Contains(msg, "2 个错误") || !strings.Contains(msg, "A1: 错误一") || !strings.Contains(msg, "A2: 错误二") {
		t.Errorf("聚合 Error() 不符：%s", msg)
	}
	// 惰性缓存
	if a.Error() != msg {
		t.Error("聚合 Error() 应缓存")
	}
}

func TestAggregateNilReceiver(t *testing.T) {
	var a *Aggregate
	if a.Error() != "" {
		t.Error("nil 聚合 Error() 应为空")
	}
	if a.Errors() != nil {
		t.Error("nil 聚合 Errors() 应为 nil")
	}
}

func TestAggregateUnwrapAndErrors(t *testing.T) {
	e1 := New(KindBusiness, "AG1", "一")
	e2 := New(KindTimeout, "AG2", "二")
	a := Join(e1, e2).(*Aggregate)

	unwrapped := a.Unwrap()
	if len(unwrapped) != 2 || unwrapped[0] != e1 || unwrapped[1] != e2 {
		t.Error("Unwrap 子错误不符")
	}
	got := a.Errors()
	if len(got) != 2 {
		t.Fatal("Errors 数量不符")
	}
	got[0] = nil // 修改拷贝不影响聚合
	if a.Errors()[0] != e1 {
		t.Error("Errors 应返回拷贝")
	}
}

func TestAggregateStdCompat(t *testing.T) {
	e1 := New(KindBusiness, "AGC1", "一")
	e2 := New(KindTimeout, "AGC2", "二")
	a := Join(e1, e2)

	if !errors.Is(a, e1) || !errors.Is(a, e2) {
		t.Error("errors.Is 应命中任一子错误")
	}
	var target *Error
	if !errors.As(a, &target) {
		t.Error("errors.As 应命中子错误")
	}
	if !Is(a, "AGC1") || !Is(a, "AGC2") {
		t.Error("errx.Is 应命中聚合内错误码")
	}
	if !Retryable(a) {
		t.Error("errx.Retryable 应命中聚合内可重试子错误")
	}
}

func TestAggregateJSON(t *testing.T) {
	agg := Join(
		New(KindBusiness, "ORDER_FAIL", "下单失败").WithField("order_id", "1"),
		New(KindInternal, "DB_DOWN", "数据库不可用"),
	)
	data, err := json.Marshal(agg)
	if err != nil {
		t.Fatalf("意外错误：%v", err)
	}
	if !strings.Contains(string(data), `"errors"`) || !strings.Contains(string(data), "ORDER_FAIL") {
		t.Errorf("聚合 JSON 缺少子错误：%s", data)
	}

	var restored Aggregate
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("聚合 Unmarshal 失败：%v", err)
	}
	if !Is(&restored, "ORDER_FAIL") || !Is(&restored, "DB_DOWN") {
		t.Error("还原后的聚合应可命中子错误码")
	}
	errs := restored.Errors()
	if len(errs) != 2 {
		t.Fatalf("还原后子错误数量不符：%d", len(errs))
	}
	if e, ok := As(errs[0]); !ok || len(e.Fields()) != 1 || e.Fields()[0].Key != "order_id" {
		t.Errorf("还原后字段不符：%v", errs[0])
	}

	var nilAgg *Aggregate
	if data, err := json.Marshal(nilAgg); err != nil || string(data) != "null" {
		t.Errorf("nil 聚合应序列化为 null：%s %v", data, err)
	}
	if data, err := nilAgg.MarshalJSON(); err != nil || string(data) != "null" {
		t.Errorf("直接调用 nil MarshalJSON 不符：%s %v", data, err)
	}

	var invalid Aggregate
	if err := invalid.UnmarshalJSON([]byte("{invalid")); err == nil {
		t.Error("非法 JSON 应返回错误")
	}
}

func TestAggregateConcurrentError(t *testing.T) {
	ag := Join(
		New(KindBusiness, "CONC_A", "a"),
		New(KindTimeout, "CONC_B", "b"),
	).(*Aggregate)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = ag.Error()
				_ = Is(ag, "CONC_A")
				_ = Retryable(ag)
				_ = ag.Errors()
			}
		}()
	}
	wg.Wait()
}
