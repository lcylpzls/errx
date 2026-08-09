# 错误码手册

> 由 errx 注册表自动生成,新增错误码请通过 `RegisterCode` 注册。

## CACHEX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `CACHEX_INVALID_CONFIG` | unknown | 配置非法 |

## CONF

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `CONF_INVALID_CONF` | unknown | CONF 解析、绑定或序列化失败 |
| `CONF_INVALID_INI` | unknown | INI 解析、绑定或序列化失败 |
| `CONF_INVALID_JSON` | unknown | JSON 解析、绑定或序列化失败 |
| `CONF_INVALID_OPTION` | unknown | 选项参数非法 |
| `CONF_INVALID_TARGET` | unknown | 配置目标必须为非空指针 |
| `CONF_INVALID_TOML` | unknown | TOML 解析、绑定或序列化失败 |
| `CONF_INVALID_YAML` | unknown | YAML 解析、绑定或序列化失败 |
| `CONF_READ_FAILED` | unknown | 配置文件读取失败 |
| `CONF_UNKNOWN_KEY` | unknown | 配置包含未声明字段 |
| `CONF_UNSUPPORTED_FORMAT` | unknown | 不支持的配置格式 |
| `CONF_WRITE_FAILED` | unknown | 配置文件写入失败 |

## DBX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `DBX_BAD_ARGUMENT` | unknown | 参数非法 |
| `DBX_CLOSE_FAILED` | unknown | 关闭数据库连接失败 |
| `DBX_DRIVER_NOT_REGISTERED` | unknown | 数据库驱动/方言未注册 |
| `DBX_DUPLICATE` | unknown | 唯一约束或重复键冲突 |
| `DBX_EXEC_FAILED` | unknown | Exec 执行失败 |
| `DBX_MIGRATION_FAILED` | unknown | 迁移执行失败 |
| `DBX_NOT_FOUND` | unknown | 查询无结果 |
| `DBX_OPEN_FAILED` | unknown | 打开数据库连接失败 |
| `DBX_QUERY_FAILED` | unknown | 查询失败 |
| `DBX_SCAN_FAILED` | unknown | 扫描或类型转换失败 |
| `DBX_TX_BEGIN_FAILED` | unknown | 开启事务失败 |
| `DBX_TX_CALLBACK_FAILED` | unknown | 事务回调失败,已回滚 |
| `DBX_TX_COMMIT_FAILED` | unknown | 提交事务失败 |
| `DBX_TX_ROLLBACK_FAILED` | unknown | 回滚事务失败 |

## HTX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `HTX_BODY_TOO_LARGE` | unknown | 响应体超过大小上限 |
| `HTX_BODY_UNREADABLE` | unknown | 请求体不可重读 |
| `HTX_DIAL_FAILED` | unknown | 建立连接失败 |
| `HTX_INVALID_CONFIG` | unknown | 配置或请求参数非法 |
| `HTX_REDIRECT_EXCEEDED` | unknown | 重定向次数超限 |
| `HTX_REDIRECT_FAILED` | unknown | 重定向地址解析或构造失败 |
| `HTX_REQUEST_FAILED` | unknown | 请求发送失败 |
| `HTX_RESPONSE_FAILED` | unknown | 读取响应失败 |
| `HTX_RETRY_EXHAUSTED` | unknown | 重试耗尽 |
| `HTX_TLS_FAILED` | unknown | TLS 握手失败 |
| `HTX_UNEXPECTED_STATUS` | unknown | 响应状态码不在期望列表 |
| `HTX_UNSUPPORTED_PROTOCOL` | unknown | 协议未注册 |

## MTRX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `MTRX_ALREADY_REGISTERED` | unknown | 指标已注册 |
| `MTRX_INVALID_CONFIG` | unknown | 配置非法 |

## RESX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `RESX_BULKHEAD_FULL` | unknown | 舱壁拒绝 |
| `RESX_CIRCUIT_OPEN` | unknown | 熔断拒绝 |
| `RESX_INVALID_CONFIG` | unknown | 配置非法 |
| `RESX_RATE_LIMITED` | unknown | 限流拒绝 |
| `RESX_WAIT_CANCELED` | unknown | 等待限流许可被取消 |

## UNKN

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `UNKNOWN` | unknown | 未知错误 |

