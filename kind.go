package errx

import (
	"fmt"
	"strings"
)

// Kind 是错误分类，用于驱动重试、告警与用户提示策略。
// 枚举对齐 Google API / gRPC 错误模型的主流场景。
type Kind uint8

const (
	// KindUnknown 未知分类。
	KindUnknown Kind = iota
	// KindInvalid 输入/参数无效，重试无意义。
	KindInvalid
	// KindNotFound 资源不存在。
	KindNotFound
	// KindAlreadyExists 资源已存在。
	KindAlreadyExists
	// KindUnauthorized 未认证（401）。
	KindUnauthorized
	// KindForbidden 已认证但无权限（403）。
	KindForbidden
	// KindConflict 状态冲突（如并发修改，409）。
	KindConflict
	// KindCancelled 操作被取消。
	KindCancelled
	// KindDeadlineExceeded 整体截止时间已过，重试无意义。
	KindDeadlineExceeded
	// KindTimeout 操作超时，可重试。
	KindTimeout
	// KindRateLimited 触发限流，稍后可重试（429）。
	KindRateLimited
	// KindQuotaExceeded 配额/资源耗尽，等待释放后可重试。
	KindQuotaExceeded
	// KindUnavailable 依赖或系统暂不可用，可重试（503）。
	KindUnavailable
	// KindInternal 内部错误，不应暴露细节，应告警（500）。
	KindInternal
	// KindNotImplemented 功能未实现（501）。
	KindNotImplemented
	// KindDataLoss 数据丢失或损坏，应告警。
	KindDataLoss
	// KindBusiness 业务规则错误，重试无意义。
	KindBusiness
)

var kindNames = map[Kind]string{
	KindUnknown:          "unknown",
	KindInvalid:          "invalid_argument",
	KindNotFound:         "not_found",
	KindAlreadyExists:    "already_exists",
	KindUnauthorized:     "unauthorized",
	KindForbidden:        "forbidden",
	KindConflict:         "conflict",
	KindCancelled:        "cancelled",
	KindDeadlineExceeded: "deadline_exceeded",
	KindTimeout:          "timeout",
	KindRateLimited:      "rate_limited",
	KindQuotaExceeded:    "quota_exceeded",
	KindUnavailable:      "unavailable",
	KindInternal:         "internal",
	KindNotImplemented:   "not_implemented",
	KindDataLoss:         "data_loss",
	KindBusiness:         "business",
}

// String 返回 Kind 的稳定小写名称，用于日志与监控打点。
func (k Kind) String() string {
	if s, ok := kindNames[k]; ok {
		return s
	}
	return "unknown"
}

// Category 是 Kind 的领域分组，便于按场景组织错误与阅读错误表。
type Category uint8

const (
	// CatInput 输入与参数校验。
	CatInput Category = iota
	// CatAuth 认证与授权。
	CatAuth
	// CatState 资源与状态。
	CatState
	// CatDependency 依赖与外部约束。
	CatDependency
	// CatSystem 系统内部。
	CatSystem
	// CatBusiness 业务规则。
	CatBusiness
)

var categoryNames = map[Category]string{
	CatInput:      "输入与参数",
	CatAuth:       "认证与授权",
	CatState:      "资源与状态",
	CatDependency: "依赖与外部",
	CatSystem:     "系统内部",
	CatBusiness:   "业务规则",
}

// String 返回分类的中文名称。
func (c Category) String() string {
	if s, ok := categoryNames[c]; ok {
		return s
	}
	return "未知"
}

// Category 返回 Kind 所属的领域分组。
func (k Kind) Category() Category {
	switch k {
	case KindInvalid:
		return CatInput
	case KindUnauthorized, KindForbidden:
		return CatAuth
	case KindNotFound, KindAlreadyExists, KindConflict, KindCancelled:
		return CatState
	case KindDeadlineExceeded, KindTimeout, KindRateLimited, KindQuotaExceeded, KindUnavailable:
		return CatDependency
	case KindInternal, KindNotImplemented, KindDataLoss:
		return CatSystem
	case KindBusiness:
		return CatBusiness
	default:
		return CatSystem
	}
}

// Policy 是 Kind 对应的错误处理策略。
type Policy struct {
	// Retryable 是否建议重试。
	Retryable bool
	// Alert 是否应触发告警。
	Alert bool
	// UserVisible 是否适合直接展示给用户。
	UserVisible bool
}

// Policy 返回该分类的错误处理策略。
func (k Kind) Policy() Policy {
	p := Policy{UserVisible: true}
	switch k {
	case KindTimeout, KindRateLimited, KindQuotaExceeded, KindUnavailable:
		p.Retryable = true
	}
	switch k {
	case KindInternal, KindDataLoss, KindUnavailable:
		p.Alert = true
	}
	switch k {
	case KindUnknown, KindInternal, KindDataLoss:
		p.UserVisible = false
	}
	return p
}

// Retryable 判断该分类是否建议重试（委托 Policy）。
func (k Kind) Retryable() bool {
	return k.Policy().Retryable
}

// KindsMarkdown 生成按领域分组、含策略标注的错误分类表 Markdown。
func KindsMarkdown() string {
	categories := []Category{CatInput, CatAuth, CatState, CatDependency, CatSystem, CatBusiness}
	var b strings.Builder
	b.WriteString("# 错误分类表\n\n")
	for _, cat := range categories {
		fmt.Fprintf(&b, "## %s\n\n", cat.String())
		b.WriteString("| Kind | 可重试 | 告警 | 用户可见 |\n")
		b.WriteString("| --- | --- | --- | --- |\n")
		for k := Kind(1); k <= KindBusiness; k++ {
			if k.Category() != cat {
				continue
			}
			p := k.Policy()
			fmt.Fprintf(&b, "| %s | %v | %v | %v |\n",
				k.String(), p.Retryable, p.Alert, p.UserVisible)
		}
		b.WriteString("\n")
	}
	return b.String()
}
