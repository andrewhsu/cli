package auth

import (
	gitCredentialCmd "github.com/andrewhsu/cli/v2/pkg/cmd/auth/gitcredential"
	authLoginCmd "github.com/andrewhsu/cli/v2/pkg/cmd/auth/login"
	authLogoutCmd "github.com/andrewhsu/cli/v2/pkg/cmd/auth/logout"
	authRefreshCmd "github.com/andrewhsu/cli/v2/pkg/cmd/auth/refresh"
	authStatusCmd "github.com/andrewhsu/cli/v2/pkg/cmd/auth/status"
	"github.com/andrewhsu/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Login, logout, and refresh your authentication",
		Long:  `Manage gh's authentication state.`,
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(authLoginCmd.NewCmdLogin(f, nil))
	cmd.AddCommand(authLogoutCmd.NewCmdLogout(f, nil))
	cmd.AddCommand(authStatusCmd.NewCmdStatus(f, nil))
	cmd.AddCommand(authRefreshCmd.NewCmdRefresh(f, nil))
	cmd.AddCommand(gitCredentialCmd.NewCmdCredential(f, nil))

	return cmd
}
