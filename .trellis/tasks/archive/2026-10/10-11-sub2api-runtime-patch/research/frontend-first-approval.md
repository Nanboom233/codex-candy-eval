# 前端从零审批

最新用户要求：剔除未经批准的前端修改，从前端开始，审批需要做的修改。
任务保持 planning。用户最新要求后续每批展示5个决策；每项仍有独立编号，整批确认可以记录各项决定。本文件是当前审批方案；上一轮全栈D/P清单暂停，不在本轮继续推进后端或安装设计。
旧实现仍是需要回退的现状，不作为默认保留项。不得把“取消旧设计”报告成已经完成产品代码回退。

## 取消旧前端设计

- adapter/ui-rules.json中8条修改全部退出方案：AppSidebar.vue的硬编码菜单/旧项过滤、CustomPageView.vue的冲突/内置回退、embedded-url.ts特殊处理。
- adapter/runtime/frontend/src/utils/codex-eval-menu.ts整份辅助文件退出方案。
- internal/eval/web中的旧API/body/响应转换、原生数据到CPA字段转换、模型目录聚合/启动前废弃请求、失效计划/开关逻辑，不按旧实现继续迁就；接口接入逐项重新审批。
- 主题重写、style-to-class、Pelican HTML重写、额外激活/状态按钮、版本/品牌改动均取消默认保留；只有单独审批的页面行为才进入新实现。
- 原四种评估的题目、算法和数据不属于本轮前端改造；原internal/plugin/web仅作布局/功能参考，不恢复其CPA兼容机制到新交付。

## 原生能力证据

源码基线来自已验证的Sub2API v0.2.16归档，路径在pristine-v0.2.16-baseline.md。

- frontend/src/components/layout/AppSidebar.vue:769-774读取已有管理员customMenuItems；849-858已生成/custom/<id>导航项。
- frontend/src/router/index.ts:389已有/custom/:id路由。
- frontend/src/views/user/CustomPageView.vue:223-236已从管理员设置读取菜单；244起使用原有嵌入页功能。
- frontend/src/views/admin/SettingsView.vue:6771起已有自定义菜单编辑；10845起新增项；11528保存custom_menu_items。
- frontend/src/utils/embedded-url.ts:23-59是原有URL构造器；本方案不修改它。它本身附加token/theme等原生嵌入参数，不能声称这些是本次新增代码。登录接入另行审批。

由现有源码可确定：展示一条管理员评估入口不需要8条原生前端补丁、不需要新Vue路由、不需要菜单冲突框架或内置菜单回退。

## 已审批：F01 页面入口

页面入口为 `/codex-eval`，不再使用 `/codex-eval/ui`。用户在菜单方案审批后回复“继续”，批准本项规划；不连带批准其他前端行为或开始实施。

使用原生自定义菜单新增一条管理员入口：

| 字段 | 内容 |
| --- | --- |
| 名称 | 账户评估 |
| 可见范围 | admin |
| 页面URL | 当前站点origin + /codex-eval |
| 页面类型 | 原生URL嵌入页，非Markdown |
| 排序 | 管理员自定义菜单末尾 |

这只是原生设置中的一个菜单项，不修改AppSidebar.vue、CustomPageView.vue、embedded-url.ts或路由系统。
原四种评估页面的具体接入尚未批准，F01不连带批准token读取、数据转换、额外按钮、缓存、主题或渲染重写。
当前只准备审批内容，没有写入外部设置，也没有宣称菜单或页面已可用。

## 已审批：F02 管理员登录接入

用户回复“是”，批准下述F02范围；记录为规划批准，产品代码尚未修改。

最小方案：评估页与管理面板同源，每次请求读取原生 `localStorage.auth_token`，发送 `Authorization: Bearer <token>`。页面不再保留CPA密钥变量、KEY_STORE、sessionStorage密钥清理或隐藏的密钥登录表单。

