# Research: wlt 约定提炼（superpowers → Trellis spec 素材）

- **Query**: 从两份既有约定文档 + 8 篇 superpowers 文档 + 真实代码中，提炼可写入 `.trellis/spec/cli/**` 与 `.trellis/spec/internal/**` 的约定；每条规则带真实 `path:line` 锚点。
- **Scope**: internal（本仓库）/ 混合（含后端同步方法论，涉及外部后端仓库位置）
- **Date**: 2026-10-16
- **权威性**: 代码为准。文档与代码冲突处单列于 §C。

> 阅读约定：行号基于 2026-10-16 `main`（HEAD `812cddb`）。标识符/路径原样保留。
> 标注 **UNVERIFIED** 的条目表示代码中找不到证据，禁止在 spec 中断言。
> 源文档简称：`sync-design` = `docs/superpowers/specs/2026-09-28-backend-sync-0923-design.md`；`stock-design`/`stock-plan` = `docs/superpowers/{specs,plans}/2026-08-07-stock-cmd-optimization*.md`；`order-design`/`order-plan` = `2026-08-21-order-prepayment-relation*.md`；`auth-design` = `2026-06-26-stateless-auth-design.md`；`backlog` = `2026-06-27-query-endpoints-backlog.md`。

---


## 本目录文件导航（拆分说明）

原始单文件 `conventions-extraction.md`（553 行 / 57 KB）超过
`context_injection.max_file_bytes`（32768 字节），注入子代理时会被截断。
按「一主题一文件」拆分为下列 6 件，无内容删减：

     1. `01-auth-contract.md` — §1 无状态鉴权契约（含退出码 0–6 定义与产出点）
     2. `02-structure-and-factories.md` — §2 命令分包与注册约定 + §3 工厂选型判据
     3. `03-flags-and-output.md` — §4 flag 与传参映射 + §5 输出协议
     4. `04-sync-methodology-and-tests.md` — §6 后端基线同步方法论 + §7 测试与验证约定
     5. `05-spec-layout-and-drift.md` — A 建议 spec 布局 + B 源映射 + C 矛盾/漂移清单
     6. `06-backlog-facts.md` — D backlog 开放任务事实 + 锚点速查附录

交叉引用对照：`§1`→01，`§2/§3`→02，`§4/§5`→03，`§6/§7`→04，`A/B/C`→05，`D/附录`→06。

拆分执行：2026-10-16，任务 `10-16-superpowers-to-trellis`。
