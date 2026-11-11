# 云效 Flow 发布流水线接入指南

本文档说明如何把仓库根目录的 `flow.yml` 接入云效 Flow，实现「推 `v*` tag → 自动构建全平台产物 → 流水线运行详情下载」的发布自动化。

产物使用与 GitHub Actions（`.github/workflows/release.yml`）、本地 `make release` 完全相同的 `.goreleaser.yml` 构建，清单与二进制内嵌版本号逻辑完全一致。流水线只构建不发布：不访问 GitHub API、不需要任何 GitHub 凭据、不产生 GitHub Release，与 GitHub Actions 互不影响。

## 一、一次性接入步骤

### 1. 创建服务连接（已有可跳过）

1. 云效首页 → 左下角「企业设置」→「服务连接管理」
2. 新建服务连接，类型选择 **CodeUp**（本企业代码库），按向导完成授权
3. 复制该服务连接的 **ID**（第 3 步占位符要用）

### 2. 新建 YAML 化流水线

1. 进入「Flow 流水线」→「我的流水线」→「新建流水线」
2. 编排方式选择 **YAML 化编排**，模板任选一个 Go 模板（模板内容随后会被覆盖）
3. 进入 YAML 编辑页，把仓库根目录 `flow.yml` 的内容 **整体粘贴** 到左侧 YAML 编辑器

### 3. 替换两处占位符

| 占位符 | 所在字段 | 替换为 |
| --- | --- | --- |
| `https://<codeup 仓库地址>.git` | `sources.repo.endpoint` | CodeUp 仓库的 HTTPS 克隆地址（代码库首页「克隆/下载」复制） |
| `<服务连接ID>` | `sources.repo.certificate.serviceConnection` | 第 1 步创建/复制的服务连接 ID |

### 4. 校验、保存并验证

1. 点击右侧「校验」按钮确认语法通过（`GolangBuild.with.goVersion` 等参数取值以校验提示和右侧「文档」面板为准，见常见问题 Q1）
2. 点击「保存并运行」手工跑一次：可以验证代码源连接、工具链与 goreleaser 安装；若默认分支 HEAD 不在某个 tag 上，goreleaser 步骤会报错退出（属预期，原因见 Q5）
3. 完整链路验证：推一个测试 tag（如 `v0.0.0-rc1`），确认流水线自动触发并产出全部制品
4. 验证通过后，正式发布直接推 `v*` tag 即可，无需再进控制台

## 二、如何触发

```bash
git tag v0.2.0
git push origin v0.2.0
```

- 仅 **`v` 开头** 的 tag 触发（`branchesFilter: '^v'`，对齐 GitHub Actions 的 `on.push.tags: v*`）；push 到分支不触发发布流水线
- tag 触发的运行会自动 checkout 到该 tag，二进制内嵌版本号即真实 tag（如 `v0.2.0`）

## 三、产物在哪下载

流水线运行成功后：进入该次 **运行详情** →「构建产物」面板，下载制品 `wlt-dist`，共 7 个文件：

| 文件 | 说明 |
| --- | --- |
| `wlt-linux-amd64.tar.gz` | Linux x86-64 |
| `wlt-linux-arm64.tar.gz` | Linux ARM64 |
| `wlt-darwin-amd64.tar.gz` | macOS Intel |
| `wlt-darwin-arm64.tar.gz` | macOS Apple Silicon |
| `wlt-windows-amd64.zip` | Windows x86-64 |
| `checksums.txt` | SHA256 校验和 |
| `wlt-skills.zip` | AI Agent Skills 文档包（`.goreleaser.yml` before hook 打包） |

下载后可按原有流程手动上传到 CodeUp「标签 → 发行版」，或通过内部文件服务器分发。

## 四、三条发布路径的关系

