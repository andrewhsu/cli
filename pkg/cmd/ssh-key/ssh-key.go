package key

import (
	cmdAdd "github.com/andrewhsu/cli/v2/pkg/cmd/ssh-key/add"
	cmdList "github.com/andrewhsu/cli/v2/pkg/cmd/ssh-key/list"
	"github.com/andrewhsu/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSSHKey(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh-key <command>",
		Short: "Manage SSH keys",
		Long:  "Manage SSH keys registered with your GitHub account",
	}

	cmd.AddCommand(cmdList.NewCmdList(f, nil))
	cmd.AddCommand(cmdAdd.NewCmdAdd(f, nil))

	return cmd
}