## VALIDX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `VALIDX_INVALID_RULE` | unknown | 校验规则语法或参数非法 |
| `VALIDX_VALIDATION_FAILED` | unknown | 字段校验失败 |

## WEBX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `WEBX_CONFIG_INVALID` | unknown | webx 配置校验失败 |
| `WEBX_CONFIG_LOAD_FAILED` | unknown | webx 配置文件加载失败 |
| `WEBX_LISTEN_FAILED` | unknown | webx 监听器创建失败 |
| `WEBX_PANIC` | unknown | webx 请求处理发生 panic |
| `WEBX_SHUTDOWN_FAILED` | unknown | webx 优雅关闭失败 |
| `WEBX_START_FAILED` | unknown | webx 服务启动失败 |

## authx

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `authx_audit_queue_full` | unknown | 审计队列已满 |
| `authx_csrf_generation_failed` | unknown | CSRF 令牌生成失败 |
| `authx_csrf_mismatch` | unknown | CSRF 校验不匹配 |
| `authx_forbidden` | unknown | 无权限 |
| `authx_mfa_config_invalid` | unknown | MFA 配置非法 |
| `authx_mfa_invalid` | unknown | MFA 校验失败 |
| `authx_oauth2_config_invalid` | unknown | OAuth2 配置非法 |
| `authx_oauth2_invalid` | unknown | OAuth2 参数非法 |
| `authx_password_config_invalid` | unknown | 哈希参数非法 |
| `authx_password_hash_invalid` | unknown | 密码哈希格式无效或损坏 |
| `authx_password_internal` | unknown | 哈希/校验过程内部失败 |
| `authx_password_mismatch` | unknown | 明文与哈希不匹配 |
| `authx_password_too_long` | unknown | 明文密码超过长度上限 |
| `authx_password_too_short` | unknown | 明文密码低于长度下限 |
| `authx_password_too_weak` | unknown | 密码强度不足 |
| `authx_rbac_cycle` | unknown | 角色继承环 |
| `authx_rbac_invalid` | unknown | RBAC 参数非法 |
| `authx_rbac_limit` | unknown | 角色/权限数量超限 |
| `authx_rbac_role_exists` | unknown | 角色已存在 |
| `authx_rbac_role_not_found` | unknown | 角色不存在 |
| `authx_refresh_token_invalid` | unknown | 刷新令牌无效 |
| `authx_security_config_invalid` | unknown | 安全配置非法 |
| `authx_session_invalid` | unknown | 会话无效 |
| `authx_session_not_found` | unknown | 会话不存在 |
| `authx_session_store_invalid` | unknown | 会话存储失败 |
| `authx_store_full` | unknown | 存储容量已满 |
| `authx_store_invalid` | unknown | 存储读写失败 |
| `authx_token_config_invalid` | unknown | 令牌配置非法 |
| `authx_token_expired` | unknown | 令牌已过期 |
| `authx_token_invalid` | unknown | 令牌格式非法或载荷非法 |
| `authx_token_missing` | unknown | 缺少令牌 |
| `authx_token_revoked` | unknown | 令牌已撤销 |
| `authx_token_signature` | unknown | 令牌签名无效 |

## jobx

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `jobx_cron_invalid` | unknown | cron 表达式非法 |
| `jobx_execution_failed` | unknown | 处理器执行失败 |
| `jobx_handler_conflict` | unknown | 处理器重复注册 |
| `jobx_handler_not_found` | unknown | 未注册处理器 |
| `jobx_invalid_config` | unknown | 配置非法 |
| `jobx_job_invalid` | unknown | 任务参数非法 |
| `jobx_job_not_found` | unknown | 任务或调度条目不存在 |
| `jobx_queue_full` | unknown | 就绪队列已满 |
| `jobx_replaced` | unknown | 同名旧任务已被替换取消 |
| `jobx_retry_exhausted` | unknown | 重试耗尽 |
| `jobx_scheduler_stopped` | unknown | 调度器已停止 |
| `jobx_shutting_down` | unknown | 关闭中拒绝新任务 |
| `jobx_skipped` | unknown | 同名任务在途,本次提交被跳过 |
| `jobx_store_invalid` | unknown | 任务存储读写失败 |
| `jobx_timeout` | unknown | 单任务执行超时 |

