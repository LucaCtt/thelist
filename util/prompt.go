//go:generate mockgen -destination=../mocks/mock_prompt.go -package=mocks github.com/lucactt/thelist/util Prompt

package util

import (
	"fmt"
	"strings"

	"github.com/lucactt/thelist/data"
	"github.com/manifoldco/promptui"
)

type Prompt interface {
	PromptShow() (string, error)
	SelectShow(shows *data.ShowList) (*data.Show, error)
}

type CliPrompt struct {
}

func (c *CliPrompt) PromptShow() (string, error) {
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

func (c *CliPrompt) SelectShow(shows *data.ShowList) (*data.Show, error) {
	if shows.TotalResults == 0 {
		return nil, fmt.Errorf("No shows found")
	}

	if shows.TotalResults == 1 {
		return shows.Results[0], nil
	}

	prompt := promptui.Select{
		Label:     "Select one",
		Items:     shows.Results,
		Templates: templates,
		Searcher:  searcher(shows.Results),
	}
	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return shows.Results[i], nil
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

func searcher(shows []*data.Show) func(string, int) bool {
	return func(input string, index int) bool {
		show := shows[index]

		// Convert string to lowercase and remove all whitespace
		name := strings.Replace(strings.ToLower(show.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
}
