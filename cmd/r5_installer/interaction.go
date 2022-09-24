package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

func gatherRunOptions(options []string) (selectedOptions []string, err error) {
	prompt := &survey.MultiSelect{
		Message: "Please select from the following options",
		Options: options,
		Default: []int{
			0,
			1,
		},
	}

	//_ = dialog.Raise("Use the arrow keys(navigation) and spacebar(toggle selection) in console to continue")
	err = survey.AskOne(prompt, &selectedOptions)
	if err != nil {
		err = fmt.Errorf("error while gathering options from user: %v", err)
	}

	return
}
