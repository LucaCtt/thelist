package cmd

import (
	"log"
	"net/http"

	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/common/store"
	"github.com/LucaCtt/thelist/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startServer() {
	dbStore, err := store.New(viper.GetString(dbPathOpt))
	defer dbStore.Close()

	if err != nil {
		log.Fatal(err)
	}

	client := client.New(viper.GetString(apiKeyOpt))

	server := server.New(dbStore, client)
	port := viper.GetString(serverPortOpt)

	log.Print(webCmdStartMsg)
	log.Fatal(http.ListenAndServe(":"+port, server))
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
