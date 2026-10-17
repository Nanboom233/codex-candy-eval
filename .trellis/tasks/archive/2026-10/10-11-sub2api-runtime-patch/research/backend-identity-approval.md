# 后端身份与初始化审批：B06—B10

状态：planning。用户在完整说明后回复“批准 继续”，批准B06—B10的规划范围。B01—B05已批准；本批不授权实施、真实模型请求、Key创建/停用或旧数据删除。固定v0.2.16 Linux/amd64、standard运行模式；不使用历史v0.2.15研究中的迁移/intent-first规则作为现行要求。

## 已核实的固定版本约束

- 原生ent/schema/usage_log.go:35要求api_key_id；ent/schema/api_key.go:39的key全局唯一，含已软删除行的唯一性仍存在。
- APIKeyService.Create:494复用用户、请求、IP、分组和创建数量/频率检查，公开expires_in_days只接受未来有效期。api_key_repo.go:45支持在一次INSERT中写入ExpiresAt，不能先创建永久有效Key再补过期。
- 原生api_key_auth.go:160/163检查用户组权限；:194及以后处理订阅、余额及配额；:246检查运行时expiry。simple模式在此前返回，用量/账单读取也跳过expiry，所以不能把过去expiry宣称为所有入口的绝对禁用。
- 当前评估已有不序列化网络凭据的process-local context capability，私有token用于文件/lookup派生而不是上游认证。全局身份及provisioning迁移并非这条内部调用的必需条件。
- RunManager.Start:315从m.lifetime派生后台任务context，不能依赖启动HTTP请求里的临时actor；必须显式保留任务执行身份。

## B06 按实际启动管理员保留一行原生记账数据，不再共享首次激活者身份

现状：全局EvalKeyIdentity绑定首次激活用户，所有评估都用它。若直接取消owner限制，会把其他管理员的请求记到该用户。原生usage外键和Key用量计数仍需要真实数据库行。
建议：每个实际发起过评估的管理员复用一行原生api_keys记账行，UserID为其原生用户ID。key使用固定保留标识__codex_eval_accounting_user_<user_id>，仅作为数据库唯一查询标识，不作为鉴权密钥；无新表、用户提供Key、随机私有token或可配置owner。Name初始为账户评估，可按原生功能改名而不妨碍执行。GroupID保持空，status初始active，ExpiresAt在首次INSERT中固定为2000-01-01T00:00:00Z，IP白名单初始仅127.0.0.1，配额/窗口初始采用原生0默认值。
创建使用小范围内部服务入口复用原生检查和插入，不复制整段Create，也不使用公开的未来expiry字段绕行；公开创建不能占用保留key命名空间。每人第一次准备占用一个原生未删除Key名额，达到原生数量/创建频率限制时返回实际错误，不额外放宽。
正常重启按固定标识读回同一行；无须文件保存user→key映射。唯一冲突只允许读回并验证该用户对应的行。其他owner、非固定expiry、非空持久group或已删除行都报告冲突，不接管、不改值、不新生成替代标识。初始创建完成但进程中断，下次使用现有行，不增加第二行。普通额度、限流或停用变更按原生规则生效，不再要求额度永远为0或名称不变。
可见行为：原生用量按真正启动者归属，API Key列表多一行已经过期的账户评估记账数据。全局历史/管理页面共享规则不变。无旧身份迁移，不接管/删除以前那一行。
限制：数据库保留标识不是秘密；B08必须阻止其作为网络凭据使用。原程序standard paid入口有过期拒绝，但原程序用量/账单读取及simple模式没有同样保证；回滚不等于撤销数据库行。不增加全站无Key计费重构。

## B07 首次有效启动准备记账，服务初始化走普通构造流程

现状：构造时Bootstrap失败关闭admitting，必须额外Activate才能启动；history/RunManager构造错误又被忽略或换成blocked通用错误。
建议：进程普通初始化时创建评估manager/历史存储、绑定既有路由，并保存实际初始化错误；不开启首次激活/repair状态机。通过原生管理员鉴权、请求参数和目标校验后，首次需要执行的/start准备B06那一行，再提交任务。准备失败不提交后台执行任务，按B05报告原错误；没有可执行目标不为它额外创建记账行。读取页面、state、历史不创建Key。对并发首次启动使用简单进程内同步和数据库原有唯一约束，唯一冲突读回核对，不新增队列、分布式锁、文件journal或repair协议。
行为：安装后登录原生后台即可直接开始，第一次自动完成必要准备，以后复用；某个用户的行异常只拒绝该用户启动，不锁全局服务。已有任务的启动/忙碌/失败摘要和取消语义保留。
删除：旧Activate/Bootstrap的首次使用流程、initializing/installed_uninitialized/repair_required框架；与B08一起删除private-token、provisioning intent、row_pending/migration/revoke状态机及旧格式迁移。不是偷偷调用旧Activate。
范围：尚未审批普通历史格式/启动恢复/存储位置的变化；保留当前正确的任务状态保存行为，具体另批审。初始化失败不能伪装成功；也不因评估初始化错误令原生其他服务一并无法启动。

## B08 内部请求携带任务身份，删除额外token和身份文件

