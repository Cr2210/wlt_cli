# 清理 superpowers 残留并将文档迁移至 Trellis 体系

## Goal

项目从 superpowers 工作流切换到 Trellis 后，`docs/superpowers/` 下遗留 8 篇文档（3 篇实现计划 + 5 篇设计/交接），
`.superpowers/`、`.flow/` 留有工作流残留目录。同时 `.trellis/spec/` 是 `trellis init` 生成的
**fullstack 通用占位模板**（`backend/` 5 个 To fill 文件 + `frontend/` 6 个前端专用文件），对本项目
（Go 编写的 CLI，无前端、无数据库 ORM）完全不适配。

本任务把 superpowers 残留清除，并将其中的持久知识提炼进改造后的 `.trellis/spec/`，
历史交付按 Trellis 任务形态归档，未完成的待办转为开放任务。

价值：`.trellis/spec/` 是 trellis-implement / trellis-check 子代理的唯一上下文来源，
填错或留空会导致子代理写出不符合项目约定的代码；把真实约定沉淀进去是一次性投入、长期收益。

---

## Background（已核实事实）

### superpowers 残留分布

| 位置 | 内容 | git 跟踪 |
|------|------|---------|
| `.superpowers/sdd/.gitignore` | 仅一行 `*` 的空壳目录 | 未跟踪 |
| `.flow/` | 空目录（无任何文件） | 未跟踪（空目录 git 本就不跟踪） |
| `docs/superpowers/plans/*.md` | 3 篇实现计划（约 118 KB） | 已跟踪 |
| `docs/superpowers/specs/*.md` | 5 篇设计/交接文档（约 37 KB） | 已跟踪 |

全仓 grep 确认：除上述文件**自身**外，没有任何其他文件引用 superpowers
（`CLAUDE.md`、`skills/`、`.claude/`、`.cursor/`、代码全部干净）。即清理是纯删除，无引用需改写。

### 8 篇文档对应的实现状态（代码证据，非文档自述）

5 篇设计文档均标注"状态:已确认,待实现"，但经代码核实**全部已实现**：

| 文档 | 实现证据 |
|------|---------|
| `specs/2026-06-26-stateless-auth-design.md` | `internal/auth/`、`cmd/auth/` 已删除；`internal/cmdutil/client.go` 含 `AuthFlags`/`RequestContext`；`internal/client/client.go` 已用 `RequestContext` |
| `specs/2026-08-07-stock-cmd-optimization-design.md` | `internal/cmdutil/crud.go:856` 有 `DocumentConfig`、`:868` 有 `DocumentCmds`；`crud_legacy.go` 已删除；commit `ab58043`…`33b4d10` |
| `specs/2026-08-21-order-prepayment-relation-design.md` | `cmd/order/order_prepayment.go` + `order_prepayment_test.go` 存在；commit `af07f87`/`bb1a8df` |
| `specs/2026-09-28-backend-sync-0923-design.md` | `cmd/settlement/logistics.go`、`cmd/sale/accounting.go`、`cmd/finance/freight_apply.go` 存在；commit `4888596`/`13e7f7e`/`67aa36a`/`40de0c9` |
| `plans/2026-06-27-*`（对应 Round 1 查询接口对接） | commit `b69f448`（见文档自述） |

`go build ./...` 当前 exit 0，仓库处于干净可构建状态。

### Trellis 工具链对 spec 目录的要求（影响改造方案）

- `.trellis/scripts/common/packages_context.py:30-40` 的 `_scan_spec_layers()`：
  扫描 `.trellis/spec/` 的**子目录**作为 layer，**显式排除 `guides`**。
  → 目录改名为 `cli/` + `internal/` 后，`get_context.py --mode packages` 会自动识别为新层，无需改脚本。
- 同文件 `:198` 期望每个 layer 下有 `index.md`（输出 `- .trellis/spec/{pkg}/{layer}/index.md`）。
- 因此**必须**保证 `spec/cli/index.md`、`spec/internal/index.md` 存在。
- `spec/guides/`（通用思维指南，已填充）不在改造范围内，保留。

