> 已被用户最新“内部密钥＋选择头，其余沿原生”指令取代，非当前方案、未获批准。当前范围见 [backend-minimal-key-header.md](backend-minimal-key-header.md)。

# B11—B15 补充调研与具体修改方案

状态：planning，五项均未批准。最新请求是“B11-15均需要补充调研 需要完整说明怎么改”。本文件替代旧 backend-dispatch-approval.md 中不准确的范围和保证；不授权实施。

调研对象：未修改 Sub2API v0.2.16，提交 01842576b3d5c1e2545b5a99ad555d36e28ed306。下列原生文件名均相对 backend/。本轮只读源码并更新规划文件，没有编译、真实上游请求或产品代码修改。

## 先纠正此前的判断

1. OpenAI 平台账户不会进入 Responses→Anthropic 转换。Forward:184 的原生注释也明确 OpenAI API Key 只落到 Responses/Chat Completions；account.go:1460、opencode_go.go:369、provider_profile.go:227 与 upstream_protocol_routing.go:32 的实际条件支持这一结论。不能因为 switch 有 Anthropic 分支就称它可达。
2. 文件名含 WS 不等于 HTTP 无用。bindHTTPResponseAccount:1473 使用原生 WS state store 保存 response_id 的账户/用户/Key 归属。这些原生方法保留；删除的是评估额外加入且不可达的 WS 转发修改。
3. 原生 wire.go 清理的十秒 context 只在 wg.Wait 和顺序关闭之后检查，不会强制终止这些函数。原生 usage worker StopAndWait 也可能超过五秒。因此 B15 只约束新增评估等待窗口，不能保证进程总退出时间。
4. 原生用量提交允许配置 drop/sample/sync。提交到 native worker 不等于已落库，更不等于全部费用一定记录。B14 不另造强制结算服务。
5. 原生 slot helper 已有利润终检；它没有完整重查账户状态/组成员/模型。保留发送前完整重查，但不要另做一套利润规则。
6. 当前 io.Pipe Body.Close 会取消 native handler。评估 parser 在 terminal/bare-error 立即 return，可能使 handler 未处理最终 usage 就被取消；必须同时修正读取收尾和 handler 完成等待。

## B11：只在实际主调度入口选取指定账户

### 调用链与最小插入点

所选分组/账户必须是 B03 已批准的 OpenAI 平台；B09 的请求局部分组副本使 routes/gateway.go:49/215 进入 OpenAIGatewayHandler.Responses，而不是 GatewayHandler.Responses。

Responses:395 固定 HTTP 入站 → :676 SelectAccountWithSchedulerForCapability → openai_account_scheduler.go:2167 包装层 → :2298 selectAccountWithSchedulerOnce。

Once 先按原生顺序装入 quota-auto-pause、group privacy、profit gate，再 NormalizeOpenAICompatiblePlatform。紧接着识别 trusted internal target 并直接返回指定账户的 selection。插入点在 guardian、previous_response_id、session sticky、legacy/load-aware、advanced scheduler 之前。包装层取得这次结果后，对内部请求直接返回，禁止随后 quarantine fail-open 第二遍选号。

外部请求没有内部 context，执行原生分支。进程内 context 不能由公开 header 创建；身份规则沿用 B08。

### 重写现有 directed helper

修改新增的 openai_eval_directed.go.in，保留一个目标资格 helper，不继续依次调用 getSchedulableAccount → snapshot refresh → DB recheck → snapshot hydration 的重复链。

具体顺序：

1. 核对 context 里的 GroupID/AccountID 和请求分组；读 accountRepo.GetByID(AccountID)，保留数据库原始 error，不改成“账户不支持”。不存在用原生 ErrAccountNotFound。
2. 核对精确 PlatformOpenAI 和原生组成员关系；不能把其他多协议供应商归一成 OpenAI 后放进来。
3. 调用原生 openAICompatibleAccountEligibilityFailureReasonBeforeProfit。它已有模型归属、平台、IsSchedulableForModelWithContext、quota pause、模型与 endpoint capability 判断；直接返回它的真实 reason。
4. 补上原生筛选链已有而上述 helper 未包含的条件：group privacy、parentHealthyForShadow、runtime block、ApplyAccountSchedulingThreshold、proxy quarantine、upstream channel restriction、profit veto。全部复用原有 predicate，不添加评估专用限制。母账户查询直接读 repository，额外保存其原始查询 error，避免原生返回 nil 后无法恢复详情。
5. 复用已有 channel pricing/transport 检查；只接收这次 HTTP Responses 所需 capability，不扩展 image/compact/WS 入口。

