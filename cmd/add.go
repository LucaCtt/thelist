package cmd

import (
	"log"

	"github.com/lucactt/thelist/constants"
	"github.com/lucactt/thelist/util"

	"github.com/lucactt/thelist/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func add(args []string, prompt util.Prompt, client util.Client, store data.Store) error {
	name := ""

	if len(args) != 0 {
		name = args[0]
	} else {
		input, err := prompt.PromptShow()
		if err != nil {
			return err
		}
		name = input
	}

	searchResult, err := client.SearchShow(name)
	if err != nil {
		return err
	}

	selectedShow, err := prompt.SelectShow(searchResult)
	if err != nil {
		return err
	}

	showID := selectedShow.ID
	err = store.Create(&data.Item{ShowID: showID})
	if err != nil {
		return err
	}

	return nil
}

var addCmd = &cobra.Command{
	Use:   constants.AddCmdUse,
	Short: constants.AddCmdShort,
	Long:  constants.AddCmdLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbStore, err := data.NewDbStore(viper.GetString(constants.DbPathOption))
		if err != nil {
			log.Fatal(err)
		}
		defer dbStore.Close()

		client := util.NewAPIClient(viper.GetString(constants.APIKeyOption))

		err = add(args, &util.CliPrompt{}, client, dbStore)
		if err != nil {
			log.Fatal(err)
		}
	},
}
