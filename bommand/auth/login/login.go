package login

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cli/browser"
	"github.com/deepsourcelabs/cli/api"
	"github.com/deepsourcelabs/cli/cmdutils"
	"github.com/deepsourcelabs/cli/deepsource"
	"github.com/deepsourcelabs/cli/deepsource/auth"

	// "github.com/deepsourcelabs/cli/internal/config"
	"github.com/deepsourcelabs/cli/config"
	"github.com/spf13/cobra"
)

// Options holds the metadata.
type LoginOptions struct {
	graphqlClient *api.DSClient
	AuthTimedOut  bool
	TokenExpired  bool
	// Config        config.ConfigData
}

const pollingInterval = 1
const pollingRetryLimit = 2

// NewCmdVersion returns the current version of cli being used
func NewCmdLogin(cf *cmdutils.CLIFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to DeepSource using Command Line Interface",
		RunE: func(cmd *cobra.Command, args []string) error {

			opts := LoginOptions{
				graphqlClient: cf.GQLClient,
				AuthTimedOut:  false,
				TokenExpired:  cf.TokenExpired,
				// Config:        cf.Config,
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

// Validate impletments the Validate method for the ICommand interface.
func (opts *LoginOptions) Validate() error {
	return nil
}

// Run executest the command.
func (opts *LoginOptions) Run() error {
	// Register the device and get a device code.
	deepsource := deepsource.New()
	ctx := context.Background()
	res, err := deepsource.RegisterDevice(ctx)
	if err != nil {
		return err
	}

	// This looks ugly.
	fmt.Printf("Please copy your one-time code: %s\n", res.Code)
	fmt.Printf("Press enter to open deepsource.io in your browser...")
	fmt.Scanln()

	browser.OpenURL(res.VerificationURIComplete)

	// Start polling for updates.
	i := 0
	var jwt *auth.JWT
	for i = 0; i < pollingRetryLimit; i++ {
		time.Sleep(pollingInterval * time.Second)
		jwt, err = deepsource.Login(ctx, res.Code)
		if err != nil {
			fmt.Print(".")
			continue
		}
	}

	if i == pollingRetryLimit {
		return errors.New("authentication attempt expired")
	}

	// Create config.
	conf := config.CLIConfig{
		Token:                 jwt.Token,
		RefreshToken:          jwt.Refreshtoken,
		RefreshTokenExpiresIn: jwt.RefreshExpiresIn,
	}
	conf.SetTokenExpiry(jwt.Payload.Exp)

	if err := conf.WriteFile(); err != nil {
		return err
	}

	return nil
}
