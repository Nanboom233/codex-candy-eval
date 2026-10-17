# 管理员测试通道重规划

状态：planning。用户已明确舍弃此前用户方案，采用管理员账户测试通道，并指定
POST /api/v1/admin/accounts/:id/codex-eval。具体S01—S05见admin-account-codex-eval-approval.md，已由用户“无误 继续”批准。C02/C03/C04/C06/C07已由用户“是 继续”批准，见core-admin-channel-approval.md；C08—C12已由用户“无误 继续”批准，见core-second-approval.md；C13—C17已由用户“是”批准，见core-third-approval.md；下一批C18—C22待审，见core-fourth-approval.md。

1. 完成具体接口及剩余core/history/delivery逐项审批。每批五项，完整说明当前代码、修改内容、消费者及变化后行为。保留已批准F01—F06，暂停旧C01—C05中的网关调用设计。用户身份/Key/过期/续期/计费归属问题退出当前规划。
2. 在现有管理员accounts路由旁新增独立codex-eval接口，基于/test发送代码接完整评估请求与响应。复制相关OpenAI发送逻辑，复用原生凭证、model mapping、代理、TLS与插件传输；不改原/test、不加选择头、不走用户账单或池内选号。详细流式/转换/状态/任务接入范围已在S01—S05批准，不暗中扩展。
3. 保留原四类评估题库、算法、数据及已批准的独立页面/管理操作。后台任务与新接口复用同一发送方法，不另造内部HTTP、JWT或退出框架。历史和核心其他差异继续单独审批。
4. 用固定纯净v0.2.16提交01842576b3d5c1e2545b5a99ad555d36e28ed306加普通补丁，在开发环境编译真实嵌入前端的Linux/amd64二进制。保留原生工具版本和锁，使用标准临时目录；不全量git clone，不在用户机器构建。输出dist/build/sub2api-codex-eval。
5. 打包预编译二进制、install.py、独立cleanup.py及必要元数据/许可到dist/build/sub2api-codex-eval.tar.gz。可丢弃环境验证替换/回滚/清理和真实fixture评估功能；不用安装成功冒充功能成功，不执行全局prune或真实业务数据清理。

退出的旧产品机制在最终实施批准后才删除，规划更新不等于产品已回退。恢复原生API-key鉴权、scheduler、sticky、计费/parser、main/shutdown；不增加race、付费评估或新的全量验收框架。
完成纯净基线证据已保留，实际失败/跳过照实报告；不重复基线调研，不把开发机检查当作N150/iStoreOS实机验证。
