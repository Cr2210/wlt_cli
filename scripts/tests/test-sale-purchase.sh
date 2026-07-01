#!/usr/bin/env bash
# scripts/tests/test-sale-purchase.sh
# wlt sale/purchase 子模块端到端测试(销售出库/采购入库 — list/page-count 新参数)
#
# 用法:
#   WLT_TOKEN=<accessToken> WLT_TENANT_ID=<tenantId> ./scripts/tests/test-sale-purchase.sh
#   ./scripts/tests/test-sale-purchase.sh             # 只跑 smoke 段(无需 token)
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
# 本次 CLI 优化点: purchase-in / sale-out 新增筛选字段
#   purchase-in  list/page-count  +supplier-id +metrics-name +(已有 product-name / start-time / end-time → inTime[0]/inTime[1])
#   sale-out     list/page-count  +customer-id +batch-no +(已有 product-name / start-time / end-time → outTime[0]/outTime[1])

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

# === Helpers ===

# run_cmd: 期望返回成功 JSON(data 字段)
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

# run_plain: 任意成功(不校验 JSON)
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
echo " wlt sale/purchase 子模块测试(销售出库/采购入库)"
echo " Binary: $WLT_BIN"
[ -n "$TOKEN" ] && echo " Tenant: $TENANT" || echo " Tenant: (无 token,只跑 smoke 段)"
echo "=================================================="
echo

# === smoke 段(无需 token) ===

echo "--- [smoke] 二进制基本可用性 ---"
run_plain "wlt --version 成功" version
run_plain "wlt purchase --help 成功" purchase --help
run_plain "wlt purchase in --help 成功" purchase in --help
run_plain "wlt purchase return --help 成功" purchase return --help
run_plain "wlt sale --help 成功" sale --help
run_plain "wlt sale out --help 成功" sale out --help
run_plain "wlt sale return --help 成功" sale return --help
echo

# --- purchase-in list --help 新参数检查 ---
echo "--- [smoke] purchase-in list 新筛选字段 ---"
run_help "list --help 列出 --warehouse-id"       "--warehouse-id"    purchase in list --help
run_help "list --help 列出 --product-id"          "--product-id"      purchase in list --help
run_help "list --help 列出 --product-name(新增)"  "--product-name"    purchase in list --help
run_help "list --help 列出 --no"                  "--no"              purchase in list --help
run_help "list --help 列出 --status"              "--status"          purchase in list --help
run_help "list --help 列出 --type"                "--type"            purchase in list --help
run_help "list --help 列出 --start-time"          "--start-time"      purchase in list --help
run_help "list --help 列出 --end-time"            "--end-time"        purchase in list --help
run_help "list --help 列出 --supplier-id(新增)"   "--supplier-id"     purchase in list --help
run_help "list --help 列出 --metrics-name(新增)"  "--metrics-name"    purchase in list --help
run_help "list --help 列出 --page-no"             "--page-no"         purchase in list --help
run_help "list --help 列出 --page-size"           "--page-size"       purchase in list --help
echo

# --- purchase-in page-count --help 新参数检查 ---
echo "--- [smoke] purchase-in page-count 新筛选字段 ---"
run_help "page-count --help 列出 --warehouse-id"      "--warehouse-id"    purchase in page-count --help
run_help "page-count --help 列出 --product-name(新增)" "--product-name"    purchase in page-count --help
run_help "page-count --help 列出 --start-time"         "--start-time"      purchase in page-count --help
run_help "page-count --help 列出 --end-time"           "--end-time"        purchase in page-count --help
run_help "page-count --help 列出 --supplier-id(新增)"  "--supplier-id"     purchase in page-count --help
run_help "page-count --help 列出 --metrics-name(新增)" "--metrics-name"    purchase in page-count --help
echo

# --- sale-out list --help 新参数检查 ---
echo "--- [smoke] sale-out list 新筛选字段 ---"
run_help "list --help 列出 --warehouse-id"       "--warehouse-id"    sale out list --help
run_help "list --help 列出 --product-id"          "--product-id"      sale out list --help
run_help "list --help 列出 --product-name(新增)"  "--product-name"    sale out list --help
run_help "list --help 列出 --no"                  "--no"              sale out list --help
run_help "list --help 列出 --status"              "--status"          sale out list --help
run_help "list --help 列出 --type"                "--type"            sale out list --help
run_help "list --help 列出 --start-time"          "--start-time"      sale out list --help
run_help "list --help 列出 --end-time"            "--end-time"        sale out list --help
run_help "list --help 列出 --customer-id(新增)"   "--customer-id"     sale out list --help
run_help "list --help 列出 --batch-no(新增)"      "--batch-no"        sale out list --help
run_help "list --help 列出 --page-no"             "--page-no"         sale out list --help
run_help "list --help 列出 --page-size"           "--page-size"       sale out list --help
echo

