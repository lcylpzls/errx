# 使用指南

## 1. 创建错误

```go
// 显式指定分类
err := errx.New(errx.KindNotFound, "USER_NOT_FOUND", "用户不存在")

// 格式化消息
err = errx.Newf(errx.KindInvalid, "BAD_ARG", "参数 %s 无效", "id")

// 基于已注册错误码自动使用分类
err = errx.NewCode("ORDER_FAIL", "下单失败")
err = errx.NewCodef("ORDER_FAIL", "订单 %s 失败", "10086")
```

## 2. 包装底层错误

```go
if err != nil {
    return errx.Wrap(err, errx.KindUnavailable, "DB_DOWN", "数据库不可用")
}
// Wrap/Wrapf 对 nil 返回 nil，可直接 return

// 基于已注册错误码包装，自动使用注册分类
return errx.WrapCode(err, "DB_DOWN", "数据库不可用")
```

## 3. 附加结构化字段

```go
err = errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").
    WithField("order_id", "10086").
    WithField("user_id", "42")

// 对任意 error 附加字段（非 errx 错误自动包装为 UNKNOWN）
err = errx.WithField(io.ErrUnexpectedEOF, "source", "upload")
```

## 4. 查询与决策

```go
if errx.Is(err, "ORDER_FAIL") { /* 业务处理 */ }
if errx.Retryable(err) { backoff.Retry(fn) }
code, ok := errx.CodeOf(err)
kind := errx.KindOf(err)
first, ok := errx.As(err)
```

## 5. 错误码注册

```go
func init() {
    errx.RegisterCode("ORDER_FAIL", "下单失败")
    errx.RegisterCode("DB_DOWN", "数据库不可用")
    errx.RegisterCodeKind("ORDER_FAIL", errx.KindBusiness)
}

// 查询说明与错误码清单
fmt.Println(errx.Describe("ORDER_FAIL"))
fmt.Println(errx.CodeKind("ORDER_FAIL"))
for _, info := range errx.Codes() {
    fmt.Printf("%s\t%s\n", info.Code, info.Description)
}
```

## 6. 输出与调试

```go
fmt.Println(err)          // ORDER_FAIL: 下单失败
fmt.Printf("%+v\n", err)  // 追加创建点调用栈

// 程序化读取调用栈（日志/监控场景）
if e, ok := errx.As(err); ok {
    for _, frame := range e.StackTrace() {
        fmt.Printf("%s:%d  %s\n", frame.File, frame.Line, frame.Function)
    }
}

// 关闭/限制栈捕获（构造频率极高时）
errx.SetStackCapture(false)
errx.SetStackDepth(16)
```

`StackTrace()` 返回 `[]errx.StackFrame`；`SetStackCapture(false)` 关闭后返回 nil。

## 7. 与 logx 集成

```go
import (
    "github.com/lcylpzls/errx"
    "github.com/lcylpzls/logx"
)

logger, _ := logx.NewBuilder().
    EnableConsole(logx.InfoLevel).
    Build()

err := errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").
    WithField("order_id", "10086")

logger.Error("业务失败", logx.FieldsFromError(err))
// err.code=ORDER_FAIL, err.kind=business, err.message=下单失败, order_id=10086
```

## 8. HTTP 状态映射与 JSON 输出

```go
import (
    "github.com/lcylpzls/errx"
    "github.com/lcylpzls/httpx"
)

// 仅取状态码
code := errx.KindHTTPStatus(errx.KindOf(err)) // not_found → 404

// 直接输出 JSON 错误响应体（code/kind/message）
httpx.WriteErrorJSON(w, err)
```

`httpx.WriteErrorJSON` 对 nil 与普通错误安全：nil 输出 500 + 未知分类，
普通错误 `code` 回退为 `UNKNOWN`，`kind` 回退为 `unknown`。

## 9. 生产建议

- 错误码启动期统一注册，供前端/API 网关映射；
- 外部可重试场景用 `KindTimeout / KindRateLimited / KindUnavailable`；
- 用户可见消息与内部消息分离：内部 `msg` 可含细节，输出给用户时自行映射；
- 错误构造频率极高且不需要栈时可 `errx.SetStackCapture(false)`；
- 需要指标观测时实现 `errx.MetricsHook` 并注入 metricsx 等外部底座。

## 10. 错误分类与策略

errx 提供三级体系：**Kind（细分枚举）→ Category（领域分组）→ Policy（处理策略）**。

```go
// Kind：17 类细分错误，对齐 Google API / gRPC 错误模型
errx.New(errx.KindAlreadyExists, "USER_EXISTS", "用户已存在")
errx.New(errx.KindForbidden, "NO_PERMISSION", "无权限")
errx.New(errx.KindQuotaExceeded, "QUOTA", "配额耗尽")

// Category：领域分组
fmt.Println(errx.KindForbidden.Category()) // 认证与授权

// Policy：处理策略（可重试 / 告警 / 用户可见）
p := errx.KindUnavailable.Policy()
if p.Retryable { backoff.Retry(fn) }
if p.Alert { alert.Send() }
```

## 11. 多错误聚合

```go
err := errx.Join(
    errx.New(errx.KindBusiness, "A1", "错误一"),
    errx.New(errx.KindTimeout, "A2", "错误二"),
)
if errx.Is(err, "A1") { /* 命中子错误 */ }
if errx.Retryable(err) { /* 聚合内存在可重试子错误 */ }

// 聚合错误同样支持 JSON 跨服务传输
data, _ := json.Marshal(err)
var restored errx.Aggregate
json.Unmarshal(data, &restored)
```

## 12. 跨服务传输

```go
// 服务端：错误直接 JSON 序列化（含原因链与字段）
data, _ := json.Marshal(err)

// 客户端：完整还原结构化错误
var restored errx.Error
json.Unmarshal(data, &restored)
if errx.Is(&restored, "ORDER_FAIL") { /* 按错误码处理 */ }
```

调用栈不跨服务传输（序列化时省略），字段与原因链完整保留。
`Aggregate` 通过 `{"errors":[...]}` 传输子错误数组，还原后
`errors.Is` 可命中任一子错误。

## 13. 观测指标

```go
type myHook struct{}

func (myHook) IncCounter(name string, labels ...string) {
    metrics.Inc(name, labels...) // 转发到 metricsx 等底座
}

errx.SetMetricsHook(myHook{})
defer errx.ResetMetricsHook()
```

`MetricsHook` 仅接收错误构造/查询事件，不影响错误语义；
未设置钩子时热路径仅多一次原子加载，无额外开销。
