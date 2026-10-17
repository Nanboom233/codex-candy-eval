# 后端逐项审批：B01—B05

状态：planning。前端 F01—F06 已记录批准，产品代码未实施。用户在完整说明后回复“批准 继续”，批准B01—B05的规划范围。记账身份或全部宿主钩子不由本批连带批准。

## B01 删除菜单与部署证明对运行就绪的耦合

现状：EvalHostRuntime.Ready/Activate/State调用checkMenu和CurrentEvalDeploymentFacts。checkMenu硬编码/codex-eval/ui、旧菜单标签和字段；前端F01已经选择/codex-eval和账户评估。菜单展示设置或安装证明文件不一致可阻止执行。
方案：删除运行入口、状态和执行路径上的菜单格式/安装来源证明检查。运行准备只检查真实依赖：初始化错误、存储、路由可用、正在退出及已审批的必要运行模式约束；记账身份准备另审。不把菜单删除/改名当作评估服务故障。版本/许可证信息保留在正常构建打包流程。
行为：修改菜单标题/排序或直接访问/codex-eval不使已准备好的服务进入repair_required。数据库、存储或执行依赖真实失败仍拒绝启动并报告原错误。
范围：这里只删除菜单/部署证明门槛，不批准删所有Ready检查、模式限制、安装备份/校验或私有分发边界。F06d手动激活移除仍依赖下一轮记账准备设计。
证据：adapter/runtime/backend/internal/service/eval_host_runtime.go.in:100/151；eval_menu.go.in:12。

## B02 使用原生管理员权限，删除首次激活者专属操作门槛

现状：接口已经位于原生admin路由下，Start又要求actor==identity.OwnerUserID；另一个合法管理员能读取/取消却不能启动。
方案：评估管理操作与原生账户管理使用同一管理员鉴权，不再有首次激活者专属权限、独立登录或owner配置。启动者ID只从服务器已认证会话取，不由请求正文决定，并作为任务执行/记账归属的输入。
行为：合法管理员可发起评估；缺少登录和缺少管理员权限仍按原生401/403处理。评估状态/历史/取消沿用管理员共享管理语义，不新增按人分割的页面。
依赖：本项批准权限规则和启动者归属要求，不批准拿首位管理员的旧记账身份给其他人扣款。当前RunManager任务上下文没有启动者身份，且执行器使用全局身份；二者需与下一轮具体记账设计一起修改，不能孤立删除一行owner判断。原生余额、订阅、限额、并发规则不得由本项绕过。
证据：adapter/runtime/backend/internal/handler/admin_codex_eval.go.in:114；eval_host_runtime.go.in:45；internal/eval/runs.go:315；原生backend/internal/server/routes/admin.go:25。

## B03 删除重复账户列表，保留服务端目标校验

现状：专用/targets和ListTargets重复聚合原生组/账户，但丢掉priority、套餐等字段；前端已批准读原生账户接口。
方案：删除/targets路由、handler、简化展示DTO和ListTargets。启动仍用原生服务读取指定group_id/account_id，核对存在性、OpenAI平台、支持的账户类型、启用状态及成员关系。实际发送时的原生调度资格/限额/并发检查保留，按最新信息处理。
行为：页面显示来自原生API；账户切组、禁用或删除后，旧页面不能靠先前列表快照启动成功。不为列表重复构建第二份目录。
范围：不删TargetResolver或精确账户约束，不从其他账户补位；不改变原生管理字段/凭据脱敏。具体调度钩子另批逐条审阅。
证据：adapter/runtime/backend/internal/service/eval_target_adapter.go.in:53/85；admin_codex_eval.go.in:87。

## B04 模型支持按原生规则判定，删除客户端目录准入

现状：模型目录并集只能用于选择，原生models接口失败时也可能返回默认目录。不能据此证明每个账户都支持，更不能把缺项等同不支持。原生Account.IsModelSupported已有model_mapping、通配符、OAuth及透传模式规则。
方案：删除废弃model_catalog输入及评估自建目录/缓存准入；定向执行复用原生IsModelSupported和实际请求资格检查，不复制规则。用户所选模型继续按原评估已有参数规则处理，不在本项改算法、默认值或模型格式。
行为：A账户支持而B不支持时，分别报告B的原生拒绝原因；不替B换账户。目录读取失败不直接判B不支持；空model_mapping也按原生规则处理，不能自行认为没有模型。
限制：原生规则允许请求仍不保证上游必定接受。上游模型错误按真实返回显示；本项不新增付费探测、启动前/models请求或重试策略。
证据：原生backend/internal/service/account.go:889；openai_account_scheduler.go:1815；adapter/runtime/backend/internal/service/openai_eval_directed.go.in:96。

## B05 恢复评估自身丢弃的原始错误

现状：handler把具体错误换成not_ready/start_failed/storage_unavailable；TargetResolver和RunManager将账户错误换成统一文本；ModelError隐藏cause，HTTP非2xx最多读64KiB后丢弃body；recover丢弃panic内容。
方案：在评估handler、目标读取、任务状态/历史、HTTP/SSE执行和panic处理各处保留原错误与错误链。现有状态码/错误类型仍用于程序判断；显示内容保留实际返回正文，不用代码替代详情。取消/超时也保留实际上下文原因；有实际panic时保留panic值与调用栈，不对普通Go错误伪造traceback。评估自己的非2xx正文不再按64KiB读取后静默丢弃，读取失败保留已读正文和读取错误。
行为：数据库连接、权限、解析、上游HTTP等失败在启动响应、状态及原结果位置提供收到的真实诊断；不新增Cause/Recovery/Next面板或另造诊断框架。原有重试分类和取消行为不在本项改变。
范围与限制：只改变本次新增评估层的错误压缩。原生网关在到达评估之前已改写的上游错误，不能在评估端凭空恢复；具体必要的定向网关错误传递钩子下一批单独审阅，不据此批准全站关闭原生错误处理。鉴权头、令牌和数据库密码不拼进诊断，已有真实诊断不因猜测而整段替换成通用错误。
证据：admin_codex_eval.go.in:36/98/114/174；eval_target_adapter.go.in:190；internal/eval/runs.go:245/284；internal/eval/model.go:36/132；internal/eval/http.go:96/103。

## 尚未决定

F06d正常初始化与记账身份的具体替代方案仍未审批。当前expired API-key行是原生usage/billing引用，不能把手动激活变成后台调用旧Activate，也不能仅删记账。下一批明确决定记账行归属/创建复用、内部身份、停用及任务身份传递，再与上述权限/就绪规则组合审查。旧兼容迁移及私有token/provisioning状态机不会因前端审批自动获准保留。

仅源码核对，无产品改动、构建、部署、真实上游调用或数据删除。