- 已登录管理员打开页面即可请求；页面可见性不替代服务端管理员鉴权。
- 缺少token或接口返回401：停止轮询，显示“请登录 Sub2API 管理面板，登录后刷新本页”。
- 接口返回403：显示“无管理员权限”，停止轮询；不把权限不足当作登录过期。
- 会话续期由原生管理面板负责；评估页不增加独立的登录、refresh-token流程或跨窗口消息同步。原生嵌入器现有的URL token参数不作为本页的鉴权来源。

原生证据：frontend/src/api/client.ts:44-46读取auth_token并发送Bearer；163起处理401刷新。frontend/src/stores/auth.ts:366保存token，475清理会话。独立HTML页没有使用原生Vue API客户端，因此不能宣称自动继承其刷新拦截器。
现状证据：internal/eval/web/app.js:25已经直接读取auth_token；63把401/403混为登录错误；276起仍有CPA key/sessionStorage残留；internal/eval/web/ui.html:17仍有隐藏密钥表单。

本项只审批鉴权来源、登录/权限错误展示及密钥残留删除。API路径、请求/响应字段转换在F03/F05审阅，不在此项决定。
验收行为：管理员登录后打开页面无需再次输入密钥；退出或收到401后停止请求并提示原生登录；403显示权限不足；本页不读写独立登录密钥或refresh token。

## 已审批：F03a 分组控件展示修正

用户回复“对 继续”，批准下述展示范围；产品代码尚未修改。

用户澄清：接口是对的，问题是F03a在当前展示中的前端设计；不据此扩展为整体前端重做。上一版孤立分组卡片方案被此展示决定替代。

源码确认的问题：

- internal/eval/web/ui.html:34把分组控件单独放在通知区后，使用field card占满一行，没有沿用原四种评估的参数工具栏布局。
- 该select缺少native-select包裹；internal/eval/web/style.css:38、45的36px高度、背景、边框、圆角、箭头和长文本处理只应用于native-select > select，因此分组没有得到与模型/推理强度相同的控件样式。
- field card不是toolbar容器，缺少工具栏已有的内边距及间距。

批准范围仅限展示：

- 移除当前孤立的整行分组卡片；分组作为每个评估页现有参数工具栏的第一项，与模型、推理强度等同排。
- 使用现有field + field-label + native-select结构，沿用既有高度、字体、边框、箭头和长文本截断。
- 四个页签共用同一个已选分组；切换页签不产生四份独立分组配置。
- 调整现有工具栏列数容纳新增字段，小屏沿用原两列折行，不另起一套控件或主题。

API能力已核实且用户明确接口正确：原生GET /api/v1/admin/groups/all支持platform=openai，admin_group.go:56使用ListActiveByPlatform。本项不再以接口选择作为审批问题。
本项不决定首次默认/记忆选择策略、账户分页、模型来源或切换组后任务行为。此前这些未批准的行为提案不作为本次展示修正的默认范围。
上述为源码核对及已审批规划，尚未修改产品代码，也未做浏览器视觉验证。

## 用户已确定：F03b 保留账户启停操作并接原生接口

用户拒绝“改成只读状态标签”的方案，明确要求通过额外原生接口保留操作。撤销只读替代方案，不再次询问是否保留开关；记录为规划决定，尚未实施。

核实到原生frontend/src/api/admin/accounts.ts:290的toggleStatus调用update(id, {status})，241使用PUT /admin/accounts/:id。backend/internal/handler/admin/account_handler.go:155接受active/inactive，1151的Update走原生账户更新服务。原生路由admin.go:379已注册PUT /accounts/:id，无需新增后端接口。

确定的接入范围：

