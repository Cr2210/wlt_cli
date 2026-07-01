#!/usr/bin/env bash
# scripts/tests/test-order.sh
# wlt order 模块端到端测试（main 主订单: purchase/sale 以及 CRUD）
#
# 用法:
#   WLT_TOKEN=<accessToken> WLT_TENANT_ID=<tenantId> ./scripts/tests/test-order.sh
#   ./scripts/tests/test-order.sh                    # 只跑 smoke 段(无需 token)
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
# 订单模块子命令:
#   主订单 main:
#     采购订单(main purchase)   list/page-count        /erp/order?type=PURCHASE
#     销售订单(main sale)       list/page-count        /erp/order?type=SALE
#     get / get-linkorder-by-orderId / create / update / delete / update-status
#     cancel / reopen / complete / link-waybill / unlink-waybill / export

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
echo " wlt order 模块测试"
echo " Binary: $WLT_BIN"
[ -n "$TOKEN" ] && echo " Tenant: $TENANT" || echo " Tenant: (无 token,只跑 smoke 段)"
echo "=================================================="
echo

# === smoke 段(无需 token) ===
echo "--- [smoke] 二进制基本可用性 ---"
run_plain "wlt --version 成功" version
run_plain "wlt order --help 成功" order --help
run_plain "wlt order main --help 成功" order main --help
run_help "wlt order main --help 列出 purchase" "purchase" order main --help
run_help "wlt order main --help 列出 sale" "sale" order main --help
echo

# --- 采购订单 smoke ---
echo "--- [smoke] 采购订单 (main purchase) ---"
run_help "list --help 列出 --no" "--no" order main purchase list --help
run_help "list --help 列出 --enterprise-id" "--enterprise-id" order main purchase list --help
run_help "list --help 列出 --product-id" "--product-id" order main purchase list --help
run_help "list --help 列出 --order-start" "--order-start" order main purchase list --help
run_help "list --help 列出 --order-end" "--order-end" order main purchase list --help
run_help "list --help 列出 --page-no" "--page-no" order main purchase list --help
run_help "list --help 列出 --page-size" "--page-size" order main purchase list --help
run_help "page-count --help 列出 --no" "--no" order main purchase page-count --help
echo

# --- 销售订单 smoke ---
echo "--- [smoke] 销售订单 (main sale) ---"
run_help "list --help 列出 --no" "--no" order main sale list --help
run_help "list --help 列出 --order-start" "--order-start" order main sale list --help
run_help "page-count --help 列出 --enterprise-id" "--enterprise-id" order main sale page-count --help
echo

# --- 通用 CRUD smoke (不区分 type) ---
echo "--- [smoke] 订单通用 CRUD (main) ---"
run_help "main get --help 列出 --id" "--id" order main get --help
run_help "main get-linkorder-by-orderId --help 列出 --order-id" "--order-id" order main get-linkorder-by-orderId --help
run_help "main create --help 列出 --data" "--data" order main create --help
run_help "main update --help 列出 --data" "--data" order main update --help
run_help "main delete --help 列出 --id" "--id" order main delete --help
run_help "main update-status --help 列出 --data" "--data" order main update-status --help
run_help "main cancel --help 列出 --data" "--data" order main cancel --help
run_help "main reopen --help 列出 --data" "--data" order main reopen --help
run_help "main complete --help 列出 --data" "--data" order main complete --help
run_help "main link-waybill --help 列出 --data" "--data" order main link-waybill --help
run_help "main unlink-waybill --help 列出 --data" "--data" order main unlink-waybill --help
run_help "main export --help 列出 --no" "--no" order main export --help
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
run_expect_fail "missing --token returns code 4" 4 order main purchase list --tenant-id "$TENANT"
run_expect_fail "missing --tenant-id returns code 4" 4 order main purchase list --token "$TOKEN"
echo

# --- 采购订单 ---
echo "--- [api] 采购订单 (main purchase) ---"
run_cmd "list (默认第1页)" order main purchase list
run_cmd "list --page-no 1 --page-size 5" order main purchase list --page-no 1 --page-size 5
run_cmd "list --no CGDD20260402000002" order main purchase list --no CGDD20260402000002
run_cmd "list --enterprise-id 2001552070697033730" order main purchase list --enterprise-id 2001552070697033730
run_cmd "list --product-id 1927670802287808513" order main purchase list --product-id 1927670802287808513
run_cmd "list --order-start --order-end" order main purchase list --order-start "2026-07-17 00:00:00" --order-end "2026-08-06 23:59:59"
run_cmd "page-count" order main purchase page-count
run_cmd "page-count --no CGDD20260402000002" order main purchase page-count --no CGDD20260402000002
echo

# --- 销售订单 ---
echo "--- [api] 销售订单 (main sale) ---"
run_cmd "list (默认第1页)" order main sale list
run_cmd "list --page-no 1 --page-size 5" order main sale list --page-no 1 --page-size 5
run_cmd "list --no XSDD20260402000001" order main sale list --no XSDD20260402000001
run_cmd "list --enterprise-id 2001494968037298178" order main sale list --enterprise-id 2001494968037298178
run_cmd "list --product-id 1927612799140261889" order main sale list --product-id 1927612799140261889
run_cmd "list --order-start --order-end" order main sale list --order-start "2026-07-14 00:00:00" --order-end "2026-08-11 23:59:59"
run_cmd "page-count" order main sale page-count
run_cmd "page-count --no XSDD20260402000001" order main sale page-count --no XSDD20260402000001
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
