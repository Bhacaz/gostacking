package stack

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
)

const stacksFile string = ".git/gostacking.json"

type Stack struct {
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

type StacksData struct {
	CurrentStack string     `json:"currentStack"`
	Stacks        []Stack   `json:"stacks"`
}

func LoadStacks() (StacksData, error) {
	var data StacksData

	jsonData, err := ioutil.ReadFile(stacksFile)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (data StacksData) GetStackByName(stackName string) (*Stack, error) {
    for i, stack := range data.Stacks {
        if stack.Name == stackName {
            return &data.Stacks[i], nil
        }
    }
    return &Stack{}, fmt.Errorf("stack with name %s not found", stackName)
}

func (data StacksData) GetBranchesByName(stackName string) ([]string, error) {
	stack, _ := data.GetStackByName(stackName)
    return stack.Branches, nil
}
