# 当前 Sub2API 修改完整清单

2026-10-14，任务保持 planning。此清单描述现有实现，不表示批准保留。除用户已要求移除的 race 门槛外，尚未批准任何删除项。C29仍待审批。

当前构建使用 adapter/rules.json、adapter/ui-rules.json、adapter/runtime 及 internal/eval 白名单；历史 v0.2.15 production patch不进入当前构建。安装包仅含二进制、install.py、cleanup.py及五份许可文件。

原生修改直接触及32个原有文件：73条AST规则、13个新增函数定义、15个文件的import注入，另有8条前端文本修改；复制16份运行时模板。下列84项是行为审计分组，部分跨层依赖重叠，不代表84个独立源码改动或84项全部必要。

## 原生服务、鉴权、调度及前端集成

| ID | 当前内容 | 主要源码 |
| --- | --- | --- |
| N01 | 在服务器二进制中加入维护能力查询、状态查询及停用命令，启动时先处理这些参数。 | [adapter/runtime/backend/cmd/server/maintenance_client.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/cmd/server/maintenance_client.go.in:14) |
| N02 | 启动维护 Unix socket 协程；同时把原来的 ListenAndServe 改成显式 net.Listen 和 Serve。 | [adapter/rules.json](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/rules.json:1) |
| N03 | 原生关机前停止评估并等待任务、流生产者退出；超时后反复新建五秒上下文，没有总等待上限。 | [adapter/runtime/backend/internal/service/eval_host_runtime.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_host_runtime.go.in:187) |
| N04 | 在 Application、服务、管理员 Handler 和 Wire 生成代码中接入评估依赖。 | [adapter/rules.json](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/rules.json:1) |
| N05 | 把真实 Gin router 绑定到评估运行时，新增 HTTP server 包装器和外部请求保护。 | [adapter/runtime/backend/internal/service/eval_host_runtime.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_host_runtime.go.in:83) |
| N06 | 新增管理员评估接口：状态、激活、目标、启动、取消、历史、清空、停用，以及评估页面路由。 | [adapter/runtime/backend/internal/handler/admin_codex_eval.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/handler/admin_codex_eval.go.in:19) |
| N07 | 新增进程内私有调用凭据；修改原生 Key 鉴权及 Google 鉴权，拒绝公开重放和伪造标记。 | [adapter/runtime/backend/internal/pkg/evalroute/evalroute.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/pkg/evalroute/evalroute.go.in:10) |
| N08 | 新增固定过期时间的原生记账 Key，以及创建、所有权验证、持久化、迁移、按条件撤销逻辑。 | [adapter/runtime/backend/internal/service/eval_key_runtime.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_key_runtime.go.in:21) |
| N09 | 修改原生 API Key DTO 映射，隐藏评估记账 Key 的 lookup 值。 | [adapter/runtime/backend/internal/service/eval_identity.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_identity.go.in:36) |
| N10 | 新增指定 group/account 的目标查询、成员关系验证和目标列表 DTO。 | [adapter/runtime/backend/internal/service/eval_target_adapter.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_target_adapter.go.in:11) |
| N11 | 修改 Responses Handler 的入口、并发等待、重新验证、重试及用量流程，锁定指定账户。 | [adapter/runtime/backend/internal/pkg/evalroute/evalroute.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/pkg/evalroute/evalroute.go.in:10) |
| N12 | 修改多个调度入口、sticky、利润控制、continuation 和 fallback 路径，阻止切换到其他账户。 | [adapter/runtime/backend/internal/service/openai_eval_directed.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/openai_eval_directed.go.in:12) |
| N13 | 修改普通转发、passthrough、Chat/Responses 转换及 WebSocket 中的评估取消和部分用量处理。 | [adapter/rules.json](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/rules.json:1) |
| N14 | 新增 Unix socket 维护协议、Linux SO_PEERCRED 对端验证、部署 ID 核对及停用确认。 | [adapter/runtime/backend/internal/service/eval_maintenance.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_maintenance.go.in:17) |
| N15 | 新增保留菜单名冲突检查，发现已有冲突菜单时禁止激活。 | [adapter/runtime/backend/internal/service/eval_menu.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_menu.go.in:11) |
| N16 | 修改侧边栏和 CustomPageView，加入内置评估菜单、重复项处理和页面回退。 | [adapter/runtime/frontend/src/utils/codex-eval-menu.ts](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/frontend/src/utils/codex-eval-menu.ts:1) |
| N17 | 修改 SPA fallback 绕过规则，让评估页面单独处理；修改同源嵌入 URL 和主题传递。 | [adapter/rules.json](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/rules.json:1) |
| N18 | 新增受保护部署 ID 文件，以及版本、源版本、源码摘要等返回字段和构建标记。 | [adapter/runtime/backend/internal/service/eval_runtime_ports.go.in](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/runtime/backend/internal/service/eval_runtime_ports.go.in:8) |
| N19 | 新增原生源码适配、运行时文件复制、前端构建、类型检查及 Linux 服务编译流程。 | [adapter/build.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/build.py:14) |

## 新增评估核心和评估页面

| ID | 当前内容 | 主要源码 |
| --- | --- | --- |
| C01 | 新增独立 internal/eval 核心，使用数字 group/account 和 ModelExecutor，替换 CPA host/C ABI 调用。 | [internal/eval/eval.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/eval.go:13) |
| C02 | 结果新增 run/group/account 身份、请求和实际 model/effort、显示名称及状态。 | [internal/eval/eval.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/eval.go:59) |
| C03 | 新增模型名后缀解析和冲突拒绝；调整 effort 集合，UI 移除原有 max 选项。 | [internal/eval/eval.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/eval.go:223) |
| C04 | 新增统一参数验证；Fingerprint 限制 low，ModelTrace 限制默认 effort，仅 Candy 支持批次。 | [internal/eval/eval.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/eval.go:95) |
| C05 | 新增固定地址的进程内 Responses HTTP 执行器，校验请求目标/body/headers 并拒绝跳转。 | [internal/eval/http.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/http.go:21) |
| C06 | 调整错误和 panic 处理，隐藏原始上游正文/URL；取消后不保存答案。 | [internal/eval/model.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/model.go:32) |
| C07 | 重写为增量 SSE 解析，要求终态、拒绝异常事件，并重建完成项中的答案。 | [internal/eval/sse.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/sse.go:52) |
| C08 | 新增请求、header、SSE 行/事件、输出、总流量和事件数量上限。 | [internal/eval/http.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/http.go:12) |
| C09 | 新增已观察到的 usage 标记，失败尝试也累计实际收到的 token 数。 | [internal/eval/sse.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/sse.go:20) |
| C10 | 新增重试前 body/header 深拷贝、父上下文控制超时和退避；缩窄可重试 HTTP 状态范围。 | [internal/eval/model.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/model.go:47) |
| C11 | 复用 Codex 请求指纹和上下文；安装命名空间改为账户 ID。上下文数据未改。 | [internal/eval/codex.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/codex.go:16) |
| C12 | Candy 增加失败/取消元数据及 error 返回；题目、21 判断和评分未改。 | [internal/eval/candy.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/candy.go:8) |
| C13 | Fingerprint 增加进度、汇总和父取消结果；原探针、算法和阈值未改。 | [internal/eval/fingerprint.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/fingerprint.go:13) |
| C14 | ModelTrace 增加 Done/Total/Valid/Errors、取消/部分结果及汇总；挑战和算法数据未改。 | [internal/eval/modeltrace.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/modeltrace.go:12) |
| C15 | Pelican 增加失败/取消元数据和列表汇总；原题目、HTML 提取规则未改。 | [internal/eval/pelican.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/pelican.go:12) |
| C16 | 新增复制式进度回调和 manager 进度存储，替代 CPA 全局状态写入。 | [internal/eval/eval.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/eval.go:33) |
| C17 | 新增精确目标解析、跨分组按账户排他、去重、busy/unavailable 结果及失败记录。 | [internal/eval/runs.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/runs.go:14) |
| C18 | 新增 Candy 批次持久化、每轮 run ID、完成计数和执行结束后释放账户占用。 | [internal/eval/runs.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/runs.go:141) |
| C19 | 新增持久化的活动任务意图、确认启动前落盘、重启后中断记录恢复。 | [internal/eval/state.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/state.go:14) |
| C20 | 新增历史/活动任务的 schema、身份、选项、状态和计数校验；保留原每类20条历史。 | [internal/eval/state.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/state.go:159) |
| C21 | 新增路径逐层拒绝符号链接、0700目录和唯一0600临时文件的原子存储。 | [internal/eval/state.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/state.go:46) |
| C22 | 新增结果深拷贝、快照、按 run ID 查记录、清空/flush及存储错误状态。 | [internal/eval/state.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/state.go:289) |
| C23 | 新增所有评估类型按 run ID 取消、manager生命周期、关闭准入及等待退出。 | [internal/eval/runs.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/runs.go:638) |
| C24 | 新增 Sub2API 静态评估页面封装、nonce和品牌替换；页脚版本目前写成字面量1。 | [internal/eval/ui.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/ui.go:10) |
| C25 | 新增原生管理员 auth_token 鉴权和旧 UI 请求到新 API 的路径、body、响应转换。 | [internal/eval/web/app.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/app.js:25) |
| C26 | 新增首次激活按钮、ready门槛和 repair_required 时禁止激活。 | [internal/eval/web/ui.html](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/ui.html:33) |
| C27 | 新增分组选择、原生账户/历史/进度到原 UI 数据格式的映射。 | [internal/eval/web/app.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/app.js:31) |
| C28 | 新增所选分组各账户模型目录聚合和缓存，替代原 CPA 配置和模型读取方式。 | [internal/eval/web/catalog.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/catalog.js:3) |
| C29 | 启动前额外请求账户模型目录；结果在发送前被删除，未传给服务端。 | [internal/eval/web/app.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/app.js:249) |
| C30 | 账户开关改成只读，但保留旧 switching、预览及无实际切换的代码。 | [internal/eval/web/components.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/components.js:145) |
| C31 | 把渲染器生成的 style 属性转换为 nonce 样式表里的 class。 | [internal/eval/web/components.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/components.js:3) |
| C32 | 重写 Pelican HTML 的 nonce及内联样式，再放入原 sandbox iframe。 | [internal/eval/web/pelican.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/pelican.js:27) |
| C33 | 更换为 Sub2API 浅/深色主题，移除 CPA 背景、动画及一条样式重置。 | [internal/eval/web/theme.css](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/theme.css:1) |
| C34 | 保留 CPA 订阅计划分级、徽章和筛选逻辑，但原生账户数据不满足其条件，当前不可达。 | [internal/eval/web/components.js](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/internal/eval/web/components.js:101) |

## 安装、构建、打包和历史工具

