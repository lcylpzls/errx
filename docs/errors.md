# 错误码手册

> 由 errx 注册表自动生成,新增错误码请通过 `RegisterCode` 注册。

## CACHEX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `CACHEX_INVALID_CONFIG` | invalid_argument | 配置非法 |

## CONF

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `CONF_INVALID_CONF` | invalid_argument | CONF 解析、绑定或序列化失败 |
| `CONF_INVALID_INI` | invalid_argument | INI 解析、绑定或序列化失败 |
| `CONF_INVALID_JSON` | invalid_argument | JSON 解析、绑定或序列化失败 |
| `CONF_INVALID_OPTION` | invalid_argument | 选项参数非法 |
| `CONF_INVALID_TARGET` | invalid_argument | 配置目标必须为非空指针 |
| `CONF_INVALID_TOML` | invalid_argument | TOML 解析、绑定或序列化失败 |
| `CONF_INVALID_YAML` | invalid_argument | YAML 解析、绑定或序列化失败 |
| `CONF_READ_FAILED` | unavailable | 配置文件读取失败 |
| `CONF_UNKNOWN_KEY` | invalid_argument | 配置包含未声明字段 |
| `CONF_UNSUPPORTED_FORMAT` | invalid_argument | 不支持的配置格式 |
| `CONF_WRITE_FAILED` | unavailable | 配置文件写入失败 |

## DBX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `DBX_BAD_ARGUMENT` | invalid_argument | 参数非法 |
| `DBX_CLOSE_FAILED` | unavailable | 关闭数据库连接失败 |
| `DBX_DRIVER_NOT_REGISTERED` | invalid_argument | 数据库驱动/方言未注册 |
| `DBX_DUPLICATE` | conflict | 唯一约束或重复键冲突 |
| `DBX_EXEC_FAILED` | unavailable | Exec 执行失败 |
| `DBX_MIGRATION_FAILED` | unavailable | 迁移执行失败 |
| `DBX_NOT_FOUND` | not_found | 查询无结果 |
| `DBX_OPEN_FAILED` | unavailable | 打开数据库连接失败 |
| `DBX_QUERY_FAILED` | unavailable | 查询失败 |
| `DBX_SCAN_FAILED` | invalid_argument | 扫描或类型转换失败 |
| `DBX_TX_BEGIN_FAILED` | unavailable | 开启事务失败 |
| `DBX_TX_CALLBACK_FAILED` | business | 事务回调失败,已回滚 |
| `DBX_TX_COMMIT_FAILED` | unavailable | 提交事务失败 |
| `DBX_TX_ROLLBACK_FAILED` | unavailable | 回滚事务失败 |

## HTX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `HTX_BODY_TOO_LARGE` | invalid_argument | 响应体超过大小上限 |
| `HTX_BODY_UNREADABLE` | invalid_argument | 请求体不可重读 |
| `HTX_DIAL_FAILED` | unavailable | 建立连接失败 |
| `HTX_INVALID_CONFIG` | invalid_argument | 配置或请求参数非法 |
| `HTX_REDIRECT_EXCEEDED` | invalid_argument | 重定向次数超限 |
| `HTX_REDIRECT_FAILED` | invalid_argument | 重定向地址解析或构造失败 |
| `HTX_REQUEST_FAILED` | unavailable | 请求发送失败 |
| `HTX_RESPONSE_FAILED` | invalid_argument | 读取响应失败 |
| `HTX_RETRY_EXHAUSTED` | unavailable | 重试耗尽 |
| `HTX_TLS_FAILED` | unavailable | TLS 握手失败 |
| `HTX_UNEXPECTED_STATUS` | invalid_argument | 响应状态码不在期望列表 |
| `HTX_UNSUPPORTED_PROTOCOL` | invalid_argument | 协议未注册 |

## MTRX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `MTRX_ALREADY_REGISTERED` | invalid_argument | 指标已注册 |
| `MTRX_INVALID_CONFIG` | invalid_argument | 配置非法 |

## RESX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `RESX_BULKHEAD_FULL` | unknown | 舱壁拒绝 |
| `RESX_CIRCUIT_OPEN` | unavailable | 熔断拒绝 |
| `RESX_INVALID_CONFIG` | invalid_argument | 配置非法 |
| `RESX_RATE_LIMITED` | rate_limited | 限流拒绝 |
| `RESX_WAIT_CANCELED` | cancelled | 等待限流许可被取消 |

## UNKN

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `UNKNOWN` | unknown | 未知错误 |

## VALIDX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `VALIDX_INVALID_RULE` | invalid_argument | 校验规则语法或参数非法 |
| `VALIDX_VALIDATION_FAILED` | invalid_argument | 字段校验失败 |

