#!/usr/bin/env bash
# scripts/tests/test-contract.sh
# wlt contract 模块端到端测试
#
# 用法:
#   WLT_TOKEN=<accessToken> WLT_TENANT_ID=<tenantId> ./scripts/tests/test-contract.sh
#   ./scripts/tests/test-contract.sh                # 只跑 smoke 段(无需 token)
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
# 合同模块子命令(4 大类 7 个子类,每个子类含 list/page-count/get/update-status/create/update/delete):
#   采购长协(purchase-long-cooperate)   /erp/contract?type=PURCHASE_LONG_COOPERATE
#   销售合同(sale-contract)             /erp/service-contract?type=SALE_CONTRACT
#   销售长协(sale-long-cooperate)       /erp/contract?type=SALE_LONG_COOPERATE
#   运输合同(transport)                 /erp/transport-contract?type=TRANSPORT
#   运输长协(transport-long)            /erp/contract?type=TRANSPORT_LONG
#   服务合同(service-contract)          /erp/provision-contract?type=SERVICE
#   服务长协(service-long)              /erp/provision-contract?type=SERVICE_LONG

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
echo " wlt contract 模块测试"
echo " Binary: $WLT_BIN"
[ -n "$TOKEN" ] && echo " Tenant: $TENANT" || echo " Tenant: (无 token,只跑 smoke 段)"
echo "=================================================="
echo

# === smoke 段(无需 token) ===
echo "--- [smoke] 二进制基本可用性 ---"
run_plain "wlt --version 成功" version
run_plain "wlt contract --help 成功" contract --help
run_help "wlt contract --help 列出 purchase-long-cooperate" "purchase-long-cooperate" contract --help
run_help "wlt contract --help 列出 sale-contract" "sale-contract" contract --help
run_help "wlt contract --help 列出 sale-long-cooperate" "sale-long-cooperate" contract --help
run_help "wlt contract --help 列出 transport" "transport" contract --help
run_help "wlt contract --help 列出 transport-long" "transport-long" contract --help
run_help "wlt contract --help 列出 service-contract" "service-contract" contract --help
run_help "wlt contract --help 列出 service-long" "service-long" contract --help
echo

# --- 采购长协 smoke ---
echo "--- [smoke] 采购长协 (purchase-long-cooperate) ---"
run_help "list --help 列出 --keyword" "--keyword" contract purchase-long-cooperate list --help
run_help "list --help 列出 --enterprise-id" "--enterprise-id" contract purchase-long-cooperate list --help
run_help "list --help 列出 --order-start" "--order-start" contract purchase-long-cooperate list --help
run_help "list --help 列出 --order-end" "--order-end" contract purchase-long-cooperate list --help
run_help "list --help 列出 --end-start" "--end-start" contract purchase-long-cooperate list --help
run_help "list --help 列出 --end-end" "--end-end" contract purchase-long-cooperate list --help
run_help "list --help 列出 --page-no" "--page-no" contract purchase-long-cooperate list --help
run_help "list --help 列出 --page-size" "--page-size" contract purchase-long-cooperate list --help
run_help "page-count --help 列出 --keyword" "--keyword" contract purchase-long-cooperate page-count --help
run_help "get --help 列出 --id" "--id" contract purchase-long-cooperate get --help
run_help "create --help 列出 --data" "--data" contract purchase-long-cooperate create --help
run_help "update --help 列出 --data" "--data" contract purchase-long-cooperate update --help
run_help "delete --help 列出 --ids" "--ids" contract purchase-long-cooperate delete --help
run_help "update-status --help 列出 --data" "--data" contract purchase-long-cooperate update-status --help
echo

# --- 销售合同 smoke ---
echo "--- [smoke] 销售合同 (sale-contract) ---"
run_help "list --help 列出 --keyword" "--keyword" contract sale-contract list --help
run_help "list --help 列出 --order-start" "--order-start" contract sale-contract list --help
run_help "page-count --help 列出 --enterprise-id" "--enterprise-id" contract sale-contract page-count --help
run_help "get --help 列出 --id" "--id" contract sale-contract get --help
echo

# --- 销售长协 smoke ---
echo "--- [smoke] 销售长协 (sale-long-cooperate) ---"
run_help "list --help 列出 --keyword" "--keyword" contract sale-long-cooperate list --help
run_help "list --help 列出 --end-end" "--end-end" contract sale-long-cooperate list --help
run_help "page-count --help 列出 --keyword" "--keyword" contract sale-long-cooperate page-count --help
echo