原生依据：openai_gateway_scheduling.go:389/1549/1608/1709；openai_account_scheduler.go:1774；account.go:191。DB GetByID:265 经 accountsToService 取得组成员与账户字段。

不能为了复用返回函数重新把已查数据库的账户换回 snapshot：直接构造 AccountSelectionResult，并调用 attachSelectionProfitGate；不调用 newSelectionResult/newAcquiredSelectionResult 的 hydrateSelectedAccount。利润门仍是原生的原门。

### 取得并发槽位与发送前重查

用原生 tryAcquireAccountSlot 抢指定 ID；未抢到则创建原生 AccountWaitPlan，继续由 acquireResponsesAccountSlot/acquireOpenAIAccountSlot:2269 的 quick-acquire、queue、wait、cancel、profit check 管理，不写新队列/等待配置。

抢槽成功后，在 Responses:799 Forward 之前再执行同一资格 helper，读同一 ID 的最新 DB 行。重新装入原生 quota/privacy context，复用 selection 携带的 pricing/profit context。成功替换本次 account 和 selection.Account；失败调用已包装的 accountReleaseFunc 一次后返回具体错误。wrapReleaseOnDone:188 已使用 sync.Once；取消与显式释放不能重复释放。

例：A 排队时被禁用或移出该组，不能用开始时的快照继续发给 A，更不能转给 B。

边界：重查不是锁住账户直至网络发送的事务；最后核对之后仍可能发生并发更新。并发上限变更按原生槽位语义执行，不新增抢槽循环来追求瞬时绝对一致。

### 删除项及理由

删除旧评估规则 61/62（两个备用 selector）和 66（WS previous-response selector）：中心分支已提前 return，实际入口无法进入它们。

删除规则 55、67—70（WS partial/cancel/read detach 修改），评估 HTTP transport 不选择 WS。删除新增 GatewayService ForwardAsChatCompletions/ForwardAsResponses context 规则 50/51；已批准的 OpenAI 组走 OpenAIGatewayService，不走这些另一服务的方法。

保留原生 Responses、OpenAI passthrough、同账户 Responses→Chat Completions；保留原生 response ownership 与会话状态。新增 helper/patch 仅固定 v0.2.16；多版本 AST/兼容配置不被重新批准。

## B12：只评估所选账户，保留原生同账户有限重试

### 精确修改三个 handler 边界

1. Responses:742 continuation 不兼容分支：仍释放 selection 已抢槽位、保留实际 continuation 错误；内部请求直接结束，不执行排除 A 后 continue。
2. :772 抢槽利润否决分支：slot helper 已释放槽位。内部请求保留原生 veto reason 后返回；不能重复释放、把 A 加入排除集后寻找 B。必要 reason 在 slot helper 产生 veto 时写进这次请求诊断，避免最终只剩聚合文案。
3. :887 failover 错误分支：保留原生安全重放判断、健康反馈、ShouldRetryNextAccount/first-output 判断及 :907 同账户 retryCount/retryDelay。只有在原生同账户重试没有继续时，内部请求在 :930 RecordOpenAIAccountSwitch 前返回最终错误。

同账户重试的 continue 会重新经过 B11 同一目标资格/并发检查；不修改有效次数、等待间隔或原四类评估自己的重试。取消仍中断 delay。语义输出已写出时保留原生禁止不安全重放条件，不能为了“必须出结果”再发一次。

继续保留 Forward 内已有同账户的请求修正/协议 fallback；没有新增备用模型或备用账号。真实账户 cooldown/quota/health side effects 不因评估免除。

### Sticky 只留一处判断

保留规则 64：OpenAIGatewayService.BindStickySessionAfterProfitAdmission 入口对内部评估 return nil。HTTP slot helper 的三个取得槽位路径均调用或可能调用该方法。

删除 63（eager binding helper）与 65（低层 setStickySessionAccountID）的重复保护。它们的其他调用来自被 B11 绕过的候选选择；HTTP 路径经 64 已终止写入。

保留原始 session/thread/Codex headers、同账户 reasoning cache，以及 bindHTTPResponseAccount 保存的 response_id 归属。后者是 response 归属/续接安全记录，不是把普通 session 改绑到账户 A 的路由 sticky；它使用 WS 命名存储也不能删。

