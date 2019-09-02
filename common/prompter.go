//go:generate mockgen -destination=../mocks/mock_prompt.go -package=mocks github.com/lucactt/thelist/util Prompt

package common

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Prompter represents an user input interface.
type Prompter interface {
	Input() (string, error)
	Select(shows []*Show) (*Show, error)
}

// CliPrompter represents a cli user input interface.
type CliPrompter struct {
}

// Input allows to ask the user for a single line input.
func (c *CliPrompter) Input() (string, error) {
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

// Select allows the user to select a single option between the available ones.
func (c *CliPrompter) Select(shows []*Show) (*Show, error) {
	if len(shows) == 0 {
		return nil, fmt.Errorf("No shows found")
	}

	if len(shows) == 1 {
		return shows[0], nil
	}

	prompt := promptui.Select{
		Label:     "Select one",
		Items:     shows,
		Templates: templates,
		Searcher:  searcher(shows),
	}
	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return shows[i], nil
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

func searcher(shows []*Show) func(string, int) bool {
	return func(input string, index int) bool {
		show := shows[index]

		// Convert string to lowercase and remove all whitespace
		name := strings.Replace(strings.ToLower(show.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
}