# --- 运输合同 smoke ---
echo "--- [smoke] 运输合同 (transport) ---"
run_help "list --help 列出 --keyword" "--keyword" contract transport list --help
run_help "list --help 列出 --order-start" "--order-start" contract transport list --help
run_help "page-count --help 列出 --enterprise-id" "--enterprise-id" contract transport page-count --help
run_help "get --help 列出 --id" "--id" contract transport get --help
run_help "create --help 列出 --data" "--data" contract transport create --help
run_help "delete --help 列出 --ids" "--ids" contract transport delete --help
run_help "update-status --help 列出 --data" "--data" contract transport update-status --help
echo

# --- 运输长协 smoke ---
echo "--- [smoke] 运输长协 (transport-long) ---"
run_help "list --help 列出 --keyword" "--keyword" contract transport-long list --help
run_help "list --help 列出 --end-start" "--end-start" contract transport-long list --help
run_help "page-count --help 列出 --keyword" "--keyword" contract transport-long page-count --help
echo

# --- 服务合同 smoke ---
echo "--- [smoke] 服务合同 (service-contract) ---"
run_help "list --help 列出 --keyword" "--keyword" contract service-contract list --help
run_help "list --help 列出 --order-end" "--order-end" contract service-contract list --help
run_help "page-count --help 列出 --enterprise-id" "--enterprise-id" contract service-contract page-count --help
run_help "get --help 列出 --id" "--id" contract service-contract get --help
echo

# --- 服务长协 smoke ---
echo "--- [smoke] 服务长协 (service-long) ---"
run_help "list --help 列出 --keyword" "--keyword" contract service-long list --help
run_help "list --help 列出 --end-end" "--end-end" contract service-long list --help
run_help "page-count --help 列出 --keyword" "--keyword" contract service-long page-count --help
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

# --- 参数校验(只校验 EnsureClient 返回 code 4 的场景;本地 flag 缺失由 cobra 返回 code 1) ---
echo "--- [api] 参数校验 ---"
run_expect_fail "missing --token returns code 4" 4 contract purchase-long-cooperate list --tenant-id "$TENANT"
run_expect_fail "missing --tenant-id returns code 4" 4 contract purchase-long-cooperate list --token "$TOKEN"
echo

# --- 采购长协 ---
echo "--- [api] 采购长协 (purchase-long-cooperate) ---"
run_cmd "list (默认第1页)" contract purchase-long-cooperate list
run_cmd "list --page-no 1 --page-size 5" contract purchase-long-cooperate list --page-no 1 --page-size 5
run_cmd "list --keyword XY20260403000001" contract purchase-long-cooperate list --keyword XY20260403000001
run_cmd "list --enterprise-id 2001489305039032322" contract purchase-long-cooperate list --enterprise-id 2001489305039032322
run_cmd "list --order-start --order-end" contract purchase-long-cooperate list --order-start "2026-07-07 00:00:00" --order-end "2026-08-13 23:59:59"
run_cmd "list --end-start --end-end" contract purchase-long-cooperate list --end-start "2026-07-13 00:00:00" --end-end "2026-08-13 23:59:59"
run_cmd "page-count" contract purchase-long-cooperate page-count
run_cmd "page-count --keyword XY20260403000001" contract purchase-long-cooperate page-count --keyword XY20260403000001
echo

# --- 销售合同 ---
echo "--- [api] 销售合同 (sale-contract) ---"
run_cmd "list (默认第1页)" contract sale-contract list
run_cmd "list --page-no 1 --page-size 5" contract sale-contract list --page-no 1 --page-size 5
run_cmd "list --keyword HT20260401000001" contract sale-contract list --keyword HT20260401000001
run_cmd "list --enterprise-id 2001494968037298178" contract sale-contract list --enterprise-id 2001494968037298178
run_cmd "list --order-start --order-end" contract sale-contract list --order-start "2026-07-07 00:00:00" --order-end "2026-08-04 23:59:59"
run_cmd "page-count" contract sale-contract page-count
run_cmd "page-count --keyword HT20260401000001" contract sale-contract page-count --keyword HT20260401000001
echo