例：A 的原生可重试 503 先按既有次数尝试 A，最后返回 A 的 503 原文。组内 B 不参与此评估。A 已产生文本后失败，按原生条件结束，不重放成功的文本。

依赖：B11 的中央入口与包装层不 fail-open；B13 最终错误映射。不能仅删 handler switch 然后继续允许 selector 偷换目标。

## B13：在原文被丢弃前保留诊断，交回现有结果错误字段

### 一个请求局部状态，不新增公开协议

在现有内部 context 包 evalroute 中加入一个小型 request-state，runtime 每次 Execute 创建；只存本次尝试的原始 cause、上游 HTTP status/body、SSE event/data、usage-observed 标记及完成信号。实际结构/命名在固定 patch 内实现，不建全局 error bus、表、API 或永久文件。

该状态只可通过 B08 trusted context 使用，不接受公开 header。普通客户端不进入诊断逻辑。结果仍是既有 Error 文本加已有状态/用量字段；B05 同时去除评估层统一错误码吞原文的行为。

### 精确捕获点

| 原生位置 | 修改内容 |
| --- | --- |
| openai_upstream_transport_error.go:117 | 在 sanitize/固定 502 replacement 前保存实际 error；原生 ops/cooldown 分类照常执行；不修改全局 UpstreamFailoverError 的公开行为。 |
| Forward:1095；passthrough:390；responses_chat_fallback:131 | 三个可达 HTTP 非成功响应点改用内部感知的读取 helper；普通请求仍原生读取。内部完整读取正文，保留已读字节和 read error，必要时回卷 Body 供原生分类/协议 fallback 再读。不能用日志截断后的字符串冒充原文。 |
| openai_gateway_response_handling.go:525；passthrough:2050 | 在 error/response.failed 的客户端规则/terminal synthesis 改写前复制原 event/data，原生 usage 解析、账户反馈仍执行。 |
| openai_gateway_cc_pipeline.go:285 scanCCStream | 内部评估捕获 CC error 对象与原始坏 JSON/解析 error；原生忽略这些内容可能生成假 completed。内部遇 error 对象先保存，继续按已有请求期限读到 DONE/EOF 收取后续 usage，结束时以 scan.Err 交回既有 partial result；坏 JSON 按真实解析失败结束。外部普通 CC parser 规则不动。 |
| Responses:676/772/799 | 保存目标查询/资格/抢槽原始 error或具体 veto reason；Forward 调用返回后标记本次实际成功/失败，不能只接收 handleFailoverExhausted 的泛化文本。 |

三个网络发送点使用 B14 已需改动的 context helper，在每次实际发送前重置上游诊断/观测状态。Responses 每次 Forward 开始也重置一次 native-attempt 状态；同账户重试成功不携带上次失败，内部协议修正后成功也不被旧 400 判失败。结算数据是每次 actual attempt 自己的，不借诊断 reset 删除已提交用量。

原生 stream 允许 bare error 后出现成功 completed，特别是 OAuth 的 codexFailureTerminal 分支。诊断中的 bare error先是暂存；依照原生最终成功/失败决定是否展示。不能见到 error event 就主动中止 native parser。response.failed 的原文保留，即使给普通客户端合成了另一份 payload。

### 回传与状态语义

native handler 全部完成后，runtime 从 request-state 合并最终诊断到既有 ModelError，保留 Unwrap 原 cause。producer 与 consumer 通过完成 channel 同步；deadline 未能等到结束时不能并发读取尚未完成的状态并声称最终诊断已齐。

- 原 transport 错误：Status=0，没有上游 HTTP response；保留 proxy/DNS/TLS cause，不把本地 502当上游502。
- 上游 HTTP 429/500：展示真实上游 status/body；已读取一半断开则同时展示半份正文和读取失败。
- HTTP 200 后 SSE failed：仍是 200 的 stream failure，显示原 event/data；没有制造“上游 HTTP500”。
- 调度/权限/转换/排队失败：说明实际本地失败来源，不说模型不支持。

既有真实凭据值脱敏继续使用；不拿“可能回显凭据”为理由删除整份错误详情。完整错误体的读取受现有请求截止/取消影响，巨大正文仍消耗内存；本项没有引入新的截断配置或跨站诊断协议。

这条 HTTP 返回管道可以继续保留普通 native response 映射，因为内部诊断从同请求 context 原文取回；不再额外改全局 handleFailoverExhausted/客户端策略。没有公开诊断接口，没有 Cause/Recovery/Next UI。

## B14：真实取消、完整收尾与已观测部分用量

