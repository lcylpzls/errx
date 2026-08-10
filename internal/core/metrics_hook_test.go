package core

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

// fakeHook 是指标钩子测试替身。
type fakeHook struct {
	mu   sync.Mutex
	recs []string
}

func (f *fakeHook) IncCounter(name string, labels ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.recs = append(f.recs, name+":"+strings.Join(labels, "|"))
}

func (f *fakeHook) count(name string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, r := range f.recs {
		if strings.HasPrefix(r, name+":") {
			n++
		}
	}
	return n
}

func TestMetricsHookConstructAndQuery(t *testing.T) {
	hook := &fakeHook{}
	SetMetricsHook(hook)
	defer ResetMetricsHook()

	_ = NewCode("HOOK_CODE", "构造一次")
	_ = New(KindTimeout, "T", "超时")
	_ = Wrap(errors.New("cause"), KindUnavailable, "U", "不可用")
	_ = Is(errors.New("x"), "NOPE")
	_, _ = CodeOf(New(KindNotFound, "N", "缺失"))

	if got := hook.count("errx.constructed"); got != 4 {
		t.Errorf("constructed 计数不符：%d", got)
	}
	if got := hook.count("errx.queried"); got != 2 {
		t.Errorf("queried 计数不符：%d", got)
	}
	foundKind := false
	hook.mu.Lock()
	for _, r := range hook.recs {
		if strings.HasPrefix(r, "errx.constructed:") && strings.HasSuffix(r, "timeout") {
			foundKind = true
		}
	}
	hook.mu.Unlock()
	if !foundKind {
		t.Errorf("constructed 应携带 kind 标签：%v", hook.recs)
	}
}

func TestMetricsHookReset(t *testing.T) {
	hook := &fakeHook{}
	SetMetricsHook(hook)
	_ = NewCode("A", "a")
	if got := hook.count("errx.constructed"); got != 1 {
		t.Fatalf("钩子未生效：%d", got)
	}
	ResetMetricsHook()
	_ = NewCode("B", "b")
	if got := hook.count("errx.constructed"); got != 1 {
		t.Errorf("重置后不应再计数：%d", got)
	}
	// 显式 nil 同样关闭。
	SetMetricsHook(hook)
	_ = NewCode("C", "c")
	SetMetricsHook(nil)
	_ = NewCode("D", "d")
	if got := hook.count("errx.constructed"); got != 2 {
		t.Errorf("nil 关闭后不应再计数：%d", got)
	}
}

func TestMetricsHookConcurrent(t *testing.T) {
	hook := &fakeHook{}
	SetMetricsHook(hook)
	defer ResetMetricsHook()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = New(KindBusiness, "B", "业务")
				_, _ = CodeOf(NewCode("C", "c"))
				if j%10 == 0 {
					SetMetricsHook(hook)
					ResetMetricsHook()
					SetMetricsHook(hook)
				}
			}
		}()
	}
	wg.Wait()
	if got := hook.count("errx.constructed"); got < 800 {
		t.Errorf("并发构造计数过少：%d", got)
	}
}
