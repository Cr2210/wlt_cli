#!/usr/bin/env bash
# scripts/tests/test-stats.sh
# wlt stats 模块端到端测试
#
# 用法:
#   WLT_TOKEN=<accessToken> WLT_TENANT_ID=<tenantId> ./scripts/tests/test-stats.sh
#   ./scripts/tests/test-stats.sh                # 只跑 smoke 段(无需 token)
#
# 环境变量:
#   WLT_BIN         wlt 二进制路径(默认 ./wlt.exe)
#   WLT_TOKEN       访问令牌(API 段必填,smoke 段无需)
#   WLT_TENANT_ID   租户 ID(API 段必填,smoke 段无需)
#   WLT_PROFILE     profile 名(默认不传,走 base_url)
#   WLT_BASE_URL    可选,覆盖 base_url
#
# 退出码 = 失败用例数(0 = 全部通过)
#
# 后续模块测试脚本按同样模板写,放 scripts/tests/test-<module>.sh
# 当前已覆盖: wlt stats overview/stock/finance/sale/purchase/produce
# (2026-07 复核:所有 stats 子命令已加 --product-id flag)

set -u

WLT_BIN="${WLT_BIN:-./wlt.exe}"
TOKEN="${WLT_TOKEN:-}"
TENANT="${WLT_TENANT_ID:-}"

PROFILE_FLAG=()
[ -n "${WLT_PROFILE:-}" ] && PROFILE_FLAG=(--profile "$WLT_PROFILE")
[ -n "${WLT_BASE_URL:-}" ] && PROFILE_FLAG+=(--base-url "$WLT_BASE_URL")

PASS=0
FAIL=0
FAILED_TESTS=()

# 前置检查(只校验 binary 存在;token/tenant 留给各段按需)
if [ ! -x "$WLT_BIN" ]; then
    echo "ERROR: WLT_BIN ($WLT_BIN) not found or not executable" >&2
    exit 1
fi

