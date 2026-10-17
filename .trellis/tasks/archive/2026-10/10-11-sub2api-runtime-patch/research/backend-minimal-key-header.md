> 历史已撤回方案：内部Key、选择头、用户记账及正常网关接入不再属于当前设计。当前采用独立管理员账户codex-eval接口。

> 历史范围收缩记录。下文逐管理员归属已经被用户撤回，不代表当前授权；当前密钥方案以 shared-key-expiry.md 的共享Key、每次调用重置5分钟为准。

> 完整当前计划见 [minimal-replan.md](minimal-replan.md)。下文是范围收缩时的记录；密钥绑定及HTTP调用方案已在新计划明确，待审批。

# 后端范围收缩：内部密钥与选择头

状态：planning。用户最新指令：“明明就两件事：一个特殊选择头，一个特殊内部密钥，检查密钥+选择头指定转发……剩下的事情sub2api本身就在干”。该指令取代未批准的B11—B15扩展方案；不是实施批准。

## 新边界

1. 内部密钥：评估执行器发送内部密钥，鉴权处验证后接入任务启动者的原生用户/Key记账身份。选择头只有通过此验证才能生效；普通Key不能使用它。沿用已批准的实际启动者计费归属、原生记账数据与请求分组，不改为共享一个管理员账单。具体密钥表示/绑定方式须在精简方案收敛时明确；不得把此前公开可推导的数据库保留标识称为秘密密钥。
2. 选择头：携带精确group_id/account_id，鉴权成功后让原生账户选择只针对指定ID。指定ID失败或被native failedAccountIDs排除，返回原生失败，不进入其他账户。不能仅把ID当sticky偏好，否则native仍可换号。

只增加身份识别与硬目标选择分支。原生资格、并发、利润检查、同账户重试、协议转换、HTTP/SSE错误处理、用量/计费、连接断开处理与退出流程沿原生代码。

## 撤回/删除此前额外方案

撤回请求局部诊断收集器、完整错误体旁路、CC parser修改、额外发送前DB资格重查、定制sticky写入保护、partial-usage builder/新提交钩子、三个上游取消context替换、producer完成屏障与定制五秒退出协调。删除现存相应eval附加修改，恢复上游原函数；既有无限退出续期也回退到原生main，而不是再添退出框架。

B05仍恢复评估层收到后被自己吞掉的错误；不再要求越过原生网关恢复它已改写的上游诊断。任务/cancel停止评估后续执行，native是否继续读上游结算遵从原生断连语义，不承诺立即停止供应商生成。

## 原生证据

v0.2.16 internal/service/openai_account_scheduler.go:72有StickyAccountID、PreserveStickyBinding、DisableStickyEscape、ExcludedIDs；selectBySessionHash:495已有资格、DB重查、并发抢槽和WaitPlan。它返回失败后通常会继续候选路径，因此不能把sticky=目标直接等同于硬定向。精简分支应复用原生已有账户检查/slot能力并终止候选回退，不复制整套资格链。

internal/handler/openai_gateway_handler.go:676每次选择传入failedAccountIDs；同账户retry不排除目标，跨账户失败把目标排除。定向分支只允许目标，因此可以保留原handler的错误/重试流程，下一次选号按native错误结束，无须在每个continue边界额外插入return。

B08先前的context-capability鉴权形式被本次用户指定的密钥+选择头方向取代；B06/B09实际启动者记账归属保留。其余已批准前端与核心原始行为要求不变。

本轮只检查源码、收缩规划。没有产品改动、编译/测试、真实请求或部署。
