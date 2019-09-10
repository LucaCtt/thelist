package cmd

const (
	appName = "thelist"

	cfgOpt      = "config"
	cfgShort    = "c"
	cfgUsage    = "Set config file"
	cfgFileName = "." + appName

	dbPathOpt = "dbpath"

	serverPortOpt     = "server"
	serverPortDefault = "9000"
	serverPortShort   = "s"
	serverPortUsage   = "REST API server port"

	clientPortOpt     = "port"
	clientPortDefault = "8000"
	clientPortShort   = "p"
	clientPortUsage   = "Web client port"

	apiKeyOpt   = "key"
	apiKeyShort = "k"
	apiKeyUsage = "TMDB api key"

	rootCmdUse   = appName
	rootCmdShort = "An app to list shows to watch."
	rootCmdLong  = "An app to list shows to watch. It has both a cli and a web interface."

	webCmdUse      = "web"
	webCmdShort    = "Start web app"
	webCmdLong     = "Start the web app interface, which is composed by a web client and a REST API."
	webCmdStartMsg = "Web app started. Press CTRL+C to stop."

	addCmdUse   = "add"
	addCmdShort = "Add a show to the list."
	addCmdLong  = "Add a show to the list. If the name is not passed as an argument, you will be prompted for it."

	listCmdUse   = "list"
	listCmdShort = "List shows"
	listCmdLong  = "Show a list of the stored shows. Here you can also mark them as watched (and delete them)"
)
