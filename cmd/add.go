package cmd

import (
	"fmt"
	"log"

	"github.com/LucaCtt/thelist/common"
	"github.com/LucaCtt/thelist/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func add(args []string, p common.Prompter, c common.Client, s common.Store) error {
	var name string

	if len(args) != 0 {
		name = args[0]
	} else {
		input, err := p.Input("Show name")
		if err != nil {
			return errors.E("prompt show name failed", err, errors.SeverityWarn)
		}
		if input == "" {
			return errors.E("invalid show name")
		}
		name = input
	}

	shows, err := searchShow(c, name)
	if err != nil {
		return errors.E("search show failed", err)
	}
	if len(shows) == 0 {
		return errors.E("no shows found")
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
			return errors.E("select option failed", err, errors.SeverityWarn)
		}
		selected = shows[i]
	}

	err = s.Create(&common.Item{ShowID: selected.ID, Type: selected.Type})
	if err != nil {
		return errors.E("create item failed", err)
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
