package errx

import (
	"fmt"
	"strings"
	"sync"
)

// Aggregate 聚合多个错误为一个。errors.Is / errors.As 可命中任一子错误。
// 构造后不可变，可安全共享与并发查询。
type Aggregate struct {
	errs []error
	once sync.Once
	msg  string
}

// Join 收集多个错误：nil 被过滤；空集合返回 nil；单个错误直接返回；
// 多个错误时返回 *Aggregate。与标准库 errors.Join 语义一致，且返回类型可展开子错误。
func Join(errs ...error) error {
	filtered := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		return &Aggregate{errs: filtered}
	}
}

// Error 返回多行错误文本（惰性缓存，重复打印零开销）。
func (a *Aggregate) Error() string {
	if a == nil {
		return ""
	}
	a.once.Do(func() {
		var b strings.Builder
		fmt.Fprintf(&b, "%d 个错误：", len(a.errs))
		for i, err := range a.errs {
			fmt.Fprintf(&b, "\n  %d) %s", i+1, err.Error())
		}
		a.msg = b.String()
	})
	return a.msg
}

// Unwrap 返回全部子错误，供 errors.Is / errors.As 展开。
func (a *Aggregate) Unwrap() []error {
	return a.errs
}

// Errors 返回子错误快照（拷贝，调用方修改不影响聚合体）。
func (a *Aggregate) Errors() []error {
	if a == nil {
		return nil
	}
	out := make([]error, len(a.errs))
	copy(out, a.errs)
	return out
}
