package errx

import (
	"errors"
	"sync"
	"testing"
)

func TestMetricsConstructed(t *testing.T) {
	ResetMetrics()

	New(KindBusiness, "M1", "a")
	Newf(KindInvalid, "M2", "%s", "b")
	Wrap(errors.New("x"), KindTimeout, "M3", "c")
	Wrapf(errors.New("x"), KindUnavailable, "M4", "%s", "d")

	m := Snapshot()
	if m.Constructed != 4 {
		t.Errorf("Constructed 不符：%d", m.Constructed)
	}
	for _, tc := range []struct {
		kind Kind
		n    uint64
	}{
		{KindBusiness, 1},
		{KindInvalid, 1},
		{KindTimeout, 1},
		{KindUnavailable, 1},
	} {
		if m.ByKind[tc.kind] != tc.n {
			t.Errorf("ByKind[%d] 不符：got %d, want %d", tc.kind, m.ByKind[tc.kind], tc.n)
		}
	}
}

func TestMetricsQueried(t *testing.T) {
	ResetMetrics()

	e := New(KindBusiness, "MQ", "q")
	_ = Is(e, "MQ")
	_ = Retryable(e)
	_, _ = CodeOf(e)
	_ = KindOf(e)

	if m := Snapshot(); m.Queried != 4 {
		t.Errorf("Queried 不符：%d", m.Queried)
	}
}

func TestResetMetrics(t *testing.T) {
	New(KindBusiness, "R", "r")
	_ = Is(New(KindBusiness, "R2", "r2"), "R2")
	ResetMetrics()
	m := Snapshot()
	if m.Constructed != 0 || m.Queried != 0 {
		t.Errorf("Reset 后应清零：%+v", m)
	}
	for i := range m.ByKind {
		if m.ByKind[i] != 0 {
			t.Errorf("ByKind[%d] 未清零", i)
		}
	}
}

func TestMetricsConcurrent(t *testing.T) {
	ResetMetrics()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				e := New(KindBusiness, "MC", "c")
				_ = Is(e, "MC")
				_ = Snapshot()
			}
		}(i)
	}
	wg.Wait()
	m := Snapshot()
	if m.Constructed != 800 || m.Queried != 800 {
		t.Errorf("并发计数不符：%+v", m)
	}
}
