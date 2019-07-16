package cmd

import (
	"log"
	"net/http"

	"github.com/LucaCtt/thelist/data"
	"github.com/LucaCtt/thelist/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startServer() {
	dbStore, err := data.NewDbStore(&data.DbOptions{
		Path: viper.GetString(dbPathOption),
	})
	defer dbStore.Close()

	if err != nil {
		log.Fatal(err)
	}

	router := router.New(dbStore)
	port := viper.GetString(serverPortOption)

	log.Printf(webCmdStartMsg, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func init() {
	webCmd.Flags().IntP(serverPortOption, serverPortShort, 0, serverPortUsage)
	webCmd.Flags().IntP(clientPortOption, clientPortShort, 0, clientPortUsage)

	viper.BindPFlag(serverPortOption, webCmd.Flags().Lookup(serverPortOption))
	viper.BindPFlag(clientPortOption, webCmd.Flags().Lookup(clientPortOption))

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
