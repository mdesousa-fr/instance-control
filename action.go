//go:generate go-enum --names
package main

import (
	"github.com/AlecAivazis/survey/v2"
)

// ENUM(start, stop)
type ActionType int

type ActionTypeAnswer struct {
	Value int `survey:"action"`
}

func ActionSelect() (ActionType, error) {
	var answer ActionTypeAnswer
	var qs = []*survey.Question{
		{
			Name: "action",
			Prompt: &survey.Select{
				Message: "What did you want to do?",
				Options: ActionTypeNames(),
				Default: ActionTypeStart.String(),
			},
		},
	}
	if err := survey.Ask(qs, &answer); err != nil {
		return ActionTypeStart, err
	}
	return ActionType(answer.Value), nil
}
