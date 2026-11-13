package finance

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// snapshotCmd 把一个子命令渲染成稳定字符串：name|short|flag:type:help;…
// flag 按名字排序，保证输出与注册顺序无关。
func snapshotCmd(c *cobra.Command) string {
	var flags []string
	c.Flags().VisitAll(func(f *pflag.Flag) {
		flags = append(flags, fmt.Sprintf("%s:%s:%s", f.Name, f.Value.Type(), f.Usage))
	})
	sort.Strings(flags)
	return fmt.Sprintf("%s|%s|%s", c.Name(), c.Short, strings.Join(flags, ";"))
}

// findGroup 在 financeCmd 下按名字找子组。
func findGroup(t *testing.T, name string) *cobra.Command {
	t.Helper()
	for _, c := range financeCmd.Commands() {
		if c.Name() == name {
			return c
		}
	}
	t.Fatalf("finance 子组 %q 不存在", name)
	return nil
}

// findSub 在子组下按名字找子命令。
func findSub(t *testing.T, group *cobra.Command, name string) *cobra.Command {
	t.Helper()
	for _, c := range group.Commands() {
		if c.Name() == name {
			return c
		}
	}
	t.Fatalf("子命令 %q %q 不存在", group.Name(), name)
	return nil
}

// migratedSnapshots 锁定由 cmdutil.Crud*Cmd 工厂生成的命令的 CLI 表面。
//
// 这些命令原先是手写的等价实现，重构为工厂后必须保持：子命令名、Short 文案、
// flag 集合（名字/类型/帮助文本）逐字不变。期望值取自重构前的真实命令树输出。
//
// 注意：不在表中的命令（如 refund get 支持 --id/--no、payment-apply update-status
// 走 PutQuery+approveStatus）是刻意保留的手写实现，工厂无法等价表达。
var migratedSnapshots = map[string]map[string]string{
	"payment-apply": {
		"get":    "get|获取付款申请详情|id:int64:付款申请 ID",
		"create": "create|创建付款申请|data:string:JSON 数据",
		"update": "update|更新付款申请|data:string:JSON 数据",
		"delete": "delete|删除付款申请|ids:string:ID 列表（逗号分隔）",
	},
	"invoice": {
		"get":    "get|获取收开票详情|id:int64:收开票 ID",
		"create": "create|创建收开票|data:string:JSON 数据",
		"update": "update|更新收开票|data:string:JSON 数据",
		"delete": "delete|删除收开票|ids:string:ID 列表（逗号分隔）",
	},
	"prepayment-apply": {
		"get":    "get|获取预付申请详情|id:int64:预付申请 ID",
		"create": "create|创建预付申请|data:string:JSON 数据",
		"update": "update|更新预付申请|data:string:JSON 数据",
		"delete": "delete|删除预付申请|ids:string:ID 列表（逗号分隔）",
	},
	"refund": {
		"create":        "create|创建退款单|data:string:JSON 数据",
		"update":        "update|更新退款单|data:string:JSON 数据",
		"delete":        "delete|删除退款单|id:int64:退款单 ID",
		"update-status": "update-status|更新退款单状态|data:string:JSON 数据（含 id 和 status）",
	},
	"settlement": {
		"get":           "get|获取结算单详情|id:int64:结算单 ID",
		"create":        "create|创建结算单|data:string:JSON 数据",
		"update":        "update|更新结算单|data:string:JSON 数据",
		"delete":        "delete|删除结算单|id:int64:结算单 ID",
		"update-status": "update-status|更新结算单状态|data:string:JSON 数据（含 id 和 status）",
	},
	"write-off": {
		"get":           "get|获取核销单详情|id:int64:核销单 ID",
		"create":        "create|创建核销单|data:string:JSON 数据",
		"update":        "update|更新核销单|data:string:JSON 数据",
		"delete":        "delete|删除核销单|id:int64:核销单 ID",
		"update-status": "update-status|更新核销单状态|data:string:JSON 数据（含 id 和 status）",
	},
	"receipt-payment": {
		"create":        "create|创建收付款|data:string:JSON 数据",
		"update":        "update|更新收付款|data:string:JSON 数据",
		"delete":        "delete|删除收付款|id:int64:收付款 ID",
		"update-status": "update-status|更新收付款状态|data:string:JSON 数据（含 id 和 status）",
	},
}

func TestMigratedCommandSurfaces(t *testing.T) {
	for group, cmds := range migratedSnapshots {
		g := findGroup(t, group)
		for name, want := range cmds {
			got := snapshotCmd(findSub(t, g, name))
			if got != want {
				t.Errorf("%s %s 的 CLI 表面与重构前不一致\n got: %s\nwant: %s", group, name, got, want)
			}
		}
	}
}
