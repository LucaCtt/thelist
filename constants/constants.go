package constants

const (
	AppName    = "thelist"
	ConfigName = "config"
	ConfigPath = "$HOME/" + AppName

	DbPathOption  = "DbPath"
	DbPathDefault = "./" + AppName + ".db"

	ServerPortOption  = "server"
	ServerPortDefault = "8080"
	ServerPortShort   = "s"
	ServerPortUsage   = "REST API port"
	ClientPortOption  = "port"
	ClientPortDefault = "8000"
	ClientPortShort   = "p"
	ClientPortUsage   = "Web client port"

	RootCmdUse   = AppName
	RootCmdShort = AppName + " is an app for listing shows to watch."
	RootCmdLong  = AppName + " is an application for listing shows to watch. It has both a cli and a web app interface."

	WebCmdUse      = "web"
	WebCmdShort    = "Start web app"
	WebCmdLong     = "Start the web app interface, which is composed by a web client and a REST API."
	WebCmdStartMsg = "Server started on port %s. Press CTRL+C to stop."

	AddCmdUse   = "add"
	AddCmdShort = "Start web app"
	AddCmdLong  = "Start the web app interface, which is composed by a web client and a REST API."

	APIKeyOption = "key"
)
