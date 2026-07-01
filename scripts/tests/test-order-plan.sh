#!/usr/bin/env bash
# scripts/tests/test-order-plan.sh
# wlt order plan 模块端到端测试（采购运输计划 / 销售运输计划 + CRUD）
#
# 用法:
#   WLT_TOKEN=<accessToken> WLT_TENANT_ID=<tenantId> ./scripts/tests/test-order-plan.sh
#   ./scripts/tests/test-order-plan.sh                # 只跑 smoke 段(无需 token)
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
# 订单计划子命令（后端统一 /erp/order-plan 端点 + type 区分子类型）:
#   采购计划(plan purchase)   list/page-count        /erp/order-plan?type=PURCHASE_TRANSPORT_PLAN
#   销售计划(plan sale)       list/page-count        /erp/order-plan?type=SALE_TRANSPORT_PLAN
#   get / create / update / delete / update-status / cancel / reopen / complete / export

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

if [ ! -x "$WLT_BIN" ]; then
    echo "ERROR: WLT_BIN ($WLT_BIN) not found or not executable" >&2
    exit 1
fi

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
echo " wlt order plan 模块测试"
echo " Binary: $WLT_BIN"
[ -n "$TOKEN" ] && echo " Tenant: $TENANT" || echo " Tenant: (无 token,只跑 smoke 段)"
echo "=================================================="
echo

# === smoke 段(无需 token) ===
echo "--- [smoke] 二进制基本可用性 ---"
run_plain "wlt --version 成功" version
run_plain "wlt order --help 成功" order --help
run_plain "wlt order plan --help 成功" order plan --help
run_help "wlt order plan --help 列出 purchase" "purchase" order plan --help
run_help "wlt order plan --help 列出 sale" "sale" order plan --help
echo

# --- 采购计划 smoke ---
echo "--- [smoke] 采购计划 (plan purchase) ---"
run_help "list --help 列出 --no" "--no" order plan purchase list --help
run_help "list --help 列出 --product-id" "--product-id" order plan purchase list --help
run_help "list --help 列出 --supplier-id" "--supplier-id" order plan purchase list --help
run_help "list --help 列出 --customer-id" "--customer-id" order plan purchase list --help
run_help "list --help 列出 --start" "--start" order plan purchase list --help
run_help "list --help 列出 --end" "--end" order plan purchase list --help
run_help "list --help 列出 --page-no" "--page-no" order plan purchase list --help
run_help "list --help 列出 --page-size" "--page-size" order plan purchase list --help
run_help "page-count --help 列出 --no" "--no" order plan purchase page-count --help
echo

# --- 销售计划 smoke ---
echo "--- [smoke] 销售计划 (plan sale) ---"
run_help "list --help 列出 --no" "--no" order plan sale list --help
run_help "list --help 列出 --start" "--start" order plan sale list --help
run_help "page-count --help 列出 --supplier-id" "--supplier-id" order plan sale page-count --help
echo

# --- 通用 CRUD smoke (不区分 type) ---
echo "--- [smoke] 计划通用 CRUD (plan) ---"
run_help "get --help 列出 --id" "--id" order plan get --help
run_help "create --help 列出 --data" "--data" order plan create --help
run_help "update --help 列出 --data" "--data" order plan update --help
run_help "delete --help 列出 --id" "--id" order plan delete --help
run_help "update-status --help 列出 --data" "--data" order plan update-status --help
run_help "cancel --help 列出 --data" "--data" order plan cancel --help
run_help "reopen --help 列出 --data" "--data" order plan reopen --help
run_help "complete --help 列出 --data" "--data" order plan complete --help
run_help "export --help 列出 --no" "--no" order plan export --help
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
run_expect_fail "missing --token returns code 4" 4 order plan purchase list --tenant-id "$TENANT"
run_expect_fail "missing --tenant-id returns code 4" 4 order plan purchase list --token "$TOKEN"
echo

# --- 采购计划 ---
echo "--- [api] 采购计划 (plan purchase) ---"
run_cmd "list (默认第1页)" order plan purchase list
run_cmd "list --page-no 1 --page-size 5" order plan purchase list --page-no 1 --page-size 5
run_cmd "list --no CGJH20260528000001" order plan purchase list --no CGJH20260528000001
run_cmd "list --product-id 1927285600675729409" order plan purchase list --product-id 1927285600675729409
run_cmd "list --supplier-id 2001552070697033730" order plan purchase list --supplier-id 2001552070697033730
run_cmd "list --start --end" order plan purchase list --start "2026-07-21 00:00:00" --end "2026-08-19 23:59:59"
run_cmd "page-count" order plan purchase page-count
run_cmd "page-count --no CGJH20260528000001" order plan purchase page-count --no CGJH20260528000001
echo

# --- 销售计划 ---
echo "--- [api] 销售计划 (plan sale) ---"
run_cmd "list (默认第1页)" order plan sale list
run_cmd "list --page-no 1 --page-size 5" order plan sale list --page-no 1 --page-size 5
run_cmd "list --no XSJH20260401000001" order plan sale list --no XSJH20260401000001
run_cmd "list --product-id 1927670802287808513" order plan sale list --product-id 1927670802287808513
run_cmd "list --customer-id 2001494968037298178" order plan sale list --customer-id 2001494968037298178
run_cmd "list --start --end" order plan sale list --start "2026-07-20 00:00:00" --end "2026-08-20 23:59:59"
run_cmd "page-count" order plan sale page-count
run_cmd "page-count --no XSJH20260401000001" order plan sale page-count --no XSJH20260401000001
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
