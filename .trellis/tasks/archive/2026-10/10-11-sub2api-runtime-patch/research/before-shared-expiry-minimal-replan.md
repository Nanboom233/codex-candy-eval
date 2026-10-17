# 最小方案重规划

状态：planning，本计划待审。此前展开的B11—B15已被用户否决。产品代码和安装包尚未回退，不能把规划删项说成已删除。

## 目标与固定边界

在Sub2API v0.2.16 Linux/amd64保留原四类账户评估，使用原生后台登录、账户API、网关转发与计费。原生网关新增职责只有内部密钥识别和选择头硬定向。用户机器只替换预编译二进制，不拉源码、不编译。

原生基线提交01842576b3d5c1e2545b5a99ad555d36e28ed306已完成纯净编译/检查调研，结果见pristine-v0.2.16-baseline.md。复用该证据，不重新跑一遍纯净基线。既有测试失败/跳过仍如实保留；race不是验收要求。

## 一、先收敛剩余改动，再开始实施

前端F01—F06和后端B01—B10已经逐项批准，按最新方向调整B08鉴权形式。保留实际发起者的记账归属，不恢复首次激活者全局身份。

已确定退出新方案的内容：原生Vue菜单/路由补丁、多版本AST适配/兼容配置、激活与维护socket协议、proof/ticket/intent/private-token文件、WS评估钩子、重复资格重查、备用selector保护、定制sticky保护、原生错误旁路、partial用量补丁、上游cancel改写、自制HTTP管道/producer屏障及退出无限续期。

未逐项批准的core/history/build/installer清单继续按每批五项审阅。每项说明保留或删除的代码、必要消费者和变化后行为。原四类题库、算法、数据、原UI功能及已批准操作不默认改动；必要宿主接口替换单独列明。不要把“简化”解释成整体删除测试或用户历史。

这一步只更新规划与对照清单。剩余决策结束后提供完整最终实施摘要，获得新批准才task.py start。

## 二、后端只接内部密钥和选择头

### 内部密钥

首次有效/start通过原生管理员认证与目标校验后，沿B06取得该管理员的原生记账行。记账行仍是过去expiry、固定DB查询标识与实际UserID；标识不是秘密密钥。

用户补充要求防止内部密钥滥用。撤回按管理员长期复用密钥：拟每次本机HTTP调用生成一次性随机密钥，内存绑定实际UserID/APIKeyID、任务context、目标、最终body摘要、允许Codex headers与原请求deadline。验证真实socket peer、固定路径、匹配请求与有效任务后原子消费；取消/结束清理，普通重试每次重新生成。无文件、公开颁发API或新配置。完整机制与限制见internal-key-abuse.md，待审批。

评估调用发送Authorization: Bearer <内部密钥>。鉴权入口识别后读取绑定的原生Key/用户，核对对应行和实际用户；只豁免已经批准的固定过去expiry，状态、IP、订阅、余额、额度等后续检查仍走原生代码。停用/删除记账行时不能自动修复或换行。

普通Key走原路径；普通请求不能提交选择头取得定向权限。公开使用可推导的记账标识仍被拒绝。内部密钥只用于POST /v1/responses，认证后移除其凭据和选择头，再进入native流水线，避免传到上游。

### 一个选择头

拟采用X-Codex-Eval-Target: <group_id>:<account_id>。字段只来自已接纳任务，两个ID都是正整数；缺失、重复或格式错误直接拒绝。

内部鉴权验证通过后，在本次Key/Group副本绑定所选分组，复用B09原生组权限与计费。选择头解析出的ID只是普通请求目标信息，不再使用capability/generation/proof协议。

主调度入口为该请求选取唯一指定AccountID。复用原生selectBySessionHash的显式StickyAccountID、PreserveStickyBinding、DisableStickyEscape及其资格/DB核对/抢槽/WaitPlan；直接调用此原生选定账户路径后返回，不调用会继续候选池的完整Select。原生已有session/thread信号沿原评估协议传入。渠道定价等入口既有检查继续执行。