# --- sale-out page-count --help 新参数检查 ---
echo "--- [smoke] sale-out page-count 新筛选字段 ---"
run_help "page-count --help 列出 --warehouse-id"      "--warehouse-id"    sale out page-count --help
run_help "page-count --help 列出 --product-name(新增)" "--product-name"    sale out page-count --help
run_help "page-count --help 列出 --start-time"         "--start-time"      sale out page-count --help
run_help "page-count --help 列出 --end-time"           "--end-time"        sale out page-count --help
run_help "page-count --help 列出 --customer-id(新增)"  "--customer-id"     sale out page-count --help
run_help "page-count --help 列出 --batch-no(新增)"     "--batch-no"        sale out page-count --help
echo

# --- purchase-in / sale-out CRUD smoke(复用现有 get/create/update/delete/update-status) ---
echo "--- [smoke] purchase-in / sale-out 通用 CRUD ---"
run_help "purchase in get --help 列出 --id"                "--id"          purchase in get --help
run_help "purchase in create --help 列出 --data"            "--data"        purchase in create --help
run_help "purchase in update --help 列出 --data"            "--data"        purchase in update --help
run_help "purchase in delete --help 列出 --ids"             "--ids"         purchase in delete --help
run_help "purchase in update-status --help 列出 --data"     "--data"        purchase in update-status --help
run_help "sale out get --help 列出 --id"                    "--id"          sale out get --help
run_help "sale out create --help 列表 --data"               "--data"        sale out create --help
run_help "sale out update --help 列表 --data"               "--data"        sale out update --help
run_help "sale out delete --help 列表 --ids"                "--ids"         sale out delete --help
run_help "sale out update-status --help 列表 --data"        "--data"        sale out update-status --help
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
run_expect_fail "purchase in list missing --token returns code 4"    4 purchase in list --tenant-id "$TENANT"
run_expect_fail "purchase in list missing --tenant-id returns code 4" 4 purchase in list --token "$TOKEN"
run_expect_fail "sale out list missing --token returns code 4"        4 sale out list    --tenant-id "$TENANT"
run_expect_fail "sale out list missing --tenant-id returns code 4"    4 sale out list    --token "$TOKEN"
echo

# --- purchase-in API ---
echo "--- [api] purchase-in (采购入库) ---"
run_cmd "list (默认第1页)"                                       purchase in list
run_cmd "list --page-no 1 --page-size 5"                          purchase in list --page-no 1 --page-size 5
run_cmd "list --no CGRK20260508000009"                           purchase in list --no CGRK20260508000009
run_cmd "list --product-name 西瓜"                               purchase in list --product-name 西瓜
run_cmd "list --supplier-id 2001489305039032322"                 purchase in list --supplier-id 2001489305039032322
run_cmd "list --metrics-name 含水"                               purchase in list --metrics-name 含水
run_cmd "list --inTime-range(自动转 inTime[0]/inTime[1])"         purchase in list --start-time "2026-07-08 00:00:00" --end-time "2026-08-05 23:59:59"
run_cmd "list 组合: supplier + product-name + warehouse"         purchase in list --supplier-id 2001489305039032322 --product-name 西瓜 --warehouse-id 7
run_cmd "page-count"                                             purchase in page-count
run_cmd "page-count --product-name 西瓜"                         purchase in page-count --product-name 西瓜
run_cmd "page-count --metrics-name 含水"                         purchase in page-count --metrics-name 含水
run_cmd "page-count --start-time + --end-time"                    purchase in page-count --start-time "2026-07-08 00:00:00" --end-time "2026-08-05 23:59:59"
echo

# --- sale-out API ---
echo "--- [api] sale-out (销售出库) ---"
run_cmd "list (默认第1页)"                                       sale out list
run_cmd "list --page-no 1 --page-size 5"                          sale out list --page-no 1 --page-size 5
run_cmd "list --no XSCK20260402000001"                           sale out list --no XSCK20260402000001
run_cmd "list --product-name 习惯"                               sale out list --product-name 习惯
run_cmd "list --batch-no PC"                                     sale out list --batch-no PC
run_cmd "list --customer-id 2001489305039032322"                 sale out list --customer-id 2001489305039032322
run_cmd "list --outTime-range(自动转 outTime[0]/outTime[1])"      sale out list --start-time "2026-07-18 00:00:00" --end-time "2026-08-12 23:59:59"
run_cmd "list 组合: batch-no + product-name + warehouse"         sale out list --batch-no PC --product-name 习惯 --warehouse-id 6
run_cmd "page-count"                                             sale out page-count
run_cmd "page-count --product-name 习惯"                         sale out page-count --product-name 习惯
run_cmd "page-count --batch-no PC"                               sale out page-count --batch-no PC
run_cmd "page-count --start-time + --end-time"                    sale out page-count --start-time "2026-07-18 00:00:00" --end-time "2026-08-12 23:59:59"
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
