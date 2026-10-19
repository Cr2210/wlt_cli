# 执行计划：superpowers 残留清理与 spec/文档迁移

> 顺序不可打乱：**先建任务、再搬运原文、然后写 spec、最后删除**。
> 删除是最后一步，保证任何中途失败都不会丢信息。

---

## Step 0 — 基线记录

- [x] 记录当前状态作为回滚参照：
      `git status --short`、`git log --oneline -3`、`python .trellis/scripts/task.py list`
- [x] 确认基线干净可构建：`go build ./...`（预期 exit 0）

## Step 1 — 建 5 个任务目录（全部 `--no-start`）

> ⚠️ **必须带 `--no-start`**：`task.py create` 默认会把新任务设为本会话 active task，
> 会覆盖本任务 `10-16-superpowers-to-trellis` 的 active 指针。

- [x] 4 个历史归档任务（`--parent 00-bootstrap-guidelines`，description 非空）：
  - `stateless-auth` — 无状态鉴权重构（2026-06-26 交付）
  - `stock-cmd-optimization` — 库存命令优化（2026-08-07 交付）
  - `order-prepayment-relation` — 订单预付关联（2026-08-21 交付）
  - `backend-sync-0923` — 后端 0923 基线同步（2026-09-28 交付）
- [x] 1 个开放任务 `query-endpoints-round2`（不加 `--parent`，独立待办）
- [x] 校验：`python .trellis/scripts/task.py list` 出现 5 个新任务，
      且 `task.py current --source` **仍指向** `10-16-superpowers-to-trellis`

## Step 2 — 搬运原文到归档任务（**必须在 Step 6 删除之前完成**）

- [x] 对 4 个归档任务，按 design.md §3 的映射表将原文**逐字**写入 `design.md` / `implement.md`：
  - `docs/superpowers/specs/2026-06-26-stateless-auth-design.md` → `stateless-auth/design.md`；
    其「代码改动/文档改动」段 → `stateless-auth/implement.md`
  - `specs/2026-08-07-...-design.md` → `stock-cmd-optimization/design.md`；
    `plans/2026-08-07-...md` → `stock-cmd-optimization/implement.md`
  - `specs/2026-08-21-...-design.md` → `order-prepayment-relation/design.md`；
    `plans/2026-08-21-...md` → `order-prepayment-relation/implement.md`
  - `specs/2026-09-28-...-design.md` → `backend-sync-0923/design.md`；
    `plans/2026-09-28-...md` → `backend-sync-0923/implement.md`
- [x] 每个文件头部加溯源块：原始路径 + 交付日期 + 对应 commit + 「历史归档，非待办」
- [x] 每个归档任务新写 `prd.md`：目标/范围/验收（从原文 Goal/目标/非目标段提炼，不臆造）
- [x] `query-endpoints-round2/prd.md`：第 2 轮端点清单逐字保留 + `-o` 方案诉求 + 2 处待决策
      + 第 1 轮已完成记录 + 未完成实测项
- [x] 用 `task.py set-meta` 为 4 个归档任务记录 `source_doc` 与 `original_date`，
      并据实改写 `task.json` 的 `createdAt` / `notes`
- [x] 校验：4 个归档任务的 design/implement 与原文字节数相当（无截断）；
      `python .trellis/scripts/task.py validate <task>` 通过（如适用）

## Step 3 — 建设新 spec 结构

- [x] 依据 `research/` 全部分片的锚点清单，创建：
      `.trellis/spec/cli/{index,directory-structure,command-factories,flags-and-params,output-protocol,backend-sync-playbook,quality-guidelines}.md`
      与 `.trellis/spec/internal/{index,auth-and-client,config,output-and-errors,quality-guidelines}.md`
      （分片导航见 `research/00-index.md`；§ 交叉引用按其中对照表定位）
- [x] **严禁照抄研究 §C 标记的漂移**，三条硬约束：不写"退出码 6 = 网络错误"（实际产出仅 1/2/4/5）、
      不写 legacy 工厂 `newCRUDSubCmd`/`crud_legacy.go`、不写"全仓分页子命令统一叫 `list`"。
      另：`CLAUDE.md` 的历史漂移（曾称 stock 用 legacy）**已自行修复**，不要再当作现存问题去"修正"
- [x] 两个 `index.md` 必须含 **Pre-Development Checklist** 与 **Quality Check** 两段 + 索引表
      （Status 列不得留 "To fill"）
- [x] 每条规则至少带一个真实 `path:line` 锚点；无法锚定的条目要么删掉，要么显式标注
      `UNVERIFIED`（不得伪装成事实）
- [x] 内容用中文；`guides/` 三个文件保持不动
- [x] 删除 `.trellis/spec/backend/`（5 文件）与 `.trellis/spec/frontend/`（6 文件）
- [x] 校验：`python .trellis/scripts/get_context.py --mode packages`
      输出 `Spec layers: cli, internal`（且不再出现 backend/frontend）

