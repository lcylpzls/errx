package errx

// Kind 是错误分类，用于驱动重试、告警与用户提示策略。
type Kind uint8

const (
	// KindUnknown 未知分类。
	KindUnknown Kind = iota
	// KindInvalid 输入/参数无效，重试无意义。
	KindInvalid
	// KindNotFound 资源不存在。
	KindNotFound
	// KindUnauthorized 未认证或无权访问。
	KindUnauthorized
	// KindConflict 状态冲突（如并发修改）。
	KindConflict
	// KindTimeout 操作超时，可重试。
	KindTimeout
	// KindRateLimited 触发限流，稍后可重试。
	KindRateLimited
	// KindUnavailable 依赖或系统暂不可用，可重试。
	KindUnavailable
	// KindInternal 内部错误，重试无意义但应告警。
	KindInternal
	// KindBusiness 业务规则错误，重试无意义。
	KindBusiness
)

var kindNames = map[Kind]string{
	KindUnknown:      "unknown",
	KindInvalid:      "invalid_argument",
	KindNotFound:     "not_found",
	KindUnauthorized: "unauthorized",
	KindConflict:     "conflict",
	KindTimeout:      "timeout",
	KindRateLimited:  "rate_limited",
	KindUnavailable:  "unavailable",
	KindInternal:     "internal",
	KindBusiness:     "business",
}

// String 返回 Kind 的稳定小写名称，用于日志与监控打点。
func (k Kind) String() string {
	if s, ok := kindNames[k]; ok {
		return s
	}
	return "unknown"
}

// Retryable 判断该分类是否建议重试。
func (k Kind) Retryable() bool {
	switch k {
	case KindTimeout, KindRateLimited, KindUnavailable:
		return true
	default:
		return false
	}
}
