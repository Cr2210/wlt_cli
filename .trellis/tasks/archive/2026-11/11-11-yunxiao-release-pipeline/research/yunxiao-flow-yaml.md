# 调研:云效 Flow YAML 流水线(发布场景)

调研日期:2026-11-11。所有结论均来自阿里云官方文档(链接见文末)。

## 1. 接入方式(关键约束)

YAML 流水线的 YAML **存放在 Flow 控制台**,不从代码库读文件:

1. Flow 控制台 → 我的流水线 → 新建流水线 → 选 **YAML 化编排** → 选模板(如 Go 模板)→ 创建
2. 进入 YAML 编辑页,左侧粘贴/编辑 YAML,右侧工具栏「文档」可查步骤参数,「校验」按钮做语法校验
3. 保存并运行

因此仓库侧交付 = 一份 `flow.yml`(唯一事实源,含接入注释)+ 接入文档;用户在控制台粘贴。

## 2. sources(代码源 + tag 触发)

```yaml
sources:
  repo:
    type: codeup
    name: 代码源展示名(≤30字符)
    endpoint: https://<codeup 仓库地址>.git     # 必填
    branch: master                              # 默认分支(手工/定时触发用)
    triggerEvents:
      - tagPush                                 # Tag 创建触发(代码推送=push)
    branchesFilter: ^v                          # 正则;"仅满足过滤规则的分支或标签可以触发"
    certificate:
      type: serviceConnection
      serviceConnection: <服务连接ID>            # 企业设置→服务连接管理 获取
```

- `tagPush` 仅部分代码源支持,codeup 支持(Tag 创建)
- `branchesFilter` 对 tag 同样生效 → 用 `^v` 实现"仅 v* tag 触发"(对齐 GitHub Actions 的 `on.push.tags: v*`)
- 克隆默认 **全部历史**(cloneDepth 默认 0)→ goreleaser 能拿到 tag 列表
- tag 触发运行时,checkout 到该 tag

## 3. jobs / steps

```yaml
stages:
  stage_id:
    name: 展示名
    jobs:
      job_id:
        name: 展示名
        runsOn: public/cn-beijing   # 云效公共构建集群(北京);容器环境,root 用户
        steps:
          step_id:
            step: GolangBuild       # 内置步骤名
            name: 展示名(≤30字符)
            workspace: repo         # 可选,默认 defaultWorkspace
            with: { ... }           # 步骤参数
```

### GolangBuild(构建分类内置步骤)
- 作用:准备 Go 工具链 + 执行自定义命令(类似 JavaBuild 的 `with.run`)
- UI 参数:Go 版本(Golang版本选择)+ 构建命令(默认 `make build`)
- Flow Go 环境默认 `GOPROXY=https://goproxy.cn`(官方文档明示)
- Go ≥1.21 且 GOTOOLCHAIN=auto 时,go.mod 声明更高版本会自动下载工具链(也走 GOPROXY)→ 环境 Go 版本偏旧也有兜底
- YAML `with.goVersion` 的具体取值格式以控制台「校验」和右侧文档面板为准

### Command(执行命令)
通用 shell 步骤,`with.run` 为脚本;备选方案(如果 GolangBuild 参数在控制台校验不过)。

### ArtifactUpload(构建物上传)
```yaml
with:
  uploadType: flowPublic      # 归档至云效公共存储空间,运行详情页可下载 ← 本任务选型
  artifact: <制品名>           # 如 wlt-dist
  filePath:                   # 项目根相对路径,可多个
    - dist/
    - wlt-skills.zip
```
另一种 `uploadType` 是组织私有通用制品仓库(Packages,需服务连接)——用户已选 flowPublic。

## 4. 内置环境变量

- `CI_COMMIT_REF_NAME`:运行分支名或 **tag 名**(tag 触发时 = tag,如 `v0.2.0`)
- `CI_COMMIT_SHA` / `CI_COMMIT_ID`(8位)/ `CI_COMMIT_TITLE`(需代码拉取步骤后有值)
- 多源时用带下标变量(`CI_COMMIT_REF_NAME_1`);单源用无下标即可
- 任务内变量传递:默认环境用 `.env`(需 `USER_` 前缀);容器/VM 环境用 `>> "$FLOW_ENV"`(必须追加模式)
- 变量不支持在 YAML 里做字符串加工,需在 shell 里赋值后再处理

