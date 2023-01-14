package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/livebotapp/publiclistener"
	handlershttp "github.com/livebotapp/publiclistener/handler/http"
	"github.com/livebotapp/publiclistener/kick"
	"github.com/livebotapp/publiclistener/redis"
)

var (
	configPath string
)

func init() {
	serverCmd.Flags().StringVar(&configPath, "config", "./config", "Defines here the config path where global conf file can be found")

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() error {
	viper.SetConfigName("global")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	return err
}

func runServer(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Initialize the configuration for the app
	err := initConfig()
	if err != nil {
		panic(err)
	}

	pusher, err := redis.NewPusher(ctx, viper.GetString("redis.addr"), viper.GetString("queues.confirmation_attempt_received"))
	if err != nil {
		panic(err)
	}

	forwarder := publiclistener.NewForwarder(ctx, pusher)
	consumer := kick.NewConsumer(ctx, viper.GetString("kickchat.host"), viper.GetString("kickchat.path"), forwarder)

	fmt.Printf("[CHAT LISTENER] - Started...\n")

	s := handlershttp.NewServer(ctx, viper.GetString("http.port"))
	if err := s.Setup(ctx); err != nil {
		panic(err)
	}

	if err = consumer.Start(ctx); err != nil {
		panic(err)
	}

	s.Start(ctx)
}

var rootCmd = &cobra.Command{
	Use:   "live",
	Short: "The solution running both the backend and the frontend of the website",
	Long:  "Live is the binary used to run the whole platform - being both the backend and the frontend",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Use this command to run the server.",
	Long:  "This command will run the server using given args and the config given as command line if any.",
	Run:   runServer,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the core.",
	Long:  "Print the current version used to run the server in the given binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Core Server Version : %s\n", publiclistener.Version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
