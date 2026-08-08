# 版本与发布

## 1. 版本策略

- 遵循语义化版本（SemVer）；
- **当前处于 v0.x 阶段**：允许破坏性变更，升级前查看 CHANGELOG；
- 功能与 API 稳定后锁定 **v1.0.0**，此后破坏性变更必须提升主版本（v2.0.0）；
- 升级前查看 CHANGELOG。

## 版本历史

| 版本 | 说明 |
| --- | --- |
| v0.3.0 | 错误码 Markdown 生成器、栈深度可配置、Fuzz 强化、性能文档 |
| v0.2.0 | 健壮性：错误链深度上限、Error() 惰性缓存、并发专项测试、基准 |
| v0.1.0 | 首个可用版本：结构化错误核心 + logx 适配 |

## 2. 发布流程

```bash
# 1) 确认 main 分支 CI 全绿
git push origin main

# 2) 更新 CHANGELOG 定版

# 3) 提交
git add CHANGELOG.md && git commit -m "chore(release): 定版 vX.Y.Z"

# 4) 打 tag 并推送（触发 release.yml：go test -race + 创建 Release）
git tag vX.Y.Z && git push origin vX.Y.Z
```

## 3. 发版检查清单

- [ ] `go test -count=1 ./...` 通过；
- [ ] `go vet ./...`、`staticcheck ./...` 零告警；
- [ ] 两包覆盖率 100%；
- [ ] GitHub CI 三平台矩阵全绿（含 Linux race）；
- [ ] CHANGELOG 已定版；
- [ ] Release 工作流成功；
- [ ] 临时模块验证 `go get github.com/lcylpzls/errx@vX.Y.Z`。
