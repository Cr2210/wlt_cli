# wlt CLI Installer for Windows
#
# Usage:
#   irm https://raw.githubusercontent.com/Cr2210/wlt_cli/main/scripts/install.ps1 | iex
#
# Environment variables (all optional):
#   WLT_INSTALL_DIR  — where to put the binary (default: ~/.local/bin)
#   WLT_VERSION      — version to install (default: latest)
#   WLT_NO_SKILLS    — set to 1 to skip skills install (default: 0)

$ErrorActionPreference = "Stop"

$Repo = "Cr2210/wlt_cli"
$BinName = "wlt"
$InstallDir = if ($env:WLT_INSTALL_DIR) { $env:WLT_INSTALL_DIR } else { Join-Path $HOME ".local\bin" }
$Version = if ($env:WLT_VERSION) { $env:WLT_VERSION } else { "latest" }
$NoSkills = $env:WLT_NO_SKILLS -eq "1"
$SkillName = "wlt"

# Agent skills 安装目标目录
$AgentDirs = @(
  ".agents\skills",
  ".claude\skills",
  ".cursor\skills",
  ".gemini\skills",
  ".codex\skills",
  ".github\skills",
  ".windsurf\skills",
  ".augment\skills",
  ".cline\skills",
  ".amp\skills",
  ".kiro\skills",
  ".trae\skills",
  ".openclaw\skills"
)

# ── Helpers ──────────────────────────────────────────────────────────────────

function Write-Say { param([string]$Message) Write-Host " $Message" }
function Write-Err { param([string]$Message) Write-Host " ❌ $Message" -ForegroundColor Red; exit 1 }

# ── Detect Architecture ──────────────────────────────────────────────────────

function Get-Arch {
  if ($env:WLT_ARCH) {
    $override = $env:WLT_ARCH.ToLower()
    if ($override -eq "amd64" -or $override -eq "arm64") { return $override }
    Write-Err "无效的 WLT_ARCH 值 '$env:WLT_ARCH'，必须是 'amd64' 或 'arm64'。"
  }
  try {
    $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    if ($arch) {
      switch ($arch.ToString()) {
        "X64"   { return "amd64" }
        "Arm64" { return "arm64" }
      }
    }
  } catch {}
  $envArch = $env:PROCESSOR_ARCHITECTURE
  if ($envArch) {
    switch ($envArch.ToUpper()) {
      "AMD64" { return "amd64" }
      "ARM64" { return "arm64" }
    }
  }
  Write-Err "无法检测系统架构。请设置 WLT_ARCH 环境变量为 'amd64' 或 'arm64'。"
}

# ── Resolve latest version from GitHub ───────────────────────────────────────

