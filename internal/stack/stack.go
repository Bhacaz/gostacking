package stack

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"
	"strings"
)

const stacksFile string = ".git/gostacking.json"

type StacksPersisting interface {
	LoadStacks(data *StacksData)
	SaveStacks(data StacksData)
}

type Stack struct {
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

type StacksData struct {
	CurrentStack    string           `json:"currentStack"`
	Stacks          []Stack          `json:"stacks"`
	StacksPersister StacksPersisting `json:"-"`
}

type StacksPersistingFile struct{}

func (s StacksPersistingFile) LoadStacks(data *StacksData) {
	jsonData, err := os.ReadFile(stacksFile)
	// If the file does not exist, return an empty data
	// Calling SaveStacks will create the file
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return
		} else {
			log.Fatal("Error reading file:", err)
		}
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}
}

func (s StacksPersistingFile) SaveStacks(data StacksData) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	err = os.WriteFile(stacksFile, jsonData, 0644)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}
}

func (data *StacksData) LoadStacks() {
	data.StacksPersister.LoadStacks(data)
}

func (data *StacksData) SaveStacks() {
	data.StacksPersister.SaveStacks(*data)
}

func (data *StacksData) GetStackByName(stackName string) (*Stack, error) {
	for i, stack := range data.Stacks {
		if stack.Name == stackName {
			return &data.Stacks[i], nil
		}
	}
	return &Stack{}, errors.New("Stack " + stackName + " not found")
}

func (data *StacksData) GetStackByBranch(branchName string) (*Stack, error) {
	for i, stack := range data.Stacks {
		if slices.Contains(stack.Branches, branchName) {
			return &data.Stacks[i], nil
		}
	}
	return &Stack{}, errors.New("Branch " + branchName + " not found")
}

func (data *StacksData) GetBranchesByName(stackName string) ([]string, error) {
	stack, err := data.GetStackByName(stackName)
	if err != nil {
		return []string{}, err
	}

	return stack.Branches, nil
}

func (data *StacksData) GetCurrentBranches() ([]string, error) {
	return data.GetBranchesByName(data.CurrentStack)
}

func (data *StacksData) SetCurrentStack(stackName string) {
	data.CurrentStack = stackName
	data.SaveStacks()
}
