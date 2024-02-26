package stack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const stacksFile string = ".git/gostacking.json"

type Stack struct {
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

type StacksData struct {
	CurrentStack string  `json:"currentStack"`
	Stacks       []Stack `json:"stacks"`
}

type StacksPersisting interface {
	LoadStacks() (StacksData, error)
	SaveStacks(StacksData)
}

type StacksPersistingFile struct{}

func (s StacksPersistingFile) LoadStacks() (StacksData, error) {
	return LoadStacksFromFile()
}

func (s StacksPersistingFile) SaveStacks(data StacksData) {
	SaveStacks(data)
}

func LoadStacksFromFile() (StacksData, error) {
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

func SaveStacks(data StacksData) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	//     fmt.Println(string(jsonData))
	// Write the JSON data to a file
	err = ioutil.WriteFile(stacksFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
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
