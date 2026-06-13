package system

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weiliantong/cli/internal/cmdutil"
	"github.com/weiliantong/cli/internal/output"
)

func init() {
	SystemCmd.AddCommand(sysUserCmd)
}

var sysUserCmd = &cobra.Command{
	Use:   "user",
	Short: "用户管理",
}

func init() {
	sysUserCmd.AddCommand(
		newSysUserListCmd(),
		newSysUserSimpleListCmd(),
		newSysUserGetCmd(),
		newSysUserCreateCmd(),
		newSysUserUpdateCmd(),
		newSysUserDeleteCmd(),
		newSysUserUpdatePasswordCmd(),
		newSysUserUpdateStatusCmd(),
		newSysUserExportCmd(),
		newSysUserGetImportTemplateCmd(),
		newSysUserImportCmd(),
	)
}

// ---- 分页查询 ----

func newSysUserListCmd() *cobra.Command {
	var pageNo, pageSize int
	var username, mobile, status, deptId string

	c := &cobra.Command{
		Use:   "list",
		Short: "分页查询用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			cmdutil.CollectStringFlags(cmd, params, "username", "mobile", "status", "dept-id")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/user/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询用户失败: %s", err), "")
			}
			return cmdutil.ParsePagedJSON(resp.Data, pageNo, pageSize)
		},
	}
	c.Flags().IntVar(&pageNo, "page-no", 1, "页码")
	c.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")
	c.Flags().StringVar(&username, "username", "", "用户名")
	c.Flags().StringVar(&mobile, "mobile", "", "手机号")
	c.Flags().StringVar(&status, "status", "", "状态")
	c.Flags().StringVar(&deptId, "dept-id", "", "部门 ID")
	return c
}

// ---- 简单列表 ----

func newSysUserSimpleListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "simple-list",
		Short: "获取用户精简列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/user/simple-list", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取用户列表失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 获取详情 ----

func newSysUserGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "get",
		Short: "获取用户详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/user/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取用户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "用户 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建 ----

func newSysUserCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "create",
		Short: "创建用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/system/user/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建用户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新 ----

func newSysUserUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update",
		Short: "更新用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/system/user/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新用户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除 ----

func newSysUserDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "delete",
		Short: "删除用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/system/user/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除用户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "用户 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 更新密码 ----

func newSysUserUpdatePasswordCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-password",
		Short: "更新用户密码",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/system/user/update-password", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新密码失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新状态 ----

func newSysUserUpdateStatusCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "update-status",
		Short: "更新用户状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/system/user/update-status", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新用户状态失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 导出 ----

func newSysUserExportCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "导出用户 Excel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/user/export-excel", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导出用户失败: %s", err), "")
			}
			fmt.Println("导出成功，返回数据：", string(resp.Data))
			return nil
		},
	}
	return c
}

// ---- 获取导入模板 ----

func newSysUserGetImportTemplateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "get-import-template",
		Short: "获取用户导入模板",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/user/get-import-template", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取导入模板失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 导入 ----

func newSysUserImportCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "import",
		Short: "导入用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/system/user/import", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("导入用户失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}