function Resolve-LatestVersion {
  if ($Version -eq "latest") {
    try {
      $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
      if ($release.tag_name) {
        $script:Version = $release.tag_name
        return
      }
    } catch {}
    try {
      $resp = Invoke-WebRequest -Uri "https://github.com/$Repo/releases/latest" `
        -MaximumRedirection 0 -UseBasicParsing 2>$null
    } catch {
      $loc = $null
      try { $loc = $_.Exception.Response.GetResponseHeader("Location") } catch {}
      if ($loc) {
        $script:Version = ($loc -split "/tag/")[-1].Trim()
        return
      }
    }
    Write-Err "无法获取最新版本。请设置 `$env:WLT_VERSION 指定版本。"
  }
}

# ── Banner ───────────────────────────────────────────────────────────────────

function Write-Banner {
  Write-Host ""
  Write-Say "┌──────────────────────────────────────┐"
  Write-Say "│  WLT Installer                       │"
  Write-Say "│  维链通 ERP 命令行工具               │"
  Write-Say "└──────────────────────────────────────┘"
  Write-Host ""
}

# ── Install Binary ───────────────────────────────────────────────────────────

function Install-Binary {
  $arch = Get-Arch
  Resolve-LatestVersion

  $archiveName = "${BinName}-windows-${arch}.zip"
  $downloadUrl = "https://github.com/$Repo/releases/download/$Version/$archiveName"

  Write-Say "⬇  下载 ${BinName} ${Version} (windows/${arch})..."

  $tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) "wlt-install-$PID"
  New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

  try {
    $archivePath = Join-Path $tmpDir $archiveName
    Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -UseBasicParsing

    # SHA256 verification
    $checksumUrl = "https://github.com/$Repo/releases/download/$Version/checksums.txt"
    try {
      $checksumPath = Join-Path $tmpDir "checksums.txt"
      Invoke-WebRequest -Uri $checksumUrl -OutFile $checksumPath -UseBasicParsing
      $checksumContent = Get-Content $checksumPath
      $expectedLine = $checksumContent | Where-Object { $_ -match [regex]::Escape($archiveName) }
      if ($expectedLine) {
        $expected = ($expectedLine -split '\s+')[0]
        $actual = (Get-FileHash -Path $archivePath -Algorithm SHA256).Hash.ToLower()
        if ($actual -ne $expected.ToLower()) {
          Write-Err "SHA256 校验失败! 期望 $expected，实际 $actual。"
        }
        Write-Say "✅ SHA256 校验通过"
      }
    } catch {
      Write-Say "⚠️  无法下载 checksums.txt，跳过验证"
    }

    # Extract
    Write-Say "📦 解压..."
    Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

    # Install
    if (!(Test-Path $InstallDir)) {
      New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $binFile = Get-ChildItem -Path $tmpDir -Recurse -Filter "${BinName}.exe" | Select-Object -First 1
    if ($null -eq $binFile) {
      Write-Err "压缩包中找不到 ${BinName}.exe。"
    }

    $destBin = Join-Path $InstallDir "${BinName}.exe"
    Copy-Item -Path $binFile.FullName -Destination $destBin -Force
    Write-Say "✅ 已安装到: $destBin"

    # Add to PATH
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($userPath -notlike "*$InstallDir*") {
      Write-Say ""
      Write-Say "⚠️  $InstallDir 不在 PATH 中，正在添加..."
      [Environment]::SetEnvironmentVariable("PATH", "$InstallDir;$userPath", "User")
      $env:PATH = "$InstallDir;$env:PATH"
      Write-Say "✅ 已添加到用户 PATH。请重启终端生效。"
    }
  } finally {
    Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
  }
}

# ── Install Skills ───────────────────────────────────────────────────────────

function Install-Skills {
  Write-Say ""
  Write-Say "📦 安装 AI Agent Skills..."

  Resolve-LatestVersion
  $skillsUrl = "https://github.com/$Repo/releases/download/$Version/wlt-skills.zip"

  $tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) "wlt-skills-$PID"
  New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

  try {
    $zipPath = Join-Path $tmpDir "wlt-skills.zip"
    try {
      Invoke-WebRequest -Uri $skillsUrl -OutFile $zipPath -UseBasicParsing
    } catch {
      Write-Say "⚠️  无法下载 skills，跳过。"
      return
    }

    $extractDir = Join-Path $tmpDir "skills"
    Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

    $installed = 0
    foreach ($agentDir in $AgentDirs) {
      $baseDir = Join-Path $HOME $agentDir
      $parentGate = Split-Path $baseDir -Parent

      if ($installed -gt 0 -and !(Test-Path $parentGate)) {
        continue
      }

      $dest = Join-Path $baseDir $SkillName
      $label = "~\$agentDir\$SkillName"

      if (Test-Path $dest) {
        Remove-Item -Path $dest -Recurse -Force
      }
      New-Item -ItemType Directory -Path $dest -Force | Out-Null

      Get-ChildItem -Path $extractDir | ForEach-Object {
        $destPath = Join-Path $dest $_.Name
        if ($_.PSIsContainer) {
          Copy-Item -Path $_.FullName -Destination $destPath -Recurse -Force
        } else {
          Copy-Item -Path $_.FullName -Destination $destPath -Force
        }
      }

      $fileCount = (Get-ChildItem -Path $dest -Recurse -File).Count

      if ($installed -eq 0) {
        Write-Say "✅ Skills → $label ($fileCount files)"
        Get-ChildItem -Path $dest | ForEach-Object {
          if ($_.PSIsContainer) {
            $subCount = (Get-ChildItem -Path $_.FullName -Recurse -File).Count
            Write-Say "   📁 $($_.Name)/ ($subCount files)"
          } else {
            Write-Say "   📄 $($_.Name)"
          }
        }
      } else {
        Write-Say "✅ Skills → $label ($fileCount files)"
      }

      $installed++
    }

    if ($installed -eq 0) {
      $fallback = Join-Path (Join-Path $HOME ".agents\skills") $SkillName
      New-Item -ItemType Directory -Path $fallback -Force | Out-Null
      Get-ChildItem -Path $extractDir | ForEach-Object {
        $destPath = Join-Path $fallback $_.Name
        if ($_.PSIsContainer) {
          Copy-Item -Path $_.FullName -Destination $destPath -Recurse -Force
        } else {
          Copy-Item -Path $_.FullName -Destination $destPath -Force
        }
      }
      $fileCount = (Get-ChildItem -Path $fallback -Recurse -File).Count
      Write-Say "✅ Skills → ~\.agents\skills\$SkillName ($fileCount files)"
    }
  } finally {
    Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
  }
}

# ── Main ─────────────────────────────────────────────────────────────────────

Write-Banner
Install-Binary
if (-not $NoSkills) {
  Install-Skills
}

Write-Host ""
Write-Say "🎉 安装完成！"
Write-Say ""
Write-Say "下一步:"
Write-Say "  wlt version        # 验证安装"
Write-Say "  wlt config init    # 初始化配置"
Write-Say "  wlt auth login     # 登录"
Write-Host ""
