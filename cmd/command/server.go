/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package command

import (
	"errors"

	"github.com/gingray/go-template/internal/api"
	"github.com/gingray/go-template/internal/consumer"
	"github.com/gingray/go-template/internal/repo"
	"github.com/gingray/go-template/pkg/app"
	"github.com/gingray/go-template/pkg/config"
	"github.com/gingray/go-template/pkg/httpserver"
	"github.com/gingray/go-template/pkg/infra"
	"github.com/gingray/go-template/pkg/kafka"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		devMode, _ := cmd.Flags().GetBool("dev")
		if !devMode {
			return
		}
		_ = godotenv.Load()
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, cfgErr := config.NewConfig()
		app, err := app.NewApp(cfg)
		if err != nil {
			err = errors.Join(cfgErr, err)
		}
		if err != nil {
			app.Logger.Error("init app", "error", err)
			return
		}

		server := httpserver.NewServer(&cfg.HTTPServiceConfig, app)
		router := api.NewRouter(repo.NewUserRepo(app.PgPool), kafka.NewProducer(app))
		router.SetupRoutes(app.HttpRouter)
		kafkaConsumer := kafka.NewConsumer(app, consumer.NewBasicConsumer(app))

		supervisor := infra.NewSupervisor(app.Logger)
		rootNode := supervisor.CreateRootNode()
		appNode := supervisor.CreateNode(app)
		serverNode := supervisor.CreateNode(server)
		consumerNode := supervisor.CreateNode(kafkaConsumer)

		rootNode.AddNode(appNode)
		appNode.AddNode(serverNode)
		appNode.AddNode(consumerNode)
		err = rootNode.Run(cmd.Context())
		if err != nil {
			app.Logger.Error("run root node", "error", err)
		}
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
