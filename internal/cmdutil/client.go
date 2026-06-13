package cmdutil

import (
	"fmt"

	"github.com/weiliantong/cli/internal/auth"
	"github.com/weiliantong/cli/internal/client"
	"github.com/weiliantong/cli/internal/config"
	"github.com/weiliantong/cli/internal/output"
)

var (
	CfgMgr    *config.Manager
	AuthMgr   *auth.Manager
	APIClient *client.Client
)

// InitManagers sets up the configuration manager.
func InitManagers(cfg *config.Manager) {
	CfgMgr = cfg
}

// EnsureClient initializes auth manager and API client on demand.
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
	AuthMgr = auth.NewManager(CfgMgr)
	APIClient = client.NewClient(profile, AuthMgr)
	return nil
}

// GetClient returns the API client.
func GetClient() *client.Client {
	return APIClient
}

// GetAuthMgr returns or creates the auth manager.
func GetAuthMgr() *auth.Manager {
	if AuthMgr != nil {
		return AuthMgr
	}
	AuthMgr = auth.NewManager(CfgMgr)
	return AuthMgr
}
