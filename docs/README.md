# errx 工业化文档

errx 是工业级、零依赖的 Go 结构化错误库（模块 `github.com/lcylpzls/errx`）。

| 文档 | 内容 |
| --- | --- |
| [architecture.md](architecture.md) | 架构、核心类型、错误链、栈捕获设计 |
| [usage.md](usage.md) | 使用指南：创建/包装/查询/日志集成/错误码注册 |
| [release.md](release.md) | 版本策略与发布流程 |
| [performance.md](performance.md) | 性能基准与调优建议 |

## 设计原则

- **标准库兼容**：`Unwrap` + `errors.Is/As`，不发明新的错误接口；
- **零依赖**：仅依赖 Go 标准库，无家族与第三方依赖；
- **不可变错误**：`WithField` 返回新实例，可安全共享；
- **可观测**：错误码、分类、字段、栈四要素齐备，`SetMetricsHook`
  可把构造/查询事件转发到外部观测底座。
