package cmd

const (
	appName    string = "thelist"
	configName string = "config"
	configPath string = "$HOME/" + appName

	dbPathOption  string = "DbPath"
	dbPathDefault string = "./" + appName + ".db"

	serverPortOption  string = "server"
	serverPortDefault string = "8080"
	serverPortShort   string = "s"
	serverPortUsage   string = "REST API port"
	clientPortOption  string = "port"
	clientPortDefault string = "8000"
	clientPortShort   string = "p"
	clientPortUsage   string = "Web client port"

	rootCmdUse   string = appName
	rootCmdShort string = appName + " is an app for listing shows to watch."
	rootCmdLong  string = appName + " is an application for listing shows to watch. It has both a cli and a web app interface."

	webCmdUse      string = "web"
	webCmdShort    string = "Start web app"
	webCmdLong     string = "Start the web app interface, which is composed by a web client and a REST API."
	webCmdStartMsg string = "Server started on port %s. Press CTRL+C to stop."
)
