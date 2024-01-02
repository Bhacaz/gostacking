package stack

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "slices"
    "github.com/gostacking/git"
)

type Stack struct {
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

type StacksData struct {
	CurrentStack string `json:"currentStack"`
	Stacks        []Stack     `json:"stacks"`
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

func WriteStacksFile(stackData StacksData) {
    jsonData, err := json.MarshalIndent(stackData, "", "    ")
        if err != nil {
            fmt.Println("Error marshaling JSON:", err)
            return
        }

    fmt.Println(string(jsonData))
    // Write the JSON data to a file
    err = ioutil.WriteFile("gostacking.json", jsonData, 0644)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
}

func LoadStacks() (StacksData, error) {
	var data StacksData

	jsonData, err := ioutil.ReadFile("gostacking.json")
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func CurrentStackName() string {
    data, err := LoadStacks()
    if err != nil {
        fmt.Println("Error loading JSON:", err)
        return ""
    }

    return data.CurrentStack
}

func CurrentStackStatus() string {
    data, err := LoadStacks()
    if err != nil {
        fmt.Println("Error loading JSON:", err)
        return ""
    }

    var displayBranches string
    branches, _ := data.GetBranchesByName(data.CurrentStack)
    for _, branch := range branches {
        displayBranches += branch + "\n\t"
    }
    return data.CurrentStack + "\n\t" + displayBranches
}

func New(stackName string) {
    stackData := StacksData{
        CurrentStack: stackName,
        Stacks: []Stack{
            Stack{
                Name: stackName,
                Branches: []string{git.CurrentBranchName()},
            },
        },
    }

    WriteStacksFile(stackData)
	fmt.Println("New stack created", stackName)
}

func Add(branchName string) {
     if branchName == "" {
        branchName = git.CurrentBranchName()
    } else {
        if !git.BranchExists(branchName) {
            fmt.Println("Branch", branchName, "does not exist")
        }
    }

    data, _ := LoadStacks()
    stack, _ := data.GetStackByName(data.CurrentStack)
    stack.Branches = append(stack.Branches, branchName)
    stack.Branches = slices.Compact(stack.Branches)
    WriteStacksFile(data)
    fmt.Println("Branch", branchName, "added to stack", data.CurrentStack)
}