现状：已经用process-bound context capability调用同进程/v1/responses，却仍维护私有token、派生lookup、全局identity及provisioning文件。
建议：任务只保存从原生认证主体取得的UserID、B06确定的APIKeyID、精确目标group_id/account_id和原任务参数。服务端在运行时生成不可由网络请求提供的context对象，携带当前任务身份进入既有HTTP Responses流水线。异步任务从manager生命周期派生context时显式携带身份，不把整个启动请求context保留到后台，也不在四类评估JSON参数中增加可由用户指定的actor。
每次模型执行读取数据库核对对应记账行和用户，使用任务绑定的那一行；不在执行时从全局“当前管理员”选身份。身份不明或行被删/改动则报实际错误，不创建替代行续跑。
公开网络请求不接受目标header/ticket/context代理。主API-key和Google公开凭据入口拒绝保留标识，必须在simple/usage等提前返回之前拒绝；公开创建不能制造保留标识。内部对已经核对的确切行只豁免固定过去expiry这一项，状态、用户、ACL、余额、配额等检查保留。不因为来自127.0.0.1就授予内部权限。
删除：私有token生成/文件读取、seed/lookup派生、全局单一identity、v1/v2迁移、旧proof/ticket和凭据隐藏兼容分支。没有token要传或隐藏，不新增外部控制协议。实际已退出的runtime不能继续生成或使用执行context；正常任务取消/退出接入另批审。
行为：两个管理员并发不会相互串账，关闭启动请求不改变已接纳任务的身份；后续评估请求仍受任务cancel控制。旧部署文件与旧Key不读取、不迁移、不自动删除。
范围：仅POST HTTP Responses执行身份；不新增评估WebSocket、Google、图像等执行功能。公开Google入口的标识拒绝只是封住已有鉴权入口，不是添加评估能力。

## B09 当前分组只在本次请求内绑定，完整复用原生计费检查

现状：内部记账行持久GroupID为空，需要在请求中装入所选分组；若直接放行或改变共享对象，会绕过用户权限或造成不同组之间串定价。
建议：新建每次请求的APIKey/Group副本，绑定当前明确group_id；不更新数据库GroupID或共享鉴权缓存。之后进入原生检查链，使用该任务启动者的实时用户状态、组权限、订阅/余额、Key额度/窗口、用户与账户并发、组倍率/渠道定价和用量记录路径。其账单更新B06那一行和该用户，AccountID仍必须是原始目标。
行为：后台能管理组不等于该管理员有该组订阅或余额。用户不具备当前组原生使用权限时，原样失败；不改记到另一个用户，不给管理员默认免单。不同组并发使用不同请求副本，不互相改绑定或定价。
范围：复用现有计费和权限能力，不另造免单/影子用户/免费Key、合成Usage ID、新收费表或独立估算账单。正常gateway部分用量、取消结算及精确路由钩子仍须下一批按实际代码审，不把所有旧billing/WS hooks连带批准。初始Key额外限额0不表示忽略用户或组限制；其后原生配置Key额度/限流也生效。

## B10 删除专用激活/停用维护协议，使用已有任务和原生Key管理

现状：除原评估start/cancel/state/history外，多出status/activate/deactivate/flush端点、维护Unix socket/client CLI，以及drain→revoke→删除身份文件状态机。安装/回滚还依赖这套协议，兼容问题多。
建议：保留已审批的任务/历史操作；删除专用激活状态接口及activate/deactivate/flush HTTP维护入口、维护socket、peer/ticket证明、CLI客户端和单例Key撤销状态机。内部正常保存/Flush能力不因删HTTP入口而删掉。单任务停止用/cancel，进程退出用正常任务取消和原生清理流程；具体等待策略另批审，不沿用无限重试。
管理员通过原生Key管理停用/限额限制自己的B06记账行，后台不得自动重启用、换行或恢复其字段。停用只按原生准入拦后续请求，不承诺立即中断已发送上游请求；要停当前任务用/cancel。软删除保留唯一标识，因此不能悄悄创建第二行绕过停用/删除意图；恢复使用需原生重新启用原行，已删除行不自动恢复。
行为：页面无需另一套激活/停用流程；二进制回滚或卸载不隐式删历史/记账行，不承诺撤销正在上游完成的请求。旧状态不迁移，普通历史保存不删除。
依赖：安装脚本现有维护调用必须在后续安装批次同步删除，不能先删server端留下installer调用。独立cleanup.py仍按用户要求打包，其范围是之前引入的构建资源，不扩展到删除数据库行或历史。不是批准今天真的停用Key、删文件或执行回滚。

## 下一批仍需审批

精确账户调度入口/失败不换账户/上游错误原文传递、必要HTTP取消和退出钩子、部分流用量结算；之后原评估核心差异/状态存储及固定预编译安装包与构建资源清理。当前既有结构修改不自动保留。

## 证据

原生路径均基于未修改v0.2.16提交01842576b3d5c1e2545b5a99ad555d36e28ed306：ent/schema/api_key.go:39/72；ent/schema/usage_log.go:35；internal/repository/api_key_repo.go:45/110；internal/service/api_key_service.go:443/494/882；internal/server/middleware/api_key_auth.go:160/179/194/246/426。
当前增加代码：adapter/runtime/backend/internal/service/eval_host_runtime.go.in:45/79/151/240/274；eval_key_runtime.go.in:56/198/233/365/428/485/540；eval_identity.go.in:15；eval_http_transport.go.in:54；internal/eval/runs.go:315。
本次仅源码核对，不是本方案已构建或运行成功的证据。
