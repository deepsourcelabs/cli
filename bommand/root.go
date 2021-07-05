package command

import (
	"fmt"

	"github.com/deepsourcelabs/cli/api"
	"github.com/deepsourcelabs/cli/bommand/auth"
	"github.com/deepsourcelabs/cli/bommand/version"
	"github.com/deepsourcelabs/cli/cmdutils"
	"github.com/deepsourcelabs/cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdRoot(cmdFactory *cmdutils.CLIFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deepsource <command> <subcommand> [flags]",
		Short: "DeepSource CLI",
		Long:  `Now ship good code directly from the command line.`,
	}
	cmd.AddCommand(version.NewCmdVersion())
	cmd.AddCommand(auth.NewCmdAuth(cmdFactory))

	return cmd
}

func Execute() error {

	var err error

	var cmdFactory cmdutils.CLIFactory

	// Config operations
	var authConfigData config.ConfigData

	// Read the config file
	// If there is a config file already, this returns its data
	// Else the fields are blank
	authConfigData, _ = config.ReadConfig()

	cmdFactory.Config = authConfigData

	// Setting defaults factory settings
	cmdFactory.HostName = "http://localhost:8000/graphql/"
	cmdFactory.TokenExpired = true

	// Creating a GraphQL client which can be picked up by any command since its in the factory
	cmdFactory.GQLClient = api.NewDSClient(cmdFactory.HostName, cmdFactory.Config.Token)

	if cmdFactory.Config.Token != "" {
		cmdFactory.TokenExpired, err = api.CheckTokenExpiry(cmdFactory.GQLClient, cmdFactory.Config.Token)
		if err != nil {
			if err == fmt.Errorf("graphql: Signature has expired") {
				fmt.Println("The token has expired. Please refresh the token or reauthenticate.")
			}
		}
	}

	// Pass configData struct to all the packages
	cmd := NewCmdRoot(&cmdFactory)
	if err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}
