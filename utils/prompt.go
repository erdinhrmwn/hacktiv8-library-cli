package utils

import (
	"errors"

	"github.com/manifoldco/promptui"
)

func AskInput(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(s string) error {
			if s == "" {
				return errors.New("input cannot be empty")
			}

			return nil
		},
	}

	return prompt.Run()
}

func SelectInput(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	idx, result, err := prompt.Run()
	return idx, result, err
}
