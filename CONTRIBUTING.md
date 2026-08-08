# 参与贡献

感谢您愿意改进 errx。请遵循以下约定，让评审和发布更顺畅。

## 基础规范

- 所有日志、打印输出、代码注释与文档均使用简体中文；
- 开发机为 Windows 操作系统，命令一律使用 PowerShell，不使用 bash 等 POSIX 命令。

## 开发环境

- Go >= 1.21（推荐 1.26.x 与本地一致）；
- 提交前通过：`gofmt`、`go vet ./...`、`staticcheck ./...`；
- 新增或修改代码必须配套测试，三个包语句覆盖率保持 100%。

## 提交规范

- 使用 Conventional Commits，消息以 `errx` 为 scope，例如：
  - `feat(errx): ...`（新能力）
  - `fix(errx): ...`（缺陷修复）
  - `docs(errx): ...` / `ci(errx): ...` / `chore(errx): ...`
- 行为变更同步更新 CHANGELOG；发布按 `docs/release.md` 流程执行。

## 代码约定

- 公开 API 注释使用简体中文，说明语义与边界行为；
- 破坏性变更必须提升主版本（自 v1.0.0 起 API 冻结，CI 的 apidiff 会拦截）；
- 核心包保持零 import 依赖；适配子包（logx/httpx）依赖独立声明。

## PR 流程

1. 从 main 新建分支，小步提交；
2. 确保 CI 全绿：三平台矩阵、race、apidiff、fuzz；
3. PR 描述说明动机、变更点与验证结果。
