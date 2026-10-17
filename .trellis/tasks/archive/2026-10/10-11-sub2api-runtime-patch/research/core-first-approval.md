> 旧C01—C05未批准；用户网关发送/内部Key方案现已明确撤回。当前先审批S01—S05管理员接口，再重新核对核心条目。

# 核心第一批审批：C01—C05

> 用户已暂停本批，转而考虑原生连通性测试发送方案。下述五项未批准；尤其C05不得作为实施输入。见native-connectivity-dispatch.md。

状态：待用户审批，仅更新规划。继续采用共享内部Key、每次评估HTTP调用前重置5分钟有效期的已确定方向。本批不决定共享Key的原生所属用户，不恢复逐管理员记账，不授权实施。

用户要求五项一批，并完整说明；因此本批一次展示五个独立决策。现状清单是证据，不是默认保留授权。

## C01：保留必要的宿主调用替换

现状：internal/eval/eval.go:13、:44和internal/eval/model.go:20—30使用Target、EvaluationRequest和单方法ModelExecutor。原internal/plugin/model.go:98起通过CPA的host.model.execute_stream调用宿主，Sub2API没有这个接口。

建议：保留group_id/account_id目标、四类评估所需输入及单方法执行器边界。执行器负责发送原有Responses请求并返回文本、token用量及错误。账户名称只用于展示。模型调用实现改为C05的本机HTTP；不引入CPA的C ABI、credential回调、inspection定时器或provider路由。

影响：四类题目、算法和请求内容保持原有语义，只替换实际发送请求的入口。执行器接口也供现有本地fixture验证使用，不是新的运行时框架。

## C02：删除重复的模型/推理强度结果字段

现状：internal/eval/eval.go:59—73同时存model/effort、requested_model/requested_effort和effective_model/effective_effort。eval.go:133和runs.go:352—355重复填写它们。这些effective值是插件本地转换结果，不是原生网关完成模型映射或供应商实际处理后的证明。

建议：结果只保留model、effort一套值，表示本次发送给原生网关的模型ID及推理强度。保留group_id、account_id、account_name、run_id和status，分别服务于结果来源、显示、任务关联和状态展示。删除requested_*、effective_*及对应生成、消费和校验代码；四类核心使用model/effort即可。真实原生模型映射及计费记录由Sub2API负责。

影响：前端原有筛选和结果显示仍读取model/effort，不新增展示列。指纹结果记录其实际固定发送的low，ModelTrace记录未指定effort。历史持久化消费者和相关测试需同步更新字段引用；具体历史格式/保存策略仍在后续C19等审批，不用兼容分支或迁移框架保住本批删除的字段，不删除真实用户文件。

## C03：模型ID与effort分开传递，撤回擅自更改的取值

现状：internal/eval/eval.go:224—267解析CPA式model(high)括号后缀，维护一套effort白名单，拒绝max，并把none转换为空。internal/eval/web/app.js:8删掉原internal/plugin/web/app.js:8里的max。原页面请求在candy.js:61和pelican.js:108本来已分别传model和effort，无须用括号承载两个字段。

建议：模型ID直接使用已审批的原生模型接口返回值，effort独立放入reasoning.effort。删除CPA括号后缀解析、后缀冲突判断和插件自己推测模型支持范围的代码。保留原有非空、长度及控制字符等基本模型输入检查。保留原页面none/low/medium/high/xhigh/max选项；none按字面发送，空值才表示未指定。max不由插件自行改成xhigh。具体模型是否接受某种effort、原生分组策略及模型映射都交给Sub2API。

原生证据：v0.2.16 internal/pkg/openai/constants.go:221—229对gpt-6.1-sol拒绝none/minimal，而非改成默认；internal/service/openai_gateway_request_body.go:1748起区分请求effort与原生映射。已有openai_responses_astra_pro_failover_test.go:128和openai_compat_model_test.go:255起覆盖max相关原生行为。这是源码证据，本轮没有运行测试，不代表所有模型都支持max或none。

影响：恢复被擅自删除的max选项，停止把显式none悄悄换成默认。原生拒绝时按B05显示评估实际收到的错误。旧CPA括号模型名不再转换；当前交付按破坏性变更，不加命名兼容器。不修改原生推理强度处理。

## C04：选项校验只保留原评估规则

现状：internal/eval/eval.go:95—132和runs.go:141—199重复校验；另加了指纹effort/后缀冲突拒绝、ModelTrace非默认effort拒绝以及负数Runs拒绝。原plugin/fingerprint.go:91起只接受mode和1—6并发并固定发送low；modeltrace.go:111起接受1—3并发并不指定effort；candy.go:64将Runs夹到1—10。

建议：任务入口统一做一份评估选项校验，执行阶段使用已经准备好的参数，不重复建立同样的规则表。保留四类类型判断、原模型基本检查、指纹已知模式及1—6并发、ModelTrace的1—3并发和原默认值。Candy次数按原行为限制在1—10，其余类型只执行一次。指纹固定low，ModelTrace不指定effort，不为这两类额外开放推理选项，也不增加针对本来不存在的effort字段的拒绝逻辑。保持原指纹/ModelTrace不使用CPA括号后缀模型的选择规则；原生模型是否可用继续由已审批的原生资格规则负责。

影响：不能进入错误模式或越界采样；不改变题库、采样数量、评分算法及并发参数。次数0或负值按原Candy行为落到1次，不另加新错误。上下文取消等执行期检查继续保留；本项不决定后续任务管理、重试和持久化机制。

## C05：替换自制进程内传输，删掉重复请求比对

现状：internal/eval/http.go:21—70包装runtime自制RoundTripper，访问虚拟codex-eval.internal地址；再次反序列化已经由插件生成的body，比对model/effort/stream并重复NormalizeModel。认证和目标靠私有context传输，而非已确定的内部Key与选择头。

建议：执行器使用标准http.Client访问原生配置对应的本机POST /v1/responses。每次发送前续期同一内部Key至当前时间+5分钟，设置Authorization及X-Codex-Eval-Target；续期失败不发送。保留原Codex请求内容和协议头，由执行器覆写内部认证/目标头。禁用环境代理和重定向，防止内部凭据被带往别处。删除虚拟域名、私有context能力、自制RoundTripper/io.Pipe及重复body反序列化/模型比对。网关中的真实鉴权与硬目标校验继续保留。

影响：评估请求进入真正的原生HTTP入口，转发、并发、协议转换、同账户重试和计费沿原生路径。取消关闭本机请求并停止后续评估工作，原生上游如何完成读取和结算仍按原生语义。HTTP/SSE错误呈现、流解析和资源大小限制分别由后续C06—C09审阅，本项不默许保留当前吞错误或重写parser行为。

## 审批边界

五项全部待审。批准本批只授权这些规划决定；共享Key归属、后续core/history/delivery及全任务最终实施摘要仍未收敛。本轮未改产品代码、未测试、未构建或部署。
