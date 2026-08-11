# errx

工业级、零依赖的 Go 结构化错误库：错误码 + 分类 + 结构化字段 + 调用栈，与标准库错误链完全兼容。

> 指标统一外置：`errx.SetMetricsHook(钩子)` 即可把错误构造/查询
> 事件转发到 metricsx 等外部观测底座（默认关闭，零额外开销）。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.26.5-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![CI](https://github.com/lcylpzls/errx/actions/workflows/ci.yml/badge.svg)](https://github.com/lcylpzls/errx/actions/workflows/ci.yml)

## 快速开始

```go
import "github.com/lcylpzls/errx"

// 注册错误码（包内 init 一次完成），之后构造无需再写分类
errx.RegisterCode("USER_NOT_FOUND", "用户不存在")
errx.RegisterCodeKind("USER_NOT_FOUND", errx.KindNotFound)

// 创建结构化错误
return errx.NewCode("USER_NOT_FOUND", "用户不存在")

// 包装底层错误
if err != nil {
    return errx.WrapCode(err, "DB_DOWN", "数据库不可用")
}

// 附加结构化字段
return errx.NewCode("ORDER_FAIL", "下单失败").
    WithField("order_id", "10086")
```

## 核心特性

- **错误码**：`Code` 字符串 + 注册表（`RegisterCode` / `Describe` / `Codes`），
  冲突检测防静默覆盖，`RegisterCodeKind` 声明分类，`NewCode` / `NewCodef` /
  `WrapCode` / `WrapCodef` 免写分类便捷构造；
- **错误分类**：`Kind` 17 类细分枚举 + `Category` 领域分组 + `Policy` 策略（可重试/告警/用户可见）；
- **多错误聚合**：`Join` 聚合多个错误，`errors.Is/As` 命中任一子错误，支持 JSON 序列化；
- **跨服务传输**：`Error` 原生 JSON 序列化/还原，`KindHTTPStatus` 直接映射 HTTP 状态码；
- **可观测**：`SetMetricsHook` 把错误构造/查询事件转发到外部观测底座（默认关闭）；
- **HTTP 映射**：`KindHTTPStatus(kind)` 直接映射 HTTP 状态码，配合 httpx 输出 JSON 错误响应体；
- **标准库兼容**：`Unwrap` / `errors.Is` / `errors.As` 全链路支持；
- **结构化字段**：`WithField` 不可变追加，随错误传递业务上下文；
- **调用栈**：创建时捕获（可全局开关），`fmt.Printf("%+v", err)` 输出完整栈，`StackTrace()` 程序化读取；
- **零依赖**：errx 仅依赖 Go 标准库，无任何家族或第三方依赖。

## 与 logx 集成

```go
import (
    "github.com/lcylpzls/logx"
    "github.com/lcylpzls/errx"
)

err := errx.NewCode("ORDER_FAIL", "下单失败").
    WithField("order_id", "10086")

logger.Error("业务失败", logx.FieldsFromError(err))
// 输出字段：err.code=ORDER_FAIL, err.kind=business, err.message=下单失败, order_id=10086
```

## HTTP 映射与输出

```go
import (
    "github.com/lcylpzls/errx"
    "github.com/lcylpzls/httpx"
)

func handler(w http.ResponseWriter, r *http.Request) {
    if err := doBusiness(); err != nil {
        // not_found → 404，响应体为 {"code","kind","message"}
        httpx.WriteErrorJSON(w, err)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

`httpx.WriteErrorJSON` 对 nil 与普通错误同样安全：nil 输出 500 + 未知分类，
普通错误的 `code` 回退为 `UNKNOWN`；仅做状态映射时可直接使用
`httpx.ErrorStatus(err)` 或 `errx.KindHTTPStatus(errx.KindOf(err))`。

## 文档

- [docs/README.md](docs/README.md) — 文档索引
- [docs/architecture.md](docs/architecture.md) — 架构与设计
- [docs/usage.md](docs/usage.md) — 使用指南
- [docs/release.md](docs/release.md) — 版本与发布

## License

MIT © [lcylpzls](https://github.com/lcylpzls)