## Step 4 — 处理 00-bootstrap-guidelines 的悬空引用

- [x] 改写其 `task.json` 的 `relatedFiles`：`.trellis/spec/backend/`、`.trellis/spec/frontend/`
      → `.trellis/spec/cli/`、`.trellis/spec/internal/`
- [x] 改写其 `prd.md` 的「Spec files to populate」表：按 cli/internal 实际文件据实列出，
      删除 frontend 相关条目（不保留假勾选框）
- [x] 校验：该 task.json 是合法 JSON 且无指向不存在路径的引用

## Step 5 — 归档 4 个历史任务

- [x] 逐个执行（**必须带 `--skip-branch-validation`**：仓库有 remote 且 `base_branch=main`、
      `branch=null`，否则会被 branch 校验拒绝）：
      `python .trellis/scripts/task.py archive <name> --skip-branch-validation --no-commit`
- [x] 校验：4 个目录出现在 `.trellis/tasks/archive/2026-10/`（月份取自归档时刻，非交付日期）；
      `python .trellis/scripts/task.py list-archive` 列出 4 条

## Step 6 — 删除 superpowers 残留 + 移除 Cursor 平台（**不可逆的一步，放最后**）

- [x] 删除前最后一次确认 Step 2/Step 3 已完成（原文已归档、约定已入 spec）
- [x] 删除 `.superpowers/`、`.flow/`、`docs/`（`docs/` 删后为空，一并移除）
- [x] **R9 移除 Cursor**（顺序：先清哈希再删目录，避免中途失败留下"目录没了但清单还在"的
      半吊子状态）：
  - [x] 从 `.trellis/.template-hashes.json` 的 `hashes` 移除全部 52 个 `.cursor/*` 键
        （保留 `__version`、其余 `hashes` 键与 JSON 格式；用脚本改写并校验可被 `json.load` 读回）
  - [x] 删除 `.cursor/` 目录
  - [x] ⚠️ `.cursor/` 与 `.template-hashes.json` **均未被 git 跟踪，删除不可由 git 恢复**
        （恢复手段仅为重跑 `trellis init --cursor`）——执行前需已获得用户明确同意
- [x] 校验：`find . -iname "*superpower*" -not -path "./.git/*"` 无输出；
      `.superpowers`、`.flow`、`docs`、`.cursor` 四个路径均不存在；
      `trellis platforms` **仅**列出 Claude Code；
      哈希清单中 `.cursor/*` 键数为 0 且文件仍是合法 JSON

## Step 7 — 全量验收

- [x] 逐条走一遍 `prd.md` 的 Acceptance Criteria，**每条给出实际命令与输出**：
      `get_context.py --mode packages`、`task.py list`、`task.py list-archive`、
      `go build ./...`、`git status --short`
- [x] 抽查 spec 锚点真实性：每层抽 3 条 `path:line`，确认文件与行号对得上
- [x] 确认 `CLAUDE.md`、`skills/`、`cmd/**`、`internal/**` **零改动**：
      `git status --short` 中不得出现这些路径

## Step 8 — 收尾

- [x] 归档 `00-bootstrap-guidelines`（其交付物即本次 spec 填充）：
      `task.py archive 00-bootstrap-guidelines --skip-branch-validation --no-commit`
- [x] 更新本任务 `prd.md` 的验收勾选；写 `task.json` 的 `meta`
- [x] **不自动提交**。向用户汇报完整 `git status` + 改动摘要，由用户决定提交粒度与信息
      （注意 `.goreleaser.yml:92` 的 changelog 过滤器排除 `^docs:` 前缀）

---

## 验证命令汇总

```bash
python .trellis/scripts/get_context.py --mode packages    # 期望 Spec layers: cli, internal
python .trellis/scripts/task.py list                      # 期望 5 个新任务，无 4 个历史任务
python .trellis/scripts/task.py list-archive              # 期望 4 条历史归档
python .trellis/scripts/task.py current --source           # 期望仍是本任务
trellis platforms                                         # 期望仅 Claude Code（无 Cursor）
go build ./...                                            # 期望 exit 0
find . -iname "*superpower*" -not -path "./.git/*"        # 期望无输出
git status --short                                        # 检查零代码改动
```

## 风险点与回滚位置

| 风险 | 回滚 |
|------|------|
| Step 1 忘了 `--no-start`，active 指针被抢 | `task.py start 10-16-superpowers-to-trellis` 复位 |
| Step 2 原文搬运不完整就进入 Step 6 | Step 6 前设硬门槛；若已误删，`git checkout -- docs/` 恢复 |
| spec 锚点失真（写了臆造路径） | Step 7 抽查；发现即改，不允许带 UNVERIFIED 混入 |
| 归档被 branch 校验拒绝 | 补 `--skip-branch-validation` |
| 误改了源码 | `git checkout -- <path>`；Step 7 的 `git status` 检查是最后防线 |