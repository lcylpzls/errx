# 版本与发布

## 1. 版本策略

- 遵循语义化版本（SemVer）；
- 家族约定：v1 之后破坏性变更统一走 minor 版本（不强制主版本升级），
  直至另行调整版本规范；
- 升级前查看 CHANGELOG。

API 文档以 `go doc -all` 导出为准，不再维护逐版本 API 快照文件。

## 版本历史

| 版本 | 说明 |
| --- | --- |
| v1.6.1 | 文档同步与历史清理（纯文档/版本元数据变更） |
| v1.6.0 | 主体下沉 internal/core、根包薄转发；删除 errx/logx 与 errx/httpx 子包；指标改为 MetricsHook 外置；零家族依赖 |
| v1.5.7 | 依赖升级：logx v1.3.4、testx v1.4.4 |
| v1.5.6 | 家族正式基线锁定 |
| v1.5.5 | 家族依赖最终对齐 v1 正式版基线 |
| v1.5.4 | 家族依赖对齐到最新基线 |
| v1.5.3 | 家族依赖对齐 |
| v1.5.2 | IP 限流统一引用 resiliencex 令牌桶（webx 联动） |
| v1.5.1 | 路由/表单/平台错误统一 errx 化（webx 联动） |
| v1.5.0 | 校验统一迁移至 validx（webx 联动） |
| v1.4.0 | 家族测试底座接入 testx |
| v1.3.x | 工程规范与文档完善 |
| v1.2.0 | StackTrace/StackFrame、Aggregate JSON、godoc 示例、fuzz |
| v1.1.x | net/http 适配与修复 |
| v1.0.0 | API 冻结：完整错误体系（17 Kind/分组/策略/聚合/JSON/HTTP/指标） |

## 2. 发布流程

```powershell
# 1) 确认 main 分支 CI 全绿
git push origin main

# 2) 更新 CHANGELOG 定版

# 3) 提交
git add CHANGELOG.md
git commit -m "chore(release): 定版 vX.Y.Z"

# 4) 打 tag 并推送（触发 release.yml：多发行版测试 + go test -race + 创建 Release）
git tag vX.Y.Z
git push origin vX.Y.Z
```

## 3. 发版检查清单

- [ ] `go test -count=1 ./...` 通过；
- [ ] `go vet ./...`、`staticcheck ./...` 零告警；
- [ ] 根包与 internal/core 语句覆盖率 100%；
- [ ] GitHub CI 三平台矩阵 + 多 Linux 发行版全绿（含 race）；
- [ ] CHANGELOG 已定版；
- [ ] Release 工作流成功；
- [ ] 临时模块验证 `go get github.com/lcylpzls/errx@vX.Y.Z`。
