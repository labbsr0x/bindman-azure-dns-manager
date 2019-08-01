package cmd

import (
	"fmt"
	"github.com/labbsr0x/bindman-azure-dns-manager/src/azure"
	"github.com/labbsr0x/bindman-azure-dns-manager/src/manager"
	"github.com/labbsr0x/bindman-azure-dns-manager/src/version"
	"github.com/labbsr0x/bindman-dns-webhook/src/hook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const basePath = "./data"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server and serves the HTTP REST API",
	Example: `  bindman-azure-dns-manager serve --dns-ttl=1h --dns-removal-delay=20m

  All command line options can be provided via environment variables by adding the prefix "BINDMAN_" 
  and converting their names to upper case and replacing punctuation and hyphen with underscores. 
  For example,

        command line option                 environment variable
        ------------------------------------------------------------------
        --dns-removal-delay                 BINDMAN_DNS_REMOVAL_DELAY
`,
	RunE: runE,
}

func runE(_ *cobra.Command, _ []string) error {
	azureBuilder := new(azure.Builder).InitFromViper(viper.GetViper())
	managerBuilder := new(manager.Builder).InitFromViper(viper.GetViper())
	nsu, err := azureBuilder.New()
	if err != nil {
		return fmt.Errorf("\n  Error occurred while setting up the DNS Manager.\n  %v", err)
	}
	azureManager, err := managerBuilder.New(nsu, basePath)
	if err != nil {
		return err
	}

	logrus.New().WithFields(logrus.Fields{
		"Version":   version.Version,
		"GitCommit": version.GitCommit,
		"BuildTime": version.BuildTime,
	}).Info("bindman-azure-dns-manager version")
	hook.Initialize(azureManager, version.Version)
	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	azure.AddFlags(serveCmd.Flags())
	manager.AddFlags(serveCmd.Flags())

	err := viper.GetViper().BindPFlags(serveCmd.Flags())
	if err != nil {
		panic(err)
	}
}
