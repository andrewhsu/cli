package delete

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andrewhsu/cli/v2/api"
	"github.com/andrewhsu/cli/v2/internal/config"
	"github.com/andrewhsu/cli/v2/pkg/cmd/gist/shared"
	"github.com/andrewhsu/cli/v2/pkg/cmdutil"
	"github.com/andrewhsu/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	IO         *iostreams.IOStreams
	Config     func() (config.Config, error)
	HttpClient func() (*http.Client, error)

	Selector string
}

func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := DeleteOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "delete {<id> | <url>}",
		Short: "Delete a gist",
		Args:  cmdutil.ExactArgs(1, "cannot delete: gist argument required"),
		RunE: func(c *cobra.Command, args []string) error {
			opts.Selector = args[0]
			if runF != nil {
				return runF(&opts)
			}
			return deleteRun(&opts)
		},
	}
	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	gistID := opts.Selector

	if strings.Contains(gistID, "/") {
		id, err := shared.GistIDFromURL(gistID)
		if err != nil {
			return err
		}
		gistID = id
	}
	client, err := opts.HttpClient()
	if err != nil {
		return err
	}

	apiClient := api.NewClientFromHTTP(client)

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	host, err := cfg.DefaultHost()
	if err != nil {
		return err
	}

	gist, err := shared.GetGist(client, host, gistID)
	if err != nil {
		return err
	}
	username, err := api.CurrentLoginName(apiClient, host)
	if err != nil {
		return err
	}

	if username != gist.Owner.Login {
		return fmt.Errorf("You do not own this gist.")
	}

	err = deleteGist(apiClient, host, gistID)
	if err != nil {
		return err
	}

	return nil
}

func deleteGist(apiClient *api.Client, hostname string, gistID string) error {
	path := "gists/" + gistID
	err := apiClient.REST(hostname, "DELETE", path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
