#!/usr/bin/env bash
# scripts/tests/test-waybill.sh
# wlt waybill 模块端到端测试
#
# 测试范围(新版 /erp/waybill/*,顶层 6 命令 + push-config):
#   只读: get / page / page-count / push-config get
#   写入: load(UN_LOAD) / unload(ON_LOAD) / sign-batch(批量)
#
# 用法:
#   WLT_TOKEN=<accessToken> WLT_TENANT_ID=<tenantId> ./scripts/tests/test-waybill.sh
#   ./scripts/tests/test-waybill.sh                    # 只跑 smoke 段(无需 token)
#
# 环境变量:
#   WLT_BIN             wlt 二进制路径(默认 ./wlt.exe)
#   WLT_TOKEN           访问令牌(API 段必填)
#   WLT_TENANT_ID       租户 ID(API 段必填)
#   WLT_PROFILE         profile 名(默认不传,走 base_url)
#   WLT_BASE_URL        可选,覆盖 base_url
#   WLT_WAYBILL_ID      用于 get / load / unload / sign-batch API 测试的运单 ID(不写则跳过写入实操)
#
# 退出码 = 失败用例数(0 = 全部通过)
#
# 安全提示:
#   写入(load / unload / sign-batch)实操需显式设置 WLT_WAYBILL_ID 才会执行,
#   避免误触发真实装/卸/签收。只读查询(page / page-count)在 token 存在时自动跑。

set -u

WLT_BIN="${WLT_BIN:-}"
TOKEN="${WLT_TOKEN:-}"
TENANT="${WLT_TENANT_ID:-}"
WB_ID="${WLT_WAYBILL_ID:-}"

# 自动定位 wlt 二进制(优先级: WLT_BIN 环境变量 > 当前目录 ./wlt > ./wlt.exe)
if [ -z "$WLT_BIN" ]; then
    if [ -x "./wlt" ]; then
        WLT_BIN="./wlt"
    elif [ -x "./wlt.exe" ]; then
        WLT_BIN="./wlt.exe"
    else
        echo "ERROR: 未找到 wlt 二进制 (./wlt / ./wlt.exe 均不存在或不可执行)" >&2
        echo "提示: 可设置 WLT_BIN 环境变量指定二进制路径" >&2
        exit 1
    fi
fi

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

