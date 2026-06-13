#!/bin/sh
# wlt CLI Installer for Linux / macOS
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/Cr2210/wlt_cli/main/scripts/install.sh | sh
#
# Environment variables (all optional):
#   WLT_INSTALL_DIR  — where to put the binary (default: ~/.local/bin)
#   WLT_VERSION      — version to install (default: latest)
#   WLT_NO_SKILLS    — set to 1 to skip skills install (default: 0)

set -eu

REPO="Cr2210/wlt_cli"
BIN_NAME="wlt"
INSTALL_DIR="${WLT_INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${WLT_VERSION:-latest}"
NO_SKILLS="${WLT_NO_SKILLS:-0}"
SKILL_NAME="wlt"

# Agent skills 安装目标目录
AGENT_DIRS="
.agents/skills
.claude/skills
.cursor/skills
.gemini/skills
.codex/skills
.github/skills
.windsurf/skills
.augment/skills
.cline/skills
.amp/skills
.kiro/skills
.trae/skills
.openclaw/skills
"

# ── Helpers ──────────────────────────────────────────────────────────────────

say() { printf ' %s\n' "$@"; }
err() { printf ' ❌ %s\n' "$@" >&2; exit 1; }

need_cmd() {
  command -v "$1" >/dev/null 2>&1
}

download() {
  url="$1" dest="$2"
  if need_cmd curl; then
    curl -fsSL "$url" -o "$dest"
  elif need_cmd wget; then
    wget -qO "$dest" "$url"
  else
    err "需要 curl 或 wget，请先安装其中一个。"
  fi
}

# ── Detect OS / Arch ─────────────────────────────────────────────────────────

detect_os() {
  os="$(uname -s)"
  case "$os" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    MINGW*|MSYS*|CYGWIN*) err "Windows 请使用 PowerShell 安装: irm <url> | iex" ;;
    *) err "不支持的操作系统: $os" ;;
  esac
}

detect_arch() {
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) err "不支持的架构: $arch" ;;
  esac
}

# ── Resolve latest version from GitHub ───────────────────────────────────────

resolve_version() {
  if [ "$VERSION" = "latest" ]; then
    if need_cmd curl; then
      VERSION="$(curl -fsSI "https://github.com/${REPO}/releases/latest" 2>/dev/null \
        | grep -i '^location:' | sed 's|.*/tag/||;s/[[:space:]]*$//')"
    elif need_cmd wget; then
      VERSION="$(wget --spider --max-redirect=0 "https://github.com/${REPO}/releases/latest" 2>&1 \
        | grep -i 'Location:' | sed 's|.*/tag/||;s/[[:space:]]*$//')"
    fi
    if [ -z "$VERSION" ]; then
      err "无法获取最新版本。请设置 WLT_VERSION 环境变量指定版本。"
    fi
  fi
}

# ── Banner ───────────────────────────────────────────────────────────────────

print_banner() {
  printf '\n'
  say "┌──────────────────────────────────────┐"
  say "│  WLT Installer                       │"
  say "│  维链通 ERP 命令行工具               │"
  say "└──────────────────────────────────────┘"
  printf '\n'
}

# ── Install Binary ───────────────────────────────────────────────────────────

