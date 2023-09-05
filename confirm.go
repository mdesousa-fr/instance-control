package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

func GetConfirm(action ActionType, instances []string) (bool, error) {
	var answer bool
	var message = fmt.Sprintf("Are you sure to %s: %s", action, instances)
	var qs = []*survey.Question{
		{
			Name: "confirm",
			Prompt: &survey.Confirm{
				Message: message,
				Default: false,
			},
		},
	}
	if err := survey.Ask(qs, &answer); err != nil {
		return false, err
	}
	return answer, nil
}