- 保留账户行现有开关、列布局及“仅启用”筛选，移除switchable始终false的限制和要求去其他页面操作的占位提示。
- 打开调用 `PUT /api/v1/admin/accounts/<id>`，仅提交 `{"status":"active"}`；关闭仅提交 `{"status":"inactive"}`，使用原生数值account id和已批准的管理员会话。
- 保留原评估页已有的测试中锁定、提交中防重复及失败还原行为；成功采用服务端返回的账户状态更新四个页签的同一账户展示。
- 当前原生账户status与schedulable是不同字段。此开关现状以status是否active表示启停，因此对接toggleStatus的接口；不同时改schedulable，也不因启用成功清除限流等原生运行限制。

源代码事实：当前stateCredentials把status !== active映射为disabled，internal/eval/web/app.js:180；原插件enableSwitch在测试中/保存中禁用，internal/plugin/web/components.js:124。上述范围为恢复原有操作并用现成原生API接入，不引入新的管理功能或兼容翻译层。

## 已审批：F03c 模型下拉接入

用户回复“对”，批准下述范围；产品代码尚未修改。

保留原评估页模型下拉的布局、原有排序、各页签模型选择和记忆行为。本项只决定模型选项的读取，不重做控件，也不连带决定启动前检查或执行接口。

核实到两个原生接口的区别：

- GET /api/v1/admin/accounts/:id/models：account_handler.go:2794起调用FetchOpenAIAccountModels，返回账户模型对象；原生自身决定发现/缓存/失败时的目录处理。frontend无需重做这些服务端逻辑。
- GET /api/v1/admin/groups/:id/model-allowlist-candidates：admin_group.go:85起只把平台默认模型与可调度账户的model_mapping键合并，不调用上游模型发现。因此不把它冒充真实账户模型目录。

建议接入：读取当前分组各账户的原生models接口，仅取模型id做去重并集，供现有四个模型下拉使用，保留原有排序和针对各评估的显示规则。并集是一个模型在组内至少一个账户的返回目录中出现，不是全账户支持保证。此必要列表合并已单独获批，不默认保留旧catalog框架。

读取结果仅保留当前页面内存中的选项，不新增持久化目录、TTL/签名/owner/revision框架或旧字段兼容转换。刷新或切组时重新读取，切组后不得把旧组响应应用到新组。读取失败使用已有提示区显示对应错误，不吞成空模型列表；部分失败显示提示，不能把该账户认定为不支持所选模型。

原界面的模型偏好值仍可保存；“不持久化目录”不意味着删除用户选中的模型偏好。本项不审批删除任何原评估算法、账户支持检查或启动/取消行为。
证据现状：internal/eval/web/catalog.js:103起已有账户原生接口调用，但catch把失败吞成空列表；随后仍有目录持久化及启动前model_catalog代码。那些旧代码不是本项自动获批的范围。

## 已审批：F04a 恢复普通评估结果的原样式渲染

用户回复“对”，批准下述范围；产品代码尚未修改。

对照原插件文件：fingerprint.js没有改动；candy.js仅把markdown结果包进nonceStyledHTML；modeltrace.js仅把概率条等结果包进该函数。components.js新增该函数，扫描所有style属性，把它们改成动态CSS类；app.js:144也在普通行渲染时调用它。该机制不是原评估功能。

直接原因：adapter/runtime/backend/internal/handler/admin_codex_eval.go.in:214自行覆盖CSP，style-src仅允许nonce；原生backend/internal/config/config.go:35的默认style-src已有unsafe-inline。额外限制导致原公式渲染器和概率条的内联样式需要重写。

批准的最小修正范围：

- 普通结果使用原Markdown/Temml及概率条输出，保留原四类评估的内容、布局和交互。
- 评估壳页的style-src沿用原生默认支持的内联样式能力，移除额外style属性转动态class的函数及其调用；原脚本nonce、同源嵌入要求和管理员接口鉴权不属于本项删除范围。
- 本项只恢复普通页的公式/概率条等样式输出，不改变主题配色、题目、算法、字段、结果或历史布局。

