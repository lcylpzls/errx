package errx

import (
	"fmt"
	"testing"
)

func FuzzErrorFormat(f *testing.F) {
	f.Add("CODE", "message", "cause")
	f.Add("", "", "")
	f.Add("中文码", "中文消息", "底层")

	f.Fuzz(func(t *testing.T, code, msg, cause string) {
		e := New(KindBusiness, Code(code), msg)
		if cause != "" {
			e = Wrap(fmt.Errorf("%s", cause), KindInternal, Code(code), msg)
		}
		_ = fmt.Sprintf("%v", e)
		_ = fmt.Sprintf("%+v", e)
		_ = fmt.Sprintf("%s", e)
		_ = fmt.Sprintf("%q", e)
	})
}

func FuzzWrapChain(f *testing.F) {
	f.Add("a", "b", "c")
	f.Add("", "x", "y")

	f.Fuzz(func(t *testing.T, a, b, c string) {
		err := New(KindBusiness, Code(a), a)
		if b != "" {
			err = Wrap(err, KindTimeout, Code(b), b)
		}
		if c != "" {
			err = Wrap(err, KindUnavailable, Code(c), c)
		}
		_ = err.Error()
		_ = Retryable(err)
		_ = Is(err, Code(c))
	})
}

func FuzzIs(f *testing.F) {
	f.Add("A", "B", int(0))
	f.Add("X", "X", int(120))

	f.Fuzz(func(t *testing.T, a, b string, depth int) {
		if depth < 0 {
			depth = -depth
		}
		depth %= maxChainDepth + 20
		err := chain(depth, Code(a))
		_ = Is(err, Code(b))
		_, _ = CodeOf(err)
		_ = KindOf(err)
	})
}

func FuzzRetryable(f *testing.F) {
	f.Add("T", int(0))
	f.Add("B", int(100))

	f.Fuzz(func(t *testing.T, code string, depth int) {
		if depth < 0 {
			depth = -depth
		}
		depth %= maxChainDepth + 20
		leaf := New(KindTimeout, Code(code), "超时")
		err := leaf
		for i := 0; i < depth; i++ {
			err = Wrap(err, KindInternal, "W", "w")
		}
		_ = Retryable(err)
	})
}
