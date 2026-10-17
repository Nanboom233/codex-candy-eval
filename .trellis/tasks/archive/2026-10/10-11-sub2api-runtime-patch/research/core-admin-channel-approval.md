# 管理员通道核心审批：C02、C03、C04、C06、C07

状态：planning，本批五项已由用户“是 继续”全部批准。用户“无误 继续”已批准S01—S05；C01必要宿主边界由S02/S05覆盖，C05旧进程内/本机用户网关执行器由S02/S05取代，无须重复审批。本批只处理仍需明确的核心差异，不授权产品实施。

## C02：结果只保留一套模型与推理参数

现状：internal/eval/eval.go:59—73和:133，同时保存model/effort、requested_model/requested_effort、effective_model/effective_effort；runs.go:352—355又重复构造。effective值是插件自己的转换结果，不是账户实际出站映射或供应商实际执行模型的证据。

建议：删除requested_*、effective_*及生成/读取/校验代码。保留model/effort，表示评估交给管理员发送方法的模型与推理参数；指纹填固定low，ModelTrace填未指定。保留group_id、account_id、account_name、run_id和status，分别用于页面分组/来源、名称、任务关联和状态。核心直接使用这一套model/effort，不增加模型证明字段。

影响：保留原页面按model/effort筛选和展示，四类结果均有来源信息，不增加结果列。模型映射继续由已批准的原生发送代码执行；不能把结果model当作供应商实际模型证明。历史字段引用和现有相关测试需同步；具体历史存储/格式在后续C19等审阅，不加入迁移兼容器，不删除真实用户历史文件。

## C03：恢复原推理选项，删除CPA命名转换

现状：eval.go:224—267的NormalizeModel解析model(high)括号后缀，自建effort白名单，拒绝max，把none改为空；web/app.js:8从原plugin/web/app.js:8删除max。原页面candy.js:61和pelican.js:108已独立提交model与effort。

建议：模型ID来自已批准的原生模型接口，作为独立model发送；reasoning.effort独立发送。删除CPA括号后缀解析、后缀冲突拒绝及插件推测模型支持范围的规则。保留原有非空/长度/控制字符等基础输入检查，原页面恢复none/low/medium/high/xhigh/max。none按字面发送，空值才表示未指定，插件不自行把max降成xhigh。旧CPA括号模型名不做转换，不加旧命名兼容分支。

影响：停止擅自改变用户选择。管理员发送方法按S02使用原生账户模型映射/必要请求准备；特定模型拒绝某种effort时按S03/C06呈现真实错误。当前选择的是管理员测试通道，不能又套回用户网关分组effort策略。原生具备max相关处理并不保证每个账户上游都接受max/none。

依据：原生internal/pkg/openai/constants.go:221—229有模型相关effort限制；原生openai_responses_astra_pro_failover_test.go:128等有max处理证据；没有在本轮运行这些测试。

## C04：集中保留原评估选项规则

现状：eval.go:95—132和runs.go:141—199重复规则；增加指纹/ModelTrace的effort拒绝、后缀冲突拒绝和负数Runs拒绝。原plugin/fingerprint.go:91起校验mode及1—6并发并固定low，modeltrace.go:111起使用1—3并发且不指定effort，candy.go:64将Runs限制到1—10。

建议：任务入口集中准备/校验一份参数，执行器使用已准备参数，不再维护两张相同规则表。保留四类类型、原模型基础检查、指纹已知模式/1—6并发及默认2、ModelTrace的1—3并发及默认3。Candy次数按原min(max(runs,1),10)，其他类型固定一次。指纹固定low，ModelTrace不指定effort，不额外开放其推理选项，也不为本来没有的选项增加拒绝逻辑。保留原指纹/ModelTrace不使用CPA括号模型名的选择规则。

影响：模式/采样参数仍有效，次数0或负数按原Candy行为执行一次；题库、采样数量、评分及并发参数不改。执行时仍检查取消。本项不连带批准现有任务恢复/调度框架，也不替代S01完整请求发送接口的必要JSON检查。

## C06：恢复实际错误，保留原异常和取消处理

现状：internal/eval/http.go:70—100读取错误body后只返回http_status，sse.go将error/failed压成固定code，model.go:32—43及:136—149把panic压成executor_panic；String/GoString另行隐藏诊断。原plugin/model.go:89起保留实际异常，:186—249保留流中断/错误消息。

建议：按已批准B05/S03保留评估实际收到的HTTP状态、错误正文/消息和SSE失败原因。保留简单的带HTTP状态错误以供原重试策略使用，删除吞正文的固定code替代与无必要String/GoString。服务端沿S03保留原生必要敏感内容脱敏，核心不再重写一个诊断或脱敏框架。

保留原plugin已有的请求边界panic恢复，转为可读异常，防止后台任务异常导致进程退出；不为panic另造新状态码或新重试类别。保留关闭响应、停止后续调用、传播context取消/超时原因和不把取消后的请求记作成功。具体重试次数/等待/usage及任务终态另在C09/C10及后续任务条目审阅。

影响：管理员能看到实际失败原因，仍保留原异常容错和取消语义；不再出现所有HTTP失败都只显示http_status。新发送通道已能拿到原始错误，不需要改用户网关或新增错误旁路。

## C07：以原判定语义读取完整Responses流

现状：internal/eval/sse.go:52—215已有增量读取，但增加坏JSON立即失败、item ID替换表、尾部没有空行一律拒绝等行为；原plugin/model.go:186—249在收完CPA流后解析完成item与terminal，保留完整错误，并在最终output为空时使用已完成items。

建议：把原Responses判定逻辑接到新通道的io.Reader，按SSE事件读取，不保留CPA stream_id/hostRead及修补CPA分块换行的逻辑。保留正常注释、事件名和多行data处理，这是实际SSE协议读取，不是命名/版本兼容层。

原终止response.completed/incomplete/failed语义保留；原生测试已经接受response.done，读取新通道时也按该原生终止事件处理。取消、未完成和失败均保留其实际状态/错误，不拿途中delta冒充成功。最终output为空时回用已完成items；只提取message中的文本，不把工具参数当答案，保留工具无直接回答及空答案错误。

删除新增item ID替换表，回到原完成item收集逻辑。未知非终止事件/无关不可解析片段沿原行为忽略，实际终止响应解析失败仍报错。尾部已经完整收到、能解析的终止事件按原行为处理，不因缺少最后一个空行额外拒绝；没有完整终止则报流中断。流大小限制由C08独立审阅，不默许本批保留其新增数值。

影响：完整回答、失败原因及用量数据能进入核心；评分和题库不变。保留增量读取所需的少量状态，不再全量保存原始网络流或引入生产者框架。原生转换函数/普通网关parser仍不改。token累积细节另在C09审阅。

## 审批与证据

本轮读取了原core/parser、原CPA对照及现有盘点，并按新仓库指令用fast-context核对原生Responses事件定义/终止行为。原生apicompat/types.go:613—658提供事件/usage定义；account_test_service.go:2947—3015接受response.completed/response.done。本轮没有运行编译、测试或实际供应商请求。

以上五项已由用户“是 继续”全部批准并记入审批账本。S01—S05已批准且记录，任务保持planning；剩余C08—C34及交付清单继续逐项审批。批准本批不等于全任务最终实施批准。
