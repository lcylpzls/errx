# errx

工业级、零依赖的 Go 结构化错误库：错误码 + 分类 + 结构化字段 + 调用栈，与标准库错误链完全兼容。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 快速开始

```go
import "github.com/lcylpzls/errx"

// 创建结构化错误
return errx.New(errx.KindNotFound, "USER_NOT_FOUND", "用户不存在")

// 包装底层错误
if err != nil {
    return errx.Wrap(err, errx.KindUnavailable, "DB_DOWN", "数据库不可用")
}

// 附加结构化字段
return errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").
    WithField("order_id", "10086")
```

## 核心特性

- **错误码**：`Code` 字符串 + 注册表（`RegisterCode` / `Describe` / `Codes`），支持文档生成；
- **错误分类**：`Kind` 17 类细分枚举 + `Category` 领域分组 + `Policy` 策略（可重试/告警/用户可见）；
- **多错误聚合**：`Join` 聚合多个错误，`errors.Is/As` 命中任一子错误；
- **跨服务传输**：`Error` 原生 JSON 序列化/还原，`HTTPStatus()` 直接映射 HTTP 状态码；
- **可观测**：`Snapshot()` 输出构造/查询计数与按 Kind 分布，零锁原子实现；
- **HTTP 适配**：`errx/httpx` 一键输出状态码与 JSON 错误响应体；
- **标准库兼容**：`Unwrap` / `errors.Is` / `errors.As` 全链路支持；
- **结构化字段**：`WithField` 不可变追加，随错误传递业务上下文；
- **调用栈**：创建时捕获（可全局开关），`fmt.Printf("%+v", err)` 输出完整栈；
- **零依赖核心**：errx 核心包仅依赖 Go 标准库。

## 与 logx 集成

```go
import (
    "github.com/lcylpzls/errx"
    errxlogx "github.com/lcylpzls/errx/logx"
    "github.com/lcylpzls/logx"
)

err := errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").
    WithField("order_id", "10086")

logger.Error("业务失败", errxlogx.Fields(err))
// 输出字段：err.code=ORDER_FAIL, err.kind=business, err.message=下单失败, order_id=10086
```

## HTTP 适配

```go
import "github.com/lcylpzls/errx/httpx"

func handler(w http.ResponseWriter, r *http.Request) {
    if err := doBusiness(); err != nil {
        // not_found → 404，响应体为 {"code","kind","message"}
        httpx.WriteJSON(w, err)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

`WriteJSON` 对 nil 与普通错误同样安全：nil 输出 500 + 未知分类，普通错误的
`code` 回退为 `UNKNOWN`。

## 文档

- [docs/architecture.md](docs/architecture.md) — 架构与设计
- [docs/usage.md](docs/usage.md) — 使用指南
- [docs/release.md](docs/release.md) — 版本与发布

## License

MIT © [lcylpzls](https://github.com/lcylpzls)
