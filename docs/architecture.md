# 架构设计

## 1. 核心类型

```text
error（标准接口）
  ▲
  │ 实现
*Error
  ├── code   Code       稳定错误码（如 USER_NOT_FOUND）
  ├── kind   Kind       分类（invalid/timeout/business...）
  ├── msg    string     人类可读消息
  ├── fields []KV       结构化上下文
  ├── cause  error      底层原因（Unwrap 链）
  └── stack  []uintptr  创建点调用栈（可关闭）
```

## 2. 错误链

`Wrap` 逐层包裹，`Unwrap()` 返回 `cause`，与标准库 `errors.Is / errors.As` 完全兼容：

```text
外层 *Error(DB_DOWN) → 内层 *Error(USER_NOT_FOUND) → errors.New("...")
```

查询函数沿链遍历：

| 函数 | 语义 |
| --- | --- |
| `As(err)` | 取链上第一个 `*Error` |
| `CodeOf(err)` | 取链上第一个错误码 |
| `KindOf(err)` | 取链上第一个分类 |
| `Is(err, code)` | 链上是否存在指定错误码 |
| `Retryable(err)` | 链上是否存在可重试分类 |

## 3. 错误码注册表

- `RegisterCode(code, desc)`：启动期注册（重复注册覆盖）；
- `Describe(code)`：查询说明；
- `Codes()`：排序快照，用于生成错误码文档与审计清单；
- 内置默认码 `CodeUnknown`。

## 4. 调用栈捕获

- 创建时 `runtime.Callers` 捕获（最多 32 帧），`%+v` 输出；
- 全局开关 `SetStackCapture(false)` 可关闭（错误构造频率极高时）；
- 栈帧在 `Format` 时才解析，构造路径无额外格式化开销。

## 5. 不可变与并发安全

- `WithField` 返回新实例，原错误不变，可安全共享/缓存；
- 注册表由 `RWMutex` 保护，注册与查询并发安全；
- `SetStackCapture` 为原子开关。

## 6. logx 适配（errx/logx）

适配子包单向依赖 `github.com/lcylpzls/logx`，将错误转为 `logx.FieldGroup`：
`err.code` / `err.kind` / `err.message` + 错误携带的 KV 字段。logx 核心保持零依赖。
