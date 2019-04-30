package main

const (
	appName           string = "thelist"
	startMsg          string = "Server started. Press CTRL+C to stop."
	configUpdatedMsg  string = "Configuration updated."
	configName        string = "config"
	configPath        string = "$HOME/" + appName
	serverPortOption  string = "port"
	serverPortDefault string = "8080"
	dbHostOption      string = "DbHost"
	dbHostDefault     string = "localhost"
	dbPortOption      string = "DbPort"
	dbPortDefault     string = "5432"
	dbUserOption      string = "DbUser"
	dbUserDefault     string = "postgres"
	dbNameOption      string = "DbName"
	dbNameDefault     string = appName
	dbPasswordOption  string = "DbPassword"
	dbPasswordDefault string = "password"
)
