package errx

import "sync/atomic"

// Metrics 是 errx 运行指标快照，可接入监控面板。
type Metrics struct {
	// Constructed 结构化错误构造次数（New/Newf/Wrap/Wrapf）。
	Constructed uint64
	// Queried 错误查询次数（Is/Retryable/CodeOf/KindOf）。
	Queried uint64
	// ByKind 按 Kind 统计的构造次数（索引为 Kind 数值）。
	ByKind [256]uint64
}

var metrics struct {
	constructed atomic.Uint64
	queried     atomic.Uint64
	byKind      [256]atomic.Uint64
}

// countConstruct 记录一次错误构造。
func countConstruct(kind Kind) {
	metrics.constructed.Add(1)
	metrics.byKind[kind].Add(1)
}

// countQuery 记录一次错误查询。
func countQuery() {
	metrics.queried.Add(1)
}

// Snapshot 返回运行指标快照。
func Snapshot() Metrics {
	m := Metrics{
		Constructed: metrics.constructed.Load(),
		Queried:     metrics.queried.Load(),
	}
	for i := range metrics.byKind {
		m.ByKind[i] = metrics.byKind[i].Load()
	}
	return m
}

// ResetMetrics 清零全部指标。
func ResetMetrics() {
	metrics.constructed.Store(0)
	metrics.queried.Store(0)
	for i := range metrics.byKind {
		metrics.byKind[i].Store(0)
	}
}
