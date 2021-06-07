package view

import (
	"fmt"
	"strings"

	"github.com/cli/browser"
	"github.com/deepsourcelabs/cli/api"
	"github.com/deepsourcelabs/cli/cmdutils"
	cliConfig "github.com/deepsourcelabs/cli/internal/config"
	"github.com/spf13/cobra"
)

type RepoViewOptions struct {
	graphqlClient *api.DSClient
	RepoArg       string
	Config        cliConfig.ConfigData
	TokenExpired  bool
	Owner         string
	RepoName      string
	VCSProvider   string
}

// NewCmdVersion returns the current version of cli being used
func NewCmdRepoView(cf *cmdutils.CLIFactory) *cobra.Command {

	opts := RepoViewOptions{
		graphqlClient: cf.GQLClient,
		RepoArg:       "",
		Config:        cf.Config,
		TokenExpired:  cf.TokenExpired,
		Owner:         "",
		RepoName:      "",
		VCSProvider:   "",
	}

	cmd := &cobra.Command{
		Use:   "view",
		Short: "Open the DeepSource dashboard of a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}

	// --repo, -r flag
	cmd.Flags().StringVarP(&opts.RepoArg, "repo", "r", "", "Open the DeepSource dashboard of the specified repository")
	return cmd
}

func (opts *RepoViewOptions) Run() error {

	if opts.RepoArg == "" {

		remotesData, err := cmdutils.ListRemotes()
		if err != nil {
			if strings.Contains(err.Error(), "exit status 128") {
				fmt.Println("This repository has not been initialized with git. Please initialize it with git using `git init`")
			}
			return err
		}

		// Condition: If only 1 remote, use it
		// If more that one remote, give the user choice which one they want to choose

		if len(remotesData) == 1 {
			for _, value := range remotesData {
				opts.Owner = value[0]
				opts.RepoName = value[1]
				opts.VCSProvider = value[2]
			}
		} else {

			var promptOpts []string
			for _, value := range remotesData {
				promptOpts = append(promptOpts, value[3])
			}

			selectedRemote, err := cmdutils.SelectFromOptions("Please select which repository you want to view?", "", promptOpts)
			if err != nil {
				return err
			}

			for _, value := range remotesData {
				if value[3] == selectedRemote {
					opts.Owner = value[0]
					opts.RepoName = value[1]
					opts.VCSProvider = value[2]
				}
			}
		}
	} else {

		// Parsing the arguments to --repo flag
		repoData, err := cmdutils.RepoArgumentResolver(opts.RepoArg)
		if err != nil {
			return err
		}

		opts.VCSProvider = repoData[0]
		opts.Owner = repoData[1]
		opts.RepoName = repoData[2]
	}

	// Making the "isActivated" query again just to confirm if the user has access to that repo
	_, err := api.GetRepoStatus(opts.graphqlClient, opts.Owner, opts.RepoName, opts.VCSProvider)
	if err != nil {
		if strings.Contains(err.Error(), "Signature has expired") {
			fmt.Println("The token has expired. Please refresh the token using the command `deepsource auth refresh`")
		}

		if strings.Contains(err.Error(), "Repository matching query does not exist") {
			fmt.Println("Unauthorized access. Please login if you haven't using the command `deepsource auth login`")
		}

		return err
	}

	var vcs string

	// Frame the URL
	switch opts.VCSProvider {
	case "GITHUB":
		vcs = "gh"
	case "GITLAB":
		vcs = "gl"
	case "BITBUCKET":
		vcs = "bb"
	}

	dashboardURL := fmt.Sprintf("https://deepsource.io/%s/%s/%s/", vcs, opts.Owner, opts.RepoName)
	fmt.Printf("Press enter to open %s in your browser...", dashboardURL)
	fmt.Scanln()
	err = browser.OpenURL(dashboardURL)
	if err != nil {
		return err
	}

	return nil

}