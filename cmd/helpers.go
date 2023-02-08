package cmd

import (
	"errors"
	"fmt"
)

func formatBoolAnswer(answer string) (bool, error) {
	switch answer {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, errors.New("wrong format")
	}
}

func askBinaryQuestion(questionText string) (bool, error) {
	var param string

	fmt.Print(questionText)
	_, err := fmt.Scanln(&param)
	if err != nil {
		return false, err
	}

	// validate

	return formatBoolAnswer(param)
}

func askStringQuestion(questionText string) (string, error) {
	var param string

	fmt.Print(questionText)
	_, err := fmt.Scanln(&param)
	if err != nil {
		return "", err
	}

	// validate

	return param, nil
}
