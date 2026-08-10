package core

import (
	"errors"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(KindBusiness, "BENCH", "基准")
	}
}

func BenchmarkWrap(b *testing.B) {
	base := errors.New("base")
	for i := 0; i < b.N; i++ {
		_ = Wrap(base, KindInternal, "BENCH", "包装")
	}
}

func BenchmarkErrorString(b *testing.B) {
	e := New(KindBusiness, "BENCH", "基准")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkIs(b *testing.B) {
	e := New(KindBusiness, "BENCH", "基准")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Is(e, "BENCH")
	}
}

func BenchmarkRetryable(b *testing.B) {
	e := New(KindTimeout, "BENCH", "超时")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Retryable(e)
	}
}
