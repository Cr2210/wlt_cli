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
	SystemCmd.AddCommand(sysDictCmd)
}

var sysDictCmd = &cobra.Command{
	Use:   "dict",
	Short: "字典管理",
}

func init() {
	sysDictCmd.AddCommand(
		newSysDictTypeListCmd(),
		newSysDictTypeGetCmd(),
		newSysDictTypeCreateCmd(),
		newSysDictTypeUpdateCmd(),
		newSysDictTypeDeleteCmd(),
		newSysDictDataListCmd(),
		newSysDictDataGetCmd(),
		newSysDictDataCreateCmd(),
		newSysDictDataUpdateCmd(),
		newSysDictDataDeleteCmd(),
	)
}

// ---- 字典类型列表 ----

func newSysDictTypeListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "type-list",
		Short: "查询字典类型列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/dict-type/list", nil)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询字典类型失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	return c
}

// ---- 字典类型详情 ----

func newSysDictTypeGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "type-get",
		Short: "获取字典类型详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/dict-type/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取字典类型失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "字典类型 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建字典类型 ----

func newSysDictTypeCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "type-create",
		Short: "创建字典类型",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/system/dict-type/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建字典类型失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新字典类型 ----

func newSysDictTypeUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "type-update",
		Short: "更新字典类型",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/system/dict-type/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新字典类型失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除字典类型 ----

func newSysDictTypeDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "type-delete",
		Short: "删除字典类型",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/system/dict-type/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除字典类型失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "字典类型 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 字典数据列表 ----

func newSysDictDataListCmd() *cobra.Command {
	var dictType, status string

	c := &cobra.Command{
		Use:   "data-list",
		Short: "查询字典数据列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			params := map[string]any{}
			cmdutil.CollectStringFlags(cmd, params, "dict-type", "status")

			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/dict-data/page", params)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("查询字典数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&dictType, "dict-type", "", "字典类型")
	c.Flags().StringVar(&status, "status", "", "状态")
	return c
}

// ---- 字典数据详情 ----

func newSysDictDataGetCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "data-get",
		Short: "获取字典数据详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Get(context.Background(), "/system/dict-data/get", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("获取字典数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "字典数据 ID")
	_ = c.MarkFlagRequired("id")
	return c
}

// ---- 创建字典数据 ----

func newSysDictDataCreateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "data-create",
		Short: "创建字典数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Post(context.Background(), "/system/dict-data/create", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("创建字典数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 更新字典数据 ----

func newSysDictDataUpdateCmd() *cobra.Command {
	var data string

	c := &cobra.Command{
		Use:   "data-update",
		Short: "更新字典数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			body, err := cmdutil.ParseJSONData(data)
			if err != nil {
				return output.NewExitError(4, fmt.Sprintf("解析 data 失败: %s", err), "data 应为 JSON 对象")
			}
			resp, err := cmdutil.GetClient().Put(context.Background(), "/system/dict-data/update", body)
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("更新字典数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().StringVar(&data, "data", "", "JSON 数据")
	_ = c.MarkFlagRequired("data")
	return c
}

// ---- 删除字典数据 ----

func newSysDictDataDeleteCmd() *cobra.Command {
	var id int64

	c := &cobra.Command{
		Use:   "data-delete",
		Short: "删除字典数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.EnsureClient(); err != nil {
				return err
			}
			resp, err := cmdutil.GetClient().Delete(context.Background(), "/system/dict-data/delete", map[string]any{"id": id})
			if err != nil {
				return output.NewExitError(5, fmt.Sprintf("删除字典数据失败: %s", err), "")
			}
			return cmdutil.OutputJSON(json.RawMessage(resp.Data))
		},
	}
	c.Flags().Int64Var(&id, "id", 0, "字典数据 ID")
	_ = c.MarkFlagRequired("id")
	return c
}
