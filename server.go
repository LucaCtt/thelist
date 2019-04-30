package main

import (
	"log"
	"net/http"

	"github.com/LucaCtt/thelist/data"
	"github.com/LucaCtt/thelist/router"
	"github.com/spf13/viper"
)

func main() {
	setupConfig()

	dbStore, err := data.NewDbStore(&data.DbOptions{
		Host:     viper.GetString(dbHostOption),
		Port:     viper.GetInt(dbPortOption),
		User:     viper.GetString(dbUserOption),
		Name:     viper.GetString(dbNameOption),
		Password: viper.GetString(dbPasswordOption),
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
	viper.SetDefault(dbHostOption, dbHostDefault)
	viper.SetDefault(dbPortOption, dbPortDefault)
	viper.SetDefault(dbUserOption, dbUserDefault)
	viper.SetDefault(dbNameOption, dbNameDefault)
	viper.SetDefault(dbPasswordOption, dbPasswordDefault)

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