### 待办（非历史，需保留）

`specs/2026-06-27-query-endpoints-backlog.md` 含第 2 轮未完成待办：
约 31 个 JSON 查询端点（invoice-item / position / screen 6 / system-account 9 / system-weight-account /
waybill 5 / service-contract / transport-contract / finance-record）、约 40 个 export-excel 端点
（阻塞于尚未设计的 `-o` 下载方案）、2 处待决策（service-contract/transport-contract 的 page 路径归属、
本轮新财务模块的写命令是否补全）。另记录了一处未完成验证：Round 1 新增约 36 个命令因 token 过期未实测返回数据。

---

## Requirements

- **R1 清除残留**：删除 `.superpowers/`、`.flow/`、`docs/superpowers/` 三处。
  删除前确保 R3/R4/R5 的知识提取与归档已完成（"先吸收再删"）。
- **R2 spec 目录改造**：`.trellis/spec/` 由 `backend/` + `frontend/` 占位层改造为
  `cli/`（命令层：`cmd/` 下各业务域）+ `internal/`（共享内部包）两层，各含 `index.md`
  （带 Pre-Development Checklist 与 Quality Check 两段 + 索引表）。
  `frontend/` 的 6 个前端专用文件全部删除；`backend/` 的 5 个通用占位文件（database/logging 等
  对本项目无意义）删除后按需重建为真实文件。
- **R3 约定提炼**：把 8 篇文档中的持久约定写入 `spec/cli/` 与 `spec/internal/`，至少覆盖：
  1. 无状态鉴权契约（`--profile` 提供什么、`--token`/`--tenant-id` 必填边界、`--base-url` 覆盖、
     缺失→退出码 4、服务端 401→退出码 5 且 message 含 `code=401`）。
     **退出码按代码现实书写**：常量 0–6 已定义，但**实际产出仅 1/2/4/5**
     （3 `authentication`、6 `network` 为保留值，全仓无产出路径；网络错误实际落 5）。
     这一条是对 `skills/SKILL.md` 与 `skills/references/api-conventions.md` 既有错误描述的**纠正**，
     spec 不得照抄"6 = 网络错误"。
  2. 命令分包与注册约定（`Register(parent)`、`init()` 注册、新增业务域步骤）
  3. 工厂选型决策（`DocumentCmds` / `AddCRUDToParent` / `Crud*Cmd` / 手写的适用条件与判据）
  4. flag 与传参映射（kebab-case flag → camelCase 参数、时间范围 → `{timeKey}[0]/[1]` 数组、
     分页 `--page-no/--page-size`、`--headers` 导出、`--data` JSON 直传、小固定请求体用离散 flag）
  5. 输出协议（stdout 数据 JSON / stderr 错误 JSON / 退出码 0-6；
     `ParsePagedJSON` vs `OutputJSON` 的选择，含"裸数组响应不可用 ParsePagedJSON"）
  6. 后端基线同步方法论（差异分析方法、A/B/排除项分类、"CLI 是薄透传层"原则）
  7. 测试与验证约定（工厂传参的表驱动单测、flag 装配测试、`go build`/`go test` 验证要求）
  每条规则必须引用仓库真实文件路径（遵守 bootstrap prd 的 "Document reality, not ideals"）。
- **R4 历史归档**：4 个已完成交付各建一个归档任务（`archive/{年-月}/`，含 `task.json` +
  `prd.md` + `design.md` + `implement.md`）：
  1. `2026-06-26` 无状态鉴权重构（仅设计文档，无独立计划）
  2. `2026-08-07` 库存命令优化（设计 + 计划）
  3. `2026-08-21` 订单预付关联（设计 + 计划）
  4. `2026-09-28` 后端 0923 基线同步（设计 + 计划）
  原文按 trellis 三件套映射存放（设计→`design.md`，计划→`implement.md`，
  目标/范围→`prd.md`），保留原文内容不重写。
