package cmd

import (
	"fmt"
	"log"

	"github.com/LucaCtt/thelist/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func add(args []string, prompter common.Prompter, client common.Client, store common.Store) error {
	var name string

	if len(args) != 0 {
		name = args[0]
	} else {
		input, err := prompter.Input("Show name")
		if err != nil {
			return fmt.Errorf("prompt show name failed: %w", err)
		}
		name = input
	}

	shows, err := client.Search(name)
	if err != nil {
		return err
	}

	options := make([]string, len(shows))
	for i, s := range shows {
		options[i] = fmt.Sprintf("%s (%d)", s.Name, s.ReleaseDate.Year())
	}

	index, err := prompter.Select(fmt.Sprintf("Found %d results", len(shows)), options)
	err = store.Create(&common.Item{ShowID: shows[index].ID})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   addCmdUse,
	Short: addCmdShort,
	Long:  addCmdLong,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbStore, err := common.NewDbStore(viper.GetString(dbPathOpt))
		if err != nil {
			log.Fatal(err)
		}
		defer dbStore.Close()

		client := common.DefaultTMDbClient(viper.GetString(apiKeyOpt))
		err = add(args, &common.CliPrompter{}, client, dbStore)
		if err != nil {
			log.Fatal(err)
		}
	},
}
