# 安全说明

## 报告漏洞

- 请勿在公开 issue 中披露漏洞细节；
- 通过邮件 [lcylpzls@qq.com](mailto:lcylpzls@qq.com) 联系维护者，并在标题注明 `[Security]`；
- 修复发布前我们会与报告者保持沟通。

## 安全使用建议

- 错误消息中不要携带密码、Token 等敏感信息；敏感上下文用 `WithField` 携带并在日志层脱敏；
- 生产环境如错误构造频率极高，可调用 `errx.SetStackCapture(false)` 关闭栈捕获以降低开销；
- 错误码注册表用于文档与审计，不要在运行时依赖其执行顺序。