### 实际上游只改三个 context 点

保留并收窄规则 54/56/58：sendCCUpstreamRequest:203、Forward:1038、forwardOpenAIPassthrough:372。内部请求用任务/request context 创建 upstream request；普通请求调用原 detachUpstreamContext。helper 的 release 必须遵守原 detach 的调用约定，不能在构造 request 后立即 cancel 真正请求。

停止/超时立即取消这些请求并关闭 pipe/upstream body，不发下一题或重试；原生 slot release 仍由 wrapReleaseOnDone 控制。

不改 Anthropic、Google、WS 或另一 GatewayService。供应商已接收请求后的执行/扣费不能撤回。

### 修正 pipe 与 executor 的完成顺序

修改 eval_http_transport.go.in、eval_host_runtime.go.in 及 internal/eval/http.go/sse.go 的必要收尾逻辑：

1. runtime 为每次 Execute 建一个独立 done channel，注册 producer 后再调用 router。done 只有 router.ServeHTTP 返回、handler defer/提交调用完成后才关闭；全局 producerDrain 仍追踪真正活着的 producer。
2. 正常 terminal 不调用会 cancel 的 Body.Close。parser 记下结果/终态后，继续消费内部 handler 的剩余输出至 pipe EOF；忽略成功 terminal 后的 DONE，而不是误判 missing_terminal。必须消费尾部，不能一边停读 io.Pipe 一边等 handler，以免 writer 阻塞造成死锁。
3. bare error 先暂存；failed terminal/EOF 按真实失败收尾，允许 native 的后续权威 terminal 和 usage 完成。流大小/解析上限不借此重新设计，本次只改完成顺序；原 core 的其他额外限制随后另批审。
4. 读到 EOF 后等本请求 done，再取最终诊断，最后释放 Body。正常收尾只证明 handler 已结束与 usage submission 已执行，不证明异步数据库落库。
5. 手动取消时马上 cancel并关闭 pipe，随后只在该请求原本的绝对截止时间以内等待 native handler 清理，不重建一个更长请求窗口。取消不能让读取已完成的普通请求误变“用户取消”。若原期限已到还没结束，返回明确取消/超时及清理未确认信息，保留 producer 跟踪，禁止报告已全部释放；B15 的全局退出等待独立受原五秒窗口限制。

当前请求期限来自已有 askCodex/executeProbe；不在本批批准新的分钟数。后续 core 审阅若修改/恢复原期限，收尾仍遵从那一次绝对截止时间。手动取消遇到不响应 cancel 的底层清理可能仍等至该截止；这不是“立即完成取消”的保证。

### 用量如何回到原生结算

原 stream parser 已有 resultWithUsage。Forward:1263 与 passthrough:477 在 error 时丢掉它，是必须修改的两处；保留且简化规则57/59，让它们把原 streamResult转成正常 OpenAIForwardResult并与原 error一起 return。

复用正常路径的 request_id、响应 headers、原模型/计费模型/上游模型、回显 tier、reasoning、duration/TTFT 及实际使用量；不要丢掉原生计费元数据或另写收费公式。只在 parser 确认实际 usage字段出现时提交部分 token用量；零 token 的真实 usage也属于观测，不能把默认全零对象当观测。B13 的请求态在原 native usage提取成功处记录这一布尔值；并发请求不共享它。其他原生实际计费项如有返回，沿既有判断处理，不增加图片/WS评估能力。

CC streamChatCompletionsAsResponses:243/259 已返回 scan.Usage 与 error，保留原函数，不重复增加 CC partial builder。只让 scanCCStream 的上述内部错误识别进入已有错误分支。

在 Responses错误处理的 failover分支补一次 submitResponsesUsage(result)，位置在任何 early-return/retry之前。该分支与原 ClientDisconnect、普通错误、成功提交分支互斥；同一 native Forward attempt只能提交一次。现有规则35需要内部作用域与已观测条件，不能扩大修改普通请求语义。不同 actual attempts不合并；保留 native request_id/payload/context去重规则，不自建全局去重表，也不声称无 request_id时具备更强的重复扣费保证。

保留 wrapUsageRecordTaskContext:265 与 native worker本身的生命周期；任务cancel不能把已提交结算使用的worker context一起取消。不增加队列、补账服务或强制同步落库。

例：A已报告usage之后断开，原result与真实error一起返回并执行native提交；没有usage时不编造0费用。队列配置drop/sample、数据库错误、deadline或强制kill仍可能使记录缺失。取消不是免费，成功提交不是已落库。

