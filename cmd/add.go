package cmd

import (
	"fmt"
	"log"

	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/common/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func add(args []string, p Prompt, c client.Client, s store.Store) error {
	var name string

	if len(args) != 0 {
		name = args[0]
	} else {
		input, err := p.Input("Show name")
		if err != nil {
			return fmt.Errorf("prompt show name failed: %w", err)
		}
		if input == "" {
			return fmt.Errorf("Invalid show name")
		}
		name = input
	}

	shows, err := searchShow(c, name)
	if err != nil {
		return fmt.Errorf("search show failed: %w", err)
	}
	if len(shows) == 0 {
		return fmt.Errorf("No shows found")
	}

	options := make([]string, len(shows))
	for i, s := range shows {
		if s.Year == 0 {
			options[i] = fmt.Sprintf("%s", s.Name)
			continue
		}
		options[i] = fmt.Sprintf("%s (%d)", s.Name, s.Year)
	}

	var selected *Show
	if len(shows) == 1 {
		selected = shows[0]
	} else {
		i, err := p.Select(fmt.Sprintf("Found %d results", len(shows)), options)
		if err != nil {
			return err
		}
		selected = shows[i]
	}

	err = s.Create(&store.Item{ShowID: selected.ID, Type: selected.Type})

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
		dbStore, err := store.New(viper.GetString(dbPathOpt))
		if err != nil {
			log.Fatal(err)
		}
		defer dbStore.Close()

		client := client.New(viper.GetString(apiKeyOpt))

		err = add(args, &CliPrompt{}, client, dbStore)
		if err != nil {
			log.Fatal(err)
		}
	},
}
