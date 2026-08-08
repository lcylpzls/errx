# 使用指南

## 1. 创建错误

```go
// 无原因
err := errx.New(errx.KindNotFound, "USER_NOT_FOUND", "用户不存在")

// 格式化消息
err = errx.Newf(errx.KindInvalid, "BAD_ARG", "参数 %s 无效", "id")
```

## 2. 包装底层错误

```go
if err != nil {
    return errx.Wrap(err, errx.KindUnavailable, "DB_DOWN", "数据库不可用")
}
// Wrap/Wrapf 对 nil 返回 nil，可直接 return
```

## 3. 附加结构化字段

```go
err = errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").
    WithField("order_id", "10086").
    WithField("user_id", "42")
```

## 4. 查询与决策

```go
if errx.Is(err, "ORDER_FAIL") { /* 业务处理 */ }
if errx.Retryable(err) { backoff.Retry(fn) }
code, _ := errx.CodeOf(err)
kind := errx.KindOf(err)
```

## 5. 错误码注册与文档生成

```go
func init() {
    errx.RegisterCode("ORDER_FAIL", "下单失败")
    errx.RegisterCode("DB_DOWN", "数据库不可用")
}

// 生成错误码清单
for _, info := range errx.Codes() {
    fmt.Printf("%s\t%s\n", info.Code, info.Description)
}
```

## 6. 输出与调试

```go
fmt.Println(err)          // ORDER_FAIL: 下单失败
fmt.Printf("%+v\n", err)  // 追加创建点调用栈
```

## 7. 与 logx 集成

```go
import (
    "github.com/lcylpzls/errx"
    errxlogx "github.com/lcylpzls/errx/logx"
    "github.com/lcylpzls/logx"
)

logger, _ := logx.NewBuilder().
    EnableConsole(logx.InfoLevel).
    Build()

err := errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").
    WithField("order_id", "10086")

logger.Error("业务失败", errxlogx.Fields(err))
// err.code=ORDER_FAIL, err.kind=business, err.message=下单失败, order_id=10086
```

## 8. 生产建议

- 错误码启动期统一注册，生成文档供前端/API 网关映射；
- 外部可重试场景用 `KindTimeout / KindRateLimited / KindUnavailable`；
- 用户可见消息与内部消息分离：内部 `msg` 可含细节，输出给用户时自行映射；
- 错误构造频率极高且不需要栈时可 `errx.SetStackCapture(false)`。

## 9. 错误分类与策略

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

// 生成按领域分组的错误分类表
fmt.Println(errx.KindsMarkdown())
```

## 10. 多错误聚合

```go
err := errx.Join(
    errx.New(errx.KindBusiness, "A1", "错误一"),
    errx.New(errx.KindTimeout, "A2", "错误二"),
)
if errx.Is(err, "A1") { /* 命中子错误 */ }
if errx.Retryable(err) { /* 聚合内存在可重试子错误 */ }
```

## 11. 跨服务传输与 HTTP 映射

```go
// 服务端：错误直接 JSON 序列化（含原因链与字段）
data, _ := json.Marshal(err)

// 客户端：完整还原结构化错误
var restored errx.Error
json.Unmarshal(data, &restored)
if errx.Is(&restored, "ORDER_FAIL") { /* 按错误码处理 */ }

// HTTP 状态映射（网关/中间件）
w.WriteHeader(err.HTTPStatus())
```

调用栈不跨服务传输（序列化时省略），字段与原因链完整保留。