不能在保留nonce-only style-src时直接删重写函数，那会破坏显示。因此该策略调整与普通结果渲染一起审批。此处只核实到原生默认策略；用户自定义CSP支持另属未审批范围，不引入多策略兼容框架。
鹈鹕sandboxHTML/脚本预览的处理单独作为F04b审阅，本项不删除其隔离或宣称预览已恢复。尚未改产品代码或做浏览器验证。

## 已审批：F04b 鹈鹕预览恢复原HTML输出

用户回复“对，接下来5个决策一批次展示”，批准下述F04b范围，并调整后续展示方式；产品代码尚未修改。

原实现证据：internal/plugin/plugin.go:215的评估页面专用CSP为script-src unsafe-inline与style-src unsafe-inline；internal/plugin/web/pelican.js:30直接把原HTML转义放进srcdoc，iframe只有sandbox=allow-scripts，没有allow-same-origin。该原实现允许模型动画及PREVIEW_FREEZE运行，隔离管理页面的DOM/存储。

新增修改：internal/eval/web/pelican.js:30起的sandboxHTML用正则处理HTML、替换nonce、把style属性转成class，再通过DOMParser序列化整份页面。这是为了适配迁移时额外引入的nonce-only script/style策略，并非原动画功能。

建议恢复：

- 删除新增sandboxHTML重写，原HTML直接用于原srcdoc预览；列表仅追加原有PREVIEW_FREEZE脚本，详情显示原HTML。
- 保留原sandbox=allow-scripts隔离，不授予allow-same-origin；继续沿用已有外部资源限制。
- 保留原列表加载约1秒后冻结动画/计时器、点击放大后播放及关闭后停止的行为。
- 仅评估壳页的script-src恢复原评估页允许运行内联脚本的策略，移除该指令的nonce-only限制，配合F04a的样式能力；原HTML及冻结脚本才不需nonce重写。原生Sub2API全站策略不修改，同源嵌入和管理员接口鉴权保留。

本项必须把删HTML重写与评估页脚本策略一起做；不能只删sandboxHTML而保留nonce-only script-src，否则模型脚本及冻结脚本会被阻止。前项F04a保留的脚本nonce策略在本项获批后，仅此评估壳页被原评估策略替代；不存在双策略兼容路径。
不引入额外预览接口、公开记录地址、URL token或HTML修复框架。此项已获规划批准，尚未浏览器验证或修改产品。

## 已审批批次：F05a—F05e 页面操作直连接口

用户在完整说明后回复“批准 继续”，批准本批5项；产品代码尚未修改。

用户要求一次展示5个决策，并补充“不要省略说明”。每项完整说明当前问题、具体修改、页面行为、影响范围与已知限制；用户该补充不是本批审批。本批保留原页面按钮、结果展示和交互，仅把调用点改成明确的新接口；不靠api函数猜测旧路径/改写旧body，不扩展后端安装或生命周期设计。
所有接口下述均相对 `/api/v1/admin/codex-eval`，使用F02的原生管理员会话。

### F05a 启动评估

保留原单账户及“测试所选/全部”操作。直接POST /start，提交type、targets数组（每项原生数值group_id/account_id）、model和对应原控件的effort/runs/mode/concurrency。参数来自原控件：糖果为model/effort/runs；指纹为model/mode/concurrency；ModelTrace为model/concurrency；鹈鹕为model/effort。各评估已有参数含义和默认值不由本项重设计。直接读取started/busy/unavailable数组显示启动结果及账户错误；不伪装成旧/run计数响应、不提交随后丢弃的model_catalog。
本项只审批页面调用，精确账户调度、账单、初始化就绪等后端行为仍另行核对；不据此移除原评估必要的账户支持检查。

### F05b 进度读取

直接GET /state，使用active里的任务id、type、target、status、progress更新原进度显示；以group_id/account_id/type匹配当前显示任务。保留原活动时2.5秒、空闲时20秒轮询和手动刷新；轮询只取状态，不反复取模型目录。状态/请求错误在现有提示区显示，保留完整返回错误文本，不自行归因为“不支持模型”。F02未登录/权限不足时停止轮询的决定继续适用。

