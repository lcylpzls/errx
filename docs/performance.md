# 性能

## 基准命令

```powershell
go test -count=1 -run '^$' -bench 'BenchmarkNew|BenchmarkWrap|BenchmarkErrorString|BenchmarkIs|BenchmarkRetryable' -benchmem -benchtime=1s ./internal/core
```

## 实测数据（Windows 10 / AMD Ryzen 5 7600 / Go 1.26.5）

| 基准 | 耗时 | 内存 |
|---|---:|---:|
| New（含栈捕获） | 251.4 ns/op | 400 B / 2 allocs |
| Wrap（含栈捕获） | 257.1 ns/op | 400 B / 2 allocs |
| Error()（惰性缓存后重复调用） | 1.244 ns/op | 0 B / 0 allocs |
| Is（单节点链） | 20.23 ns/op | 16 B / 1 alloc |
| Retryable（单节点链） | 7.155 ns/op | 0 B / 0 allocs |

## 性能说明

- 错误**构造**路径的分配来自调用栈捕获（32 帧）与错误实例本身；
  对构造频率极敏感的场景可 `SetStackCapture(false)` 或 `SetStackDepth(n)` 调低；
- 错误**输出**与 `Retryable` 路径零分配；`Is` 在聚合展开路径有一次分配；
- 基准数字与机器强相关，发布声明时请在本机复跑。