## WEBX

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `WEBX_CONFIG_INVALID` | invalid_argument | webx 配置校验失败 |
| `WEBX_CONFIG_LOAD_FAILED` | unavailable | webx 配置文件加载失败 |
| `WEBX_LISTEN_FAILED` | unavailable | webx 监听器创建失败 |
| `WEBX_PANIC` | internal | webx 请求处理发生 panic |
| `WEBX_SHUTDOWN_FAILED` | unavailable | webx 优雅关闭失败 |
| `WEBX_START_FAILED` | invalid_argument | webx 服务启动失败 |

## authx

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `authx_audit_queue_full` | rate_limited | 审计队列已满 |
| `authx_csrf_generation_failed` | unavailable | CSRF 令牌生成失败 |
| `authx_csrf_mismatch` | forbidden | CSRF 校验不匹配 |
| `authx_forbidden` | forbidden | 无权限 |
| `authx_mfa_config_invalid` | invalid_argument | MFA 配置非法 |
| `authx_mfa_invalid` | invalid_argument | MFA 校验失败 |
| `authx_oauth2_config_invalid` | invalid_argument | OAuth2 配置非法 |
| `authx_oauth2_invalid` | unauthorized | OAuth2 参数非法 |
| `authx_password_config_invalid` | invalid_argument | 哈希参数非法 |
| `authx_password_hash_invalid` | invalid_argument | 密码哈希格式无效或损坏 |
| `authx_password_internal` | unavailable | 哈希/校验过程内部失败 |
| `authx_password_mismatch` | unauthorized | 明文与哈希不匹配 |
| `authx_password_too_long` | invalid_argument | 明文密码超过长度上限 |
| `authx_password_too_short` | invalid_argument | 明文密码低于长度下限 |
| `authx_password_too_weak` | invalid_argument | 密码强度不足 |
| `authx_rbac_cycle` | invalid_argument | 角色继承环 |
| `authx_rbac_invalid` | invalid_argument | RBAC 参数非法 |
| `authx_rbac_limit` | invalid_argument | 角色/权限数量超限 |
| `authx_rbac_role_exists` | conflict | 角色已存在 |
| `authx_rbac_role_not_found` | invalid_argument | 角色不存在 |
| `authx_refresh_token_invalid` | unauthorized | 刷新令牌无效 |
| `authx_security_config_invalid` | invalid_argument | 安全配置非法 |
| `authx_session_invalid` | invalid_argument | 会话无效 |
| `authx_session_not_found` | unauthorized | 会话不存在 |
| `authx_session_store_invalid` | unavailable | 会话存储失败 |
| `authx_store_full` | unavailable | 存储容量已满 |
| `authx_store_invalid` | unavailable | 存储读写失败 |
| `authx_token_config_invalid` | invalid_argument | 令牌配置非法 |
| `authx_token_expired` | unauthorized | 令牌已过期 |
| `authx_token_invalid` | unauthorized | 令牌格式非法或载荷非法 |
| `authx_token_missing` | unauthorized | 缺少令牌 |
| `authx_token_revoked` | unauthorized | 令牌已撤销 |
| `authx_token_signature` | unauthorized | 令牌签名无效 |

## idgenx

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `idgenx_clock_backward` | unavailable | 检测到时钟回拨 |
| `idgenx_collision` | conflict | 短 ID 碰撞重试耗尽 |
| `idgenx_invalid_config` | invalid_argument | 配置非法 |
| `idgenx_invalid_id` | invalid_argument | ID 解析失败 |
| `idgenx_node_invalid` | invalid_argument | 节点 ID 越界 |
| `idgenx_rand_failure` | unavailable | 随机源失败 |
| `idgenx_timestamp_overflow` | unavailable | 时间戳超出位宽范围 |
| `idgenx_wait_timeout` | timeout | 等待时钟追平超时 |

## jobx

| 错误码 | 分类 | 说明 |
| --- | --- | --- |
| `jobx_cron_invalid` | invalid_argument | cron 表达式非法 |
| `jobx_execution_failed` | internal | 处理器执行失败 |
| `jobx_handler_conflict` | already_exists | 处理器重复注册 |
| `jobx_handler_not_found` | not_found | 未注册处理器 |
| `jobx_id_generate_failed` | unavailable | 任务 ID 生成失败 |
| `jobx_invalid_config` | invalid_argument | 配置非法 |
| `jobx_job_cancelled` | cancelled | 任务已取消 |
| `jobx_job_invalid` | invalid_argument | 任务参数非法 |
| `jobx_job_not_found` | not_found | 任务或调度条目不存在 |
| `jobx_queue_full` | rate_limited | 就绪队列已满 |
| `jobx_replaced` | conflict | 同名旧任务已被替换取消 |
| `jobx_retry_exhausted` | internal | 重试耗尽 |
| `jobx_scheduler_stopped` | unavailable | 调度器已停止 |
| `jobx_shutting_down` | unavailable | 关闭中拒绝新任务 |
| `jobx_skipped` | already_exists | 同名任务在途,本次提交被跳过 |
| `jobx_store_invalid` | unavailable | 任务存储读写失败 |
| `jobx_timeout` | timeout | 单任务执行超时 |

