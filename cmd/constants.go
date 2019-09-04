package cmd

const (
	appName    = "thelist"
	configName = "config"
	configPath = "$HOME/" + appName

	dbPathOpt     = "DbPath"
	dbPathDefault = "./" + appName + ".db"

	serverPortOpt     = "server"
	serverPortDefault = "8080"
	serverPortShort   = "s"
	serverPortUsage   = "REST API port"

	clientPortOpt     = "port"
	clientPortDefault = "8000"
	clientPortShort   = "p"
	clientPortUsage   = "Web client port"

	apiKeyOpt   = "key"
	apiKeyShort = "k"
	apiKeyUsage = "TMDB api key"

	rootCmdUse   = appName
	rootCmdShort = appName + " is an app for listing shows to watch."
	rootCmdLong  = appName + " is an application for listing shows to watch. It has both a cli and a web app interface."

	webCmdUse      = "web"
	webCmdShort    = "Start web app"
	webCmdLong     = "Start the web app interface, which is composed by a web client and a REST API."
	webCmdStartMsg = "Server started on port %s. Press CTRL+C to stop."

	addCmdUse   = "add"
	addCmdShort = "Start web app"
	addCmdLong  = "Start the web app interface, which is composed by a web client and a REST API."
)
