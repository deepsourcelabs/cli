package status

import (
	"fmt"

	"github.com/deepsourcelabs/cli/cmdutils"
	"github.com/deepsourcelabs/cli/internal/config"
	"github.com/spf13/cobra"
)

type AuthStatusOptions struct {
	TokenExpired  bool
	Authenticated bool
	Config        config.ConfigData
}

// NewCmdVersion returns the current version of cli being used
func NewCmdStatus(cf *cmdutils.CLIFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Refresh stored authentication credentials",
		Run: func(cmd *cobra.Command, args []string) {
			opts := AuthStatusOptions{
				TokenExpired:  cf.TokenExpired,
				Authenticated: false,
				Config:        cf.Config,
			}
			opts.Run()
		},
	}
	return cmd
}

func (opts *AuthStatusOptions) Run() {

	if opts.TokenExpired == false {
		opts.Authenticated = true

		fmt.Printf("Logged in to DeepSource (deepsource.io) as %s\n", opts.Config.User)
	} else {
		fmt.Println("You are not logged into DeepSource. Run \"gh auth login\" to authenticate.")
	}

}
