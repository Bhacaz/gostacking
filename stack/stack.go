package stack

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
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
    // read the json file called gostacking.json in the current directory
    // return the currentStack value
    data, err := LoadStacks()
    if err != nil {
        fmt.Println("Error loading JSON:", err)
        return ""
    }

    return data.CurrentStack
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
    }
    fmt.Println("Stack module add called", branchName)
}
