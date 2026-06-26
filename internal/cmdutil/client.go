package cmdutil

import (
	"fmt"

	"github.com/weiliantong/cli/internal/client"
	"github.com/weiliantong/cli/internal/config"
	"github.com/weiliantong/cli/internal/output"
)

var (
	CfgMgr    *config.Manager
	APIClient *client.Client
	// AuthFlags holds the per-call authentication flags parsed from the
	// command line. Auth is stateless: the caller (e.g. an AI agent) supplies
	// a fresh token and tenant-id on every invocation — nothing is persisted.
	AuthFlags struct {
		Token    string
		TenantID string
		BaseURL  string // optional override of the profile's base_url
	}
)

// InitManagers sets up the configuration manager.
func InitManagers(cfg *config.Manager) {
	CfgMgr = cfg
}

// SetAuthFlags stores the auth-related flags (called from root's PersistentPreRunE).
func SetAuthFlags(token, tenantID, baseURL string) {
	AuthFlags.Token = token
	AuthFlags.TenantID = tenantID
	AuthFlags.BaseURL = baseURL
}

// EnsureClient initializes the API client on demand from the active profile
// plus the per-call auth flags. It fails fast (exit code 4) if the required
// --token / --tenant-id flags are missing.
func EnsureClient() error {
	if APIClient != nil {
		return nil
	}
	if CfgMgr == nil {
		return fmt.Errorf("配置未加载")
	}
	profile, err := CfgMgr.ActiveProfile()
	if err != nil {
		return output.NewExitError(2, fmt.Sprintf("当前环境配置不存在: %s", err), "运行 wlt config init 初始化")
	}
	if AuthFlags.Token == "" || AuthFlags.TenantID == "" {
		return output.NewExitError(4, "缺少必填鉴权参数",
			"请通过 --token 与 --tenant-id 传入鉴权信息(对应 Authorization 与 tenant-id 请求头)")
	}
	baseURL := profile.BaseURL
	if AuthFlags.BaseURL != "" {
		baseURL = AuthFlags.BaseURL
	}
	APIClient = client.NewClient(client.RequestContext{
		BaseURL:        baseURL,
		APIPrefix:      profile.APIPrefix,
		TenantID:       AuthFlags.TenantID,
		EnterpriseType: profile.EnterpriseType,
		Token:          AuthFlags.Token,
	})
	return nil
}

// GetClient returns the API client.
func GetClient() *client.Client {
	return APIClient
}