install_binary() {
  os="$(detect_os)"
  arch="$(detect_arch)"
  resolve_version

  archive_name="${BIN_NAME}-${os}-${arch}.tar.gz"
  download_url="https://github.com/${REPO}/releases/download/${VERSION}/${archive_name}"

  say "⬇  下载 ${BIN_NAME} ${VERSION} (${os}/${arch})..."

  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT INT TERM

  download "$download_url" "$tmpdir/$archive_name"

  # ── SHA256 校验 ────────────────────────────────────────────────────────────
  checksum_url="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
  if download "$checksum_url" "$tmpdir/checksums.txt" 2>/dev/null; then
    expected="$(awk -v file="$archive_name" '$2 == file {print $1; exit}' "$tmpdir/checksums.txt")"
    if [ -n "$expected" ]; then
      if need_cmd sha256sum; then
        actual="$(sha256sum "$tmpdir/$archive_name" | awk '{print $1}')"
      elif need_cmd shasum; then
        actual="$(shasum -a 256 "$tmpdir/$archive_name" | awk '{print $1}')"
      else
        actual=""
      fi
      if [ -n "$actual" ] && [ "$actual" != "$expected" ]; then
        err "SHA256 校验失败! 期望 ${expected}，实际 ${actual}。"
      fi
      if [ -n "$actual" ]; then
        say "✅ SHA256 校验通过"
      else
        say "⚠️  无法计算校验和，跳过验证"
      fi
    fi
  else
    say "⚠️  无法下载 checksums.txt，跳过验证"
  fi

  # ── 解压 + 安装 ────────────────────────────────────────────────────────────
  say "📦 解压..."
  tar xzf "$tmpdir/$archive_name" -C "$tmpdir"
  mkdir -p "$INSTALL_DIR"

  if [ -f "$tmpdir/$BIN_NAME" ]; then
    cp "$tmpdir/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
  else
    found="$(find "$tmpdir" -name "$BIN_NAME" -type f | head -1)"
    if [ -n "$found" ]; then
      cp "$found" "$INSTALL_DIR/$BIN_NAME"
    else
      err "压缩包中找不到 ${BIN_NAME} 二进制文件。"
    fi
  fi

  chmod +x "$INSTALL_DIR/$BIN_NAME"
  say "✅ 已安装到: ${INSTALL_DIR}/${BIN_NAME}"

  # ── PATH 检查 ──────────────────────────────────────────────────────────────
  case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
      say ""
      say "⚠️  ${INSTALL_DIR} 不在 PATH 中。"
      say "   添加方式:"
      say "     export PATH=\"${INSTALL_DIR}:\$PATH\""
      say "   或将此行加入 ~/.bashrc / ~/.zshrc"
      ;;
  esac
}

# ── Install Skills ───────────────────────────────────────────────────────────

install_skills() {
  say ""
  say "📦 安装 AI Agent Skills..."

  resolve_version
  skills_url="https://github.com/${REPO}/releases/download/${VERSION}/wlt-skills.zip"

  tmpdir_skills="$(mktemp -d)"

  if ! download "$skills_url" "$tmpdir_skills/wlt-skills.zip" 2>/dev/null; then
    say "⚠️  无法下载 skills，跳过。"
    rm -rf "$tmpdir_skills"
    return
  fi

  if ! need_cmd unzip; then
    say "⚠️  需要 unzip 命令解压 skills，跳过。"
    rm -rf "$tmpdir_skills"
    return
  fi

  unzip -q "$tmpdir_skills/wlt-skills.zip" -d "$tmpdir_skills/skills"

  skill_src="$tmpdir_skills/skills"
  installed=0

  for agent_dir in $AGENT_DIRS; do
    base_dir="$HOME/$agent_dir"
    if [ "$installed" -gt 0 ] && [ ! -e "$(dirname "$base_dir")" ]; then
      continue
    fi

    dest="$base_dir/$SKILL_NAME"
    label="~/$agent_dir/$SKILL_NAME"

    rm -rf "$dest"
    mkdir -p "$dest"
    cp -R "$skill_src/"* "$dest/" 2>/dev/null || cp -r "$skill_src/"* "$dest/" 2>/dev/null || true

    file_count="$(find "$dest" -type f | wc -l | tr -d ' ')"

    if [ "$installed" -eq 0 ]; then
      say "✅ Skills → $label ($file_count files)"
      for entry in "$dest"/*; do
        entry_name="$(basename "$entry")"
        if [ -d "$entry" ]; then
          say "   📁 ${entry_name}/"
        else
          say "   📄 ${entry_name}"
        fi
      done
    else
      say "✅ Skills → $label ($file_count files)"
    fi

    installed=$((installed + 1))
  done

  if [ "$installed" -eq 0 ]; then
    dest="$HOME/.agents/skills/$SKILL_NAME"
    mkdir -p "$dest"
    cp -R "$skill_src/"* "$dest/" 2>/dev/null || cp -r "$skill_src/"* "$dest/" 2>/dev/null || true
    file_count="$(find "$dest" -type f | wc -l | tr -d ' ')"
    say "✅ Skills → ~/.agents/skills/$SKILL_NAME ($file_count files)"
  fi

  rm -rf "$tmpdir_skills"
}

# ── Main ─────────────────────────────────────────────────────────────────────

main() {
  print_banner
  install_binary
  if [ "$NO_SKILLS" != "1" ]; then
    install_skills
  fi
  printf '\n'
  say "🎉 安装完成！"
  say ""
  say "下一步:"
  say "  wlt version        # 验证安装"
  say "  wlt config init    # 初始化配置"
  say "  wlt auth login     # 登录"
  printf '\n'
}

main
