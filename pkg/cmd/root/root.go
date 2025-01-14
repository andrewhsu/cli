package root

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	actionsCmd "github.com/andrewhsu/cli/v2/pkg/cmd/actions"
	aliasCmd "github.com/andrewhsu/cli/v2/pkg/cmd/alias"
	apiCmd "github.com/andrewhsu/cli/v2/pkg/cmd/api"
	authCmd "github.com/andrewhsu/cli/v2/pkg/cmd/auth"
	browseCmd "github.com/andrewhsu/cli/v2/pkg/cmd/browse"
	completionCmd "github.com/andrewhsu/cli/v2/pkg/cmd/completion"
	configCmd "github.com/andrewhsu/cli/v2/pkg/cmd/config"
	extensionCmd "github.com/andrewhsu/cli/v2/pkg/cmd/extension"
	"github.com/andrewhsu/cli/v2/pkg/cmd/factory"
	gistCmd "github.com/andrewhsu/cli/v2/pkg/cmd/gist"
	issueCmd "github.com/andrewhsu/cli/v2/pkg/cmd/issue"
	prCmd "github.com/andrewhsu/cli/v2/pkg/cmd/pr"
	releaseCmd "github.com/andrewhsu/cli/v2/pkg/cmd/release"
	repoCmd "github.com/andrewhsu/cli/v2/pkg/cmd/repo"
	creditsCmd "github.com/andrewhsu/cli/v2/pkg/cmd/repo/credits"
	runCmd "github.com/andrewhsu/cli/v2/pkg/cmd/run"
	secretCmd "github.com/andrewhsu/cli/v2/pkg/cmd/secret"
	sshKeyCmd "github.com/andrewhsu/cli/v2/pkg/cmd/ssh-key"
	versionCmd "github.com/andrewhsu/cli/v2/pkg/cmd/version"
	workflowCmd "github.com/andrewhsu/cli/v2/pkg/cmd/workflow"
	"github.com/andrewhsu/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *cmdutil.Factory, version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gh <command> <subcommand> [flags]",
		Short: "GitHub CLI",
		Long:  `Work seamlessly with GitHub from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ gh issue create
			$ gh repo clone cli/cli
			$ gh pr checkout 321
		`),
		Annotations: map[string]string{
			"help:feedback": heredoc.Doc(`
				Open an issue using 'gh issue create -R github.com/cli/cli'
			`),
			"help:environment": heredoc.Doc(`
				See 'gh help environment' for the list of supported environment variables.
			`),
		},
	}

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(f, cmd, args)
	})
	cmd.SetUsageFunc(rootUsageFunc)
	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	formattedVersion := versionCmd.Format(version, buildDate)
	cmd.SetVersionTemplate(formattedVersion)
	cmd.Version = formattedVersion
	cmd.Flags().Bool("version", false, "Show gh version")

	// Child commands
	cmd.AddCommand(versionCmd.NewCmdVersion(f, version, buildDate))
	cmd.AddCommand(actionsCmd.NewCmdActions(f))
	cmd.AddCommand(aliasCmd.NewCmdAlias(f))
	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(configCmd.NewCmdConfig(f))
	cmd.AddCommand(creditsCmd.NewCmdCredits(f, nil))
	cmd.AddCommand(gistCmd.NewCmdGist(f))
	cmd.AddCommand(completionCmd.NewCmdCompletion(f.IOStreams))
	cmd.AddCommand(extensionCmd.NewCmdExtension(f))
	cmd.AddCommand(secretCmd.NewCmdSecret(f))
	cmd.AddCommand(sshKeyCmd.NewCmdSSHKey(f))

	// the `api` command should not inherit any extra HTTP headers
	bareHTTPCmdFactory := *f
	bareHTTPCmdFactory.HttpClient = bareHTTPClient(f, version)

	cmd.AddCommand(apiCmd.NewCmdApi(&bareHTTPCmdFactory, nil))

	// below here at the commands that require the "intelligent" BaseRepo resolver
	repoResolvingCmdFactory := *f
	repoResolvingCmdFactory.BaseRepo = factory.SmartBaseRepoFunc(f)

	cmd.AddCommand(browseCmd.NewCmdBrowse(&repoResolvingCmdFactory, nil))
	cmd.AddCommand(prCmd.NewCmdPR(&repoResolvingCmdFactory))
	cmd.AddCommand(issueCmd.NewCmdIssue(&repoResolvingCmdFactory))
	cmd.AddCommand(releaseCmd.NewCmdRelease(&repoResolvingCmdFactory))
	cmd.AddCommand(repoCmd.NewCmdRepo(&repoResolvingCmdFactory))
	cmd.AddCommand(runCmd.NewCmdRun(&repoResolvingCmdFactory))
	cmd.AddCommand(workflowCmd.NewCmdWorkflow(&repoResolvingCmdFactory))

	// Help topics
	cmd.AddCommand(NewHelpTopic("environment"))
	cmd.AddCommand(NewHelpTopic("formatting"))
	cmd.AddCommand(NewHelpTopic("mintty"))
	referenceCmd := NewHelpTopic("reference")
	referenceCmd.SetHelpFunc(referenceHelpFn(f.IOStreams))
	cmd.AddCommand(referenceCmd)

	cmdutil.DisableAuthCheck(cmd)

	// this needs to appear last:
	referenceCmd.Long = referenceLong(cmd)
	return cmd
}

func bareHTTPClient(f *cmdutil.Factory, version string) func() (*http.Client, error) {
	return func() (*http.Client, error) {
		cfg, err := f.Config()
		if err != nil {
			return nil, err
		}
		return factory.NewHTTPClient(f.IOStreams, cfg, version, false)
	}
}
