package common

import (
	"github.com/AlecAivazis/survey/v2"
)

// Prompter represents an user input interface.
type Prompter interface {
	Input(label string) string
	Select(label string, options []string) int
	MultiSelect(label string, options []string) []int
}

// CliPrompter represents a cli user input interface.
type CliPrompter struct {
}

// Input allows to ask the user for a single line input.
func (c *CliPrompter) Input(label string) string {
	prompt := &survey.Input{
		Message: label,
	}

	var show string
	survey.AskOne(prompt, &show)
	return show
}

// Select allows the user to select a single option between the available ones.
// Will return the index of the selected option.
func (c *CliPrompter) Select(label string, options []string) int {
	prompt := &survey.Select{
		Message: label,
		Options: options,
	}

	var index int
	survey.AskOne(prompt, &index)
	return index
}

// MultiSelect allows the user to select multiple options between the available ones.
// Will return a slice containig the indexes of the selected options.
func (c *CliPrompter) MultiSelect(label string, options []string) []int {
	prompt := &survey.MultiSelect{
		Message: label,
		Options: options,
	}

	var indexes []int
	survey.AskOne(prompt, &indexes)
	return indexes
}
