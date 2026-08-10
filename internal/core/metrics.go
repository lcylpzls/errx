package core

import "sync/atomic"

// MetricsHook 是可选的全局指标接收器。
// 实现方通常来自 metricsx/adapters/errx 适配器；未设置时
// 热路径仅多一次原子加载，不产生额外开销。
type MetricsHook interface {
	// IncCounter 增加一个计数指标。
	IncCounter(name string, labels ...string)
}

// metricsHookState 承载指标钩子，配合 atomic.Pointer 无锁读取。
type metricsHookState struct {
	hook MetricsHook
}

// metricsHook 是全局指标钩子，nil 表示关闭。
var metricsHook atomic.Pointer[metricsHookState]

// Metrics 是 errx 运行指标快照，可接入监控面板。

var metrics struct {
	constructed atomic.Uint64
	queried     atomic.Uint64
	byKind      [256]atomic.Uint64
}

// countConstruct 记录一次错误构造。
func countConstruct(kind Kind) {
	metrics.constructed.Add(1)
	metrics.byKind[kind].Add(1)
	if h := loadMetricsHook(); h != nil {
		h.IncCounter("errx.constructed", kind.String())
	}
}

// countQuery 记录一次错误查询。
func countQuery() {
	metrics.queried.Add(1)
	if h := loadMetricsHook(); h != nil {
		h.IncCounter("errx.queried")
	}
}

// SetMetricsHook 设置全局指标钩子；传 nil 或调用 ResetMetricsHook 关闭。
// 钩子仅用于外部观测，不影响错误构造与查询语义。
func SetMetricsHook(hook MetricsHook) {
	if hook == nil {
		ResetMetricsHook()
		return
	}
	metricsHook.Store(&metricsHookState{hook: hook})
}

// ResetMetricsHook 清空全局指标钩子。
func ResetMetricsHook() {
	metricsHook.Store(nil)
}

// loadMetricsHook 返回当前指标钩子（无锁）。
func loadMetricsHook() MetricsHook {
	st := metricsHook.Load()
	if st == nil {
		return nil
	}
	return st.hook
}

// Snapshot 返回运行指标快照。

// ResetMetrics 清零全部指标。