| ID | 当前内容 | 主要源码 |
| --- | --- | --- |
| T01 | 新增预编译二进制安装入口和固定 v0.2.16 amd64 安装包，安装时不编译。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:26) |
| T02 | 新增本地完整包识别、固定 release HTTPS 下载及有大小/成员限制的解包。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:148) |
| T03 | 新增 Docker 容器列表、强制明确选择、管道下终端输入和不可变容器 ID 重验。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:347) |
| T04 | 新增容器运行状态、账单模式、持久化挂载、UID/GID、可执行路径及路径安全预检。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:399) |
| T05 | 新增临时 Alpine/git 工具镜像获取和缓存，供 nc 探测使用。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:60) |
| T06 | 新增 nc 子容器 HTTP 请求、报文 framing 解析、超时和进程/容器清理。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:69) |
| T07 | 新增进程及 socket 归属查找、监听端口发现、提示回退和健康/setup探测。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:547) |
| T08 | 新增固定 amd64、/app/data 及已知版本冲突限制。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:794) |
| T09 | 新增私有安装目录、schema3安装日志、原始二进制备份及原子写入。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:392) |
| T10 | 新增原子二进制替换、只重启选定容器、替换后检查及失败恢复。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:809) |
| T11 | 新增补丁源码摘要和安装facts，匹配时短路重复安装；不比较已装二进制摘要。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:26) |
| T12 | 新增维护状态核对、撤销记账Key后回滚/卸载，以及可恢复阶段流程。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:813) |
| T13 | 新增完整 stdout/stderr、traceback、实时子进程输出、等待显示和终止处理。 | [install.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/install.py:130) |
| T14 | 新增独立 cleanup.py，列出并按选择清理可归属的旧助手容器/临时工作区。 | [cleanup.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/cleanup.py:14) |
| T15 | cleanup另提供未使用的 Go/Node 工具容器及镜像候选清理。 | [cleanup.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/cleanup.py:121) |
| T16 | 新增8文件白名单打包、固定 dist/build 输出、规范化 tar/gzip 及原子归档替换。 | [scripts/package-sub2api.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/scripts/package-sub2api.py:15) |
| T17 | 新增功能源码白名单、功能摘要及打包facts耦合；仍依赖部分旧工具函数。 | [scripts/universal_payload.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/scripts/universal_payload.py:1) |
| T18 | 新增开发端 v0.2.16 tag/archive 获取、构建profile；当前 CI 使用 Node24。 | [scripts/prebuilt_dev.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/scripts/prebuilt_dev.py:14) |
| T19 | 新增真正的前端/Go builder、数据目录facts、release/source/digest标记及锁文件核对。 | [adapter/build.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/build.py:14) |
| T20 | 新增 AST 修改引擎，匹配声明/作用域/相邻语句，再格式化和写入宿主源码。 | [adapter/main.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/main.go:45) |
| T21 | 新增额外 compiler-export 类型/方法形状门禁，最终 Go 编译前再编译并检查依赖图。 | [adapter/typecheck.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/typecheck.go:15) |
| T22 | 保留开发用 AST规则导出器；当前构建不编译或使用它。 | [adapter/export_dev.go](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/export_dev.go:136) |
| T23 | 新增独立测试工作区、host/core测试覆盖补丁及生产源码字节一致性比较。 | [scripts/prebuilt_test.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/scripts/prebuilt_test.py:13) |
| T24 | 新增本地开发检查入口；UI检查仍读取 CPA 页面，原 race 调用已移除。 | [scripts/test-sub2api.sh](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/scripts/test-sub2api.sh:20) |
| T25 | 改写 CI/release 流程，构建Sub2API、检查、打包、容器生命周期验证和上传归档。 | [.github/workflows/release.yml](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/.github/workflows/release.yml:12) |
| T26 | 保留旧 v0.2.15 兼容清单、拉源码/应用patch/工作区/来源证明工具；只有少数helper仍被当前打包调用。 | [compatibility/sub2api-v0.2.15.json](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/compatibility/sub2api-v0.2.15.json:1) |
| T27 | 保留历史组合patch和production patch；当前构建不应用它们，test-overlay仍使用。 | [patches/sub2api/v0.2.15/production/001-codex-eval-auth.patch](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/patches/sub2api/v0.2.15/production/001-codex-eval-auth.patch:1) |
| T28 | 保留旧 build/export/Wire/pnpm脚本；部分入口已废弃或说明过时。 | [scripts/build-sub2api.sh](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/scripts/build-sub2api.sh:1) |
| T29 | 保留当前不执行的 Dockerfile，却仍把它纳入功能源码摘要。 | [adapter/Dockerfile](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/Dockerfile:1) |
| T30 | 保留旧“两最新版本” universal CI runner；当前工作流不调用，且其source-package假设已过时。 | [tests/installer/universal_ci.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/tests/installer/universal_ci.py:17) |
| T31 | 保留旧安装期源码构建/进度测试；当前测试套件仍选择性引用其中一些用例。 | [tests/installer/current_suite.py](/Users/minami/Code/C/cpa-plugin-codex-candy-eval/tests/installer/current_suite.py:3) |

## 直接改动的32个宿主文件

路径相对于 Sub2API 的 backend；frontend路径保留原前缀。

| 文件 | AST规则数 | 注入函数数 |
| --- | --- | --- |
| `cmd/server/main.go` | 3 | 0 |
| `cmd/server/wire.go` | 2 | 0 |
| `cmd/server/wire_gen.go` | 23 | 0 |
| `frontend/src/components/layout/AppSidebar.vue` | 0 | 0 |
| `frontend/src/utils/embedded-url.ts` | 0 | 0 |
| `frontend/src/views/user/CustomPageView.vue` | 0 | 0 |
| `internal/handler/dto/mappers.go` | 1 | 0 |
| `internal/handler/gateway_handler_responses.go` | 1 | 0 |
| `internal/handler/handler.go` | 1 | 0 |
| `internal/handler/openai_gateway_handler.go` | 6 | 0 |
| `internal/handler/wire.go` | 3 | 0 |
| `internal/repository/api_key_repo.go` | 0 | 4 |
| `internal/server/http.go` | 2 | 1 |
| `internal/server/middleware/api_key_auth.go` | 4 | 0 |
| `internal/server/middleware/api_key_auth_google.go` | 1 | 0 |
| `internal/server/router.go` | 1 | 0 |
| `internal/server/routes/admin.go` | 1 | 0 |
| `internal/service/api_key_service.go` | 1 | 2 |
| `internal/service/gateway_forward_as_chat_completions.go` | 1 | 0 |
| `internal/service/gateway_forward_as_responses.go` | 1 | 0 |
| `internal/service/gateway_usage_billing.go` | 0 | 2 |
| `internal/service/openai_account_scheduler.go` | 2 | 0 |
| `internal/service/openai_gateway_cc_pipeline.go` | 1 | 0 |
| `internal/service/openai_gateway_forward.go` | 3 | 1 |
| `internal/service/openai_gateway_passthrough.go` | 3 | 3 |
| `internal/service/openai_gateway_scheduling.go` | 2 | 0 |
| `internal/service/openai_profit_control.go` | 2 | 0 |
| `internal/service/openai_sticky_compat.go` | 1 | 0 |
| `internal/service/openai_ws_forwarder_support.go` | 1 | 0 |
| `internal/service/openai_ws_forwarder_v2.go` | 4 | 0 |
| `internal/service/wire.go` | 1 | 0 |
| `internal/web/embed_on.go` | 1 | 0 |

## 16份新增运行时模板

- `adapter/runtime/backend/cmd/server/maintenance_client.go.in`
- `adapter/runtime/backend/internal/handler/admin_codex_eval.go.in`
- `adapter/runtime/backend/internal/pkg/evalroute/evalroute.go.in`
- `adapter/runtime/backend/internal/server/middleware/api_key_auth_eval.go.in`
- `adapter/runtime/backend/internal/service/eval_host_runtime.go.in`
- `adapter/runtime/backend/internal/service/eval_http_transport.go.in`
- `adapter/runtime/backend/internal/service/eval_identity.go.in`
- `adapter/runtime/backend/internal/service/eval_key_runtime.go.in`
- `adapter/runtime/backend/internal/service/eval_maintenance.go.in`
- `adapter/runtime/backend/internal/service/eval_maintenance_linux.go.in`
- `adapter/runtime/backend/internal/service/eval_maintenance_other.go.in`
- `adapter/runtime/backend/internal/service/eval_menu.go.in`
- `adapter/runtime/backend/internal/service/eval_runtime_ports.go.in`
- `adapter/runtime/backend/internal/service/eval_target_adapter.go.in`
- `adapter/runtime/backend/internal/service/openai_eval_directed.go.in`
- `adapter/runtime/frontend/src/utils/codex-eval-menu.ts`

## 逐条原生修改的原始定义

下面完整保留当前73条规则及新增声明/import；没有省略源码替换内容。

