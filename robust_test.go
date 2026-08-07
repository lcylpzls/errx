package errx

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// 错误链深度上限
// ---------------------------------------------------------------------------

func chain(depth int, leafCode Code) error {
	err := New(KindBusiness, leafCode, "叶子")
	for i := 1; i < depth; i++ {
		err = Wrap(err, KindInternal, Code(fmt.Sprintf("W%d", i)), "w")
	}
	return err
}

func TestIsChainDepthLimit(t *testing.T) {
	// 100 层内应命中叶子
	if !Is(chain(maxChainDepth, "DEEP_IN"), "DEEP_IN") {
		t.Error("深度上限内应命中")
	}
	// 第 101 层不命中（防御死循环，不 panic）
	outer := Wrap(chain(maxChainDepth, "DEEP_OUT"), KindInternal, "TOO_DEEP", "超出")
	if Is(outer, "DEEP_OUT") {
		t.Error("超过深度上限不应命中")
	}
	if !Is(outer, "TOO_DEEP") {
		t.Error("最外层应命中")
	}
}

func TestRetryableChainDepthLimit(t *testing.T) {
	// 可重试节点在深度上限内
	inner := New(KindTimeout, "R_IN", "超时")
	err := inner
	for i := 1; i < maxChainDepth; i++ {
		err = Wrap(err, KindInternal, Code(fmt.Sprintf("W%d", i)), "w")
	}
	if !Retryable(err) {
		t.Error("深度上限内应可重试")
	}
	// 可重试节点超过深度上限（叶子在第 101 层，外层均不可重试）
	leaf := New(KindTimeout, "R_101", "超时")
	deep := leaf
	for i := 0; i < maxChainDepth; i++ {
		deep = Wrap(deep, KindInternal, Code(fmt.Sprintf("W%d", i)), "w")
	}
	if Retryable(deep) {
		t.Error("可重试节点超出深度上限后不应判定可重试")
	}
}

// ---------------------------------------------------------------------------
// Error() 惰性缓存
// ---------------------------------------------------------------------------

func TestErrorStringCached(t *testing.T) {
	e := New(KindBusiness, "CACHE", "缓存")
	first := e.Error()
	second := e.Error()
	if first != second {
		t.Errorf("Error() 两次结果不一致：%q %q", first, second)
	}
	// 字段不影响 Error() 文本；WithField 重建后缓存重新计算且结果一致
	e2 := e.WithField("k", "v")
	if e2.Error() != first {
		t.Errorf("WithField 后 Error() 文本不应变化：%q", e2.Error())
	}
}

// ---------------------------------------------------------------------------
// 并发专项测试（配合 CI 的 -race 验证）
// ---------------------------------------------------------------------------

func TestConcurrentRegisterAndDescribe(t *testing.T) {
	const workers = 8
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code := Code(fmt.Sprintf("CONC_%d", i))
			for j := 0; j < 200; j++ {
				RegisterCode(code, "并发注册")
				_ = Describe(code)
				_ = Codes()
			}
		}(i)
	}
	wg.Wait()
	if got := Describe("CONC_0"); got != "并发注册" {
		t.Errorf("并发注册后描述不符：%s", got)
	}
}

func TestConcurrentSetStackCapture(t *testing.T) {
	const workers = 8
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(enable bool) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				SetStackCapture(enable)
				e := New(KindInternal, "CSC", "并发栈")
				_ = e.Error()
				_ = fmt.Sprintf("%+v", e)
				_ = Is(e, "CSC")
				_ = Retryable(e)
			}
		}(i%2 == 0)
	}
	wg.Wait()
	SetStackCapture(true)
}

func TestConcurrentErrorUsage(t *testing.T) {
	const workers = 8
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				e := New(KindBusiness, "CONC_ERR", "并发错误").
					WithField("i", i).
					WithField("j", j)
				_ = e.Error()
				_ = fmt.Sprintf("%+v", e)
				_, _ = CodeOf(e)
				_ = KindOf(e)
				_ = Is(e, "CONC_ERR")
				_ = Retryable(e)
				_ = WithField(errors.New("普通"), "k", "v")
			}
		}(i)
	}
	wg.Wait()
}
