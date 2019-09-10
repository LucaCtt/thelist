package cmd

import (
	"fmt"
	"log"

	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/common/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func list(p Prompt, c client.Client, s store.Store) error {
	items, err := s.All()
	if err != nil {
		return fmt.Errorf("get items from store failed: %w", err)
	}
	if len(items) == 0 {
		return fmt.Errorf("The list is empty")
	}

	shows := make([]*Show, len(items))
	for i, item := range items {
		show, err := getShow(c, item.ShowID, item.Type)
		if err != nil {
			return fmt.Errorf("get show from client failed: %w", err)
		}
		shows[i] = show
	}

	options := make([]string, len(shows))
	for i, s := range shows {
		if s.Year == 0 {
			options[i] = fmt.Sprintf("%s", s.Name)
			continue
		}
		options[i] = fmt.Sprintf("%s (%d)", s.Name, s.Year)
	}

	watched, err := p.MultiSelect("Shows", options)
	if err != nil {
		return fmt.Errorf("prompt shows list failed: %w", err)
	}

	for _, w := range watched {
		err = s.Delete(items[w].ID)
		if err != nil {
			return fmt.Errorf("delete watched shows failed: %w", err)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   listCmdUse,
	Short: listCmdShort,
	Long:  listCmdLong,
	Run: func(cmd *cobra.Command, args []string) {
		dbStore, err := store.New(viper.GetString(dbPathOpt))
		if err != nil {
			log.Fatal(err)
		}
		defer dbStore.Close()

		client := client.New(viper.GetString(apiKeyOpt))

		err = list(&CliPrompt{}, client, dbStore)
		if err != nil {
			log.Fatal(err)
		}
	},
}
