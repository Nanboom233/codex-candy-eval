# 管理员账户完整评估发送接口：S01—S05

状态：planning。用户明确选择管理员测试通道，舍弃之前的用户方案，并指定新入口/api/v1/admin/accounts/:id/codex-eval。该方向取代内部Key、五分钟续期、用户记账及池内定向方案。用户回复“无误 继续”批准S01—S05全部具体规划；不授权产品实施。

## 已确定的方向

- 新增POST /api/v1/admin/accounts/:id/codex-eval，以现有/test的账户发送代码为基础；原/test无需改成完整接口。
- 只用原生管理员认证，不准备内部用户/Key/记账行，不添加选择头或网关鉴权特例。
- accountID来自路径，发送只使用该账户；不选择其他账户，不进入普通用户网关计费和账户池调度。
- 保留此前已批准的页面、四类评估、原题库/算法/数据与固定版本预编译交付。
- 舍弃方案不等于已删除旧代码、安装包、密钥或用户数据；本轮只改规划。

## S01：单次请求契约

建议在现有admin/accounts路由组新增POST /:id/codex-eval，由AccountHandler的新方法处理。原生adminAuth、面板管理权限与审计仍负责对外访问。它是单次模型请求发送入口，不是另一套评估任务start/state/cancel接口。

请求体直接使用完整Responses JSON，不再包model_id/prompt/mode，不固定hi，不从body再读另一个account_id或group_id。保留原评估生成的input、reasoning、tools、include、prompt_cache_key、client_metadata等原始内容；账户model mapping和必要的原生OAuth请求准备继续复用。必要Codex协议头随HTTP请求传入，服务端上游鉴权/目标/代理来自账户配置，管理员Bearer不转发给供应商。

本批建议明确支持原四类评估使用的stream:true完整请求；不默默把“完整”扩大为非流式聚合、compact、图片、音频或任意其他OpenAI端点。非流式范围未被当前/test实现证明，不能声称完整OpenAI兼容。这一具体边界已由用户批准。

影响：既有前端评估任务接口继续承担任务管理；新入口可以作为管理员对某个账户发送完整评估请求的直接接口。页面不增加登录、激活、密钥或新的菜单步骤。

原生依据：admin.go:29—36及:393，admin/account_handler.go:1265—1317。

## S02：按原测试逻辑复制指定账户发送

建议新增独立account_codex_eval_handler.go / account_codex_eval_service.go，以现有OpenAI账户测试发送逻辑为基础，使用现有AccountHandler.accountTestService依赖。新service方法挂在现有AccountTestService，不增加身份服务/DI框架，不为了这一接口重构原/test。

复制/复用的内容为：GetByID、凭证影子解析、OAuth或API-key出站鉴权、账户model mapping、上游URL校验、账户头覆盖、代理、TLS与doOpenAIAccountTestUpstream。改变的内容为：将固定测试payload换成调用者完整评估body，将简化测试事件输出换成S03的标准响应。请求失败仍只有该账户，没有候选池或跨账户fallback。

保持当前评估目标的OpenAI OAuth和API-key范围。仅支持Chat Completions的OpenAI账户用原生ResponsesToChatCompletionsRequest及ChatCompletionsChunkToResponsesEvents等函数转换完整请求/响应；不能只取最后一段prompt冒充等价。原生转换函数源码不改，未实际fixture验证前不宣称转换已经通过。不顺便扩展Anthropic、Gemini或媒体测试。

影响：业务网关的API-key中间件、账户scheduler、sticky策略、计费、parser、ListenAndServe及shutdown无需为评估添加补丁。原/test、原生插件系统与原生转换库保持原职责。

依据：account_test_service.go:800—1000、:2118—2186；openai_plugin_transport.go:23—48；apicompat/chatcompletions_responses_bridge.go:39、:1534、:1637、:1687、:1800。

## S03：返回完整流、用量和真实错误

Responses上游返回标准状态码和SSE事件，建议新入口保留完整response.completed/output/usage/error信息，不转换成content/test_complete/error摘要，不先写200再请求上游。非成功HTTP返回保留真实上游状态及响应错误内容，沿用原生Agent Identity敏感内容脱敏，不增加诊断旁路或另一个错误接口。

Chat Completions账户由原生转换函数输出统一Responses事件，保留已有转换语义和可获得的usage；供应商不返回用量时不编造。核心现有流解析的逐项审批仍待后续完成，本项不授权保留当前吞错误或新增另一套parser。

影响：Candy/指纹/ModelTrace/鹈鹕能接收到真实回答、完成状态、错误和可用token数据。token展示属于管理员评估记录，不创建用户usage账单、不扣Sub2API用户余额，供应商请求仍消耗供应商额度。

依据：当前account_test_service.go:2947—3015的简化输出不携带usage；原生发送入口已经返回标准http.Response。

## S04：明确账户状态副作用

建议沿用原测试发送中的账户快照更新、401标记及429限流信息处理，不新增评估专用恢复/修复规则。新接口不因每次评估回答成功而调用RecoverAccountAfterSuccessfulTest去清除全部error/rate-limit/overload/temp-unsched/model-rate-limit状态；普通/test仍保持它原本的恢复行为。

理由：新接口返回标准响应流，既有测试处理器的成功恢复依赖其简化流处理成功返回。不能把io.Copy正常EOF当作完整评估成功，进而改变运营状态；为此另添终止事件观察器和恢复框架也不是本次发送需求。401/429及header快照已经位于现有发送代码中，可以按原测试逻辑复用。

影响：新接口不是完全只读，仍有上述原生错误/额度快照写入；但一次评估成功不会作为显式清除所有限制的操作。不修改普通/test及原用户网关的状态策略。

依据：admin/account_handler.go:1313—1317、ratelimit_service.go:2202，account_test_service.go:966—998。

## S05：四类后台任务复用同一发送方法

建议新HTTP入口和后台四类评估执行器调用同一AccountTestService发送方法。后台任务已由原生管理员Start接口接纳，直接调用service方法；不绕本机HTTP、不保留管理员JWT、不伪造Gin请求或新增内部鉴权。

任务仍由现有start/state/cancel/records/clear-history接口管理，四类算法和采样参数沿原行为。分组只服务于已批准的页面筛选/目标展示，不绑定用户账单。保留原评估自己的采样并发，不套普通用户网关Redis抢槽或再发明并发管理器。

取消传递任务context，关闭本次上游响应并停止后续采样；保持原生HTTP传输管理，不承诺供应商收到取消后立刻停止消耗。流式读取/转换不得重新引入已撤回的io.Pipe、producer注册表或自制退出协调器。恢复原生server启动和退出。

影响：HTTP入口与四类后台任务共享真正的发送实现，减少两条发送链漂移。原页面/UI功能仍保持；历史持久化和任务管理复杂度继续按后续清单逐项审批，不由本项连带获准。

## 需要退出现行方案的内容

内部用户/owner、真实或虚拟记账Key、过期及续期逻辑、私有capability/proof/token、维护socket和激活流程、API-key特殊鉴权、目标选择头、scheduler硬定向、计费/partial-usage补丁、额外sticky/资格/诊断/cancel/shutdown修改全部退出。只在实施批准后删相关产品代码，不能执行真实数据库清理或假称已删除。

## 状态

新管理员通道和路径是用户明确指定的方向。S01—S05是已批准的具体方案，包括流式范围、原生完整转换、状态副作用和任务接入。C01—C05旧批次仍未批准，涉及网关调用的部分已被新方向取代。最终task-wide实施审批仍在后续core/history/build/installer收敛之后。