目标在native ExcludedIDs中或原生选择失败就返回native失败。原handler的同账户retry不会排除目标，仍能执行；跨账户失败会排除目标，下次只能得到无候选，不能转给B。无需在每个continue边界追加拦截。

原生quarantine策略、利润终检、同账户重试、session记录、protocol转换、错误、用量及连接断开处理均不另改。删除旧freshDirectedEvaluationAccount资格链和post-slot评估重查。

### 使用原生HTTP入口

删除evalServeHTTPTransport/io.Pipe，执行器使用本机原生监听端口POST /v1/responses。地址由已加载的原生Server配置构造，不采用客户端Host/Origin，也不新增用户配置。当前部署SERVER_HOST=0.0.0.0，使用127.0.0.1与原生端口，调用不经过环境HTTP代理、不跟随重定向。

这让原生HTTP server负责连接、断连和Shutdown。评估取消停止后续工作并关闭本次连接；原生可能继续读取上游完成计费，按原生语义报告，不承诺立即停止供应商。

主要修改位置：runtime密钥准备/执行器、API-key鉴权入口、主账户选择入口。server接线、管理员评估API与页面资源仍是评估功能必要的普通接入；“两件事”指网关扩展职责，不是声称整个移植只改两个文件。

## 三、恢复已审批的前端与原评估行为

删除8条原生Vue修改及菜单辅助模块。/codex-eval通过原生管理员custom menu接入，页面沿F02读取同源auth_token；内部密钥不交给浏览器。

按F03/F05/F06接原生分组、分页账户、状态开关和模型列表，以及评估start/state/cancel/records/clear-history。保留原工具栏、结果与Pelican渲染、过滤、轮询和最小主题映射；不新增按钮、catalog或第二套登录。

四类评估保留原题库、算法、数据和协议行为。B05只恢复评估自己吞掉的实际响应/error，不建立原生网关诊断旁路。core/history其他差异先完成第一步的逐项审批。

## 四、固定补丁，在开发环境构建

用纯净v0.2.16源码加一份固定版本普通补丁和必要评估资源，替换多版本AST规则、兼容manifest及分支。不全量git clone、不修改上游依赖锁，不让安装端参与构建。

按已经核实的原生工具版本与构建流程，在开发环境生成真实嵌入前端的Linux/amd64二进制。开发临时目录用标准tempfile/mktemp，最终仅输出dist/build/sub2api-codex-eval和dist/build/sub2api-codex-eval.tar.gz。

复用原生required checks与已有相关验证。重点验证内部密钥/选择头组合、A失败不发B、普通Key行为不变、实际启动者记账、原评估/UI操作和安装替换回滚。无race加码、付费在线评估或新的全量验收框架。运行fixture请求，与真实供应商请求分开说明。

## 五、安装与独立清理脚本打包

安装包包含预编译二进制、install.py、独立cleanup.py和必要版本/许可文件。install.py明确选择容器，备份、替换、重启并检查必要原生可用状态；不再调用激活/维护socket、拉源码、拉构建镜像或本机编译。

cleanup.py只针对旧评估构建引入的可确认资源，不执行全局prune，不删除其他项目镜像/缓存、业务volume、数据库、记账行或历史。菜单使用已批准的原生设置项，不修改Vue。实际外部设置写入与目标部署必须有明确授权和地址。

以最终包内脚本在可丢弃的Docker部署检查安装、回滚及清理边界。编译成功、测试成功、安装成功和实际功能可用分别报告；没有N150/iStoreOS实机验证就如实标注，不把开发机测试当作实机证明。

## 执行依赖与审批

保留同一Trellis任务：这五步共同交付一个安装包，不另起独立功能任务或工作流框架。

依赖顺序：剩余逐项审批→最小后端→前端/core接入→固定构建→安装包。实施获准后按文件所有权交给Trellis implement/check agents；规划阶段不分派实施。

本轮拟确定的新增技术方案是内存一次性且绑定实际请求的随机密钥、一个group/account选择头、原生本机HTTP调用与原生指定账户复用。它们替换旧B08/B11—B15方案，等待用户审阅。之后仍需完成未审批core/history/delivery清单，才给全任务最终实施审批。
