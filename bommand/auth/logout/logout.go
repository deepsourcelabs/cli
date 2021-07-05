package logout

import (
	"github.com/deepsourcelabs/cli/cmdutils"
	"github.com/deepsourcelabs/cli/internal/config"
	"github.com/spf13/cobra"
)

type LogoutOptions struct {
}

// NewCmdVersion returns the current version of cli being used
func NewCmdLogout(cf *cmdutils.CLIFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from the DeepSource authentication on Command Line",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := LogoutOptions{}
			err := opts.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func (opts *LogoutOptions) Run() error {

	logoutConfirmationMsg := "Are you sure you want to log out of DeepSource account?"
	helpText := ""
	response, err := cmdutils.ConfirmFromUser(logoutConfirmationMsg, helpText)
	if err != nil {
		return err
	}

	if response == true {
		err := config.DeleteConfigFile()
		if err != nil {
			return err
		}
	}

	return nil
}