| | 云效 Flow（`flow.yml`，推荐） | GitHub Actions（`.github/workflows/release.yml`） | 本地 `make release` |
| --- | --- | --- | --- |
| 触发方式 | 推 `v*` tag 到 CodeUp | 推 `v*` tag 到 GitHub | 本地手工执行 |
| 构建命令 | `goreleaser release --clean --skip=publish` | `goreleaser release --clean` | `goreleaser release --clean` |
| 发布目标 | 不发布，归档为流水线构建产物 | GitHub Releases | 本地 `dist/`，手动上传云效 |
| 需要 GitHub 凭据 | 否 | 是（GITHUB_TOKEN） | 是（发布到 GitHub 时） |
| 适用场景 | 团队内部分发（云效主导） | GitHub 开源分发 | 离线/调试 |

- 三条路径共用同一份 `.goreleaser.yml`，产物清单一致
- 云效 Flow 与 GitHub Actions 并存互不影响；本地构建保留为后备方案

## 五、常见问题

### Q1：GolangBuild 的 Go 版本（goVersion）怎么填？

`GolangBuild.with.goVersion` 的确切取值格式以 Flow 控制台「校验」提示和右侧「文档」面板为准。**即使集群预装 Go 版本偏旧也没关系**：脚本已设置 `GOTOOLCHAIN=auto`，且项目 `go.mod` 声明 `go 1.26.4`（≥1.21），Go 会按声明自动下载对应工具链；工具链与模块下载统一走 `GOPROXY=https://goproxy.cn,direct`（Flow 的 Go 环境默认 GOPROXY 即为 goproxy.cn）。

### Q2：为什么用 `--skip=publish`？为什么禁用 `--snapshot`？

- **`--skip=publish`**：本流水线只负责构建与归档。skip 后 goreleaser 不访问 GitHub API，因此不需要任何 GitHub 凭据，也不产生 GitHub Release。注意 goreleaser v2 的统一写法是 `--skip=publish`（v1 的 `--skip-publish` 已废弃）。
- **禁用 `--snapshot`**：snapshot 模式下 `.goreleaser.yml` 的 `version_template: "{{ incpatch .Version }}-snapshot"` 会把二进制内嵌版本号变成 `0.2.1-snapshot` 这类错位版本号。release 模式（tag 触发 + 全量克隆 + 工作区干净）才能内嵌真实 tag。

### Q3：goreleaser 想固定版本怎么办？

`flow.yml` 默认安装最新版：

```bash
go install github.com/goreleaser/goreleaser/v2@latest
```

固定到指定版本只需把 `@latest` 换成具体版本号（v2 的 module path 固定，`@latest` 不会跨大版本）：

```bash
go install github.com/goreleaser/goreleaser/v2@v2.12.5
```

> 具体版本号以 goreleaser 官方 releases 列表为准（示例仅供参考，照抄不存在的版本号会导致安装失败）。

### Q4：zip 安装失败怎么办？

`wlt-skills.zip` 由 `.goreleaser.yml` 的 before hook（`cd skills && zip -r ../wlt-skills.zip .`）打包，依赖容器内有 `zip` 命令。脚本已做防御性安装：缺失时执行 `apt-get update -y && apt-get install -y zip`（公共构建集群容器为 root，apt 源为阿里云镜像）。

若仍失败（如镜像源暂时不可用）：

1. 查看流水线运行日志中 apt 的具体报错，稍后重试一次运行
2. 治本方案：改用私有构建集群，自定义镜像预装 `zip`
3. 不推荐的做法：从归档清单里去掉 `wlt-skills.zip`（会导致 Skills 文档包缺失）

### Q5：手工运行失败、tag 触发才成功，正常吗？

goreleaser release 模式要求 HEAD 恰好位于某个 tag 上、工作区干净、git 历史含 tags：

- **tag 触发**（tagPush）：Flow 自动 checkout 到该 tag，且默认全量克隆（含全部 tag 历史），前提全部满足
- **手工运行**：checkout `branch` 指定的默认分支，若其 HEAD 不在 tag 上，goreleaser 会在版本解析阶段报错退出——属预期行为。完整验证请推一个测试 tag（如 `v0.0.0-rc1`）
