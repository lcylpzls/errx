package core

// KV 是错误携带的结构化键值对，便于日志与监控输出。
type KV struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