### F05c 停止评估

保留原已有停止按钮，直接POST /cancel，提交按钮所属账户、分组和评估类型对应的run_id。响应cancelled只表示取消请求结果；显示正在取消并用后续状态确认终止，不直接宣称已停止。移除原“按四类running优先顺序猜任务”的做法。当前manager会取消run context，所以页面提示不再承诺已发请求一定继续完成；后端取消机制不在此项重新设计。

### F05d 历史详情

保留原展开历史和详情弹窗/动画预览操作。列表与历史概要由/state读取，详情直接GET /records/<type>/<account_id>/<record_id>，使用服务端记录标识；按接口定义的type及candy/fingerprint/modeltrace/pelican字段读取对应结果。使用新接口字段不等于旧CPA响应兼容转换。保留加载提示、失败原文和关闭/切换详情后忽略过时响应；不在通用api里解析旧/record?auth_id路径。

### F05e 清空历史

保留原按评估类型清空及“清空全部”的按钮和确认弹窗。确认后直接POST /clear-history，scope为candy/fingerprint/modeltrace/pelican/all。范围明确为相应类型在所有分组/账户的记录，沿用原全账户清空语义；不能让分组筛选误导为仅清空当前组。服务端成功后刷新；失败保留页面记录并显示返回错误。清空历史不取消运行任务，之后完成的结果仍会保存。此规划不授权实际删除用户记录。

### 本批源码证据与限制

- adapter/runtime/backend/internal/handler/admin_codex_eval.go.in:98/114/143/157/174对应State/Start/Cancel/Record/ClearHistory。
- internal/eval/runs.go:25声明StartRequest；52的StartSummary是数组；62的ManagerState包含active及各类型历史；527按类型/账户/记录id读取；566的ClearHistory清空类型的全账户映射而保留活动任务；638的Cancel调用run.cancel。
- internal/plugin/web/app.js:236是原2.5秒/20秒轮询；470起已有按类型/全部的清空弹窗。internal/eval/web/app.js:35-70则集中翻译旧路径/body/响应；这批把明确调用写到操作入口而不是保留该翻译层。
- 当前后端handler仍将部分错误压成start_failed/storage_unavailable等代码。本批只能展示返回文本；不能声称已经保留被后端丢弃的详细错误。完整错误输出要求仍保留，后端错误处理另行审批。

## 已审批批次：F06a—F06e 前端剩余行为

用户在完整说明后回复“继续 这不是能想明白吗 之前搞这么垃圾”，批准本批5项。仅记录规划，产品代码尚未修改。

### F06a 账户读取与真实信息

当前stateCredentials仅拿targets的简化账户字段，丢掉priority等原生信息，随后credentialView把缺失priority显示为P0。建议直接使用原生GET /api/v1/admin/accounts，查询platform=openai、group=<组id>、lite=1，按分页结果读齐当前组，保持原整组列表展示，不新增分页控件。使用原数值id、name、type、priority、status、schedulable及已有运行时错误字段；优先级按真实值显示，缺失则不伪造0。保留原名称、标签、脱敏、行展开和启停操作。运行/历史数据按账户与任务标识关联，不伪造成CPA credential或auth_file。只读取原生管理端已提供的展示数据，不在页面保存整份账户目录到持久存储。原生handler account_handler.go:643支持group/lite；dto/types.go:337的AccountListItem包含所需字段。

### F06b 类型与订阅筛选

当前credentialPlan只认source=auth_files/provider=codex，但新账户被标成source=accounts/provider=<type>，因此套餐徽章/筛选始终失效；providerType仍包含CPA配置提供者规则。建议保留原账户类型、订阅筛选、Pro ALL、订阅排序和徽章样式，直接以原生type驱动类型选项，以credentials.plan_type或原生parent_plan_type驱动已有套餐显示。后一字段是Sub2API原生影子账户关系，按其固定版本账户页规则读取，不是版本兼容机制。没有套餐数据时不猜测、不硬补Plus/Pro，也不按字段缺失删除类型/订阅功能。继续沿用原“仅启用”过滤，status与schedulable分别保留，不把开关开启当作没有限流或禁调度。证据：AccountsView.vue:1652的getAccountPlanType使用上述非Grok字段；dto/mappers.go已脱敏credentials。

