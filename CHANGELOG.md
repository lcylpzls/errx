# Changelog

本项目遵循语义化版本（SemVer）。值得记录的变更统一维护在此文件。

## [v1.5.6] - 2026-08-10

### 变更

- 家族正式基线锁定：依赖统一指向 v1 基线已发布版本（errx v1.5.5 / logx v1.3.2 / testx v1.4.3 / validx v1.2.4 / cryptox v1.0.2 / confx v1.0.2 / webx v1.5.4 等），此后家族依赖不再前进。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.5.5] - 2026-08-10

### 变更

- 家族依赖最终对齐到 v1 正式版基线（errx v1.5.4 / logx v1.3.1 / testx v1.4.2 / validx v1.2.3 / confx v1.0.1 / cryptox v1.0.1 等），无 API 变更。

### 质量

- 全部库包语句覆盖率保持 100%；race / vet / staticcheck / fuzz / govulncheck 全绿。

## [v1.5.4] - 2026-08-10

### 变更

- go 指令与 CI/Release 工作流统一为 Go 1.26.5；
- README Go 版本徽章同步更新。

## [v1.5.3] - 2026-08-10

### 变更

- 家族统一 Go 1.21：全部 go.mod 与 CI/Release 工作流版本号对齐 1.21。

## [v1.5.2] - 2026-08-10

### 修复

- 依赖升级 `testx v1.2.1` 后 go 指令正式恢复 `go 1.21`
  （此前被 v1.2.0 的 go 1.26.5 门槛顶回），CI Go 1.21 矩阵可用。

## [v1.5.1] - 2026-08-10

### 修复

- go.mod 的 go 指令恢复 `go 1.21`（testx v1.2.1 已同步降级），
  修复 CI Go 1.21 矩阵失败；
- TestVersion 期望值同步为 v1.5.0。

## [v1.5.0] - 2026-08-10

### 变更

- 家族测试底座接入：httpx 与 logx 子包测试改用 testx 断言；
  根包测试因 testx 依赖 errx 存在循环依赖，保持原生 testing 写法；
- 测试依赖新增 `testx v1.2.0`（仅子包测试引用）。

### 质量

- 根包与子包语句覆盖率 100%；race / vet / staticcheck 全绿。

## [v1.4.0] - 2026-08-10

### 新增

- `SetMetricsHook(MetricsHook)` / `ResetMetricsHook()`：可选全局
  指标钩子，构造与查询事件转发给 metricsx 等外部底座
  （`errx.constructed` 带 kind 标签 / `errx.queried`）；
- 默认关闭：热路径仅多一次原子加载，零分配、零额外开销；
- 配套 `metricsx/adapters/errx` 适配器（家族统一接入）。

### 质量

- 新增钩子转发/重置/并发切换测试；覆盖率保持 100%，
  race / vet / staticcheck 全绿。

## [v1.3.0] - 2026-08-09

### 新增

- `RegisterCode` 冲突检测：同码同描述幂等，同码不同描述 panic，
  防止不同模块静默覆盖同一错误码；
- `RegisterCodeKind` / `CodeKind`：错误码注册分类元数据，
  `CodeInfo` 携带 Kind；
- `NewCode` / `NewCodef`：基于注册错误码便捷构造，
  自动使用注册分类，消灭 `New(Kind, Code, ...)` 样板；
- `CodesMarkdown()`：生成全库错误码手册（按前缀分组，
  含分类与说明），配套 `docs/errors.md` 自动生成文档；
- `Version` 常量。

### 生态

- jobx / authx 补齐错误码注册（15 + 33 个），
  全库 102 个错误码全部进入全局注册表，0 冲突。

### 兼容性

- v1.x 向后兼容：`RegisterCode` 签名不变，冲突场景由静默覆盖
  改为 panic（全库审计无冲突，升级安全）。

## [v1.2.0] - 2026-08-08

### 新增

- `StackTrace()` 与 `StackFrame`：程序化读取创建点调用栈，不再局限于 `%+v` 格式化输出；
- `Aggregate` 支持 JSON 序列化/还原，多错误可跨服务传输；
- godoc 示例：`ExampleNew` / `ExampleWrap` / `ExampleJoin` / `ExampleWriteJSON`。

### 修复

- `WithField` 包装普通错误时不再重复输出原错误文本（如 `UNKNOWN: 普通错误: 普通错误` → `UNKNOWN: 普通错误`）；
- 空错误码统一归一为 `CodeUnknown`，避免 `": 消息"` 式前导冒号。

### 工程

- CI 新增 apidiff API 兼容检查（对比 v1.0.0 冻结基线）与 Fuzz 短跑；
- GitHub Actions 升级至 checkout@v7 / setup-go@v7 / upload-artifact@v7 / action-gh-release@v3，消除 Node 20 弃用警告；
- CI 的 Staticcheck 步骤显式使用 `GOTOOLCHAIN=auto`，兼容 Go 1.21 矩阵（staticcheck 2026.1 要求 Go >= 1.25）；
- 新增 `.gitattributes` 统一行尾处理，README 补充 CI 徽章并修正零依赖措辞；
- 新增 CONTRIBUTING 基础规范（简体中文 / PowerShell）与 issue 模板（bug / feature）；
- `.gitignore` 增加 AI 工具本地文件（`.agents/`、`.codex/`、`AGENTS.md`），避免误推送。

## [v1.1.1] - 2026-08-08

### 修复

- `httpx.WriteJSON` nil 安全：err 为 nil 时输出 500 + 未知分类，不再 panic；
- `httpx.WriteJSON` 普通错误的 `code` 回退为 `UNKNOWN`，与 `kind` 回退语义一致；
- 补充 `httpx` 对应测试用例，三个包语句覆盖率保持 100%。

## [v1.1.0] - 2026-08-08

### 新增

- **net/http 适配子包 `errx/httpx`**：`Status(err)` 映射 HTTP 状态码，`WriteJSON(w, err)` 输出统一 JSON 错误响应体（code/kind/message），零第三方依赖；
- staticcheck 升级至 2026.1（支持 Go 1.26 的 net/http 依赖分析）。

## [v1.0.0] - 2026-08-08

### API 冻结

- 自本版本起冻结 API：破坏性变更必须提升主版本（v2.0.0）；
- API 基线见 `docs/api-v1.0.0.md`，后续可用 apidiff 检测破坏；
- 自 v0.1.0 以来累计能力：Error 结构化错误（错误码/分类/字段/栈）、Kind 17 类 + Category 分组 + Policy 策略、Join 多错误聚合、JSON 跨服务传输、HTTP 状态映射、观测指标、logx 适配。

## [v0.6.0] - 2026-08-08

### 新增

- **观测指标**：`Snapshot()` 返回构造数 / 查询数 / 按 Kind 构造分布（原子计数，并发安全）；
- `ResetMetrics()` 清零指标，便于压测与统计窗口；
- 构造（New/Newf/Wrap/Wrapf）与查询（Is/Retryable/CodeOf/KindOf）自动打点。

## [v0.5.0] - 2026-08-08

### 新增

- **跨服务传输**：`Error` 实现 `json.Marshaler/Unmarshaler`，错误码/分类/消息/字段/原因链可跨 RPC 传递并完整还原；
- **HTTP 状态映射**：`KindHTTPStatus(kind)` / `Error.HTTPStatus()` 覆盖 400/401/403/404/409/422/429/499/501/503/504；
- KV 增加 JSON 字段名（`key`/`value`）。

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