# --- 销售长协 ---
echo "--- [api] 销售长协 (sale-long-cooperate) ---"
run_cmd "list (默认第1页)" contract sale-long-cooperate list
run_cmd "list --page-no 1 --page-size 5" contract sale-long-cooperate list --page-no 1 --page-size 5
run_cmd "list --keyword XY20260401000001" contract sale-long-cooperate list --keyword XY20260401000001
run_cmd "list --order-start --order-end" contract sale-long-cooperate list --order-start "2026-07-17 00:00:00" --order-end "2026-08-03 23:59:59"
run_cmd "list --end-start --end-end" contract sale-long-cooperate list --end-start "2026-07-27 00:00:00" --end-end "2026-08-18 23:59:59"
run_cmd "page-count" contract sale-long-cooperate page-count
run_cmd "page-count --keyword XY20260401000001" contract sale-long-cooperate page-count --keyword XY20260401000001
echo

# --- 运输合同 ---
echo "--- [api] 运输合同 (transport) ---"
run_cmd "list (默认第1页)" contract transport list
run_cmd "list --page-no 1 --page-size 5" contract transport list --page-no 1 --page-size 5
run_cmd "list --keyword HT20260319000001" contract transport list --keyword HT20260319000001
run_cmd "list --enterprise-id 2001552070697033730" contract transport list --enterprise-id 2001552070697033730
run_cmd "list --order-start --order-end" contract transport list --order-start "2026-07-09 00:00:00" --order-end "2026-08-05 23:59:59"
run_cmd "list --end-start --end-end" contract transport list --end-start "2026-07-07 00:00:00" --end-end "2026-08-12 23:59:59"
run_cmd "page-count" contract transport page-count
run_cmd "page-count --keyword HT20260319000001" contract transport page-count --keyword HT20260319000001
echo

# --- 运输长协 ---
echo "--- [api] 运输长协 (transport-long) ---"
run_cmd "list (默认第1页)" contract transport-long list
run_cmd "list --page-no 1 --page-size 5" contract transport-long list --page-no 1 --page-size 5
run_cmd "list --keyword XY20260515000002" contract transport-long list --keyword XY20260515000002
run_cmd "list --order-start --order-end" contract transport-long list --order-start "2026-07-15 00:00:00" --order-end "2026-08-18 23:59:59"
run_cmd "list --end-start --end-end" contract transport-long list --end-start "2026-07-13 00:00:00" --end-end "2026-08-13 23:59:59"
run_cmd "page-count" contract transport-long page-count
run_cmd "page-count --keyword XY20260515000002" contract transport-long page-count --keyword XY20260515000002
echo

# --- 服务合同 ---
echo "--- [api] 服务合同 (service-contract) ---"
run_cmd "list (默认第1页)" contract service-contract list
run_cmd "list --page-no 1 --page-size 5" contract service-contract list --page-no 1 --page-size 5
run_cmd "list --keyword HT20260402000001" contract service-contract list --keyword HT20260402000001
run_cmd "list --enterprise-id 2001552070697033730" contract service-contract list --enterprise-id 2001552070697033730
run_cmd "list --order-start --order-end" contract service-contract list --order-start "2026-07-15 00:00:00" --order-end "2026-08-12 23:59:59"
run_cmd "list --end-start --end-end" contract service-contract list --end-start "2026-07-16 00:00:00" --end-end "2026-08-12 23:59:59"
run_cmd "page-count" contract service-contract page-count
run_cmd "page-count --keyword HT20260402000001" contract service-contract page-count --keyword HT20260402000001
echo

# --- 服务长协 ---
echo "--- [api] 服务长协 (service-long) ---"
run_cmd "list (默认第1页)" contract service-long list
run_cmd "list --page-no 1 --page-size 5" contract service-long list --page-no 1 --page-size 5
run_cmd "list --keyword XY20260409000001" contract service-long list --keyword XY20260409000001
run_cmd "list --enterprise-id 2003369636046413825" contract service-long list --enterprise-id 2003369636046413825
run_cmd "list --order-start --order-end" contract service-long list --order-start "2026-07-15 00:00:00" --order-end "2026-08-12 23:59:59"
run_cmd "list --end-start --end-end" contract service-long list --end-start "2026-07-14 00:00:00" --end-end "2026-08-15 23:59:59"
run_cmd "page-count" contract service-long page-count
run_cmd "page-count --keyword XY20260409000001" contract service-long page-count --keyword XY20260409000001
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