### F06c 分组记忆与切换

F03a仅批准控件展示，初始选择/切换行为还未定。建议保留当前已有行为：记住最后所选组；仍存在于启用OpenAI组时恢复，否则选返回列表第一组；没有组时显示明确空状态并禁用启动。四个页签共用同一个组，切组清空四类账户勾选、旧账户/模型显示并忽略旧组未完成的读取响应，随后读取新组。各评估原模型偏好继续保留，只有模型实际出现在新组选项中时才恢复。切组只切换页面操作目标，不取消后台已启动任务、不清空历史。只保存组id和已有UI偏好，不恢复已移除的持久目录/owner/revision框架。证据：app.js:31/508及catalog.js:98已有组记忆、首组选择和清空勾选。

### F06d 移除额外激活卡片与手动前置步骤

现有ui.html:33插入额外激活卡片；app.js:328又把hostReady当作全页测试开关；原plugin页面没有这套手动激活UI。建议删除该卡片、按钮、页面hostState/hostReady状态机和启动前额外/status轮询，页面保留已审批的登录、分组、模型、账户与任务操作；初始化失败仍显示真实返回错误。
此项有明确依赖：当前后端Start会检查Ready和owner identity，直接删按钮会让未初始化安装无法启动。因此“删除手动前置步骤”必须与后端确保启动所需准备可由正常流程完成的具体方案配套；后端方案仍需下一轮单独审批。不得先删按钮、再默认调用原/activate冒充去掉激活，也不自动批准旧key/owner/repair框架。当前只审批页面去掉额外手动激活步骤，不声称生命周期问题已解决。

### F06e 最小主题接入

当前theme.css用Sub2API基础色重写整份主题，但把premium/plus/go/business全部映射为普通灰色，丢掉原套餐视觉区分。建议保持原评估布局、控件尺寸、图标、交互和套餐徽章差异，只对必要的页面背景/文字/边框/主色变量使用Sub2API原生light/dark颜色；原项目套餐强调色继续保留，不重做整页样式。
嵌入时读取原生已附加的theme参数，直接打开时读取同站点原生localStorage.theme；切换主题沿用CustomPageView现有更新iframe URL的行为，不修改原生Vue源码，不新增主题通信协议或CPA主题探测分支。没有显式主题时按原系统明暗选择。证据：原生tailwind.config.js的primary与dark颜色；main.ts:26/AppSidebar.vue:869保存读取theme；CustomPageView.vue:245/428构造URL并观察原生主题变化。此项不新增主题选项、不改变算法或额外装饰。

本批已获规划批准；源码核对不等于浏览器验证。原生账户接口与主题字段已确认，F06d的后端依赖仍需单独审批。

## 后续审批顺序

- F03b：保留账户启停并接原生status更新接口；用户已确定，等待整体实施阶段。
- 账户其余标签/筛选、分页及选择另行审阅。
- F03c：原模型下拉的原生目录接入；已审批，等待整体实施阶段。
- F04a：普通结果的原样式渲染；已审批，等待整体实施阶段。
- F04b：鹈鹕原HTML预览和对应评估页脚本策略；已审批，等待整体实施阶段。
- 四种评估其余结果、主题、账户信息仍逐项核对，不由F04a连带获批。
- F05a—F05e：5项页面调用；已审批，等待整体实施阶段。
- F06a—F06e：上述前端剩余5项；已审批，等待整体实施阶段。

前端F01—F06审批已记录。下一轮继续后端逐项审批，不从前端批准推导旧生命周期/身份框架获批。
