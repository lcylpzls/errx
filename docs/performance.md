# 性能

## 基准命令

```bash
go test -bench=. -benchmem -benchtime=1s ./...
```

## 实测数据（Windows 10 / AMD Ryzen 5 7600 / Go 1.26.5）

| 基准 | 耗时 | 内存 |
|---|---:|---:|
| New（含栈捕获） | 212.6 ns/op | 256 B / 1 alloc |
| Wrap（含栈捕获） | 213.1 ns/op | 256 B / 1 alloc |
| Error()（惰性缓存后重复调用） | 1.27 ns/op | 0 B / 0 alloc |
| Is（单节点链） | 2.71 ns/op | 0 B / 0 alloc |
| Retryable（单节点链） | 1.67 ns/op | 0 B / 0 alloc |

## 性能说明

- 错误**构造**路径的 1 次分配来自调用栈捕获（32 帧）；对构造频率极敏感的场景可 `SetStackCapture(false)` 或 `SetStackDepth(n)` 调低；
- 错误**查询与输出**路径零分配；
- 基准数字与机器强相关，发布声明时请在本机复跑。
