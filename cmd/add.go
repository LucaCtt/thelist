package cmd

import (
	"log"
	"strconv"

	"github.com/LucaCtt/thelist/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   addCmdUse,
	Short: addCmdShort,
	Long:  addCmdLong,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbStore, err := data.NewDbStore(&data.DbOptions{
			Path: viper.GetString(dbPathOption),
		})
		defer dbStore.Close()

		showID, _ := strconv.Atoi(args[0])
		dbStore.CreateShow(&data.Show{ShowID: &showID})
		if err != nil {
			log.Fatal(err)
		}
	},
}
