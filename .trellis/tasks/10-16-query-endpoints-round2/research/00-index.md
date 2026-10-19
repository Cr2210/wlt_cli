# 研究索引 — Round 2 端点盘点（2026-10-19）

> 方法：以后端仓库 `D:\Projects\weiliantong\weiliantong-api`（分支 `ge-dev-1.8`，2026-10-19 核对，
> HEAD `7011df50`）的 Controller 注解为**唯一事实源**，逐条核对 prd A 类清单与 wlt_cli 当前代码。
> ERP 模块 Controller 根目录：`wlclouds-module-erp/src/main/java/com/wlclouds/module/erp/controller/admin/`
> （下文锚点中的 `$C` 均指该目录；CLI 侧锚点相对 `D:\Projects\wlt_cli`）。

## 分片

| 文件 | 内容 |
|---|---|
| [01-endpoint-inventory.md](01-endpoint-inventory.md) | A 类端点逐条盘点：后端存在性 / CLI 现状 / 结论（待做·已完成·漂移） |
| [02-decisions-facts.md](02-decisions-facts.md) | 待决策 1（合同路径归属）与待决策 2（财务写命令）的事实与建议结论 |
| [03-export-inventory.md](03-export-inventory.md) | export 全景：后端 65 端点分布 + CLI 23 处 `fmt.Println` 临时形态，`-o` 方案设计输入 |

## 一句话结论

- prd 清单约 31 个端点中，**约 10 个已在后续开发中完成**（合同 4 + 运单 3 + 预付关联 3），实际待做约 20 个
- **待决策 1 已被合同域重构顺带解决**（CLI 已切独立端点），只剩 token 实测
- 发现 **1 处端点漂移**：CLI `finance invoice page-count` 调的 `/erp/finance-invoice/page-count` 后端已不存在
- `finance-record/init` 实为无参初始化**写操作**，不是查询，建议移出 A 类
