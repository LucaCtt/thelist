package cmd

import (
	"log"
	"net/http"

	"github.com/LucaCtt/thelist/common"
	"github.com/LucaCtt/thelist/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startServer() {
	dbStore, err := common.NewDbStore(viper.GetString(dbPathOpt))
	defer dbStore.Close()

	if err != nil {
		log.Fatal(err)
	}

	router := router.New(dbStore)
	port := viper.GetString(serverPortOpt)

	log.Printf(webCmdStartMsg, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func init() {
	webCmd.Flags().IntP(serverPortOpt, serverPortShort, 0, serverPortUsage)
	webCmd.Flags().IntP(clientPortOpt, clientPortShort, 0, clientPortUsage)

	viper.BindPFlag(serverPortOpt, webCmd.Flags().Lookup(serverPortOpt))
	viper.BindPFlag(clientPortOpt, webCmd.Flags().Lookup(clientPortOpt))

	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   webCmdUse,
	Short: webCmdShort,
	Long:  webCmdLong,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}