- **R5 待办转开放任务**：`2026-06-27-query-endpoints-backlog.md` 转为 1 个开放 Trellis 任务，
  其 `prd.md` 保留第 2 轮完整端点清单（约 31 JSON + 约 40 export-excel）、两处待决策、
  `-o` 下载方案的设计诉求，以及"第 1 轮已完成"的实现记录与未完成的实测项作为背景。
  该任务状态保持 `planning`，本任务**不实施**其内容。
- **R6 工具链兼容**：改造后 `get_context.py --mode packages` 列出 `cli, internal`；
  `task.py list` / `list-archive` / `validate` 正常工作。
- **R6b 记录真实的不一致，不粉饰**：研究已证实现实中存在若干并非"整齐统一"的约定，
  spec 必须**如实记录**而非按理想改写（bootstrap prd 的 "Document reality, not ideals"）：
  1. 分页子命令主名并非全仓统一 `list`：`finance`(6)/`produce`(2)/`quality`(1)/`waybill`(1)
     仍以 `page` 为主名；`page` 仅在 stock 四个单据上是 `list` 的别名
  2. `DocumentCmds` vs `AddCRUDToParent` **无强制规则**，域间并存且不一致
     （`product`/`stock_warehouse` 用后者，`product_unit`/`product_metrics` 用 `NewCRUDGroup`）
  3. `--headers` 只部分命令支持；`settlement main export` 无、logistics/销售核算/运费申请有；
     `DocumentCmds` 的 page-count 刻意不收集
  4. export 类命令输出**不是 JSON 信封**（21 处 `fmt.Println`），是待 `-o` 方案统一的已知例外
  5. 时间范围有两套并存实现：通用 `AddTimeRangeFlags/CollectTimeRangeFlags` 与
     `crud.go` 内 6 组域内 helper；`stock query/record` 的 `in-time`/`create-time` 是**单值**，
     与 `stock in/out` 的 `{timeKey}[0]/[1]` 数组契约不同
- **R7 文档语言**：spec 与归档产物用**中文**编写，与项目既有 `CLAUDE.md`、`skills/`、
  历史文档一致。（`trellis init` 生成的 `spec/backend/index.md` 尾部有 "All documentation should
  be written in English"，但该行是模板默认值，与项目现实冲突；bootstrap prd 自身要求
  "Document reality, not ideals" 并"优先从现有约定文档导入"，因此以中文为准。）
- **R8 范围外**：不修改 `CLAUDE.md`、`skills/`、任何 `cmd/**`、`internal/**` Go 代码。
  研究已发现 `skills/` 存在两处事实性漂移（退出码 6 从未产出；`api-conventions.md:62-71`
  仍在描述已删除的 `newCRUDSubCmd`），**本轮只记录不修改**，作为后续独立任务的依据。
- **R9 移除 Cursor 平台配置（用户新增要求）**：项目只需要 Claude Code，移除 Cursor 相关的
  全部跟踪。**不能只删目录**：`trellis` 判定"已配置平台"的唯一依据是
  `.trellis/.template-hashes.json`（`configurators/index.js:103-133` 的 `getConfiguredPlatforms()`
  遍历 `PLATFORM_IDS`，检查该平台 `configDir` 下是否有被哈希跟踪的模板文件），
  且 `trellis update` 正是遍历该结果来管理平台文件（`commands/update.js:435`）。
  因此必须：
  1. 删除 `.cursor/` 目录（52 个文件，**未跟踪**，删后不可由 git 恢复，但可由
     `trellis init --cursor` 重新生成）
  2. 从 `.trellis/.template-hashes.json` 的 `hashes` 中移除全部 52 个 `.cursor/*` 键
     （trellis 自带 `pruneOrphanManifestKeys` 工具，说明清理悬空清单键是其认可的操作）
  3. 验证 `trellis platforms` 仅列出 Claude Code
  `AGENTS.md` **保留**：它是跨工具共享的入口说明（非 Cursor 专属），且不影响平台判定。

---

## Acceptance Criteria

