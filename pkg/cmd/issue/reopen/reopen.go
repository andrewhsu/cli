package reopen

import (
	"fmt"
	"net/http"

	"github.com/andrewhsu/cli/v2/api"
	"github.com/andrewhsu/cli/v2/internal/config"
	"github.com/andrewhsu/cli/v2/internal/ghrepo"
	"github.com/andrewhsu/cli/v2/pkg/cmd/issue/shared"
	"github.com/andrewhsu/cli/v2/pkg/cmdutil"
	"github.com/andrewhsu/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type ReopenOptions struct {
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (ghrepo.Interface, error)

	SelectorArg string
}

func NewCmdReopen(f *cmdutil.Factory, runF func(*ReopenOptions) error) *cobra.Command {
	opts := &ReopenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "reopen {<number> | <url>}",
		Short: "Reopen issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// support `-R, --repo` override
			opts.BaseRepo = f.BaseRepo

			if len(args) > 0 {
				opts.SelectorArg = args[0]
			}

			if runF != nil {
				return runF(opts)
			}
			return reopenRun(opts)
		},
	}

	return cmd
}

func reopenRun(opts *ReopenOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	issue, baseRepo, err := shared.IssueFromArg(apiClient, opts.BaseRepo, opts.SelectorArg)
	if err != nil {
		return err
	}

	if issue.State == "OPEN" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Issue #%d (%s) is already open\n", cs.Yellow("!"), issue.Number, issue.Title)
		return nil
	}

	err = api.IssueReopen(apiClient, baseRepo, *issue)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.ErrOut, "%s Reopened issue #%d (%s)\n", cs.SuccessIconWithColor(cs.Green), issue.Number, issue.Title)

	return nil
}