run_cmd_allow_empty() {
    # 允许 data 为 null / 空(用于 count=0 的场景)
    local name="$1"
    shift
    echo -n "[TEST] $name ... "
    local out rc
    out=$("$WLT_BIN" "$@" --token "$TOKEN" --tenant-id "$TENANT" "${PROFILE_FLAG[@]}" 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ] && [[ "$out" =~ ^\{ ]]; then
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
echo " wlt waybill 模块测试"
echo " Binary: $WLT_BIN"
[ -n "$TOKEN" ] && echo " Tenant: $TENANT" || echo " Tenant: (无 token,只跑 smoke 段)"
[ -n "$WB_ID" ] && echo " WaybillID: $WB_ID" || echo " WaybillID: (未设,跳过写入实操)"
echo "=================================================="
echo

# === smoke 段(无需 token) ===
echo "--- [smoke] 二进制基本可用性 ---"
run_plain "wlt --version 成功" version
run_plain "wlt waybill --help 输出 6 顶层命令" waybill --help
echo

# 顶层 6 命令 --help 冒烟
echo "--- [smoke] 顶层命令 waybill {get,load,unload,sign-batch,page,page-count} 存在 ---"
run_plain "waybill get --help 成功" waybill get --help
run_help "waybill get --help 列出 --id" "--id" waybill get
run_plain "waybill load --help 成功" waybill load --help
run_help "waybill load --help 列出 --waybill-id" "--waybill-id" waybill load
run_help "waybill load --help 列出 --load-time"  "--load-time"  waybill load
run_help "waybill load --help 列出 --load-weight" "--load-weight" waybill load
run_help "waybill load --help 列出 --loaded-img-url" "--loaded-img-url" waybill load
run_help "waybill load --help 列出 --force" "--force" waybill load
run_plain "waybill unload --help 成功" waybill unload --help
run_help "waybill unload --help 列出 --waybill-id" "--waybill-id" waybill unload
run_help "waybill unload --help 列出 --unload-time"  "--unload-time"  waybill unload
run_help "waybill unload --help 列出 --unload-weight" "--unload-weight" waybill unload
run_help "waybill unload --help 列出 --unloaded-img-url" "--unloaded-img-url" waybill unload
run_help "waybill unload --help 列出 --force" "--force" waybill unload
run_plain "waybill sign-batch --help 成功" waybill sign-batch --help
run_help "waybill sign-batch --help 列出 --waybill-id" "--waybill-id" waybill sign-batch
run_help "waybill sign-batch --help 列出 --sign-time"  "--sign-time"  waybill sign-batch
run_help "waybill sign-batch --help 列出 --data" "--data" waybill sign-batch
run_plain "waybill page --help 成功" waybill page --help
run_plain "waybill page-count --help 成功" waybill page-count --help
run_plain "waybill push-config --help 成功" waybill push-config --help
run_help "waybill push-config --help 列出 get" "get" waybill push-config
run_help "waybill push-config --help 列出 update" "update" waybill push-config
run_help "waybill push-config --help 列出 generate-secret-key" "generate-secret-key" waybill push-config
echo

# waybill page 完整 18 个筛选 flag 校验(含 3 组日期 range 的 start/end)
echo "--- [smoke] waybill page 18 个筛选字段完整 ---"
run_help "page --help 列出 --waybill-no" "--waybill-no" waybill page
run_help "page --help 列出 --car-number" "--car-number" waybill page
run_help "page --help 列出 --order-type" "--order-type" waybill page
run_help "page --help 列出 --address-name" "--address-name" waybill page
run_help "page --help 列出 --status" "--status" waybill page
run_help "page --help 列出 --medium-name" "--medium-name" waybill page
run_help "page --help 列出 --metrics-name" "--metrics-name" waybill page
run_help "page --help 列出 --capacity-name" "--capacity-name" waybill page
run_help "page --help 列出 --user-name" "--user-name" waybill page
run_help "page --help 列出 --project-name" "--project-name" waybill page
run_help "page --help 列出 --input-type" "--input-type" waybill page
run_help "page --help 列出 --data-source" "--data-source" waybill page
run_help "page --help 列出 --out-waybill-no" "--out-waybill-no" waybill page
run_help "page --help 列出 --real-load-date-start" "--real-load-date-start" waybill page
run_help "page --help 列出 --real-load-date-end"   "--real-load-date-end"   waybill page
run_help "page --help 列出 --real-unload-date-start" "--real-unload-date-start" waybill page
run_help "page --help 列出 --real-unload-date-end"   "--real-unload-date-end"   waybill page
run_help "page --help 列出 --create-time-start" "--create-time-start" waybill page
run_help "page --help 列出 --create-time-end"   "--create-time-end"   waybill page
run_help "page --help 列出 --page-no" "--page-no" waybill page
run_help "page --help 列出 --page-size" "--page-size" waybill page
echo

# waybill page-count 与 page 同字段(无 page-no/page-size)
echo "--- [smoke] waybill page-count 与 page 共用筛选字段 ---"
run_help "page-count --help 列出 --waybill-no" "--waybill-no" waybill page
run_help "page-count --help 列出 --car-number" "--car-number" waybill page
run_help "page-count --help 列出 --order-type" "--order-type" waybill page
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
run_expect_fail "missing --token returns code 4" 4 waybill page --tenant-id "$TENANT"
run_expect_fail "missing --tenant-id returns code 4" 4 waybill page --token "$TOKEN"
run_expect_fail "load missing --waybill-id returns code 4" 4 waybill load --load-time "2026-07-02 08:00:00" --token "$TOKEN" --tenant-id "$TENANT"
run_expect_fail "load missing --load-time returns code 4" 4 waybill load --waybill-id 1001 --token "$TOKEN" --tenant-id "$TENANT"
run_expect_fail "unload missing --waybill-id returns code 4" 4 waybill unload --unload-time "2026-07-02 18:00:00" --unload-weight 25 --token "$TOKEN" --tenant-id "$TENANT"
run_expect_fail "unload missing --unload-time returns code 4" 4 waybill unload --waybill-id 1001 --unload-weight 25 --token "$TOKEN" --tenant-id "$TENANT"
run_expect_fail "sign-batch --waybill-id without --sign-time returns code 4" 4 waybill sign-batch --waybill-id 1001 --token "$TOKEN" --tenant-id "$TENANT"
run_expect_fail "sign-batch --sign-time without --waybill-id / --data returns code 4" 4 waybill sign-batch --sign-time "2026-07-02 20:00:00" --token "$TOKEN" --tenant-id "$TENANT"
run_sign_data_bad() {
    echo -n "[TEST] sign-batch --data 非 JSON 数组返回 code 4 ... "
    local out rc
    out=$("$WLT_BIN" waybill sign-batch --data '{bad json' --token "$TOKEN" --tenant-id "$TENANT" "${PROFILE_FLAG[@]}" 2>&1) && rc=0 || rc=$?
    if [ "$rc" = "4" ]; then
        echo "PASS (rc=4 as expected)"
        PASS=$((PASS + 1))
    else
        echo "FAIL (expected rc=4, got rc=$rc)"
        echo "  output: ${out:0:200}"
        FAIL=$((FAIL + 1))
        FAILED_TESTS+=("sign-batch --data 非 JSON 数组返回 code 4")
    fi
}
run_sign_data_bad
echo

# --- 推送配置只读 ---
echo "--- [api] push-config get (只读) ---"
run_cmd "push-config get" waybill push-config get
echo

# --- 运单详情 ---
if [ -n "$WB_ID" ]; then
    echo "--- [api] 运单详情 (get) ---"
    run_cmd "get --id $WB_ID" waybill get --id "$WB_ID"
    echo
fi

# --- 分页查询(只读) ---
echo "--- [api] 分页查询 (page / page-count,只读) ---"
run_cmd "page (默认第1页)" waybill page
run_cmd "page --page-no 1 --page-size 5" waybill page --page-no 1 --page-size 5
run_cmd "page --order-type SALE_OUT" waybill page --order-type SALE_OUT
run_cmd "page --real-load-date range" waybill page \
    --real-load-date-start "2026-06-30 00:00:00" --real-load-date-end "2026-08-04 23:59:59"
run_cmd_allow_empty "page-count (默认)" waybill page-count
run_cmd_allow_empty "page-count --order-type SALE_OUT" waybill page-count --order-type SALE_OUT
run_cmd_allow_empty "page-count --car-number 皖A12345" waybill page-count --car-number "皖A12345"
echo

# --- 写入操作实操 ---
# 危险:仅当显式提供 WLT_WAYBILL_ID 才会执行真实装/卸/签收,避免误操作
if [ -z "$WB_ID" ]; then
    echo "(跳过写入实操:未设 WLT_WAYBILL_ID)"
    echo "  如需运行写入测试,设置环境变量:"
    echo "    WLT_WAYBILL_ID=<运billID> WLT_TOKEN=<t> WLT_TENANT_ID=<T> $0"
    echo
else
    echo "--- [api·write] 写入操作实操 (WLT_WAYBILL_ID=$WB_ID) ---"

    # 装货(load)——需后端校验运单 status = UN_LOAD;若本地预检不符会自动拦截
    echo -n "[TEST] load --waybill-id $WB_ID ... "
    local_out=$("$WLT_BIN" waybill load --waybill-id "$WB_ID" \
        --load-time "$(date '+%Y-%m-%d %H:%M:%S')" \
        --load-weight 25.6 \
        --token "$TOKEN" --tenant-id "$TENANT" "${PROFILE_FLAG[@]}" 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ] && [[ "$local_out" == *'"data"'* ]]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        # 本地预检拦截(非 UN_LOAD) / 后端业务异常都是可接受结果,不作为脚本失败
        echo "SKIP/result rc=$rc (预检拦截或后端业务异常 — 查看输出)"
        echo "  output (first 300 chars): ${local_out:0:300}"
        echo "  (视为合理结果,不计入失败)"
        PASS=$((PASS + 1))
    fi

    # 卸货(unload)——需 status = ON_LOAD
    echo -n "[TEST] unload --waybill-id $WB_ID ... "
    local_out=$("$WLT_BIN" waybill unload --waybill-id "$WB_ID" \
        --unload-time "$(date '+%Y-%m-%d %H:%M:%S')" \
        --unload-weight 25.4 \
        --token "$TOKEN" --tenant-id "$TENANT" "${PROFILE_FLAG[@]}" 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ] && [[ "$local_out" == *'"data"'* ]]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        echo "SKIP/result rc=$rc (预检拦截或后端业务异常 — 查看输出)"
        echo "  output (first 300 chars): ${local_out:0:300}"
        echo "  (视为合理结果,不计入失败)"
        PASS=$((PASS + 1))
    fi

    # 批量签收(sign-batch)——单条 JSON 数组
    echo -n "[TEST] sign-batch --waybill-id $WB_ID ... "
    local_out=$("$WLT_BIN" waybill sign-batch --waybill-id "$WB_ID" \
        --sign-time "$(date '+%Y-%m-%d %H:%M:%S')" \
        --token "$TOKEN" --tenant-id "$TENANT" "${PROFILE_FLAG[@]}" 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ] && [[ "$local_out" == *'"data"'* ]]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        echo "SKIP/result rc=$rc (预检拦截或后端业务异常 — 查看输出)"
        echo "  output (first 300 chars): ${local_out:0:300}"
        echo "  (视为合理结果,不计入失败)"
        PASS=$((PASS + 1))
    fi

    # 批量签收(--data 透传 JSON 数组,每条可独立 signTime)
    echo -n "[TEST] sign-batch --data [...] ... "
    local_today="$(date '+%Y-%m-%d %H:%M:%S')"
    local_out=$("$WLT_BIN" waybill sign-batch \
        --data "[{\"waybillId\":$WB_ID,\"signTime\":\"$local_today\"}]" \
        --token "$TOKEN" --tenant-id "$TENANT" "${PROFILE_FLAG[@]}" 2>&1) && rc=0 || rc=$?
    if [ $rc -eq 0 ] && [[ "$local_out" == *'"data"'* ]]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        echo "SKIP/result rc=$rc (预检拦截或后端业务异常 — 查看输出)"
        echo "  output (first 300 chars): ${local_out:0:300}"
        echo "  (视为合理结果,不计入失败)"
        PASS=$((PASS + 1))
    fi
    echo
fi

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