- [x] `find . -iname "*superpower*" -not -path "./.git/*"` **除本任务自身目录外**无输出；
      `.superpowers/`、`.flow/`、`docs/superpowers/` 三个路径均不存在
      （判据修正说明：本任务目录名 `10-16-superpowers-to-trellis` 本身含 "superpowers"，
      故该 find 命令**必然**命中它一次；判据的本意是"无 superpowers 残留物"，据此措辞修正。
      已核实：命中项只有本任务目录，无任何 superpowers 工件。）
- [x] `python .trellis/scripts/get_context.py --mode packages` 输出 `Spec layers: cli, internal`，
      **不再**出现 backend / frontend
- [x] `.trellis/spec/cli/index.md` 与 `.trellis/spec/internal/index.md` 存在，
      均含 Pre-Development Checklist 与 Quality Check 段落，索引表 Status 列无 "To fill"
- [x] R3 列出的 7 类约定在 spec 中各有落点，且每条规则至少引用一个真实存在的仓库路径
      （抽查：鉴权契约引 `internal/cmdutil/client.go`；工厂选型引 `internal/cmdutil/crud.go`；
      输出协议引 `internal/cmdutil/output.go` 与 `internal/output/`）
- [x] spec 中**不出现**以下两条已知错误：不得写"退出码 6 = 网络错误"（无产出路径），
      不得写 legacy 工厂 `NewCRUDSubCmd`/`crud_legacy.go`（已删除）
- [x] R6b 的 5 处不一致在 spec 中均有如实记载（不写成"全仓统一"的理想描述）
- [x] `python .trellis/scripts/task.py list-archive` 能列出 **4 个**归档的历史交付任务
- [x] `python .trellis/scripts/task.py list` 能列出 backlog 转换出的开放任务，
      其 prd 含第 2 轮完整端点清单 + export-excel `-o` 方案 + 2 处待决策 + Round 1 未实测项
- [x] `go build ./...` 仍 exit 0（本任务不动代码，作为回归确认）
- [x] `trellis platforms` **仅**列出 Claude Code（不再出现 Cursor）；
      `.trellis/.template-hashes.json` 中 `.cursor/*` 键数为 0，且 `.cursor/` 不存在
- [x] `.trellis/spec/guides/` 未被破坏（其 index.md 与 2 个 guide 文件仍在）
- [x] `00-bootstrap-guidelines` 父任务的校验项按实际完成情况更新（`backend/frontend` 相关项
      因结构变更而失效，需据实改写而非留假勾）

---

## Out of Scope

- `CLAUDE.md` 与 `skills/` 的内容整合/去重（用户已决定本轮不动；spec 与二者会有部分重复，
  留待后续单独任务收敛）。**含研究发现的 `skills/` 两处事实性漂移**
  （`SKILL.md`/`api-conventions.md` 的"退出码 6 = 网络错误"、`api-conventions.md:62-71` 的
  已删除 legacy 工厂描述）——本轮只记录进 spec 与研究报告，不在本轮修 `skills/`
- 第 2 轮后端端点的**实际对接实现**（仅转成开放任务，不在本任务实施）
- 修改任何 Go 代码或命令行为
- 安装 Trellis statusLine（`trellis init` 时默认未装）

---

## Key Decisions

| 决定 | 选择 | 理由 |
|------|------|------|
| 归档粒度 | 每个交付一个归档任务（4 个） | 用户拍板；与 Trellis"一个交付 = 一个 task"惯例一致，`task.py list-archive` 可按月检索 |
| 旧文档处置 | 先吸收再删 | 用户拍板；持久约定进 spec，原文进归档任务，删除后才不丢信息 |
| spec 目录结构 | `cli/` + `internal/` 两层 | 用户拍板；替换 fullstack 占位层，且被 `_scan_spec_layers` 自动识别 |
| `CLAUDE.md` / `skills/` | 本轮不动 | 用户拍板；与 spec 的重复留待后续任务收敛 |
| spec 语言 | 中文 | 项目现实（CLAUDE.md/skills/docs/commit 全中文），模板的 English 行是默认值而非项目决定 |

无遗留阻塞性开放问题。