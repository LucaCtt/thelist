package cmd

import (
	"log"
	"net/http"

	"github.com/lucactt/thelist/constants"
	"github.com/lucactt/thelist/data"
	"github.com/lucactt/thelist/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startServer() {
	dbStore, err := data.NewDbStore(&data.DbOptions{
		Path: viper.GetString(constants.DbPathOption),
	})
	defer dbStore.Close()

	if err != nil {
		log.Fatal(err)
	}

	router := router.New(dbStore)
	port := viper.GetString(constants.ServerPortOption)

	log.Printf(constants.WebCmdStartMsg, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func init() {
	webCmd.Flags().IntP(constants.ServerPortOption, constants.ServerPortShort, 0, constants.ServerPortUsage)
	webCmd.Flags().IntP(constants.ClientPortOption, constants.ClientPortShort, 0, constants.ClientPortUsage)

	viper.BindPFlag(constants.ServerPortOption, webCmd.Flags().Lookup(constants.ServerPortOption))
	viper.BindPFlag(constants.ClientPortOption, webCmd.Flags().Lookup(constants.ClientPortOption))

	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   constants.WebCmdUse,
	Short: constants.WebCmdShort,
	Long:  constants.WebCmdLong,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}
