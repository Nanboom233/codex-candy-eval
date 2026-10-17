> 已被用户最新“内部密钥＋选择头，其余沿原生”指令取代，非当前方案、未获批准。当前范围见 [backend-minimal-key-header.md](backend-minimal-key-header.md)。

# B11—B15 当前待审批方案

状态：planning，B01—B10 已批准；B11—B15 均未批准，产品实施仍暂停。

用户要求：“B11-15均需要补充调研 需要完整说明怎么改”。补充调研已完成，完整修改顺序、原生函数定位、旧规则保留/删除对照与实际限制见 [backend-dispatch-supplemental.md](backend-dispatch-supplemental.md)。该文是本批唯一完整方案，替代本文件旧版。

## 五项决策

- B11：在主调度 Once 的原生 quota/privacy/profit context 准备后、所有候选路径前定向；读取指定账户最新 DB 行并复用原生资格谓词，构造 selection 时不换回 snapshot。取得原生槽位后再核对，同账户失败释放一次。删除实际 HTTP 路径不可达的备用 selector、WS 与另一 GatewayService 修改。
- B12：continuation、profit-veto、跨账户 failover 边界结束定向任务；原生同账户 retry 在跨账户禁止分支之前执行。只留 post-admission sticky 方法的一处评估保护，保留 response_id 归属、session/thread 与原生 reasoning 状态。
- B13：本次 Execute 内部 context 保存原始 transport cause、完整 HTTP error body/read error、SSE payload 与 CC error/解析失败；在失败详情丢弃前捕获。尝试开始重置，成功恢复不携带旧失败；handler 完成后合并到原有错误结果。普通客户端错误策略不改。
- B14：只改 Responses、OpenAI passthrough、同账户 CC 三处上游 context。正常 terminal 后读完内部管道并等 handler done，再关闭 body；手动 cancel立即取消生成，清理等待不超过原请求绝对期限。HTTP stream partial 与实际 observed usage 返回 native handler，各实际 attempt只提交一次，保留原生异步结算/overflow策略；没有额外队列或费用估算。
- B15：恢复 ListenAndServe，信号时立即停止接收评估并 cancel，runtime 与 HTTP server并行共用原一次五秒 context。等 manager 与实际 producers，超时记录未完成后进入原生 cleanup，不无限续期。

## 必须纠正的旧判断

OpenAI 平台不会进入 Responses→Anthropic；原生 HTTP response归属使用 WS 命名存储，不能删掉；用量队列允许配置丢弃；native cleanup 的十秒 context不强制终止各Stop函数。不得把本方案描述为全部费用必定落库或整个进程必定五秒退出。

## 审批边界

本批仍需用户审批，不是所有现存native/core/tooling修改的批准，更不是最终实施批准。之后继续原四类评估/core、history/restart、固定构建交付与独立cleanup脚本审阅。本轮没有运行测试、真实上游调用、部署或用户数据删除，也未追加race/全量测试要求。