删除规则60的评估非stream partial，因为四类请求stream=true；原生非stream能力保持原样。B11已删除WS partial。保留已存在的正常计费算法，不为评估扩展媒介计费框架。

## B15：恢复监听；新增退出等待只用原五秒期限

### main与runtime改法

恢复 cmd/server/main.go:173 的 ListenAndServe，删除为维护socket引入的手工 listener/Serve与相关路径。Application/runtime 的普通DI接线保留，B10删除维护协议不等于删manager。

SIGINT/SIGTERM到达后：

1. 同步关闭评估 admission并标记stopping；锁内仅更新状态/取得cancel列表，锁外逐个cancel。
2. manager停止接收并cancel runs，runtime cancel所有注册producer。
3. 创建原生那一次context.WithTimeout(...,5s)，让 runtime.Shutdown与Server.Shutdown共享它并行执行。
4. 等两者完成或该context过期，分别记录实际error；任何timeout都不renew。buffered completion channel保证main停止等候后发送方不因报告结果卡住。
5. 进入原生Application.Cleanup，不在eval补丁中再添一轮等待或改写所有服务stop顺序。

runtime.Shutdown等manager drained和producerDrain两个条件。manager因storage error结束不代表producer已结束；不能因为manager返回error直接跳过producer检查。所有等待使用同一个deadline，不持runtime mutex等待/IO。

### 为什么不能只等http.Server

内部router.ServeHTTP不属于net/http接收的真实connection，Server.Shutdown不会追踪它。manager算法返回也不必然说明native handler已提交账单或释放完defer；因此B14 done/producerDrain是必须的，不是另一套维护状态机。

### 真实边界

五秒内完成时才报告评估已drained；过期则准确报告未完成。仍活着的producer不能从计数器里假删除。

原生wire.go会并行Stop服务，再关Redis/Ent；其中usage worker.Stop:217调用StopAndWait，文件history Save也不是context可强制中断操作。保留这些事实，不保证全进程五秒/十秒退出、不保证deadline前所有历史和费用落库、不假装kill可以保存未收到的结果。

主进程继续原生cleanup可能与期限后残留producer重叠，Docker强制kill仍可能中断它们。本批不为此重建整站cleanup框架或延长Docker停止时间。重启history如何显示中断任务由后续任务状态批次审，不自动续跑/补发。

## 旧评估规则处置对照

编号是当前 adapter/rules.json.Rules 的零基索引，只供审计定位；不是继续批准AST适配框架。

| 规则 | 处置 | 归属/依据 |
| --- | --- | --- |
| 31 | 保留入口识别，随B08 context简化 | B11/B12 |
| 32/33 | 保留终止跨账户边界；保存具体错误 | B12/B13 |
| 34 | 重写fresh DB recheck/具体错误；release一次 | B11 |
| 35 | 收窄为内部observed partial提交一次 | B14 |
| 36 | 保留同账户retry之后、跨账户switch之前的返回 | B12 |
| 50/51 | 删除另一GatewayService的eval context修改 | B11固定OpenAI route不可达 |
| 52/53 | 保留中心选号/禁止quarantine第二遍 | B11 |
| 54/56/58 | 保留三个可达HTTP upstream context点 | B14 |
| 55/66/67/68/69/70 | 删除eval WS修改 | HTTP transport不可达 |
| 57/59 | 简化为真实HTTP stream partial与native元数据 | B14 |
| 60 | 删除eval nonstream partial | stream=true |
| 61/62 | 删除备用selector拦截 | 中心direct return不可达 |
| 63/65 | 删除重复sticky拦截 | HTTP由64覆盖，其余被中心入口绕过 |
| 64 | 保留一处post-admission sticky保护 | HTTP slot实际可达 |

本批还需要上表原先没有的原文捕获/CC错误识别/请求done收尾；这些是解决已明确缺陷的具体点，不能只“删到更少”然后继续吞错误或提前cancel。

## 审批关系与后续

B11依赖已批准B03/B08/B09；B12依赖B11；B13依赖B05与B08；B14与B13共享一次请求的小状态，依赖B12单次attempt控制；B15依赖B10去掉维护协议和B14真实producer完成信号。

五项仍待用户逐批批准。原四类算法/题库、其他core改动、history/restart、固定构建交付与cleanup仍另批审。没有追加race、付费验证、全量测试或构建要求；本轮源码调研不能写作运行验证成功。
