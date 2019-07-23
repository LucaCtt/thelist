package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/LucaCtt/thelist/api"
	"github.com/LucaCtt/thelist/constants"

	"github.com/LucaCtt/thelist/data"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func promptShow(args []string) (string, error) {
	if len(args) != 0 {
		return args[0], nil
	}

	prompt := promptui.Prompt{
		Label: "Show name",
	}

	show, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return show, nil
}

func selectShow(shows *api.ShowSearchResult) (int, error) {
	if shows.TotalResults == 0 {
		return 0, fmt.Errorf("No shows found")
	}

	if shows.TotalResults == 1 {
		return shows.Results[0].ID, nil
	}

	prompt := promptui.Select{
		Label:     "Select one",
		Items:     shows.Results,
		Templates: templates,
		Searcher:  searcher(shows.Results),
	}
	i, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return i, nil
}

var templates = &promptui.SelectTemplates{
	Label:    "{{ . }}:",
	Active:   "\u2B95 {{ .Name | cyan }}",
	Inactive: "  {{ .Name | cyan }}",
	Selected: "\u2B95 {{ .Name | red | cyan }}",
	Details: `
{{ "Name:" | faint }}	{{ .Name }}
{{ "Release Date:" | faint }}	{{ .ReleaseDate }}
{{ "Popularity:" | faint }}	{{ .Popularity }}
{{ "Vote Average:" | faint }}	{{ .VoteAverage}}`,
}

func searcher(shows []*api.ShowSearchInfo) func(string, int) bool {
	return func(input string, index int) bool {
		show := shows[index]

		// Convert string to lowercase and remove all whitespace
		name := strings.Replace(strings.ToLower(show.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
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