```json
{
  "Rules": [
    {
      "File": "cmd/server/main.go",
      "Decl": ".main",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if len(os.Args) > 1 && os.Args[1] == \"--codex-eval-maintenance-capabilities\" {\n os.Stdout.WriteString(\"codex-eval-maintenance-v1\\n\")\n os.Exit(0)\n}\nif len(os.Args) > 1 && os.Args[1] == \"--codex-eval-maintenance\" {\n os.Exit(runMaintenanceClient(os.Args[2:]))\n}",
      "Before": "",
      "After": "logger.InitBootstrap()",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/main.go",
      "Decl": ".runMainServer",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "go func() {\n\tif err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {\n\t\tlog.Fatalf(\"Failed to start server: %v\", err)\n\t}\n}()",
      "New": "listener, err := net.Listen(\"tcp\", app.Server.Addr)\nif err != nil {\n\tlog.Fatalf(\"Failed to bind server: %v\", err)\n}\nmaintenanceCtx, stopMaintenance := context.WithCancel(context.Background())\nmaintenanceDone := make(chan struct{})\ndefer func() { stopMaintenance(); <-maintenanceDone }()\ngo func() {\n\tdefer close(maintenanceDone)\n\tif err := app.CodexEval.ServeMaintenance(maintenanceCtx); err != nil {\n\t\tlog.Print(\"Evaluation maintenance is unavailable\")\n\t}\n}()\ngo func() {\n\tif err := app.Server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {\n\t\tlog.Fatalf(\"Failed to start server: %v\", err)\n\t}\n}()",
      "Before": "if app.PromptAudit != nil {\n\tif err := app.PromptAudit.Start(context.Background()); err != nil {\n\n\t\tlog.Printf(\"Prompt Audit started in degraded state: %v\", err)\n\t}\n}",
      "After": "log.Printf(\"Server started on %s\", app.Server.Addr)",
      "Scope": "",
      "Comma": false,
      "Count": 1
    },
    {
      "File": "cmd/server/main.go",
      "Decl": ".runMainServer",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "for {\n\tdrainErr := app.CodexEval.Shutdown(ctx)\n\tif !errors.Is(drainErr, context.DeadlineExceeded) && !errors.Is(drainErr, context.Canceled) {\n\t\tif drainErr != nil {\n\t\t\tlog.Print(\"Evaluation drain finished with a storage error\")\n\t\t}\n\t\tbreak\n\t}\n\tcancel()\n\tctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)\n}\ncancel()\nctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)\ndefer cancel()",
      "Before": "defer cancel()",
      "After": "if err := app.Server.Shutdown(ctx); err != nil {\n\tlog.Printf(\"Server forced to shutdown: %v\", err)\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire.go",
      "Decl": "type:Application",
      "Kind": "*ast.FieldList",
      "Field": "List",
      "Old": "",
      "New": "CodexEval *service.EvalHostRuntime",
      "Before": "",
      "After": "Server *http.Server",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.CallExpr",
      "Field": "Args",
      "Old": "",
      "New": "\"CodexEval\"",
      "Before": "new(Application)",
      "After": "\"Server\"",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "apiKeyRepository := repository.NewAPIKeyRepository(client, db)",
      "Before": "if err != nil {\n\treturn nil, err\n}",
      "After": "userRepository := repository.NewUserRepository(client, db)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "redeemCodeRepository := repository.NewRedeemCodeRepository(client)\nredisClient := repository.ProvideRedis(configConfig)\nrefreshTokenCache := repository.NewRefreshTokenCache(redisClient)\nsettingRepository := repository.NewSettingRepository(client)",
      "New": "",
      "Before": "userRepository := repository.NewUserRepository(client, db)",
      "After": "groupRepository := repository.NewGroupRepository(client, db)",
      "Scope": "",
      "Comma": false,
      "Count": 4
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "proxyRepository := repository.NewProxyRepository(client, db)\nsettingService := service.ProvideSettingService(settingRepository, groupRepository, proxyRepository, configConfig)\nemailCache := repository.NewEmailCache(redisClient)\nemailService := service.NewEmailService(settingRepository, emailCache)\nturnstileVerifier := repository.NewTurnstileVerifier()\nturnstileService := service.NewTurnstileService(settingService, turnstileVerifier)\ntencentCaptchaVerifier := repository.NewTencentCaptchaVerifier()\ntencentCaptchaService := service.NewTencentCaptchaService(settingService, tencentCaptchaVerifier)\naliyunCaptchaVerifier := repository.NewAliyunCaptchaVerifier()\naliyunCaptchaService := service.NewAliyunCaptchaService(settingService, aliyunCaptchaVerifier)\nemailQueueService := service.ProvideEmailQueueService(emailService)\npromoCodeRepository := repository.NewPromoCodeRepository(client)\nbillingCache := repository.NewBillingCache(redisClient)",
      "New": "",
      "Before": "groupRepository := repository.NewGroupRepository(client, db)",
      "After": "userSubscriptionRepository := repository.NewUserSubscriptionRepository(client)",
      "Scope": "",
      "Comma": false,
      "Count": 13
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "apiKeyRepository := repository.NewAPIKeyRepository(client, db)\nuserRPMCache := repository.NewUserRPMCache(redisClient)",
      "New": "",
      "Before": "userSubscriptionRepository := repository.NewUserSubscriptionRepository(client)",
      "After": "userGroupRateRepository := repository.NewUserGroupRateRepository(db)",
      "Scope": "",
      "Comma": false,
      "Count": 2
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "redisClient := repository.ProvideRedis(configConfig)\napiKeyCache := repository.NewAPIKeyCache(redisClient)\nbillingCache := repository.NewBillingCache(redisClient)\nuserRPMCache := repository.NewUserRPMCache(redisClient)",
      "Before": "userGroupRateRepository := repository.NewUserGroupRateRepository(db)",
      "After": "userPlatformQuotaRepository := repository.NewUserPlatformQuotaRepository(client)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "apiKeyCache := repository.NewAPIKeyCache(redisClient)",
      "New": "",
      "Before": "billingCacheService := service.ProvideBillingCacheService(billingCache, userRepository, userSubscriptionRepository, apiKeyRepository, userRPMCache, userGroupRateRepository, configConfig, serviceUserPlatformQuotaRepository)",
      "After": "concurrencyCache := repository.ProvideConcurrencyCache(redisClient, configConfig)",
      "Scope": "",
      "Comma": false,
      "Count": 1
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "adminGroupRepository := repository.NewAdminGroupRepository(client, db)\nadminAccountRepository := repository.NewAdminAccountRepository(client, db, schedulerCache)\nproxyRepository := repository.NewProxyRepository(client, db)\nredeemCodeRepository := repository.NewRedeemCodeRepository(client)\nproxyExitInfoProber := repository.NewProxyExitInfoProber(configConfig)\nproxyLatencyCache := repository.NewProxyLatencyCache(redisClient)",
      "Before": "apiKeyService := service.ProvideAPIKeyService(apiKeyRepository, userRepository, groupRepository, userSubscriptionRepository, userGroupRateRepository, apiKeyCache, configConfig, billingCacheService, concurrencyService)",
      "After": "apiKeyAuthCacheInvalidator := service.ProvideAPIKeyAuthCacheInvalidator(apiKeyService)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "promoService := service.NewPromoService(promoCodeRepository, userRepository, billingCacheService, client, apiKeyAuthCacheInvalidator)",
      "New": "settingRepository := repository.NewSettingRepository(client)\nsettingService := service.ProvideSettingService(settingRepository, groupRepository, proxyRepository, configConfig)",
      "Before": "apiKeyAuthCacheInvalidator := service.ProvideAPIKeyAuthCacheInvalidator(apiKeyService)",
      "After": "subscriptionService := service.NewSubscriptionService(groupRepository, userSubscriptionRepository, billingCacheService, client, configConfig)",
      "Scope": "",
      "Comma": false,
      "Count": 1
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "affiliateRepository := repository.NewAffiliateRepository(client, db)\naffiliateService := service.NewAffiliateService(affiliateRepository, settingService, apiKeyAuthCacheInvalidator, billingCacheService)\nauthService := service.ProvideAuthService(client, userRepository, redeemCodeRepository, refreshTokenCache, configConfig, settingService, emailService, turnstileService, tencentCaptchaService, aliyunCaptchaService, emailQueueService, promoService, subscriptionService, affiliateService, serviceUserPlatformQuotaRepository)\nuserService := service.NewUserService(userRepository, settingRepository, apiKeyAuthCacheInvalidator, billingCache)\nredeemCache := repository.NewRedeemCache(redisClient)\nredeemService := service.NewRedeemService(redeemCodeRepository, userRepository, subscriptionService, redeemCache, billingCacheService, client, apiKeyAuthCacheInvalidator, affiliateService)\nsecretEncryptor, err := repository.NewAESEncryptor(configConfig)\nif err != nil {\n\treturn nil, err\n}\ntotpCache := repository.NewTotpCache(redisClient)\ntotpService := service.NewTotpService(userRepository, secretEncryptor, totpCache, settingService, emailService, emailQueueService)\nuserAttributeDefinitionRepository := repository.NewUserAttributeDefinitionRepository(client)\nuserAttributeValueRepository := repository.NewUserAttributeValueRepository(client)\nuserAttributeService := service.NewUserAttributeService(userAttributeDefinitionRepository, userAttributeValueRepository)\nauthHandler := handler.NewAuthHandler(configConfig, authService, userService, settingService, promoService, redeemService, totpService, userAttributeService)\nuserHandler := handler.NewUserHandler(userService, authService, emailService, emailCache, affiliateService, serviceUserPlatformQuotaRepository)\napiKeyHandler := handler.NewAPIKeyHandler(apiKeyService)",
      "New": "privacyClientFactory := providePrivacyClientFactory()",
      "Before": "subscriptionService := service.NewSubscriptionService(groupRepository, userSubscriptionRepository, billingCacheService, client, configConfig)",
      "After": "usageLogRepository := repository.NewUsageLogRepository(client, db)",
      "Scope": "",
      "Comma": false,
      "Count": 16
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "usageService := service.NewUsageService(usageLogRepository, userRepository, client, apiKeyAuthCacheInvalidator)\nopsRepository := repository.NewOpsRepository(db)",
      "New": "",
      "Before": "usageLogRepository := repository.NewUsageLogRepository(client, db)",
      "After": "usageBillingRepository := repository.NewUsageBillingRepository(client, db)",
      "Scope": "",
      "Comma": false,
      "Count": 2
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "secretEncryptor, err := repository.NewAESEncryptor(configConfig)\nif err != nil {\n\treturn nil, err\n}",
      "Before": "httpUpstream := repository.NewHTTPUpstream(configConfig)",
      "After": "leaderLockCache := repository.NewLeaderLockCache(redisClient)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "identityCache := repository.NewIdentityCache(redisClient)\nidentityService := service.NewIdentityService(identityCache)",
      "New": "",
      "Before": "rateLimitService := service.ProvideRateLimitService(accountRepository, usageLogRepository, configConfig, geminiQuotaService, tempUnschedCache, timeoutCounterCache, openAI403CounterCache, settingService, compositeTokenCacheInvalidator, ollamaCloudUsageService)",
      "After": "timingWheelService, err := service.ProvideTimingWheelService()",
      "Scope": "",
      "Comma": false,
      "Count": 2
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "openAIOAuthClient := repository.NewOpenAIOAuthClient()\nopenAIOAuthService := service.ProvideOpenAIOAuthService(proxyRepository, openAIOAuthClient, privacyClientFactory)\noAuthRefreshAPI := service.ProvideOAuthRefreshAPI(accountRepository, geminiTokenCache)\nopenAITokenProvider := service.ProvideOpenAITokenProvider(accountRepository, geminiTokenCache, openAIOAuthService, oAuthRefreshAPI)\ngrokOAuthClient := repository.NewGrokOAuthClient()\ngrokOAuthService := service.ProvideGrokOAuthService(proxyRepository, grokOAuthClient, configConfig, redisClient)\ngrokTokenProvider := service.ProvideGrokTokenProvider(accountRepository, geminiTokenCache, grokOAuthService, oAuthRefreshAPI, tempUnschedCache)\nchannelRepository := repository.NewChannelRepository(db)\nchannelCachePubSub := repository.NewChannelCache(redisClient)\nchannelService := service.NewChannelService(channelRepository, groupRepository, apiKeyAuthCacheInvalidator, pricingService, channelCachePubSub)\nmodelPricingResolver := service.NewModelPricingResolver(channelService, billingService)\nemailCache := repository.NewEmailCache(redisClient)\nemailService := service.NewEmailService(settingRepository, emailCache)\nnotificationEmailService := service.NewNotificationEmailService(settingRepository, emailService)\nbalanceNotifyService := service.ProvideBalanceNotifyService(emailService, settingRepository, accountRepository, notificationEmailService)\nopenAIGatewayService := service.NewOpenAIGatewayService(accountRepository, usageLogRepository, usageBillingRepository, userRepository, userSubscriptionRepository, userGroupRateRepository, gatewayCache, configConfig, schedulerSnapshotService, concurrencyService, billingService, rateLimitService, billingCacheService, httpUpstream, deferredService, openAITokenProvider, grokTokenProvider, modelPricingResolver, channelService, balanceNotifyService, settingService, serviceUserPlatformQuotaRepository)\naffiliateRepository := repository.NewAffiliateRepository(client, db)\naffiliateService := service.NewAffiliateService(affiliateRepository, settingService, apiKeyAuthCacheInvalidator, billingCacheService)\ncompositeModelRouteRepository := repository.NewCompositeModelRouteRepository(client)\ncompositeRouteResolver := service.NewCompositeRouteResolver(compositeModelRouteRepository)\nadminService := service.NewAdminService(configConfig, userRepository, adminGroupRepository, adminAccountRepository, proxyRepository, apiKeyRepository, redeemCodeRepository, userGroupRateRepository, userRPMCache, billingCacheService, proxyExitInfoProber, proxyLatencyCache, apiKeyAuthCacheInvalidator, client, settingService, subscriptionService, userSubscriptionRepository, privacyClientFactory, openAIGatewayService, affiliateService, compositeModelRouteRepository, compositeRouteResolver, channelService)\nserviceBuildInfo := provideServiceBuildInfo(buildInfo)\nevalHostRuntime := service.NewEvalHostRuntime(configConfig, apiKeyService, adminService, groupRepository, settingRepository, serviceBuildInfo)\nrefreshTokenCache := repository.NewRefreshTokenCache(redisClient)\nturnstileVerifier := repository.NewTurnstileVerifier()\nturnstileService := service.NewTurnstileService(settingService, turnstileVerifier)\ntencentCaptchaVerifier := repository.NewTencentCaptchaVerifier()\ntencentCaptchaService := service.NewTencentCaptchaService(settingService, tencentCaptchaVerifier)\naliyunCaptchaVerifier := repository.NewAliyunCaptchaVerifier()\naliyunCaptchaService := service.NewAliyunCaptchaService(settingService, aliyunCaptchaVerifier)\nemailQueueService := service.ProvideEmailQueueService(emailService)\npromoCodeRepository := repository.NewPromoCodeRepository(client)\npromoService := service.NewPromoService(promoCodeRepository, userRepository, billingCacheService, client, apiKeyAuthCacheInvalidator)\nauthService := service.ProvideAuthService(client, userRepository, redeemCodeRepository, refreshTokenCache, configConfig, settingService, emailService, turnstileService, tencentCaptchaService, aliyunCaptchaService, emailQueueService, promoService, subscriptionService, affiliateService, serviceUserPlatformQuotaRepository)\nuserService := service.NewUserService(userRepository, settingRepository, apiKeyAuthCacheInvalidator, billingCache)\nredeemCache := repository.NewRedeemCache(redisClient)\nredeemService := service.NewRedeemService(redeemCodeRepository, userRepository, subscriptionService, redeemCache, billingCacheService, client, apiKeyAuthCacheInvalidator, affiliateService)\ntotpCache := repository.NewTotpCache(redisClient)\ntotpService := service.NewTotpService(userRepository, secretEncryptor, totpCache, settingService, emailService, emailQueueService)\nuserAttributeDefinitionRepository := repository.NewUserAttributeDefinitionRepository(client)\nuserAttributeValueRepository := repository.NewUserAttributeValueRepository(client)\nuserAttributeService := service.NewUserAttributeService(userAttributeDefinitionRepository, userAttributeValueRepository)\nauthHandler := handler.NewAuthHandler(configConfig, authService, userService, settingService, promoService, redeemService, totpService, userAttributeService)\nuserHandler := handler.NewUserHandler(userService, authService, emailService, emailCache, affiliateService, serviceUserPlatformQuotaRepository)\napiKeyHandler := handler.NewAPIKeyHandler(apiKeyService)\nusageService := service.NewUsageService(usageLogRepository, userRepository, client, apiKeyAuthCacheInvalidator)\nopsRepository := repository.NewOpsRepository(db)\nidentityCache := repository.NewIdentityCache(redisClient)\nidentityService := service.NewIdentityService(identityCache)",
      "Before": "deferredService := service.ProvideDeferredService(accountRepository, timingWheelService)",
      "After": "claudeOAuthClient := repository.NewClaudeOAuthClient()",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "oAuthRefreshAPI := service.ProvideOAuthRefreshAPI(accountRepository, geminiTokenCache)",
      "New": "",
      "Before": "oAuthService := service.NewOAuthService(proxyRepository, claudeOAuthClient)",
      "After": "claudeTokenProvider := service.ProvideClaudeTokenProvider(accountRepository, geminiTokenCache, oAuthService, oAuthRefreshAPI)",
      "Scope": "",
      "Comma": false,
      "Count": 1
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "channelRepository := repository.NewChannelRepository(db)\nchannelCachePubSub := repository.NewChannelCache(redisClient)\nchannelService := service.NewChannelService(channelRepository, groupRepository, apiKeyAuthCacheInvalidator, pricingService, channelCachePubSub)\nmodelPricingResolver := service.NewModelPricingResolver(channelService, billingService)\ncompositeModelRouteRepository := repository.NewCompositeModelRouteRepository(client)\ncompositeRouteResolver := service.NewCompositeRouteResolver(compositeModelRouteRepository)\nnotificationEmailService := service.NewNotificationEmailService(settingRepository, emailService)\nbalanceNotifyService := service.ProvideBalanceNotifyService(emailService, settingRepository, accountRepository, notificationEmailService)",
      "New": "",
      "Before": "tlsFingerprintProfileService := service.NewTLSFingerprintProfileService(tlsFingerprintProfileRepository, tlsFingerprintProfileCache)",
      "After": "gatewayService := service.NewGatewayService(accountRepository, groupRepository, usageLogRepository, usageBillingRepository, userRepository, userSubscriptionRepository, userGroupRateRepository, gatewayCache, configConfig, schedulerSnapshotService, concurrencyService, billingService, rateLimitService, billingCacheService, identityService, httpUpstream, deferredService, claudeTokenProvider, sessionLimitCache, rpmCache, digestSessionStore, settingService, tlsFingerprintProfileService, channelService, modelPricingResolver, compositeRouteResolver, balanceNotifyService, serviceUserPlatformQuotaRepository)",
      "Scope": "",
      "Comma": false,
      "Count": 8
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "openAIOAuthClient := repository.NewOpenAIOAuthClient()\nprivacyClientFactory := providePrivacyClientFactory()\nopenAIOAuthService := service.ProvideOpenAIOAuthService(proxyRepository, openAIOAuthClient, privacyClientFactory)\nopenAITokenProvider := service.ProvideOpenAITokenProvider(accountRepository, geminiTokenCache, openAIOAuthService, oAuthRefreshAPI)\ngrokOAuthClient := repository.NewGrokOAuthClient()\ngrokOAuthService := service.ProvideGrokOAuthService(proxyRepository, grokOAuthClient, configConfig, redisClient)\ngrokTokenProvider := service.ProvideGrokTokenProvider(accountRepository, geminiTokenCache, grokOAuthService, oAuthRefreshAPI, tempUnschedCache)\nopenAIGatewayService := service.NewOpenAIGatewayService(accountRepository, usageLogRepository, usageBillingRepository, userRepository, userSubscriptionRepository, userGroupRateRepository, gatewayCache, configConfig, schedulerSnapshotService, concurrencyService, billingService, rateLimitService, billingCacheService, httpUpstream, deferredService, openAITokenProvider, grokTokenProvider, modelPricingResolver, channelService, balanceNotifyService, settingService, serviceUserPlatformQuotaRepository)",
      "New": "",
      "Before": "gatewayService := service.NewGatewayService(accountRepository, groupRepository, usageLogRepository, usageBillingRepository, userRepository, userSubscriptionRepository, userGroupRateRepository, gatewayCache, configConfig, schedulerSnapshotService, concurrencyService, billingService, rateLimitService, billingCacheService, identityService, httpUpstream, deferredService, claudeTokenProvider, sessionLimitCache, rpmCache, digestSessionStore, settingService, tlsFingerprintProfileService, channelService, modelPricingResolver, compositeRouteResolver, balanceNotifyService, serviceUserPlatformQuotaRepository)",
      "After": "geminiOAuthClient := repository.NewGeminiOAuthClient(configConfig)",
      "Scope": "",
      "Comma": false,
      "Count": 8
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "adminCodexEvalHandler := handler.NewAdminCodexEvalHandler(evalHostRuntime)",
      "Before": "channelMonitorV2Handler := handler.NewChannelMonitorV2Handler(channelMonitorV2Service, apiKeyService)",
      "After": "dashboardAggregationRepository := repository.NewDashboardAggregationRepository(db)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "adminGroupRepository := repository.NewAdminGroupRepository(client, db)\nadminAccountRepository := repository.NewAdminAccountRepository(client, db, schedulerCache)\nproxyExitInfoProber := repository.NewProxyExitInfoProber(configConfig)\nproxyLatencyCache := repository.NewProxyLatencyCache(redisClient)\nadminService := service.NewAdminService(configConfig, userRepository, adminGroupRepository, adminAccountRepository, proxyRepository, apiKeyRepository, redeemCodeRepository, userGroupRateRepository, userRPMCache, billingCacheService, proxyExitInfoProber, proxyLatencyCache, apiKeyAuthCacheInvalidator, client, settingService, subscriptionService, userSubscriptionRepository, privacyClientFactory, openAIGatewayService, affiliateService, compositeModelRouteRepository, compositeRouteResolver, channelService)",
      "New": "",
      "Before": "dashboardHandler := admin.NewDashboardHandler(dashboardService, dashboardAggregationService)",
      "After": "adminUserHandler := admin.NewUserHandler(adminService, concurrencyService, serviceUserPlatformQuotaRepository, billingCache, totpService, userService, settingService)",
      "Scope": "",
      "Comma": false,
      "Count": 5
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "serviceBuildInfo := provideServiceBuildInfo(buildInfo)",
      "New": "",
      "Before": "gitHubReleaseClient := repository.ProvideGitHubReleaseClient(configConfig)",
      "After": "updateService := service.ProvideUpdateService(updateCache, gitHubReleaseClient, serviceBuildInfo)",
      "Scope": "",
      "Comma": false,
      "Count": 1
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.CallExpr",
      "Field": "Args",
      "Old": "",
      "New": "adminCodexEvalHandler",
      "Before": "",
      "After": "dashboardHandler",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "engine := server.ProvideRouter(configConfig, handlers, jwtAuthMiddleware, optionalJWTAuthMiddleware, adminAuthMiddleware, apiKeyAuthMiddleware, auditLogMiddleware, stepUpAuthMiddleware, apiKeyService, subscriptionService, opsService, settingService, compositeRouteResolver, redisClient)\nhttpServer := server.ProvideHTTPServer(configConfig, engine)",
      "New": "engine := server.ProvideRouter(configConfig, handlers, jwtAuthMiddleware, optionalJWTAuthMiddleware, adminAuthMiddleware, apiKeyAuthMiddleware, auditLogMiddleware, stepUpAuthMiddleware, apiKeyService, subscriptionService, opsService, settingService, compositeRouteResolver, redisClient, evalHostRuntime)\nhttpServer := server.ProvideEvalHTTPServer(configConfig, engine, evalHostRuntime)",
      "Before": "stepUpAuthMiddleware := middleware.NewStepUpAuthMiddleware(totpService, userService, settingService)",
      "After": "opsMetricsCollector := service.ProvideOpsMetricsCollector(opsRepository, settingRepository, accountRepository, concurrencyService, db, redisClient, configConfig)",
      "Scope": "",
      "Comma": false,
      "Count": 2
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": ".initializeApplication",
      "Kind": "*ast.CompositeLit",
      "Field": "Elts",
      "Old": "",
      "New": "CodexEval: evalHostRuntime",
      "Before": "",
      "After": "Server: httpServer",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "cmd/server/wire_gen.go",
      "Decl": "type:Application",
      "Kind": "*ast.FieldList",
      "Field": "List",
      "Old": "",
      "New": "CodexEval *service.EvalHostRuntime",
      "Before": "",
      "After": "Server *http.Server",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/dto/mappers.go",
      "Decl": ".APIKeyFromService",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if service.IsEvalOwnedAPIKeyID(k.ID) {\n\tout.Key = \"\"\n}",
      "Before": "out := &APIKey{\n\tID:\t\t\tk.ID,\n\tUserID:\t\t\tk.UserID,\n\tKey:\t\t\tk.Key,\n\tName:\t\t\tk.Name,\n\tGroupID:\t\tk.GroupID,\n\tStatus:\t\t\tk.Status,\n\tIPWhitelist:\t\tk.IPWhitelist,\n\tIPBlacklist:\t\tk.IPBlacklist,\n\tLastUsedAt:\t\tk.LastUsedAt,\n\tLastUsedIP:\t\tk.LastUsedIP,\n\tQuota:\t\t\tk.Quota,\n\tQuotaUsed:\t\tk.QuotaUsed,\n\tExpiresAt:\t\tk.ExpiresAt,\n\tCreatedAt:\t\tk.CreatedAt,\n\tUpdatedAt:\t\tk.UpdatedAt,\n\tCurrentConcurrency:\tk.CurrentConcurrency,\n\tRateLimit5h:\t\tk.RateLimit5h,\n\tRateLimit1d:\t\tk.RateLimit1d,\n\tRateLimit7d:\t\tk.RateLimit7d,\n\tUsage5h:\t\tk.EffectiveUsage5h(),\n\tUsage1d:\t\tk.EffectiveUsage1d(),\n\tUsage7d:\t\tk.EffectiveUsage7d(),\n\tWindow5hStart:\t\tk.Window5hStart,\n\tWindow1dStart:\t\tk.Window1dStart,\n\tWindow7dStart:\t\tk.Window7dStart,\n\tUser:\t\t\tUserFromServiceShallow(k.User),\n\tGroup:\t\t\tGroupFromServiceShallow(k.Group),\n}",
      "After": "if k.Window5hStart != nil && !service.IsWindowExpired(k.Window5hStart, service.RateLimitWindow5h) {\n\tt := k.Window5hStart.Add(service.RateLimitWindow5h)\n\tout.Reset5hAt = &t\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/gateway_handler_responses.go",
      "Decl": "*GatewayHandler.Responses",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(requestCtx); directed {\n\th.responsesErrorResponse(c, http.StatusForbidden, \"invalid_request_error\", \"Directed evaluation routing is unsupported on this endpoint\")\n\treturn\n}",
      "Before": "sessionHash := h.gatewayService.GenerateSessionHash(parsedReq)",
      "After": "fs := NewFailoverState(h.maxAccountSwitches, false)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/handler.go",
      "Decl": "type:AdminHandlers",
      "Kind": "*ast.FieldList",
      "Field": "List",
      "Old": "",
      "New": "CodexEval *AdminCodexEvalHandler",
      "Before": "",
      "After": "Dashboard *admin.DashboardHandler",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/openai_gateway_handler.go",
      "Decl": "*OpenAIGatewayHandler.Responses",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "directedTarget, directedRequest := evalroute.FromContext(c.Request.Context())",
      "Before": "requireCompact := legacyCompact",
      "After": "maxAccountSwitches := h.maxAccountSwitches",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/openai_gateway_handler.go",
      "Decl": "*OpenAIGatewayHandler.Responses",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if directedRequest {\n\th.handleFailoverExhausted(c, lastFailoverErr, streamStarted)\n\treturn\n}",
      "Before": "reqLog.Debug(\"openai.account_skipped_http_continuation_unsupported\",\n\tzap.Int64(\"account_id\", account.ID),\n\tzap.String(\"account_type\", account.Type),\n)",
      "After": "continue",
      "Scope": "/for /if ;previousResponseID != \"\" && requestPlatform == service.PlatformOpenAI && !account.IsOpenAIApiKey()/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/openai_gateway_handler.go",
      "Decl": "*OpenAIGatewayHandler.Responses",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if directedRequest {\n\th.handleOpenAIProfitVetoExhausted(c, streamStarted, reqLog, profitVetoCount+1)\n\treturn\n}",
      "Before": "",
      "After": "if !recordOpenAIProfitVeto(failedAccountIDs, account.ID, &profitVetoCount) {\n\th.handleOpenAIProfitVetoExhausted(c, streamStarted, reqLog, profitVetoCount)\n\treturn\n}",
      "Scope": "/for /if ;slotResult == openAISlotAcquireProfitVetoed/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/openai_gateway_handler.go",
      "Decl": "*OpenAIGatewayHandler.Responses",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if directedRequest {\n\tadmissionCtx := service.ContextWithSelectionProfitGate(c.Request.Context(), selection)\n\tlatest, recheckErr := h.gatewayService.RecheckDirectedEvaluationAccount(\n\t\tadmissionCtx, apiKey.GroupID, directedTarget,\n\t\tforwardModel, requireCompact,\n\t\tservice.OpenAIUpstreamTransportAny, requiredCapability,\n\t)\n\tif recheckErr != nil || latest == nil {\n\t\tif accountReleaseFunc != nil {\n\t\t\taccountReleaseFunc()\n\t\t\taccountReleaseFunc = nil\n\t\t}\n\t\treqLog.Warn(\"openai.eval_directed_post_slot_recheck_failed\",\n\t\t\tzap.Int64(\"account_id\", directedTarget.AccountID),\n\t\t\tzap.Error(recheckErr),\n\t\t)\n\t\th.handleStreamingAwareError(c, http.StatusServiceUnavailable, \"api_error\", \"Evaluation target account is unavailable\", streamStarted)\n\t\treturn\n\t}\n\taccount = latest\n\tselection.Account = latest\n}",
      "Before": "if slotResult != openAISlotAcquireOK {\n\treturn\n}",
      "After": "service.SetOpsLatencyMs(c, service.OpsRoutingLatencyMsKey, time.Since(routingStart).Milliseconds())",
      "Scope": "/for ",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/openai_gateway_handler.go",
      "Decl": "*OpenAIGatewayHandler.Responses",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "submitResponsesUsage(result)",
      "Before": "",
      "After": "if failoverClientGone(c) {\n\treqLog.Info(\"openai.failover_aborted_client_disconnected\",\n\t\tzap.Int64(\"account_id\", account.ID),\n\t\tzap.Int(\"upstream_status\", failoverErr.StatusCode),\n\t)\n\treturn\n}",
      "Scope": "/for /if ;err != nil/then/if ;result != nil && result.ImageCount > 0/else/if ;errors.As(err, &failoverErr)/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/openai_gateway_handler.go",
      "Decl": "*OpenAIGatewayHandler.Responses",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if directedRequest {\n\th.handleFailoverExhausted(c, failoverErr, streamStarted)\n\treturn\n}",
      "Before": "if failoverErr.RetryableOnSameAccount {\n\tretryLimit := effectiveSameAccountRetryLimit(failoverErr, account)\n\tif sameAccountRetryAllowed(failoverErr, sameAccountRetryCount[account.ID], retryLimit) {\n\t\tsameAccountRetryCount[account.ID]++\n\t\tretryDelay := sameAccountRetryDelayFor(failoverErr, sameAccountRetryCount[account.ID])\n\t\treqLog.Warn(\"openai.pool_mode_same_account_retry\",\n\t\t\tzap.Int64(\"account_id\", account.ID),\n\t\t\tzap.Int(\"upstream_status\", failoverErr.StatusCode),\n\t\t\tzap.Int(\"retry_limit\", retryLimit),\n\t\t\tzap.Int(\"retry_count\", sameAccountRetryCount[account.ID]),\n\t\t\tzap.Duration(\"retry_delay\", retryDelay),\n\t\t)\n\t\tselect {\n\t\tcase <-c.Request.Context().Done():\n\t\t\treturn\n\t\tcase <-time.After(retryDelay):\n\t\t}\n\t\tcontinue\n\t}\n}",
      "After": "h.gatewayService.RecordOpenAIAccountSwitch()",
      "Scope": "/for /if ;err != nil/then/if ;result != nil && result.ImageCount > 0/else/if ;errors.As(err, &failoverErr)/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/wire.go",
      "Decl": ".ProvideAdminHandlers",
      "Kind": "*ast.FieldList",
      "Field": "List",
      "Old": "",
      "New": "codexEval *AdminCodexEvalHandler",
      "Before": "",
      "After": "dashboardHandler *admin.DashboardHandler",
      "Scope": "",
      "Comma": true,
      "Count": 0
    },
    {
      "File": "internal/handler/wire.go",
      "Decl": ".ProvideAdminHandlers",
      "Kind": "*ast.CompositeLit",
      "Field": "Elts",
      "Old": "",
      "New": "CodexEval: codexEval",
      "Before": "",
      "After": "Dashboard: dashboardHandler",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/handler/wire.go",
      "Decl": "var:ProviderSet",
      "Kind": "*ast.CallExpr",
      "Field": "Args",
      "Old": "",
      "New": "NewAdminCodexEvalHandler",
      "Before": "",
      "After": "NewAuthHandler",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/server/http.go",
      "Decl": "var:ProviderSet",
      "Kind": "*ast.Ident",
      "Field": "",
      "Old": "ProvideHTTPServer",
      "New": "ProvideEvalHTTPServer",
      "Before": "",
      "After": "",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/server/http.go",
      "Decl": ".ProvideRouter",
      "Kind": "*ast.FieldList",
      "Field": "List",
      "Old": "",
      "New": "evalRuntime *service.EvalHostRuntime",
      "Before": "redisClient *redis.Client",
      "After": "",
      "Scope": "",
      "Comma": true,
      "Count": 0
    },
    {
      "File": "internal/server/middleware/api_key_auth.go",
      "Decl": ".apiKeyAuthWithSubscription",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "if rejectInvalidAuthAbuse(c, apiKeyService) {\n\tAbortWithError(c, http.StatusTooManyRequests, \"INVALID_AUTH_RATE_LIMITED\", \"Too many invalid authentication attempts; retry later\")\n\treturn\n}\nif apiKeyHeadersTooLarge(c) {\n\trecordInvalidAuthFailure(c, apiKeyService)\n\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\tAbortWithError(c, http.StatusUnauthorized, \"INVALID_API_KEY\", \"Invalid API key\")\n\treturn\n}\nqueryKey := strings.TrimSpace(c.Query(\"key\"))\nqueryApiKey := strings.TrimSpace(c.Query(\"api_key\"))\nif queryKey != \"\" || queryApiKey != \"\" {\n\trecordInvalidAuthFailure(c, apiKeyService)\n\tMarkIngressRejected(c, IngressRejectQueryAPIKeyDeprecated)\n\tAbortWithError(c, 400, \"api_key_in_query_deprecated\", \"API key in query parameter is deprecated. Please use Authorization header instead.\")\n\treturn\n}\nauthHeader := c.GetHeader(\"Authorization\")\nvar apiKeyString string\nif authHeader != \"\" {\n\n\tparts := strings.SplitN(authHeader, \" \", 2)\n\tif len(parts) == 2 && strings.EqualFold(parts[0], \"Bearer\") {\n\t\tapiKeyString = strings.TrimSpace(parts[1])\n\t}\n}\nif apiKeyString == \"\" {\n\tapiKeyString = c.GetHeader(\"x-api-key\")\n}\nif len(apiKeyString) > service.MaxAPIKeyCredentialBytes {\n\trecordInvalidAuthFailure(c, apiKeyService)\n\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\tAbortWithError(c, http.StatusUnauthorized, \"INVALID_API_KEY\", \"Invalid API key\")\n\treturn\n}\nif apiKeyString == \"\" {\n\tapiKeyString = c.GetHeader(\"x-goog-api-key\")\n}\nif apiKeyString == \"\" {\n\trecordInvalidAuthFailure(c, apiKeyService)\n\tif hasAPIKeyCredentialInput(c) {\n\t\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\t} else {\n\t\tMarkIngressRejected(c, IngressRejectAPIKeyRequired)\n\t}\n\tAbortWithError(c, 401, \"API_KEY_REQUIRED\", \"API key is required in Authorization header (Bearer scheme), x-api-key header, or x-goog-api-key header\")\n\treturn\n}\napiKey, err := apiKeyService.GetByKey(c.Request.Context(), apiKeyString)\nif err != nil {\n\tif errors.Is(err, service.ErrAPIKeyNotFound) {\n\t\trecordInvalidAuthFailure(c, apiKeyService)\n\t\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\t\tAbortWithError(c, 401, \"INVALID_API_KEY\", \"Invalid API key\")\n\t\treturn\n\t}\n\tif errors.Is(err, service.ErrAPIKeyAuthOverloaded) {\n\t\tMarkIngressRejected(c, IngressRejectAPIKeyAuthOverloaded)\n\t\tAbortWithError(c, http.StatusServiceUnavailable, \"API_KEY_AUTH_OVERLOADED\", \"API key authentication is temporarily unavailable\")\n\t\treturn\n\t}\n\tAbortWithError(c, 500, \"INTERNAL_ERROR\", \"Failed to validate API key\")\n\treturn\n}",
      "New": "apiKey, target, internal, err := apiKeyService.AuthenticateEvalCapability(c.Request)\nevalDec := evalIngressDecision{private: internal, directed: internal, target: target}\nif internal && err != nil {\n\tAbortWithError(c, http.StatusForbidden, \"EVAL_UNAVAILABLE\", \"Evaluation is not ready\")\n\treturn\n}\nif !internal {\n\n\tif rejectInvalidAuthAbuse(c, apiKeyService) {\n\t\tAbortWithError(c, http.StatusTooManyRequests, \"INVALID_AUTH_RATE_LIMITED\", \"Too many invalid authentication attempts; retry later\")\n\t\treturn\n\t}\n\n\tif apiKeyHeadersTooLarge(c) {\n\t\trecordInvalidAuthFailure(c, apiKeyService)\n\t\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\t\tAbortWithError(c, http.StatusUnauthorized, \"INVALID_API_KEY\", \"Invalid API key\")\n\t\treturn\n\t}\n\n\tqueryKey := strings.TrimSpace(c.Query(\"key\"))\n\tqueryApiKey := strings.TrimSpace(c.Query(\"api_key\"))\n\tif queryKey != \"\" || queryApiKey != \"\" {\n\t\trecordInvalidAuthFailure(c, apiKeyService)\n\t\tMarkIngressRejected(c, IngressRejectQueryAPIKeyDeprecated)\n\t\tAbortWithError(c, 400, \"api_key_in_query_deprecated\", \"API key in query parameter is deprecated. Please use Authorization header instead.\")\n\t\treturn\n\t}\n\n\tauthHeader := c.GetHeader(\"Authorization\")\n\tvar apiKeyString string\n\n\tif authHeader != \"\" {\n\n\t\tparts := strings.SplitN(authHeader, \" \", 2)\n\t\tif len(parts) == 2 && strings.EqualFold(parts[0], \"Bearer\") {\n\t\t\tapiKeyString = strings.TrimSpace(parts[1])\n\t\t}\n\t}\n\n\tif apiKeyString == \"\" {\n\t\tapiKeyString = c.GetHeader(\"x-api-key\")\n\t}\n\tif len(apiKeyString) > service.MaxAPIKeyCredentialBytes {\n\t\trecordInvalidAuthFailure(c, apiKeyService)\n\t\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\t\tAbortWithError(c, http.StatusUnauthorized, \"INVALID_API_KEY\", \"Invalid API key\")\n\t\treturn\n\t}\n\n\tif apiKeyString == \"\" {\n\t\tapiKeyString = c.GetHeader(\"x-goog-api-key\")\n\t}\n\n\tif apiKeyString == \"\" {\n\t\trecordInvalidAuthFailure(c, apiKeyService)\n\t\tif hasAPIKeyCredentialInput(c) {\n\t\t\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\t\t} else {\n\t\t\tMarkIngressRejected(c, IngressRejectAPIKeyRequired)\n\t\t}\n\t\tAbortWithError(c, 401, \"API_KEY_REQUIRED\", \"API key is required in Authorization header (Bearer scheme), x-api-key header, or x-goog-api-key header\")\n\t\treturn\n\t}\n\n\tresolvedKey, decision, evalAborted := resolveEvalIngress(c, apiKeyService, apiKeyString, cfg)\n\tif evalAborted {\n\t\treturn\n\t}\n\tapiKeyString = resolvedKey\n\tevalDec = decision\n\n\tapiKey, err = apiKeyService.GetByKey(c.Request.Context(), apiKeyString)\n\tif err != nil {\n\t\tif errors.Is(err, service.ErrAPIKeyNotFound) {\n\t\t\trecordInvalidAuthFailure(c, apiKeyService)\n\t\t\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\t\t\tAbortWithError(c, 401, \"INVALID_API_KEY\", \"Invalid API key\")\n\t\t\treturn\n\t\t}\n\t\tif errors.Is(err, service.ErrAPIKeyAuthOverloaded) {\n\t\t\tMarkIngressRejected(c, IngressRejectAPIKeyAuthOverloaded)\n\t\t\tAbortWithError(c, http.StatusServiceUnavailable, \"API_KEY_AUTH_OVERLOADED\", \"API key authentication is temporarily unavailable\")\n\t\t\treturn\n\t\t}\n\t\tAbortWithError(c, 500, \"INTERNAL_ERROR\", \"Failed to validate API key\")\n\t\treturn\n\t}\n\n}",
      "Before": "",
      "After": "SetOpsFallbackAPIKey(c, apiKey)",
      "Scope": "",
      "Comma": false,
      "Count": 14
    },
    {
      "File": "internal/server/middleware/api_key_auth.go",
      "Decl": ".apiKeyAuthWithSubscription",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if evalDec.directed {\n\tresolvedKey, aborted := applyEvalDirectedTarget(c, apiKeyService, apiKey, evalDec.target)\n\tif aborted {\n\t\treturn\n\t}\n\tapiKey = resolvedKey\n}",
      "Before": "if !apiKey.User.IsActive() {\n\tMarkIngressRejected(c, IngressRejectUserInactive)\n\tAbortWithError(c, 401, \"USER_INACTIVE\", \"User account is not active\")\n\treturn\n}",
      "After": "if abortIfAPIKeyGroupUnavailable(c, apiKey) {\n\treturn\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/server/middleware/api_key_auth.go",
      "Decl": ".apiKeyAuthWithSubscription",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if evalDec.directed {\n\tc.Request = c.Request.WithContext(evalroute.WithTarget(c.Request.Context(), evalDec.target))\n\tevalroute.StripTargetHeaders(c.Request.Header)\n}",
      "Before": "c.Request = c.Request.WithContext(ctx)",
      "After": "billingInfoRequest := c.Request.URL.Path == \"/v1/sub2api/billing\"",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/server/middleware/api_key_auth.go",
      "Decl": ".apiKeyAuthWithSubscription",
      "Kind": "*ast.CallExpr",
      "Field": "",
      "Old": "apiKey.IsExpired()",
      "New": "!evalDec.private && apiKey.IsExpired()",
      "Before": "",
      "After": "",
      "Scope": "/if ;!skipBilling/then/if ;apiKey.IsExpired()",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/server/middleware/api_key_auth_google.go",
      "Decl": ".APIKeyAuthWithSubscriptionGoogle",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if service.IsEvalStoredLookupCredential(apiKeyString) || service.IsEvalPrivateTokenCredential(apiKeyString) {\n\trecordInvalidAuthFailure(c, apiKeyService)\n\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\tabortWithGoogleError(c, 401, \"Invalid API key\")\n\treturn\n}\nif evalroute.HasTargetHeaders(c.Request.Header) {\n\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\tabortWithGoogleError(c, 403, \"Reserved evaluation headers are not allowed\")\n\treturn\n}",
      "Before": "if len(apiKeyString) > service.MaxAPIKeyCredentialBytes {\n\trecordInvalidAuthFailure(c, apiKeyService)\n\tMarkIngressRejected(c, IngressRejectInvalidAPIKey)\n\tabortWithGoogleError(c, 401, \"Invalid API key\")\n\treturn\n}",
      "After": "apiKey, err := apiKeyService.GetByKey(c.Request.Context(), apiKeyString)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/server/router.go",
      "Decl": ".registerRoutes",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "r.GET(\"/codex-eval/ui\", h.Admin.CodexEval.UI)",
      "Before": "routes.RegisterAdminRoutes(v1, h, adminAuth, auditLog, stepUpAuth, settingService, panelRateLimiter)",
      "After": "routes.RegisterGatewayRoutes(r, h, apiKeyAuth, apiKeyService, subscriptionService, opsService, settingService, compositeResolver, cfg)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/server/routes/admin.go",
      "Decl": ".RegisterAdminRoutes",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "evaluation := admin.Group(\"/codex-eval\")\nevaluation.Use(func(c *gin.Context) { c.Header(\"Cache-Control\", \"no-store\"); c.Next() })\nevaluation.GET(\"/status\", h.Admin.CodexEval.Status)\nevaluation.POST(\"/activate\", h.Admin.CodexEval.Activate)\nevaluation.GET(\"/targets\", h.Admin.CodexEval.Targets)\nevaluation.GET(\"/state\", h.Admin.CodexEval.State)\nevaluation.POST(\"/start\", h.Admin.CodexEval.Start)\nevaluation.POST(\"/cancel\", h.Admin.CodexEval.Cancel)\nevaluation.GET(\"/records/:type/:account_id/:run_id\", h.Admin.CodexEval.Record)\nevaluation.POST(\"/clear-history\", h.Admin.CodexEval.ClearHistory)\nevaluation.POST(\"/flush\", h.Admin.CodexEval.Flush)\nevaluation.POST(\"/deactivate\", h.Admin.CodexEval.Deactivate)",
      "Before": "",
      "After": "registerAdminComplianceRoutes(admin, h)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/api_key_service.go",
      "Decl": "*APIKeyService.Create",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "if err := validateCreateAPIKeyRequest(req); err != nil {\n\treturn nil, err\n}\nuser, err := s.userRepo.GetByID(ctx, userID)\nif err != nil {\n\treturn nil, fmt.Errorf(\"get user: %w\", err)\n}\nif len(req.IPWhitelist) > 0 {\n\tif invalid := ip.ValidateIPPatterns(req.IPWhitelist); len(invalid) > 0 {\n\t\treturn nil, fmt.Errorf(\"%w: %v\", ErrInvalidIPPattern, invalid)\n\t}\n}\nif len(req.IPBlacklist) > 0 {\n\tif invalid := ip.ValidateIPPatterns(req.IPBlacklist); len(invalid) > 0 {\n\t\treturn nil, fmt.Errorf(\"%w: %v\", ErrInvalidIPPattern, invalid)\n\t}\n}\nif req.GroupID != nil {\n\tgroup, err := s.groupRepo.GetByID(ctx, *req.GroupID)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"get group: %w\", err)\n\t}\n\n\tif !s.canUserBindGroup(ctx, user, group) {\n\t\treturn nil, ErrGroupNotAllowed\n\t}\n}\nvar key string\nif req.CustomKey != nil && *req.CustomKey != \"\" {\n\n\tif err := s.checkAPIKeyRateLimit(ctx, userID); err != nil {\n\t\treturn nil, err\n\t}\n\n\tif err := s.ValidateCustomKey(*req.CustomKey); err != nil {\n\t\treturn nil, err\n\t}\n\n\texists, err := s.apiKeyRepo.ExistsByKey(ctx, *req.CustomKey)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"check key exists: %w\", err)\n\t}\n\tif exists {\n\n\t\ts.incrementAPIKeyErrorCount(ctx, userID)\n\t\treturn nil, ErrAPIKeyExists\n\t}\n\n\tkey = *req.CustomKey\n} else {\n\n\tvar err error\n\tkey, err = s.GenerateKey()\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"generate key: %w\", err)\n\t}\n}\nif err := s.checkAPIKeyCreateLimits(ctx, userID); err != nil {\n\treturn nil, err\n}\napiKey := &APIKey{\n\tUserID:\t\tuserID,\n\tKey:\t\tkey,\n\tName:\t\thtml.EscapeString(req.Name),\n\tGroupID:\treq.GroupID,\n\tStatus:\t\tStatusActive,\n\tIPWhitelist:\treq.IPWhitelist,\n\tIPBlacklist:\treq.IPBlacklist,\n\tQuota:\t\treq.Quota,\n\tQuotaUsed:\t0,\n\tRateLimit5h:\treq.RateLimit5h,\n\tRateLimit1d:\treq.RateLimit1d,\n\tRateLimit7d:\treq.RateLimit7d,\n}\nif req.ExpiresInDays != nil && *req.ExpiresInDays > 0 {\n\texpiresAt := time.Now().AddDate(0, 0, *req.ExpiresInDays)\n\tapiKey.ExpiresAt = &expiresAt\n}\nif err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {\n\treturn nil, fmt.Errorf(\"create api key: %w\", err)\n}\ns.InvalidateAuthCacheByKey(ctx, apiKey.Key)\ns.compileAPIKeyIPRules(apiKey)\nreturn apiKey, nil",
      "New": "return s.createWithExpiration(ctx, userID, req, nil)",
      "Before": "",
      "After": "",
      "Scope": "",
      "Comma": false,
      "Count": 15
    },
    {
      "File": "internal/service/gateway_forward_as_chat_completions.go",
      "Decl": "*GatewayService.ForwardAsChatCompletions",
      "Kind": "*ast.Ident",
      "Field": "",
      "Old": "detachStreamUpstreamContext",
      "New": "streamUpstreamRequestContext",
      "Before": "",
      "After": "",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/gateway_forward_as_responses.go",
      "Decl": "*GatewayService.ForwardAsResponses",
      "Kind": "*ast.Ident",
      "Field": "",
      "Old": "detachStreamUpstreamContext",
      "New": "streamUpstreamRequestContext",
      "Before": "",
      "After": "",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_account_scheduler.go",
      "Decl": "*OpenAIGatewayService.selectAccountWithScheduler",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed {\n\treturn selection, decision, err\n}",
      "Before": "selection, decision, err := s.selectAccountWithSchedulerOnce(ctx, groupID, previousResponseID, sessionHash, requestedModel, excludedIDs, requiredTransport, requiredCapability, requiredImageCapability, requireCompact, platform, previousResponseCanMove, useUpstreamTokenCost)",
      "After": "if err == nil || openAIProxyStreamQuarantineBypassed(ctx) {\n\treturn selection, decision, err\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_account_scheduler.go",
      "Decl": "*OpenAIGatewayService.selectAccountWithSchedulerOnce",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if target, directed := evalroute.FromContext(ctx); directed {\n\tif platform != PlatformOpenAI {\n\t\treturn nil, OpenAIAccountScheduleDecision{Layer: openAIAccountScheduleLayerEvalDirected}, ErrEvalDirectedRoutingUnsupported\n\t}\n\treturn s.selectDirectedEvaluationAccount(ctx, groupID, requestedModel, excludedIDs, requiredTransport, requiredCapability, requiredImageCapability, requireCompact, target)\n}",
      "Before": "platform = NormalizeOpenAICompatiblePlatform(platform)",
      "After": "decision := OpenAIAccountScheduleDecision{}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_cc_pipeline.go",
      "Decl": "*OpenAIGatewayService.sendCCUpstreamRequest",
      "Kind": "*ast.Ident",
      "Field": "",
      "Old": "detachUpstreamContext",
      "New": "upstreamRequestContext",
      "Before": "",
      "After": "",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_forward.go",
      "Decl": "*OpenAIGatewayService.Forward",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed && wsResult != nil {\n\treturn wsResult, wsErr\n}",
      "Before": "s.writeOpenAIWSFallbackErrorResponse(c, account, wsErr)",
      "After": "return nil, wsErr",
      "Scope": "/if ;wsDecision.Transport == OpenAIUpstreamTransportResponsesWebsocketV2/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_forward.go",
      "Decl": "*OpenAIGatewayService.Forward",
      "Kind": "*ast.Ident",
      "Field": "",
      "Old": "detachUpstreamContext",
      "New": "upstreamRequestContext",
      "Before": "",
      "After": "",
      "Scope": "/for ",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_forward.go",
      "Decl": "*OpenAIGatewayService.Forward",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if partial := openAIDirectedStreamPartialResult(ctx, c, resp, originalModel, billingModel, upstreamModel, serviceTier, reasoningEffort, reqStream, startTime, streamResult, imageSizeTier, imageInputSize, imageBillingModel); partial != nil {\n\treturn partial, err\n}",
      "Before": "if signal, ok := asOpenAICompactFallbackSignal(err); ok {\n\tif retryBody, fallbackModel, retry := s.prepareOpenAICompactFallbackRetry(\n\t\tc, account, requestedModel, body, http.StatusBadRequest, signal.message, signal.payload, compactModelFallbackRetried,\n\t); retry {\n\t\ts.appendOpenAICompactFallbackRetryOps(c, account, resp, signal.payload, signal.message, false)\n\t\tbody = retryBody\n\t\trequestView = newOpenAIRequestView(body)\n\t\tupstreamModel = fallbackModel\n\t\tcompactModelFallbackRetried = true\n\t\tSetOpsUpstreamModel(c, fallbackModel)\n\t\tcontinue\n\t}\n\tif resp.Body != nil {\n\t\t_ = resp.Body.Close()\n\t}\n\tcompactResp, compactBody := openAICompactFallbackErrorResponse(resp, signal)\n\tif s.shouldFailoverOpenAIUpstreamResponse(account, compactResp.StatusCode, signal.message, compactBody) {\n\t\tappendOpsUpstreamError(c, OpsUpstreamErrorEvent{\n\t\t\tProxyID:\t\topsUpstreamProxyID(account),\n\t\t\tProxyName:\t\topsUpstreamProxyName(account),\n\t\t\tPlatform:\t\taccount.Platform,\n\t\t\tAccountID:\t\taccount.ID,\n\t\t\tAccountName:\t\taccount.Name,\n\t\t\tUpstreamStatusCode:\tcompactResp.StatusCode,\n\t\t\tUpstreamRequestID:\tcompactResp.Header.Get(\"x-request-id\"),\n\t\t\tKind:\t\t\t\"failover\",\n\t\t\tMessage:\t\tsignal.message,\n\t\t})\n\t\tshouldDisable := s.handleFailoverSideEffects(ctx, compactResp, account, compactBody, upstreamModel)\n\t\treturn nil, s.newOpenAIAccountFailoverError(\n\t\t\taccount, compactResp.StatusCode, compactResp.Header, compactBody, signal.message, shouldDisable,\n\t\t\t!shouldDisable && account.IsPoolMode() && (account.IsPoolModeRetryableStatus(compactResp.StatusCode) || isOpenAITransientProcessingError(compactResp.StatusCode, signal.message, compactBody)),\n\t\t)\n\t}\n\treturn s.handleErrorResponse(ctx, compactResp, c, account, body, resolveOpenAIErrorSchedulingModel(billingModel, upstreamModel))\n}",
      "After": "return nil, err",
      "Scope": "/for /if ;reqStream/then/if ;err != nil/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_passthrough.go",
      "Decl": "*OpenAIGatewayService.forwardOpenAIPassthrough",
      "Kind": "*ast.Ident",
      "Field": "",
      "Old": "detachUpstreamContext",
      "New": "upstreamRequestContext",
      "Before": "",
      "After": "",
      "Scope": "/for ",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_passthrough.go",
      "Decl": "*OpenAIGatewayService.forwardOpenAIPassthrough",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if partial := openAIDirectedPassthroughStreamPartialResult(ctx, c, resp, reqModel, upstreamPassthroughModel, extractOpenAIServiceTierFromBody(body), reasoningEffort, reqStream, startTime, result, imageSizeTier, imageInputSize, imageBillingModel); partial != nil {\n\treturn partial, handleErr\n}",
      "Before": "_ = resp.Body.Close()",
      "After": "return nil, handleErr",
      "Scope": "/for /if ;reqStream/then/if ;handleErr != nil/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_passthrough.go",
      "Decl": "*OpenAIGatewayService.forwardOpenAIPassthrough",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if partial := openAIDirectedPassthroughNonStreamPartialResult(ctx, c, resp, reqModel, upstreamPassthroughModel, extractOpenAIServiceTierFromBody(body), reasoningEffort, reqStream, startTime, result, imageSizeTier, imageInputSize, imageBillingModel); partial != nil {\n\treturn partial, handleErr\n}",
      "Before": "_ = resp.Body.Close()",
      "After": "return nil, handleErr",
      "Scope": "/for /if ;reqStream/else/if ;handleErr != nil/then",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_scheduling.go",
      "Decl": "*OpenAIGatewayService.selectAccountForModelWithExclusionsStickyHit",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed {\n\treturn nil, false, ErrEvalDirectedRoutingUnsupported\n}",
      "Before": "platform = NormalizeOpenAICompatiblePlatform(platform)",
      "After": "if s.checkChannelPricingRestriction(ctx, groupID, requestedModel) {\n\tslog.Warn(\"channel pricing restriction blocked request\",\n\t\t\"group_id\", derefGroupID(groupID),\n\t\t\"model\", requestedModel)\n\treturn nil, false, fmt.Errorf(\"%w supporting model: %s (channel pricing restriction)\", ErrNoAvailableAccounts, requestedModel)\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_gateway_scheduling.go",
      "Decl": "*OpenAIGatewayService.selectAccountWithLoadAwareness",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed {\n\treturn nil, ErrEvalDirectedRoutingUnsupported\n}",
      "Before": "platform = NormalizeOpenAICompatiblePlatform(platform)",
      "After": "if s.checkChannelPricingRestriction(ctx, groupID, requestedModel) {\n\tslog.Warn(\"channel pricing restriction blocked request\",\n\t\t\"group_id\", derefGroupID(groupID),\n\t\t\"model\", requestedModel)\n\treturn nil, fmt.Errorf(\"%w supporting model: %s (channel pricing restriction)\", ErrNoAvailableAccounts, requestedModel)\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_profit_control.go",
      "Decl": "*OpenAIGatewayService.bindOpenAIStickySessionDuringSelection",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed {\n\treturn nil\n}",
      "Before": "",
      "After": "if gatewayProfitControlGateActive(ctx) || preserveOpenAIGuardianParentBinding(ctx, sessionHash) {\n\treturn nil\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_profit_control.go",
      "Decl": "*OpenAIGatewayService.BindStickySessionAfterProfitAdmission",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed {\n\treturn nil\n}",
      "Before": "",
      "After": "if sessionHash == \"\" || accountID <= 0 {\n\treturn nil\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_sticky_compat.go",
      "Decl": "*OpenAIGatewayService.setStickySessionAccountID",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed {\n\treturn nil\n}",
      "Before": "if s == nil || s.cache == nil || accountID <= 0 {\n\treturn nil\n}",
      "After": "primaryKey := s.openAISessionCacheKey(sessionHash)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_ws_forwarder_support.go",
      "Decl": "*OpenAIGatewayService.selectAccountByPreviousResponseIDForCapability",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "if _, directed := evalroute.FromContext(ctx); directed {\n\treturn nil, ErrEvalDirectedRoutingUnsupported\n}",
      "Before": "if s == nil {\n\treturn nil, nil\n}",
      "After": "accountID, account, responseID, store := s.resolveAccountByPreviousResponseIDForCapability(ctx, groupID, previousResponseID, requestedModel, excludedIDs, requiredCapability, requireCompact)",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_ws_forwarder_v2.go",
      "Decl": "*OpenAIGatewayService.forwardOpenAIWSV2",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "_, directedEval := evalroute.FromContext(ctx)",
      "Before": "readTimeout := s.openAIWSReadTimeout()",
      "After": "upstreamReadCtx := ctx",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_ws_forwarder_v2.go",
      "Decl": "*OpenAIGatewayService.forwardOpenAIWSV2",
      "Kind": "*ast.BlockStmt",
      "Field": "List",
      "Old": "",
      "New": "var upstreamReadCancel context.CancelFunc\nif directedEval {\n\tupstreamReadCtx, upstreamReadCancel = context.WithCancel(ctx)\n\tdefer upstreamReadCancel()\n}",
      "Before": "upstreamReadDetached := false",
      "After": "clientRequestCanceled := func() bool {\n\treturn ctx != nil && errors.Is(ctx.Err(), context.Canceled)\n}",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_ws_forwarder_v2.go",
      "Decl": "*OpenAIGatewayService.forwardOpenAIWSV2",
      "Kind": "*ast.IfStmt",
      "Field": "",
      "Old": "if !upstreamReadDetached {\n\tupstreamReadCtx = context.WithoutCancel(ctx)\n\tupstreamReadDetached = true\n}",
      "New": "if directedEval {\n\tif upstreamReadCancel != nil {\n\t\tupstreamReadCancel()\n\t}\n} else if !upstreamReadDetached {\n\tupstreamReadCtx = context.WithoutCancel(ctx)\n\tupstreamReadDetached = true\n}",
      "Before": "",
      "After": "",
      "Scope": "/if ;!upstreamReadDetached",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/openai_ws_forwarder_v2.go",
      "Decl": "*OpenAIGatewayService.forwardOpenAIWSV2",
      "Kind": "*ast.UnaryExpr",
      "Field": "",
      "Old": "!readUsedDetachedContext",
      "New": "!readUsedDetachedContext && !directedEval",
      "Before": "",
      "After": "",
      "Scope": "/for /if ;readErr != nil/then/if ;clientDisconnected/then/if ;!readUsedDetachedContext && errors.Is(readErr, context.Canceled) && clientRequestCanceled()",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/service/wire.go",
      "Decl": "var:ProviderSet",
      "Kind": "*ast.CallExpr",
      "Field": "Args",
      "Old": "",
      "New": "NewEvalHostRuntime",
      "Before": "",
      "After": "ProvideAuthService",
      "Scope": "",
      "Comma": false,
      "Count": 0
    },
    {
      "File": "internal/web/embed_on.go",
      "Decl": ".shouldBypassEmbeddedFrontend",
      "Kind": "*ast.CallExpr",
      "Field": "",
      "Old": "strings.HasPrefix(trimmed, \"/api/\")",
      "New": "trimmed == \"/codex-eval/ui\" || strings.HasPrefix(trimmed, \"/api/\")",
      "Before": "",
      "After": "",
      "Scope": "",
      "Comma": false,
      "Count": 0
    }
  ],
  "Imports": {
    "cmd/server/main.go": [
      "\"net\""
    ],
    "internal/handler/gateway_handler_responses.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/handler/openai_gateway_handler.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/repository/api_key_repo.go": [
      "\"github.com/Wei-Shaw/sub2api/ent/predicate\""
    ],
    "internal/server/middleware/api_key_auth.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/server/middleware/api_key_auth_google.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/gateway_usage_billing.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_account_scheduler.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_gateway_forward.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_gateway_passthrough.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_gateway_scheduling.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_profit_control.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_sticky_compat.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_ws_forwarder_support.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ],
    "internal/service/openai_ws_forwarder_v2.go": [
      "\"github.com/Wei-Shaw/sub2api/internal/pkg/evalroute\""
    ]
  },
  "Declarations": {
    "internal/repository/api_key_repo.go": [
      "func evalAccountingPredicates(row *service.APIKey) []predicate.APIKey {\n\treturn []predicate.APIKey{apikey.IDEQ(row.ID), apikey.UserIDEQ(row.UserID), apikey.KeyEQ(row.Key), apikey.StatusEQ(service.StatusActive), apikey.NameEQ(row.Name), apikey.DeletedAtIsNil(), apikey.GroupIDIsNil(), apikey.QuotaEQ(row.Quota), apikey.RateLimit5hEQ(row.RateLimit5h), apikey.RateLimit1dEQ(row.RateLimit1d), apikey.RateLimit7dEQ(row.RateLimit7d), func(s *entsql.Selector) {\n\t\twhitelist, _ := json.Marshal(row.IPWhitelist)\n\t\ts.Where(entsql.EQ(s.C(apikey.FieldIPWhitelist), string(whitelist)))\n\t\ts.Where(entsql.Or(entsql.IsNull(s.C(apikey.FieldIPBlacklist)), entsql.EQ(s.C(apikey.FieldIPBlacklist), \"[]\")))\n\t}}\n}",
      "func (r *apiKeyRepository) ExpireEvalAccountingKey(ctx context.Context, row *service.APIKey, expiry time.Time) error {\n\tfilters := append(evalAccountingPredicates(row), apikey.ExpiresAtIsNil())\n\tn, err := clientFromContext(ctx, r.client).APIKey.Update().Where(filters...).SetExpiresAt(expiry).SetUpdatedAt(time.Now()).Save(ctx)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif n != 1 {\n\t\treturn service.ErrEvalKeyExpiryConflict\n\t}\n\treturn nil\n}",
      "func (r *apiKeyRepository) RevokeEvalAccountingKey(ctx context.Context, row *service.APIKey) error {\n\tif row.ExpiresAt == nil {\n\t\treturn service.ErrEvalKeyExpiryConflict\n\t}\n\tfilters := append(evalAccountingPredicates(row), apikey.ExpiresAtEQ(*row.ExpiresAt))\n\tn, err := clientFromContext(ctx, r.client).APIKey.Update().Where(filters...).SetKey(fmt.Sprintf(\"__codex_eval_revoked__%d\", row.ID)).SetDeletedAt(time.Now()).SetUpdatedAt(time.Now()).Save(ctx)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif n != 1 {\n\t\treturn service.ErrEvalStorageCorrupt\n\t}\n\treturn nil\n}",
      "func (r *apiKeyRepository) EvalAccountingKeyRevoked(ctx context.Context, id, owner int64, expiry time.Time) (bool, error) {\n\n\treturn r.client.APIKey.Query().Where(apikey.IDEQ(id), apikey.UserIDEQ(owner), apikey.KeyEQ(fmt.Sprintf(\"__codex_eval_revoked__%d\", id)), apikey.DeletedAtNotNil(), apikey.ExpiresAtEQ(expiry)).Exist(mixins.SkipSoftDelete(ctx))\n}"
    ],
    "internal/server/http.go": [
      "func ProvideEvalHTTPServer(cfg *config.Config, router *gin.Engine, evalRuntime *service.EvalHostRuntime) *http.Server {\n\tserver := ProvideHTTPServer(cfg, router)\n\tevalRuntime.BindRouter(server.Handler)\n\tserver.Handler = evalRuntime.Guard(server.Handler)\n\treturn server\n}"
    ],
    "internal/service/api_key_service.go": [
      "func (s *APIKeyService) createEvalAccountingKey(ctx context.Context, userID int64, lookup string, sources []string) (*APIKey, error) {\n\texpiry := evalAccountingExpiry()\n\treturn s.createWithExpiration(ctx, userID, CreateAPIKeyRequest{Name: EvalAPIKeyName, CustomKey: &lookup, IPWhitelist: sources}, &expiry)\n}",
      "func (s *APIKeyService) createWithExpiration(ctx context.Context, userID int64, req CreateAPIKeyRequest, fixedExpiry *time.Time) (*APIKey, error) {\n\tif err := validateCreateAPIKeyRequest(req); err != nil {\n\t\treturn nil, err\n\t}\n\n\tuser, err := s.userRepo.GetByID(ctx, userID)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"get user: %w\", err)\n\t}\n\n\tif len(req.IPWhitelist) > 0 {\n\t\tif invalid := ip.ValidateIPPatterns(req.IPWhitelist); len(invalid) > 0 {\n\t\t\treturn nil, fmt.Errorf(\"%w: %v\", ErrInvalidIPPattern, invalid)\n\t\t}\n\t}\n\n\tif len(req.IPBlacklist) > 0 {\n\t\tif invalid := ip.ValidateIPPatterns(req.IPBlacklist); len(invalid) > 0 {\n\t\t\treturn nil, fmt.Errorf(\"%w: %v\", ErrInvalidIPPattern, invalid)\n\t\t}\n\t}\n\n\tif req.GroupID != nil {\n\t\tgroup, err := s.groupRepo.GetByID(ctx, *req.GroupID)\n\t\tif err != nil {\n\t\t\treturn nil, fmt.Errorf(\"get group: %w\", err)\n\t\t}\n\n\t\tif !s.canUserBindGroup(ctx, user, group) {\n\t\t\treturn nil, ErrGroupNotAllowed\n\t\t}\n\t}\n\n\tvar key string\n\n\tif req.CustomKey != nil && *req.CustomKey != \"\" {\n\n\t\tif err := s.checkAPIKeyRateLimit(ctx, userID); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\n\t\tif err := s.ValidateCustomKey(*req.CustomKey); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\n\t\texists, err := s.apiKeyRepo.ExistsByKey(ctx, *req.CustomKey)\n\t\tif err != nil {\n\t\t\treturn nil, fmt.Errorf(\"check key exists: %w\", err)\n\t\t}\n\t\tif exists {\n\n\t\t\ts.incrementAPIKeyErrorCount(ctx, userID)\n\t\t\treturn nil, ErrAPIKeyExists\n\t\t}\n\n\t\tkey = *req.CustomKey\n\t} else {\n\n\t\tvar err error\n\t\tkey, err = s.GenerateKey()\n\t\tif err != nil {\n\t\t\treturn nil, fmt.Errorf(\"generate key: %w\", err)\n\t\t}\n\t}\n\n\tif err := s.checkAPIKeyCreateLimits(ctx, userID); err != nil {\n\t\treturn nil, err\n\t}\n\n\tapiKey := &APIKey{\n\t\tUserID:\t\tuserID,\n\t\tKey:\t\tkey,\n\t\tName:\t\thtml.EscapeString(req.Name),\n\t\tGroupID:\treq.GroupID,\n\t\tStatus:\t\tStatusActive,\n\t\tIPWhitelist:\treq.IPWhitelist,\n\t\tIPBlacklist:\treq.IPBlacklist,\n\t\tQuota:\t\treq.Quota,\n\t\tQuotaUsed:\t0,\n\t\tRateLimit5h:\treq.RateLimit5h,\n\t\tRateLimit1d:\treq.RateLimit1d,\n\t\tRateLimit7d:\treq.RateLimit7d,\n\t}\n\n\tif req.ExpiresInDays != nil && *req.ExpiresInDays > 0 {\n\t\texpiresAt := time.Now().AddDate(0, 0, *req.ExpiresInDays)\n\t\tapiKey.ExpiresAt = &expiresAt\n\t}\n\n\tif fixedExpiry != nil {\n\t\texpiry := *fixedExpiry\n\t\tapiKey.ExpiresAt = &expiry\n\t}\n\tif err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {\n\t\treturn nil, fmt.Errorf(\"create api key: %w\", err)\n\t}\n\n\ts.InvalidateAuthCacheByKey(ctx, apiKey.Key)\n\ts.compileAPIKeyIPRules(apiKey)\n\n\treturn apiKey, nil\n}"
    ],
    "internal/service/gateway_usage_billing.go": [
      "func upstreamRequestContext(ctx context.Context) (context.Context, context.CancelFunc) {\n\tif _, directed := evalroute.FromContext(ctx); directed {\n\t\treturn ctx, func() {}\n\t}\n\treturn detachUpstreamContext(ctx)\n}",
      "func streamUpstreamRequestContext(ctx context.Context, stream bool) (context.Context, context.CancelFunc) {\n\tif _, directed := evalroute.FromContext(ctx); directed {\n\t\tif ctx == nil {\n\t\t\treturn context.Background(), func() {}\n\t\t}\n\t\treturn ctx, func() {}\n\t}\n\treturn detachStreamUpstreamContext(ctx, stream)\n}"
    ],
    "internal/service/openai_gateway_forward.go": [
      "func openAIDirectedStreamPartialResult(\n\tctx context.Context,\n\tc *gin.Context,\n\tresp *http.Response,\n\toriginalModel, billingModel, upstreamModel string,\n\tserviceTier *string,\n\treasoningEffort *string,\n\treqStream bool,\n\tstartTime time.Time,\n\tstreamResult *openaiStreamingResult,\n\timageSizeTier, imageInputSize, imageBillingModel string,\n) *OpenAIForwardResult {\n\tif streamResult == nil {\n\t\treturn nil\n\t}\n\tif _, directed := evalroute.FromContext(ctx); !directed {\n\t\treturn nil\n\t}\n\tusage := OpenAIUsage{}\n\tif streamResult.usage != nil {\n\t\tusage = *streamResult.usage\n\t}\n\tresult := &OpenAIForwardResult{\n\t\tRequestID:\t\t\tresp.Header.Get(\"x-request-id\"),\n\t\tUpstreamHeaders:\t\tresp.Header,\n\t\tResponseID:\t\t\tstreamResult.responseID,\n\t\tUsage:\t\t\t\tusage,\n\t\tModel:\t\t\t\toriginalModel,\n\t\tBillingModel:\t\t\tbillingModel,\n\t\tUpstreamModel:\t\t\tupstreamModel,\n\t\tUpstreamResponseModel:\t\tobservedUpstreamResponseModel(c),\n\t\tUpstreamResponseModelConflict:\tobservedUpstreamResponseModelConflict(c),\n\t\tUpstreamResponseServiceTier:\tobservedUpstreamResponseServiceTier(c),\n\t\tServiceTier:\t\t\tresolvedOpenAIUpstreamServiceTier(c, serviceTier),\n\t\tReasoningEffort:\t\treasoningEffort,\n\t\tStream:\t\t\t\treqStream,\n\t\tDuration:\t\t\ttime.Since(startTime),\n\t\tFirstTokenMs:\t\t\tstreamResult.firstTokenMs,\n\t\tClientDisconnect:\t\tctx != nil && ctx.Err() != nil,\n\t}\n\tif streamResult.imageCount > 0 {\n\t\tresult.ImageCount = streamResult.imageCount\n\t\tresult.ImageSize = imageSizeTier\n\t\tresult.ImageInputSize = imageInputSize\n\t\tresult.ImageOutputSizes = streamResult.imageOutputSizes\n\t\tresult.BillingModel = imageBillingModel\n\t}\n\tif streamResult.searchCount > 0 {\n\t\tresult.SearchCount = streamResult.searchCount\n\t}\n\treturn result\n}"
    ],
    "internal/service/openai_gateway_passthrough.go": [
      "func openAIDirectedPassthroughStreamPartialResult(\n\tctx context.Context,\n\tc *gin.Context,\n\tresp *http.Response,\n\treqModel, upstreamPassthroughModel string,\n\tserviceTier *string,\n\treasoningEffort *string,\n\treqStream bool,\n\tstartTime time.Time,\n\tresult *openaiStreamingResultPassthrough,\n\timageSizeTier, imageInputSize, imageBillingModel string,\n) *OpenAIForwardResult {\n\tif result == nil {\n\t\treturn nil\n\t}\n\treturn openAIDirectedPassthroughPartialResult(ctx, c, resp, reqModel, upstreamPassthroughModel, serviceTier, reasoningEffort, reqStream, startTime, result.usage, result.responseID, result.firstTokenMs, result.imageCount, result.imageOutputSizes, imageSizeTier, imageInputSize, imageBillingModel)\n}",
      "func openAIDirectedPassthroughNonStreamPartialResult(\n\tctx context.Context,\n\tc *gin.Context,\n\tresp *http.Response,\n\treqModel, upstreamPassthroughModel string,\n\tserviceTier *string,\n\treasoningEffort *string,\n\treqStream bool,\n\tstartTime time.Time,\n\tresult *openaiNonStreamingResultPassthrough,\n\timageSizeTier, imageInputSize, imageBillingModel string,\n) *OpenAIForwardResult {\n\tif result == nil {\n\t\treturn nil\n\t}\n\treturn openAIDirectedPassthroughPartialResult(ctx, c, resp, reqModel, upstreamPassthroughModel, serviceTier, reasoningEffort, reqStream, startTime, result.usage, result.responseID, nil, result.imageCount, result.imageOutputSizes, imageSizeTier, imageInputSize, imageBillingModel)\n}",
      "func openAIDirectedPassthroughPartialResult(\n\tctx context.Context,\n\tc *gin.Context,\n\tresp *http.Response,\n\treqModel, upstreamPassthroughModel string,\n\tserviceTier *string,\n\treasoningEffort *string,\n\treqStream bool,\n\tstartTime time.Time,\n\tusage *OpenAIUsage,\n\tresponseID string,\n\tfirstTokenMs *int,\n\timageCount int,\n\timageOutputSizes []string,\n\timageSizeTier, imageInputSize, imageBillingModel string,\n) *OpenAIForwardResult {\n\tif _, directed := evalroute.FromContext(ctx); !directed {\n\t\treturn nil\n\t}\n\tusageValue := OpenAIUsage{}\n\tif usage != nil {\n\t\tusageValue = *usage\n\t}\n\tresult := &OpenAIForwardResult{\n\t\tRequestID:\t\t\tresp.Header.Get(\"x-request-id\"),\n\t\tUpstreamHeaders:\t\tresp.Header,\n\t\tResponseID:\t\t\tresponseID,\n\t\tUsage:\t\t\t\tusageValue,\n\t\tModel:\t\t\t\treqModel,\n\t\tUpstreamModel:\t\t\tupstreamPassthroughModel,\n\t\tUpstreamResponseModel:\t\tobservedUpstreamResponseModel(c),\n\t\tUpstreamResponseModelConflict:\tobservedUpstreamResponseModelConflict(c),\n\t\tUpstreamResponseServiceTier:\tobservedUpstreamResponseServiceTier(c),\n\t\tServiceTier:\t\t\tresolvedOpenAIUpstreamServiceTier(c, serviceTier),\n\t\tReasoningEffort:\t\treasoningEffort,\n\t\tStream:\t\t\t\treqStream,\n\t\tDuration:\t\t\ttime.Since(startTime),\n\t\tFirstTokenMs:\t\t\tfirstTokenMs,\n\t\tClientDisconnect:\t\tctx != nil && ctx.Err() != nil,\n\t}\n\tif imageCount > 0 {\n\t\tresult.ImageCount = imageCount\n\t\tresult.ImageSize = imageSizeTier\n\t\tresult.ImageInputSize = imageInputSize\n\t\tresult.ImageOutputSizes = imageOutputSizes\n\t\tresult.BillingModel = imageBillingModel\n\t}\n\treturn result\n}"
    ]
  }
}
```

