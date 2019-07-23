package cmd

import (
	"fmt"
	"strings"

	"github.com/LucaCtt/thelist/api"
	"github.com/manifoldco/promptui"
)

func promptShow(args []string) (string, error) {
	if len(args) != 0 {
		return args[0], nil
	}

	prompt := promptui.Prompt{
		Label: "Show name",
		Templates: &promptui.PromptTemplates{
			Success: fmt.Sprintf("%s {{ . | bold | green }} ", promptui.IconGood),
		},
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
	Active:   fmt.Sprintf("%s {{ .Name | cyan | underline }}", promptui.IconSelect),
	Inactive: "  {{ .Name | cyan }}",
	Selected: fmt.Sprintf("%s {{ .Name | bold | green }}", promptui.IconGood),
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
