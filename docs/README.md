# errx 工业化文档

errx 是工业级、零依赖的 Go 结构化错误库（模块 `github.com/lcylpzls/errx`）。

| 文档 | 内容 |
| --- | --- |
| [architecture.md](architecture.md) | 架构、核心类型、错误链、栈捕获设计 |
| [usage.md](usage.md) | 使用指南：创建/包装/查询/日志集成/错误码注册 |
| [release.md](release.md) | 版本策略与发布流程 |

## 设计原则

- **标准库兼容**：`Unwrap` + `errors.Is/As`，不发明新的错误接口；
- **零依赖核心**：核心包仅用标准库；logx 适配在独立子包；
- **不可变错误**：`WithField` 返回新实例，可安全共享；
- **可观测**：错误码、分类、字段、栈四要素齐备，直接对接日志与监控。
