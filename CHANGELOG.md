# Changelog

本项目遵循语义化版本（SemVer）。值得记录的变更统一维护在此文件。

## [v0.4.0] - 2026-08-08

### 新增

- **Kind 扩展至 17 类**：新增 `AlreadyExists` / `Forbidden` / `DeadlineExceeded` / `QuotaExceeded` / `NotImplemented` / `DataLoss`，对齐 Google API / gRPC 错误模型；
- **Category 领域分组**：`Kind.Category()` 将错误分为输入/认证授权/资源状态/依赖外部/系统内部/业务规则六类；
- **Policy 策略抽象**：`Kind.Policy()` 输出可重试 / 告警 / 用户可见三维策略，`Retryable()` 委托 Policy；
- **KindsMarkdown()**：按领域分组生成带策略标注的错误分类表；
- **多错误聚合**：`Join(errs...)` 返回 `*Aggregate`，`errors.Is/As` 可命中任一子错误；
- **Is / Retryable 基于 errors.Is**：支持聚合多错误展开，与标准库语义一致；
- **logx 适配增强**：新增 `err.retryable`、`err.code_desc` 字段。

## [v0.3.0] - 2026-08-08

### 新增

- 错误码 Markdown 文档生成器：`Markdown()` 输出排序错误码表，交付前端/网关/审计；
- 栈捕获深度可配置：`SetStackDepth(n)`（默认 32，`<=0` 恢复默认）；
- Fuzz 强化：`FuzzIs` / `FuzzRetryable`（随机链深度，覆盖深度上限边界）；
- 性能文档与基准（构造/查询/输出路径）。

## [v0.2.0] - 2026-08-08

### 健壮性

- 错误链遍历（`Is` / `Retryable`）增加 100 层深度上限，防御意外成环导致的死循环；
- `Error()` 惰性缓存（`sync.Once`），重复打印零额外开销（实测 1.27 ns/op）；
- `WithField` 显式重建新实例，规避 `sync.Once` 复制风险；
- 新增并发专项测试：注册表并发读写、栈捕获开关并发切换、并发错误构造/查询；
- 新增基准：New / Wrap / ErrorString / Is / Retryable。

## [v0.1.0] - 2026-08-08

### 新增

- `Error` 结构化错误：错误码 + 分类 + 消息 + 结构化字段 + 可选调用栈；
- 错误码注册表：`RegisterCode` / `Describe` / `Codes`（排序快照，支持文档生成）；
- 错误分类：`Kind` 枚举与 `Retryable()`；
- 包装链：`New` / `Newf` / `Wrap` / `Wrapf` / `WithField`，`Wrap(nil)` 返回 nil；
- 标准库兼容：`Unwrap` / `Cause`，`errors.Is` / `errors.As` 全链路；
- 查询辅助：`As` / `CodeOf` / `KindOf` / `Is` / `Retryable`；
- 调用栈：创建时捕获，`%+v` 输出，`SetStackCapture` 全局开关；
- `errx/logx` 适配子包：结构化错误转 logx 字段组；
- 质量：两包语句覆盖率 100%、Fuzz、三平台 CI（含 race）。

### 工程

- MIT LICENSE、CHANGELOG、SECURITY、CODEOWNERS、PR/Issue 模板；
- 工业化文档集：架构 / 使用 / 发布。
