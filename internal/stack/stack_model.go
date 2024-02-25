package stack

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "github.com/Bhacaz/gostacking/internal/git"
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

type StacksLoader struct {
    load func() (StacksData, error)
}

func stacksLoaderFromFile() StacksLoader {
    return StacksLoader{
        load: LoadStacksFromFile,
    }
}

func (loader StacksLoader) LoadStacks() (StacksData, error) {
    return loader.load()
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

func New(stackName string) {
    newStack := Stack{
                    Name: stackName,
                    Branches: []string{git.CurrentBranchName()},
                }

    data, _ := stacksLoaderFromFile().load()
    data.CurrentStack = stackName
    data.Stacks = append(data.Stacks, newStack)

    data.SaveStacks()
	fmt.Println("New stack created", stackName)
}

func (data StacksData) SaveStacks() {
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