## 8条前端修改的原始定义

```json
{
  "frontend/src/components/layout/AppSidebar.vue": [
    {
      "before": "import { useAdminSettingsStore, useAppStore, useAuthStore, useOnboardingStore } from '@/stores'\n",
      "old": "",
      "after": "import VersionBadge from '@/components/common/VersionBadge.vue'\n",
      "new": "import { isOwnedCodexEvalMenu } from '@/utils/codex-eval-menu'\n"
    },
    {
      "before": "  return adminSettingsStore.customMenuItems\n",
      "old": "    .filter((item) => item.visibility === 'admin')\n",
      "after": "    .sort((a, b) => a.sort_order - b.sort_order)\n",
      "new": "    .filter((item) => item.visibility === 'admin' && !isOwnedCodexEvalMenu(item))\n"
    },
    {
      "before": "    { path: '/admin/accounts', label: t('nav.accounts'), icon: GlobeIcon },\n",
      "old": "",
      "after": "    { path: '/admin/plugins', label: t('nav.plugins'), icon: PluginIcon, featureFlag: flagPluginManagement },\n",
      "new": "    { path: '/custom/codex-eval', label: '账户评估', icon: ChartIcon, hideInSimpleMode: true },\n"
    }
  ],
  "frontend/src/views/user/CustomPageView.vue": [
    {
      "before": "      <div class=\"card flex-1 min-h-0 overflow-hidden\">\n",
      "old": "        <div v-if=\"loading\" class=\"flex h-full items-center justify-center py-12\">\n",
      "after": "          <div\n",
      "new": "        <div v-if=\"evalMenuConflict\" class=\"flex h-full items-center justify-center p-10\" role=\"alert\">账户评估菜单冲突。请检查管理面板的自定义菜单。</div>\n        <div v-else-if=\"loading\" class=\"flex h-full items-center justify-center py-12\">\n"
    },
    {
      "before": "import { buildEmbeddedUrl, detectTheme } from '@/utils/embedded-url'\n",
      "old": "",
      "after": "import { marked } from 'marked'\n",
      "new": "import { CODEX_EVAL_MENU_ID, builtInCodexEvalMenu, codexEvalMenuConflict } from '@/utils/codex-eval-menu'\n"
    },
    {
      "before": "\n",
      "old": "",
      "after": "const menuItem = computed(() => {\n",
      "new": "const evalMenuConflict = computed(() => authStore.isAdmin && menuItemId.value === CODEX_EVAL_MENU_ID\n  && codexEvalMenuConflict([...(appStore.cachedPublicSettings?.custom_menu_items ?? []), ...adminSettingsStore.customMenuItems]))\n\n"
    },
    {
      "before": "  const id = menuItemId.value\n",
      "old": "",
      "after": "  const publicItems = appStore.cachedPublicSettings?.custom_menu_items ?? []\n",
      "new": "  if (id === CODEX_EVAL_MENU_ID) return authStore.isAdmin && !evalMenuConflict.value ? builtInCodexEvalMenu() : null\n"
    }
  ],
  "frontend/src/utils/embedded-url.ts": [
    {
      "before": "    const url = new URL(baseUrl)\n",
      "old": "",
      "after": "    if (userId) {\n",
      "new": "    if (typeof window !== 'undefined' && url.origin === window.location.origin && url.pathname === '/codex-eval/ui') {\n      url.search = ''\n      url.hash = ''\n      url.searchParams.set(EMBEDDED_THEME_QUERY_KEY, theme)\n      if (lang) url.searchParams.set(EMBEDDED_LANG_QUERY_KEY, lang)\n      return url.toString()\n    }\n"
    }
  ]
}
```

## 未改变的原评估内容

Codex上下文和探针/基线/ModelTrace数据相同；Fingerprint和ModelTrace纯算法仅包名不同。原题目、评分、挑战以及相应数据不应被本清单误计为重新设计。详情见已有core/tooling审计的逐项比较。

## 证据来源

- [原生审计](patch-scope-native-inventory.md)及[逐条结构记录](patch-scope-native-inventory.json)。
- [核心与工具链审计](patch-scope-core-tooling-inventory.md)及[逐项结构记录](patch-scope-core-tooling-inventory.json)。
- [审批记录](patch-scope-approval-ledger.json)。
