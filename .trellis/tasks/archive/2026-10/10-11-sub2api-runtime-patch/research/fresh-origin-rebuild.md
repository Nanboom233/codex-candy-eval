# 从原始版本重做

用户明确要求从原始版本开始，停止修补旧移植；此前“继续，如此这类剔除冗余逻辑的无需额外提问”批准C18—C22及后续同类冗余剔除。已批准功能仍有效。

最新决定：直接使用已完成原始编译的Sub2API v0.2.16源码，以及当前Git仓库的internal/plugin。无需固定新来源、导出指定插件提交、再下载/拉Git或重复校验归档/文件身份。fresh-origin-source-manifest.json仅保留为历史调查证据，不是施工输入或验收门槛。

从原插件移植四评估题库/算法/数据/原请求/原独立页面到新后端包；新增已批准原生管理员/api/v1/admin/accounts/:id/codex-eval，与后台评估共用原生发送。只替换必需CPA宿主接口，保留原布局/renderers/Pelican/历史/进度及已批准原生操作。原简单内存任务和四类历史，无新增任务账本/Key/计费/恢复/激活/维护。

旧internal/eval、adapter/runtime/rules、compatibility、旧patch/scripts/install.py/cleanup.py只作过去错误对照，不作新代码输入。旧逐项清单不再是施工单。使用普通补丁、简单构建/打包及二进制替换install.py；独立cleanup已取消，移除其交付、打包和测试要求，不执行真实清理。不继承AST/type-export/profiles/hash-gates/overlays/维护socket/部署状态机。

输出仍dist/build/sub2api-codex-eval及同名.tar.gz，Linux/amd64原生嵌入前端，用户端不编译。复用已完成原始编译的前端输出和缓存，前端构建输入未改不重建。保留明确容器选择/备份/回滚/进度、必要许可；不依赖docker --format、nc/buildx/工具镜像。

原始编译和测试沿用已完成报告，保留真实失败/跳过说明，不重跑基线。检查仅针对实际改动，优先复用现有测试；未覆盖的原生发送/页面对接与替换失败回滚做必要本地验证。同一证据不重复验收，不新增race、覆盖率、全量套件、架构矩阵、资源校验或独立cleanup测试。历史研究和旧spec的额外检查不成为当前task要求。

工作区、旧canonical产物、真实业务/历史资源当前保留。旧规划已保存before-fresh-origin-*，整体替换旧三个规划文件时按本地钩子先移入macOS废纸篓，副本完整可恢复；没有git reset、产品删除、重编译或部署。准确新文件映射见fresh-origin-port-map.md，当前需求/设计/执行清单见../prd.md、../design.md、../implement.md。

不再询问逐项冗余删除。仅真正未解决功能取舍需提问；用户已批准执行，task为in_progress；按Trellis实现代理、检查代理、spec更新、提交及收尾顺序进行。
