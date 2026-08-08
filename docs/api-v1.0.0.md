package errx // import "github.com/lcylpzls/errx"


FUNCTIONS

func Describe(code Code) string
    Describe 返回错误码的注册说明；未注册时返回空字符串。

func Is(err error, code Code) bool
    Is 判断错误链中是否存在指定错误码（支持单链与 Aggregate 多错误展开）。

func Join(errs ...error) error
    Join 收集多个错误：nil 被过滤；空集合返回 nil；单个错误直接返回； 多个错误时返回 *Aggregate。与标准库 errors.Join
    语义一致，且返回类型可展开子错误。

func KindHTTPStatus(kind Kind) int
    KindHTTPStatus 返回 Kind 对应的 HTTP 状态码。

func KindsMarkdown() string
    KindsMarkdown 生成按领域分组、含策略标注的错误分类表 Markdown。

func Markdown() string
    Markdown 生成错误码注册表的 Markdown 文档（按错误码排序）， 可直接交付前端、API 网关或审计使用。

func RegisterCode(code Code, description string)
    RegisterCode 注册错误码及其说明。重复注册以最后一次为准。 建议在包 init 或程序启动阶段完成注册。

func ResetMetrics()
    ResetMetrics 清零全部指标。

func Retryable(err error) bool
    Retryable 判断错误链中是否存在可重试分类（支持单链与 Aggregate 多错误展开）。

func SetStackCapture(enabled bool)
    SetStackCapture 全局开关调用栈捕获。生产环境如对错误构造频率敏感可关闭。

func SetStackDepth(depth int)
    SetStackDepth 设置栈捕获的最大帧数；depth <= 0 时恢复默认 32。

func WithField(err error, key string, val any) error
    WithField 为任意错误附加结构化字段： 已是 *Error 时返回新实例；否则包装为 UNKNOWN 错误并保留原因链。


TYPES

type Aggregate struct {
	// Has unexported fields.
}
    Aggregate 聚合多个错误为一个。errors.Is / errors.As 可命中任一子错误。 构造后不可变，可安全共享与并发查询。

func (a *Aggregate) Error() string
    Error 返回多行错误文本（惰性缓存，重复打印零开销）。

func (a *Aggregate) Errors() []error
    Errors 返回子错误快照（拷贝，调用方修改不影响聚合体）。

func (a *Aggregate) Unwrap() []error
    Unwrap 返回全部子错误，供 errors.Is / errors.As 展开。

type Category uint8
    Category 是 Kind 的领域分组，便于按场景组织错误与阅读错误表。

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
func (c Category) String() string
    String 返回分类的中文名称。

type Code string
    Code 是稳定的业务错误码，建议使用大写下划线风格（如 "USER_NOT_FOUND"）。

const CodeUnknown Code = "UNKNOWN"
    CodeUnknown 是未指定错误码时的默认值。

func CodeOf(err error) (Code, bool)
    CodeOf 返回错误链中第一个结构化错误的错误码。

type CodeInfo struct {
	Code        Code
	Description string
}
    CodeInfo 描述一个错误码及其含义，用于文档生成与审计。

func Codes() []CodeInfo
    Codes 返回全部已注册错误码的快照，按错误码排序。

type Error struct {
	// Has unexported fields.
}
    Error 是结构化错误：携带错误码、分类、消息、结构化字段与可选调用栈。 通过 errors.Is / errors.As 与标准库错误链完全兼容。

func As(err error) (*Error, bool)
    As 从错误链中取出第一个 *Error。

func New(kind Kind, code Code, msg string) *Error
    New 创建一个无原因的结构化错误。

func Newf(kind Kind, code Code, format string, args ...any) *Error
    Newf 创建带格式化消息的结构化错误。

func Wrap(err error, kind Kind, code Code, msg string) *Error
    Wrap 包装一个底层错误并附加分类与错误码。 当 err 为 nil 时返回 nil，便于直接 return errx.Wrap(err,
    ...) 的写法。

func Wrapf(err error, kind Kind, code Code, format string, args ...any) *Error
    Wrapf 包装底层错误并附加格式化消息。

func (e *Error) Cause() error
    Cause 返回被包装的底层错误（兼容 pkg/errors 习惯）。

func (e *Error) Code() Code
    Code 返回错误码。

func (e *Error) Error() string
    Error 返回格式为 "CODE: message: cause" 的文本；空字段自动省略。 结果惰性缓存，重复打印零额外开销。

func (e *Error) Fields() []KV
    Fields 返回错误携带的结构化字段快照。

func (e *Error) Format(f fmt.State, verb rune)
    Format 支持 %s / %q / %v；%+v 额外输出创建时捕获的调用栈。

func (e *Error) HTTPStatus() int
    HTTPStatus 返回该错误对应的 HTTP 状态码。

func (e *Error) Is(target error) bool
    Is 支持 errors.Is 按错误码或可重试分类匹配（沿错误链与聚合子错误展开）。

func (e *Error) Kind() Kind
    Kind 返回错误分类。

func (e *Error) MarshalJSON() ([]byte, error)
    MarshalJSON 将 Error 序列化为跨服务可传输的 JSON（含原因链与字段）。

func (e *Error) Message() string
    Message 返回错误消息（不含错误码与原因）。

func (e *Error) UnmarshalJSON(data []byte) error
    UnmarshalJSON 从 JSON 恢复 Error。调用栈不跨服务传输。

func (e *Error) Unwrap() error
    Unwrap 返回被包装的底层错误，支持 errors.Is / errors.As 链路。

func (e *Error) WithField(key string, val any) *Error
    WithField 返回携带附加字段的新错误（不可变风格，原错误不受影响）。

type KV struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
    KV 是错误携带的结构化键值对，便于日志与监控输出。

type Kind uint8
    Kind 是错误分类，用于驱动重试、告警与用户提示策略。 枚举对齐 Google API / gRPC 错误模型的主流场景。

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
func KindOf(err error) Kind
    KindOf 返回错误链中第一个结构化错误的分类；无结构化错误时返回 KindUnknown。

func (k Kind) Category() Category
    Category 返回 Kind 所属的领域分组。

func (k Kind) Policy() Policy
    Policy 返回该分类的错误处理策略。

func (k Kind) Retryable() bool
    Retryable 判断该分类是否建议重试（委托 Policy）。

func (k Kind) String() string
    String 返回 Kind 的稳定小写名称，用于日志与监控打点。

type Metrics struct {
	// Constructed 结构化错误构造次数（New/Newf/Wrap/Wrapf）。
	Constructed uint64
	// Queried 错误查询次数（Is/Retryable/CodeOf/KindOf）。
	Queried uint64
	// ByKind 按 Kind 统计的构造次数（索引为 Kind 数值）。
	ByKind [256]uint64
}
    Metrics 是 errx 运行指标快照，可接入监控面板。

func Snapshot() Metrics
    Snapshot 返回运行指标快照。

type Policy struct {
	// Retryable 是否建议重试。
	Retryable bool
	// Alert 是否应触发告警。
	Alert bool
	// UserVisible 是否适合直接展示给用户。
	UserVisible bool
}
    Policy 是 Kind 对应的错误处理策略。


<!-- v1.0.0 API 冻结基线；生成方式：go doc -all . -->