# === Helpers ===
# run_cmd: 期望返回成功 JSON
run_cmd() {
    local name="$1"
    shift
    echo -n "[TEST] $name ... "
    local out rc
    out=$("$WLT_BIN" "$@" --token "$TOKEN" --tenant-id "$TENANT" "${PROFILE_FLAG[@]}" 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ] && [[ "$out" =~ ^\{ ]] && [[ "$out" == *'"data"'* ]]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        echo "FAIL (rc=$rc)"
        echo "  output (first 300 chars): ${out:0:300}"
        FAIL=$((FAIL + 1))
        FAILED_TESTS+=("$name")
    fi
}

# run_expect_fail: 期望返回特定退出码
run_expect_fail() {
    local name="$1"
    local expected_rc="$2"
    shift 2
    echo -n "[TEST] $name ... "
    local out rc
    out=$("$WLT_BIN" "$@" 2>&1) && rc=0 || rc=$?
    if [ "$rc" = "$expected_rc" ]; then
        echo "PASS (rc=$rc as expected)"
        PASS=$((PASS + 1))
    else
        echo "FAIL (expected rc=$expected_rc, got rc=$rc)"
        echo "  output: ${out:0:200}"
        FAIL=$((FAIL + 1))
        FAILED_TESTS+=("$name")
    fi
}

# run_help: 期望 --help 输出含指定子串(不需 token)
run_help() {
    local name="$1"
    local expect="$2"
    shift 2
    echo -n "[TEST] $name ... "
    local out rc
    out=$("$WLT_BIN" "$@" --help 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ] && [[ "$out" == *"$expect"* ]]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        echo "FAIL (rc=$rc, missing '$expect')"
        echo "  output: ${out:0:200}"
        FAIL=$((FAIL + 1))
        FAILED_TESTS+=("$name")
    fi
}

# run_plain: 任意成功(不校验 JSON),用于 --version 等
run_plain() {
    local name="$1"
    shift
    echo -n "[TEST] $name ... "
    local out rc
    out=$("$WLT_BIN" "$@" 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        echo "FAIL (rc=$rc)"
        echo "  output: ${out:0:200}"
        FAIL=$((FAIL + 1))
        FAILED_TESTS+=("$name")
    fi
}

echo "=================================================="
echo " wlt stats 模块测试"
echo " Binary: $WLT_BIN"
[ -n "$TOKEN" ] && echo " Tenant: $TENANT" || echo " Tenant: (无 token,只跑 smoke 段)"
echo "=================================================="
echo

# === smoke 段(无需 token) ===
echo "--- [smoke] 二进制基本可用性 ---"
run_plain "wlt --version 成功" version
run_plain "wlt stats --help 成功" stats --help
run_help "wlt stats --help 列出 overview" "overview" stats --help
run_help "wlt stats --help 列出 stock" "stock" stats --help
run_help "wlt stats --help 列出 finance/sale/purchase/produce" "finance" stats --help
run_help "wlt stats overview --help 列出 --type" "--type" stats overview --help
run_help "wlt stats overview --help 列出 --start-time" "--start-time" stats overview --help
run_help "wlt stats overview --help 列出 --product-id" "--product-id" stats overview --help
run_help "wlt stats finance data-overview --help 列出 --product-id" "--product-id" stats finance data-overview --help
echo

# === API 段(需 token + tenant) ===
if [ -z "$TOKEN" ] || [ -z "$TENANT" ]; then
    echo "(跳过 API 段:未设 WLT_TOKEN / WLT_TENANT_ID)"
    echo
    echo "=================================================="
    echo " 结果: $PASS 通过, $FAIL 失败 (仅 smoke 段)"
    echo "=================================================="
    exit $FAIL
fi

# --- 参数校验 ---
echo "--- [api] 参数校验 ---"
run_expect_fail "missing --token returns code 4" 4 stats overview --tenant-id "$TENANT"
run_expect_fail "missing --tenant-id returns code 4" 4 stats overview --token "$TOKEN"
echo

# --- 经营总览 ---
echo "--- [api] 经营总览 (overview) ---"
run_cmd "stats overview (default type=month)" stats overview
run_cmd "stats overview --type year" stats overview --type year
run_cmd "stats overview --type month --start-time 2026-07-01" stats overview --type month --start-time "2026-07-01 00:00:00"
run_cmd "stats overview --product-id 1" stats overview --product-id 1
run_cmd "stats overview --sort-by amount" stats overview --sort-by amount
echo

# --- 库存分析 ---
echo "--- [api] 库存分析 (stock) ---"
run_cmd "stats stock --product-id 1" stats stock --product-id 1
run_cmd "stats stock --warehouse-id 1" stats stock --warehouse-id 1
run_cmd "stats stock --product-id 1 --warehouse-id 1" stats stock --product-id 1 --warehouse-id 1
echo

# --- 财务分析(5 个)---
echo "--- [api] 财务分析 (finance, 5 个) ---"
run_cmd "stats finance data-overview --product-id 1" stats finance data-overview --product-id 1
run_cmd "stats finance receivable-rankings --product-id 1" stats finance receivable-rankings --product-id 1
run_cmd "stats finance overdue-receivable-rankings --product-id 1" stats finance overdue-receivable-rankings --product-id 1
run_cmd "stats finance payable-rankings --product-id 1" stats finance payable-rankings --product-id 1
run_cmd "stats finance overdue-payable-rankings --product-id 1" stats finance overdue-payable-rankings --product-id 1
run_cmd "stats finance data-overview --sort-by enterprise" stats finance data-overview --product-id 1 --sort-by enterprise
echo

# --- 销售分析(5 个)---
echo "--- [api] 销售分析 (sale, 5 个) ---"
run_cmd "stats sale data-overview --product-id 1" stats sale data-overview --product-id 1
run_cmd "stats sale customer-rankings --product-id 1" stats sale customer-rankings --product-id 1
run_cmd "stats sale product-rankings --product-id 1" stats sale product-rankings --product-id 1
run_cmd "stats sale employee-rankings --product-id 1" stats sale employee-rankings --product-id 1
run_cmd "stats sale region-rankings --product-id 1" stats sale region-rankings --product-id 1
echo

# --- 采购分析(5 个)---
echo "--- [api] 采购分析 (purchase, 5 个) ---"
run_cmd "stats purchase data-overview --product-id 1" stats purchase data-overview --product-id 1
run_cmd "stats purchase supplier-rankings --product-id 1" stats purchase supplier-rankings --product-id 1
run_cmd "stats purchase product-rankings --product-id 1" stats purchase product-rankings --product-id 1
run_cmd "stats purchase employee-rankings --product-id 1" stats purchase employee-rankings --product-id 1
run_cmd "stats purchase region-rankings --product-id 1" stats purchase region-rankings --product-id 1
echo

# --- 生产分析(1 个)---
echo "--- [api] 生产分析 (produce, 1 个) ---"
run_cmd "stats produce data-overview --product-id 1" stats produce data-overview --product-id 1
echo

echo "=================================================="
echo " 结果: $PASS 通过, $FAIL 失败"
if [ $FAIL -gt 0 ]; then
    echo " 失败用例:"
    for t in "${FAILED_TESTS[@]}"; do
        echo "   - $t"
    done
fi
echo "=================================================="
exit $FAIL

