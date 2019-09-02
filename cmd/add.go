package cmd

import (
	"log"

	"github.com/lucactt/thelist/common"
	"github.com/lucactt/thelist/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func add(args []string, prompter common.Prompter, client common.Client, store common.Store) error {
	var name string

	if len(args) != 0 {
		name = args[0]
	} else {
		input, err := prompter.Input()
		if err != nil {
			return err
		}
		name = input
	}

	shows, err := client.Search(name)
	if err != nil {
		return err
	}

	selected, err := prompter.Select(shows)
	if err != nil {
		return err
	}

	err = store.Create(&common.Item{ShowID: selected.ID})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   constants.AddCmdUse,
	Short: constants.AddCmdShort,
	Long:  constants.AddCmdLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbStore, err := common.NewDbStore(viper.GetString(constants.DbPathOption))
		if err != nil {
			log.Fatal(err)
		}
		defer dbStore.Close()

		client := common.DefaultTMDbClient(viper.GetString(constants.APIKeyOption))
		err = add(args, &common.CliPrompter{}, client, dbStore)
		if err != nil {
			log.Fatal(err)
		}
	},
}
