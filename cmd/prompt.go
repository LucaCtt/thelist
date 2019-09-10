package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Prompt represents an user input interface.
type Prompt interface {
	Input(label string) (string, error)
	Select(label string, options []string) (int, error)
	MultiSelect(label string, options []string) ([]int, error)
}

// CliPrompt represents a cli user input interface.
type CliPrompt struct {
}

// Input allows to ask the user for a single line input.
func (c *CliPrompt) Input(label string) (string, error) {
	prompt := &survey.Input{
		Message: label,
	}

	var input string
	err := survey.AskOne(prompt, &input)
	if err != nil {
		return "", fmt.Errorf("prompt input failed: %w", err)
	}

	return input, nil
}

// Select allows the user to select a single option between the available ones.
// Will return the index of the selected option.
func (c *CliPrompt) Select(label string, options []string) (int, error) {
	prompt := &survey.Select{
		Message: label,
		Options: options,
	}

	var index int
	err := survey.AskOne(prompt, &index)
	if err != nil {
		return 0, fmt.Errorf("prompt select failed: %w", err)
	}

	return index, nil
}

// MultiSelect allows the user to select multiple options between the available ones.
// Will return a slice containig the indexes of the selected options.
func (c *CliPrompt) MultiSelect(label string, options []string) ([]int, error) {
	prompt := &survey.MultiSelect{
		Message: label,
		Options: options,
	}

	var indexes []int
	err := survey.AskOne(prompt, &indexes)
	if err != nil {
		return nil, fmt.Errorf("prompt multiselect failed: %w", err)
	}

	return indexes, nil
}
