# 云效 Flow 发布流水线(YAML 化,对齐 GitHub Actions release)

## Goal

提供一套云效 Flow YAML 流水线:`v*` tag 推送到 Codeup 后自动触发,用 GoReleaser 构建与 GitHub Actions 完全一致的全平台产物,归档到流水线运行产物(公共存储空间)供下载;仓库保留 `flow.yml` 作为唯一事实源,并附 Flow 控制台一次性接入说明。

## Requirements

- **R1 触发**:codeup 代码源 `tagPush` 事件触发;`branchesFilter: ^v` 仅 `v` 开头的 tag 触发(对齐 GitHub Actions `on.push.tags: v*`)。push 到分支不触发发布流水线。
- **R2 构建一致性**:使用现有 `.goreleaser.yml`(不改动),产物与 CLAUDE.md 发布流程列出的清单一致:`wlt-{linux,darwin}-{amd64,arm64}.tar.gz`、`wlt-windows-amd64.zip`、`checksums.txt`、`wlt-skills.zip`。二进制内嵌版本号必须是真实 tag(如 `v0.2.0`),因此用 `goreleaser release --clean --skip=publish`,禁用 snapshot 模式。
- **R3 无 GitHub 依赖**:流水线不需要 GITHUB_TOKEN,不产生 GitHub Release、不影响 `.github/workflows/release.yml`(两套并存)。
- **R4 构建环境策略**(云效公共构建集群,国内网络):
  - Go 模块/工具链代理用 goproxy.cn;
  - goreleaser 用 `go install github.com/goreleaser/goreleaser/v2@latest` 安装(不直连 GitHub 下载二进制);
  - `zip` 命令缺失时防御性 apt 安装(before hook 依赖它打 wlt-skills.zip)。
- **R5 产物上传**:`ArtifactUpload` + `uploadType: flowPublic`(归档云效公共存储空间,运行详情页可下载),路径含 `dist/` 与 `wlt-skills.zip`。
- **R6 仓库交付物**:
  - 根目录 `flow.yml`:完整流水线 YAML,头部中文注释说明用途、触发方式、接入方式、需替换的占位符(endpoint / serviceConnection id);
  - 接入文档(新增 `docs/yunxiao-release.md`):控制台一次性接入步骤、产物下载位置、与 GitHub Actions/本地 make release 的关系、常见问题(Go 版本工具链、zip、goreleaser 固定版本);
  - `CLAUDE.md` 发布流程章节补一小节:tag 推送后云效流水线自动构建,产物从流水线运行详情下载(替代手动 `make snapshot` + 上传,本地构建保留为后备)。

## Acceptance Criteria

- [ ] `flow.yml` 字段均符合云效 YAML 语法且有官方文档依据(`tagPush`、`branchesFilter`、`GolangBuild`、`ArtifactUpload`/`flowPublic`),占位符清晰标注
- [ ] 构建命令为 `goreleaser release --clean --skip=publish`;脚本中不出现 GITHUB_TOKEN / 发布 GitHub 的行为
- [ ] 不使用 `--snapshot`(版本号错位风险已在注释/文档中说明原因)
- [ ] GOPROXY 指向 goproxy.cn;goreleaser 经 `go install` 安装;zip 缺失时有安装兜底
- [ ] 文档完整覆盖:接入步骤(含服务连接)、tag 触发方式、产物下载位置、故障排查要点
- [ ] `.goreleaser.yml`、`.github/`、Go 代码零改动;`go build ./...` 与 `go vet ./...` 通过(工作区无 Go 代码变更)

## Notes

- 轻量任务,PRD-only;技术选型与官方文档依据见 `research/yunxiao-flow-yaml.md`(决策 D1–D4)。
- Flow YAML 存在于控制台而非仓库(官方机制),`flow.yml` 仅为仓库侧事实源,接入时需粘贴到控制台。
- `GolangBuild.with.goVersion` 的确切取值格式以 Flow 控制台「校验」提示为准,文档中注明。
