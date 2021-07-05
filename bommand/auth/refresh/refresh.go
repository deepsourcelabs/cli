package refresh

import (
	"fmt"
	"time"

	"github.com/deepsourcelabs/cli/api"
	"github.com/deepsourcelabs/cli/cmdutils"
	"github.com/deepsourcelabs/cli/internal/config"
	"github.com/spf13/cobra"
)

type RefreshOptions struct {
	graphqlClient *api.DSClient
	Config        config.ConfigData
}

// NewCmdVersion returns the current version of cli being used
func NewCmdRefresh(cf *cmdutils.CLIFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh stored authentication credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := RefreshOptions{
				graphqlClient: cf.GQLClient,
				Config:        cf.Config,
			}
			err := opts.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func (opts *RefreshOptions) Run() error {

	if opts.Config.Token != "" || opts.Config.RefreshToken != "" {

		refreshedConfigData, err := api.RefreshAuthCreds(opts.graphqlClient, opts.Config.RefreshToken)
		if err != nil {
			return err
		}

		// Convert incoming config into the ConfigData format
		finalConfig := config.ConfigData{
			User:                refreshedConfigData.Refreshtoken.Payload.Email,
			Token:               refreshedConfigData.Refreshtoken.Token,
			TokenExpiry:         refreshedConfigData.Refreshtoken.Payload.Exp,
			RefreshToken:        refreshedConfigData.Refreshtoken.Refreshtoken,
			OrigIAT:             refreshedConfigData.Refreshtoken.Payload.Origiat,
			RefreshTokenExpiry:  refreshedConfigData.Refreshtoken.Refreshexpiresin,
			RefreshTokenSetTime: time.Now().Unix(),
		}

		err = config.WriteConfigToFile(finalConfig)
		if err != nil {
			fmt.Println("Error in writing authentication data to a file. Exiting...")
			return err
		}
	} else {
		return fmt.Errorf("User not authenticated. Please login.")
	}
	return nil
}
