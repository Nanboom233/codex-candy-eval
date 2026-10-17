> 后续用户已选择管理员测试语义，并要求独立/accounts/:id/codex-eval接口，而非修改原/test。本报告保留源码证据；候选方案部分已由admin-account-codex-eval-approval.md取代。

# 使用原生账户连通性测试发送评估请求

状态：planning / 调研方案待选择。用户要求暂停C01—C05，考虑直接修改Sub2API测试连通性接口承担请求发送。本报告不把这一探索视为撤销旧审批或实施授权。

## 结论

可行的方向是扩展已有账户测试的请求发送能力，接收完整评估请求并保留完整响应，而不是将评估题目直接塞入现有prompt字段。现有测试服务已经直接获取指定账户并处理出站鉴权、代理、TLS和原生插件发送路径。若评估采用管理员测试语义，特殊内部Key、5分钟续期、选择头和账户池定向补丁都可以退出方案；是否接受相应计费/调度语义变化仍需用户决定。

## 原生事实（固定v0.2.16）

源码提交01842576b3d5c1e2545b5a99ad555d36e28ed306，复用已验证的纯净源码。本轮仅静态阅读，无实际供应商请求、构建或测试。

| 原生位置 | 事实 | 对方案的含义 |
| --- | --- | --- |
| internal/server/routes/admin.go:29—36、:393 | /api/v1/admin/accounts/:id/test属于adminAuth保护的管理路由。 | 对外权限复用管理员认证，无须新增内部调用密钥。后台评估服务调用共享发送函数时，其任务仍须先经已批准的管理员Start鉴权。 |
| internal/handler/admin/account_handler.go:1265—1273 | 请求仅含model_id、prompt、mode和媒体选项。 | 当前接口不能表达原Codex完整input、reasoning、tools、metadata和协议头。 |
| internal/service/account_test_service.go:361—399 | 按accountID读取账户，直接分流到平台测试方法。 | 账户是硬目标，不经过候选池，不需要选择头。原测试允许诊断异常账户，不能宣称继承正常网关的所有资格限制。 |
| 同文件:800—889、:2770—2796 | Responses测试构建固定hi、默认instructions；这一分支没有把prompt放入payload。 | 只修改prompt无效，而且会改变原四类评估的请求上下文。需要完整评估请求body输入。 |
| 同文件:827—984 | 已有凭证影子解析、OAuth/API-key鉴权、model mapping、账号头覆盖、代理及TLS/插件发送。 | 应复用这些原生能力，避免把上游凭据带进评估核心或另写发送客户端。 |
| internal/service/openai_plugin_transport.go:23—48 | doOpenAIAccountTestUpstream先尝试原生OAuth插件发送，否则复用HTTPUpstream DoWithTLS/Do，返回标准http.Response。 | 已有可复用的发送入口；不是新的插件开发或再研究一套插件API。 |
| account_test_service.go:2947—3015 | 流转换为content/test_complete/error；未保留Responses usage及完整terminal数据。 | 当前简化事件不足以保留原评估回答、用量和终止语义；评估分支应取得标准响应流。 |
| internal/handler/admin/account_handler.go:1313—1317；internal/service/ratelimit_service.go:2202 | 普通测试成功会调用RecoverAccountAfterSuccessfulTest，恢复error/rate-limit/overload/temp-unsched/model-rate-limit。 | 直接复用整个HTTP处理器会继承状态恢复副作用。评估分支是否执行恢复需在选定方向后明确，不能无声决定。 |
| account_test_service.go:966—998 | 测试路径更新部分账户快照，429会协调限流状态，401可标记账户错误。 | 管理员测试仍可能改变账户运行状态；不能描述为完全只读或完全没有副作用。 |
| internal/repository/http_upstream.go:194—240、:556—633 | accountConcurrency参与连接池大小；inFlight防止活跃客户端淘汰，enforceLimit限制缓存客户端数量。 | 这是原生出站传输管理，不等于正常网关的Redis业务并发槽、WaitPlan或配额检查。 |

所读OpenAI测试调用链不进入网关API-key认证、账户池scheduler或用户RecordUsage/扣费链路。供应商请求本身仍消耗供应商账户额度，获取/展示token用量也不等于产生Sub2API用户账单。

## 最小候选改法

1. 在已有账户测试服务中抽出共用的指定账户发送函数，复用原生凭证、URL/model mapping、代理、TLS和插件发送代码。函数接收完整请求body和必要Codex协议头，返回标准响应；上游Authorization、Host、URL始终由服务端账户配置决定。
2. 原生/test接口增加明确的评估请求分支，接收完整body，而非仅prompt；该分支保留标准响应流、用量和错误。原普通测试请求继续使用其现有hi/简化测试事件行为，原生Vue测试弹窗无需修改。具体字段及分支仍待下一轮审批，不先添加新路由或兼容协议。
3. 四类后台评估任务调用同一个发送函数，无须绕一圈本机HTTP、复制管理员JWT或伪造Gin路由。任务的启动/状态/取消/历史及原算法保持现有已审批边界。实际错误/响应解析仍须后续逐项审阅。
4. 如果选择管理员测试语义，则撤出特殊内部Key/记账行/归属准备/续期、选择头、API-key鉴权修改和池内硬选号修改；分组仍用于已批准的页面筛选，不暗中用于用户定价或计费。

这不是“一处改prompt就完成”：完整payload和完整response是保留评估原语义的必要改动。无需复制完整网关调度和计费框架来装成另一套普通网关。

## 尚需决定的行为

首要选择：是否将评估定义为管理员账户测试，接受它不走Sub2API用户扣费、分组计费策略与正常账户池调度。建议选该方向以缩小补丁，但不能保留原计划的“正常网关全套计费/并发保证”宣传。

选择后再给五项一批的具体修改审批：完整请求契约、原生发送复用、响应/用量契约、账户状态副作用、协议与并发范围。Responses/OAuth有清楚的复用证据；仅支持Chat Completions的账户不能只发最后一段prompt，必须核对完整请求转换，复用现成原生转换能力，避免另写转换规则。已发现internal/pkg/apicompat/chatcompletions_responses_bridge.go:39—45提供Responses到Chat请求转换入口，但本轮没有证明接入后所有模型/工具/流语义等价。

## 当前审批状态

C01—C05暂停，未批准，尤其C05原生本机HTTP提案不能继续当作实施输入。五分钟续期是此前方向已批准决定；在用户选择新的管理员测试方向之前只暂停评估其必要性，不假称已删除。产品代码、安装包及外部设置均未更改。
