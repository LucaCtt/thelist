package cmd

import (
	"log"
	"strconv"

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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbStore, err := data.NewDbStore(&data.DbOptions{
			Path: viper.GetString(constants.DbPathOption),
		})
		defer dbStore.Close()

		showID, _ := strconv.Atoi(args[0])
		dbStore.CreateShow(&data.Show{ShowID: &showID})
		if err != nil {
			log.Fatal(err)
		}
	},
}
