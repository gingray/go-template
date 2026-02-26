/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package command

import (
	"github.com/gingray/go-template/pkg/common"
	"github.com/gingray/go-template/pkg/component"
	"github.com/gingray/go-template/pkg/server"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		devMode, err := cmd.Flags().GetBool("dev")
		if err != nil {
			return err
		}
		if !devMode {
			return nil
		}

		err = godotenv.Load()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := common.NewConfig()
		if err != nil {
			return err
		}
		app, err := common.NewApp(cfg)
		if err != nil {
			return err
		}
		rootNode := component.NewSupervisor(app.Logger)
		appNode := rootNode.AddComponent(app)
		appNode.AddComponent(server.NewServer(&cfg.HTTPServiceConfig, app))
		return rootNode.Run(cmd.Context())
	},
}

func init() {
	serverCmd.Flags().BoolP("dev", "t", false, "Run server in dev mode")
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
