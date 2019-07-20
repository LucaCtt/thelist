package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/LucaCtt/thelist/constants"
	"github.com/ryanbradynd05/go-tmdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   constants.RootCmdUse,
	Short: constants.RootCmdShort,
	Long:  constants.RootCmdLong,
}

func init() {
	// Viper configuration and defaults.
	viper.SetDefault(constants.ServerPortOption, constants.ServerPortDefault)
	viper.SetDefault(constants.ClientPortOption, constants.ClientPortDefault)
	viper.SetDefault(constants.DbPathOption, constants.DbPathDefault)

	viper.SetConfigName(constants.ConfigName)
	viper.AddConfigPath(constants.ConfigPath)
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(strings.ToUpper(constants.AppName))
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

var APIClient = tmdb.Init(tmdb.Config{
	APIKey: viper.GetString(constants.APIKeyOption),
})
