package cmd

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
)

func formatBoolAnswer(answer string) (bool, error) {
	switch answer {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, errors.New("wrong format")
	}
}

func askBinaryQuestion(questionText string) (bool, error) {

	var param string

	prompt := promptui.Select{
		Label: questionText,
		Items: []string{"yes", "no"},
	}

	_, param, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false, err
	}

	// validate

	return formatBoolAnswer(param)
}

func askStringQuestion(questionText string) (string, error) {

	var param string

	prompt := promptui.Prompt{
		Label: questionText,
		// Validate: ,
	}

	param, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	// validate

	return param, nil
}

func askCustomSelectable(questionText string, choices []string) (string, error) {

	var param string

	prompt := promptui.Select{
		Label: questionText,
		Items: choices,
	}

	_, param, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	// validate

	return param, nil
}
