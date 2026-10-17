> 用户已否决按管理员划分与逐请求一次性密钥。下文仅保留为未批准的历史提案，不得实施。当前方案见 [shared-key-expiry.md](shared-key-expiry.md)。

# 内部密钥滥用边界

状态：planning，待审批。用户补充：“还需要思考如何避免内部密钥滥用”。这修正上一版按管理员长期复用内部密钥的提议，不改变“特殊密钥＋选择头，其余原生”的网关边界。

## 实际缺口

仅验证随机密钥与选择头，会让拿到该管理员密钥的人提交任意body、改目标并重复请求。只存在内存或使用127.0.0.1并不能封住这种用途。

原生API-key middleware检查用户/组/余额/额度，但不知道某段body是否由评估任务生成。原生密钥的固定过去expiry被内部入口豁免，也不能替代内部凭据自己的生命周期。

建议撤回“每管理员一把可长期复用的内部密钥”，改成“每次本机HTTP调用一把一次性内部密钥”。还是同一个Authorization字段和同一个选择头，没有额外ticket/header/API/文件/数据库表。

## 一次性密钥如何工作

1. 执行器已取得真实任务身份，并构造完本次body和Codex headers后，使用crypto/rand生成256-bit随机密钥。
2. 现有runtime内存映射只保存待接纳调用：密钥→真实UserID/APIKeyID、任务context、目标group/account、body的SHA-256、允许的Codex headers副本和该请求原有deadline。
3. 密钥不交给浏览器，不接收客户端actor/APIKeyID，不进入history、配置、Key列表或日志。数据库里的记账标识仍只是查找标识。
4. API-key鉴权入口确认POST /v1/responses、无query/path别名及单一目标头。读取Request.RemoteAddr的真实socket peer并要求loopback，不用转发IP头授予内部权限。
5. 查内存记录，确认任务/request context未取消、deadline未过、选择头目标完全匹配；在原生请求体限制内读取body，核对SHA-256与允许Codex headers，然后恢复Body供native读取。
6. 最后在runtime锁下重新核对记录仍有效并删除，保证并发相同密钥最多一个被接纳。比较body与数据库查询都不能持这把锁。
7. 接入绑定的原生用户/记账行与所选分组，保留B06/B09原生约束。只豁免该确切记账行固定过去expiry；不是豁免内部凭据自己的deadline。
8. 移除Authorization和选择头再进native forwarding。无论本地发送成功、失败、取消或超时，执行器defer移除未消费记录；取消context立即让未接纳凭据无效，即使defer尚未运行。

native内部同账户retry不再经过入口认证，不需要再次使用密钥。原评估再次执行HTTP尝试时，执行器再生成一把。一次性指一次HTTP接纳，不承诺一个供应商请求；原生有限retry仍保留。

任务正常结束、取消和进程重启后，不残留可复用凭据。密钥没有独立可配置TTL或刷新流程；期限沿已有本次请求deadline。

## 这防止什么

| 滥用 | 结果 |
| --- | --- |
| 普通用户带选择头 | 拒绝，不能靠目标头进入内部身份。 |
| 拿Key后改group/account/模型/prompt/tools等body | 目标或body/header校验失败。不能把一次评估调用变成任意代理调用。 |
| 重放或并发复制相同请求 | 已消费Key拒绝；并发最多一个接纳。 |
| 任务取消后继续使用未消费Key | task/request context检查拒绝。 |
| 调用/models、/usage、其他别名或任意路径 | 内部凭据不受理这些入口。固定过去expiry不作为这项隔离保证。 |
| 公网直接带Key或伪造X-Forwarded-For | socket peer检查拒绝公网直连；伪造IP头不产生内部权限。 |

没有新增免费用户、无限额度、目标外账号或身份恢复机制。原生余额、组权限、Key/user/account限制继续生效。

## 不能宣称的保证

loopback不是可信调用者证明：本机进程、SSRF和本机反向代理都可能连到原生监听。它只作为来源限制，随机密钥与完整请求绑定才阻止任意用法。

若攻击者同时拿到尚未消费密钥和完整请求，仍可能抢先发出这一次相同的已授权调用，导致真正调用被拒绝。一次性不能证明谁先提交；不能声称完全阻止这种抢占。

若攻击者能读取/修改Sub2API进程内存或已有服务端代码执行权限，本应用内方案不能隔离他。不要为此增加签名服务、维护socket、另一身份框架或新配置。

任务cancel不会撤销已经接纳的native请求，也不会取消供应商已发生费用。原生后续行为不改。

## 修改位置和原生证据

仍只在runtime/executor准备与API-key内部鉴权内完成，硬目标选择保持原计划。没有原生CC/SSE parser、forward/error/usage/Shutdown改动。

v0.2.16 routes/gateway.go:181在API-key auth之前已有RequestBodyLimit；request_body_limit.go:13使用http.MaxBytesReader，可在此限制内核对body并恢复读取。

api_key_auth.go:120—175已有Key状态、IP、用户及组检查；:179有simple提前返回，usage/billing也会跳过部分计费条件。特殊密钥的路径/来源/一次性校验必须在这些早退之前完成，不能指望expired记账行替代它。

ip/ip.go:164的GetSecurityClientIP可能受配置与forwarded headers影响，不能用它给内部请求授予信任。内部来源直接解析Request.RemoteAddr。

internal/eval/model.go:100每次评估HTTP重试分别调用executor，适合每次生成独立凭据；原生同账户retry留在同一次已认证HTTP处理内。

## 最小验证

在用户批准后，使用本地fixture核对：改目标/body/允许headers拒绝、重复/并发消费最多一次、取消后拒绝、其他路径及转发IP伪造拒绝、正常调用和native同账户retry可用。

这是用户明确提出的密钥滥用验证，不追加全站安全审计、race、付费线上检查或永久凭据管理框架。当前仅源码核对与规划，不是已经实现或测试通过。
