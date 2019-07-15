package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/LucaCtt/thelist/data"
	"github.com/LucaCtt/thelist/router"
	"github.com/spf13/viper"
)

func main() {
	setupConfig()

	dbStore, err := data.NewDbStore(&data.DbOptions{
		Path: viper.GetString(dbPathOption),
	})
	defer dbStore.Close()

	if err != nil {
		log.Fatal(err)
	}

	router := router.New(dbStore)

	log.Println(startMsg)
	log.Fatal(http.ListenAndServe(":"+viper.GetString(serverPortOption), router))
}

func setupConfig() {
	viper.SetDefault(serverPortOption, serverPortDefault)
	viper.SetDefault(dbPathOption, dbPathDefault)
	viper.SetDefault(dbUserOption, dbUserDefault)
	viper.SetDefault(dbPasswordOption, dbPasswordDefault)

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(strings.ToUpper(appName))
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
