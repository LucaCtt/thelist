package main

const (
	appName           string = "thelist"
	startMsg          string = "Server started. Press CTRL+C to stop."
	configUpdatedMsg  string = "Configuration updated."
	configName        string = "config"
	configPath        string = "$HOME/" + appName
	serverPortOption  string = "port"
	serverPortDefault string = "8080"
	dbPathOption      string = "DbPath"
	dbPathDefault     string = "./" + appName + ".db"
)
