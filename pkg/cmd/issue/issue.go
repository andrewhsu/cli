package issue

import (
	"github.com/MakeNowJust/heredoc"
	cmdClose "github.com/andrewhsu/cli/v2/pkg/cmd/issue/close"
	cmdComment "github.com/andrewhsu/cli/v2/pkg/cmd/issue/comment"
	cmdCreate "github.com/andrewhsu/cli/v2/pkg/cmd/issue/create"
	cmdDelete "github.com/andrewhsu/cli/v2/pkg/cmd/issue/delete"
	cmdEdit "github.com/andrewhsu/cli/v2/pkg/cmd/issue/edit"
	cmdList "github.com/andrewhsu/cli/v2/pkg/cmd/issue/list"
	cmdReopen "github.com/andrewhsu/cli/v2/pkg/cmd/issue/reopen"
	cmdStatus "github.com/andrewhsu/cli/v2/pkg/cmd/issue/status"
	cmdTransfer "github.com/andrewhsu/cli/v2/pkg/cmd/issue/transfer"
	cmdView "github.com/andrewhsu/cli/v2/pkg/cmd/issue/view"
	"github.com/andrewhsu/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage issues",
		Long:  `Work with GitHub issues`,
		Example: heredoc.Doc(`
			$ gh issue list
			$ gh issue create --label bug
			$ gh issue view --web
		`),
		Annotations: map[string]string{
			"IsCore": "true",
			"help:arguments": heredoc.Doc(`
				An issue can be supplied as argument in any of the following formats:
				- by number, e.g. "123"; or
				- by URL, e.g. "https://github.com/OWNER/REPO/issues/123".
			`),
		},
	}

	cmdutil.EnableRepoOverride(cmd, f)

	cmd.AddCommand(cmdClose.NewCmdClose(f, nil))
	cmd.AddCommand(cmdCreate.NewCmdCreate(f, nil))
	cmd.AddCommand(cmdList.NewCmdList(f, nil))
	cmd.AddCommand(cmdReopen.NewCmdReopen(f, nil))
	cmd.AddCommand(cmdStatus.NewCmdStatus(f, nil))
	cmd.AddCommand(cmdView.NewCmdView(f, nil))
	cmd.AddCommand(cmdComment.NewCmdComment(f, nil))
	cmd.AddCommand(cmdDelete.NewCmdDelete(f, nil))
	cmd.AddCommand(cmdEdit.NewCmdEdit(f, nil))
	cmd.AddCommand(cmdTransfer.NewCmdTransfer(f, nil))

	return cmd
}
