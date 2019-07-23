package cmd

import (
	"log"

	"github.com/LucaCtt/thelist/api"
	"github.com/LucaCtt/thelist/constants"

	"github.com/LucaCtt/thelist/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   constants.AddCmdUse,
	Short: constants.AddCmdShort,
	Long:  constants.AddCmdLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbStore, err := data.NewDbStore(&data.DbOptions{
			Path: viper.GetString(constants.DbPathOption),
		})
		if err != nil {
			log.Fatal(err)
		}
		defer dbStore.Close()

		client := api.NewClient(viper.GetString(constants.APIKeyOption))

		show, err := promptShow(args)
		if err != nil {
			log.Fatal(err)
		}

		shows, err := client.SearchShow(show)
		if err != nil {
			log.Fatal(err)
		}

		i, err := selectShow(shows)
		if err != nil {
			log.Fatal(err)
		}

		showID := shows.Results[i].ID
		err = dbStore.CreateShow(&data.Show{ShowID: &showID})
		if err != nil {
			log.Fatal(err)
		}
	},
}