## 5. GoReleaser 在云效上的运行策略(核心决策)

### 决策 D1:用 `goreleaser release --clean --skip=publish`,不用 snapshot
- `.goreleaser.yml` 的 `release:` 段指向 GitHub(Cr2210/wlt_cli),直接 full release 需要 GITHUB_TOKEN 且会发 GitHub Release → 必须 skip publish
- snapshot 模式的 `version_template: "{{ incpatch .Version }}-snapshot"` 会把内嵌版本变成 `0.2.1-snapshot` 这类(版本号错位)→ 不用
- `--skip=publish` 是 goreleaser v2 统一 flag(v1 的 `--skip-publish` 已废弃),skip 后不访问 GitHub API、不需要 token
- full release 模式要求:HEAD 恰好位于 tag 上(Flow tagPush 运行满足)、工作区干净(仓库保持 tidy,与 GitHub Actions/本地 make release 相同前提)、git 历史含 tags(Flow 默认全量克隆满足)

### 决策 D2:goreleaser 安装走 `go install` + goproxy.cn
- 公共构建集群直连 GitHub release 下载不稳定(用户网络环境对国外源不友好)
- `go install github.com/goreleaser/goreleaser/v2@latest` 走 goproxy.cn 稳;v2 module path 固定,@latest 不会跨 major
- 编译耗时约 1–3 分钟,可接受;文档中说明如何固定版本(如 `@v2.x.y`)

### 决策 D3:zip 依赖防御性安装
- `.goreleaser.yml` before hook `cd skills && zip -r ../wlt-skills.zip .` 依赖 `zip` 命令
- 构建脚本前置 `command -v zip || apt-get install -y zip`(公共集群容器为 root,apt 源为阿里云镜像)
- 不改 `.goreleaser.yml`(它同时服务 GitHub Actions 和本地构建)

### 决策 D4:产物上传清单
- `dist/` 整目录(含 5 平台包 + checksums.txt)+ 项目根 `wlt-skills.zip`(hook 生成在根目录)
- 产物文件名与 CLAUDE.md 发布流程章节一致:wlt-linux-amd64.tar.gz / wlt-linux-arm64.tar.gz / wlt-darwin-amd64.tar.gz / wlt-darwin-arm64.tar.gz / wlt-windows-amd64.zip / checksums.txt / wlt-skills.zip

## 6. 现状事实

- `go.mod`: module github.com/weiliantong/cli, go 1.26.4
- `.goreleaser.yml`: version 2;before hooks = `go mod tidy` + skills zip;builds 3 OS × 2 Arch,CGO_ENABLED=0,-trimpath,ldflags 注入 `internal/build.Version={{.Version}}`;archives `{{.Binary}}-{{.Os}}-{{.Arch}}`(windows→zip);checksums sha256;snapshot version_template `{{ incpatch .Version }}-snapshot`
- `.github/workflows/release.yml`:tag `v*` → setup-go(读 go.mod)→ `goreleaser release --clean`(发布 GitHub Releases)——保持不动,两套并存
- Makefile:`make snapshot`/`make release` 本地构建,产物手动传 CodeUp 发行版(现状,本任务提供自动化替代)
- CLAUDE.md「发布流程(云效 CodeUp)」章节需补充云效流水线入口

## 7. 参考文档

- YAML 流水线总览: https://help.aliyun.com/zh/yunxiao/user-guide/yaml-preliminary-experience
- sources 语法(tagPush/branchesFilter/certificate): https://help.aliyun.com/zh/yunxiao/user-guide/pipeline-sources (help.aliyun.com → YAML语法介绍 → 语法细则 → 流水线源 sources)
- steps 语法: https://help.aliyun.com/zh/yunxiao/user-guide/pipeline-steps
- 内置步骤清单(GolangBuild/Command/ArtifactUpload): https://help.aliyun.com/zh/yunxiao/user-guide/step-steps-list
- Go 测试构建(Go 环境/GOPROXY 默认/上传方式): https://help.aliyun.com/zh/yunxiao/user-guide/go-test-build
- 环境变量: https://help.aliyun.com/zh/yunxiao/user-guide/environment-variables
- 新建 YAML 流水线: https://help.aliyun.com/zh/yunxiao/user-guide/new-yaml-pipeline
- goreleaser v2 `--skip` flag: https://goreleaser.com/resources/deprecations (v1 --skip-publish → v2 --skip=publish)
