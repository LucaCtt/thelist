package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   rootCmdUse,
	Short: rootCmdShort,
	Long:  rootCmdLong,
}

func init() {
	// Viper configuration and defaults.
	viper.SetDefault(serverPortOption, serverPortDefault)
	viper.SetDefault(clientPortOption, clientPortDefault)
	viper.SetDefault(dbPathOption, dbPathDefault)

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(strings.ToUpper(appName))
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

// Execute starts the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
