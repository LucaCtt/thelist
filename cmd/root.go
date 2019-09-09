package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   rootCmdUse,
	Short: rootCmdShort,
	Long:  rootCmdLong,
}

var cfgFile string

func getHomedir() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

func init() {
	cobra.OnInitialize(initConfig, checkAPIKey)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, cfgOpt, cfgShort, "", cfgUsage)
	rootCmd.PersistentFlags().StringP(apiKeyOpt, apiKeyShort, "", apiKeyUsage)
	viper.BindPFlag(apiKeyOpt, rootCmd.PersistentFlags().Lookup(apiKeyOpt))

	viper.SetDefault(serverPortOpt, serverPortDefault)
	viper.SetDefault(clientPortOpt, clientPortDefault)
	viper.SetDefault(dbPathOpt, fmt.Sprintf("%s/%s.db", getHomedir(), appName))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(getHomedir())
		viper.SetConfigName(cfgFileName)
	}

	viper.SetEnvPrefix(strings.ToUpper(appName))
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func checkAPIKey() {
	if viper.GetString(apiKeyOpt) == "" {
		log.Fatalf("API key is unset")
	}
}

// Execute starts the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
