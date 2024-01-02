package stack

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "os"
    "github.com/gostacking/git"
)

func CurrentStackName() string {
    fmt.Println("Stack module currentStackName called")
    // read the json file called gostacking.json in the current directory
    // return the currentStack value

    // Open our jsonFile
    jsonFile, err := os.Open("gostacking.json")
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Successfully Opened gostacking.json")
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    // read our opened jsonFile as a byte array.
    byteValue, _ := ioutil.ReadAll(jsonFile)

    // we initialize our Users array
    var currentStack map[string]string

    // we unmarshal our byteArray which contains our
    // jsonFile's content into 'users' which we defined above
    json.Unmarshal(byteValue, &currentStack)

    return currentStack["currentStack"]
}

func New(stackName string) {
    fmt.Println("Stack module new called")
    // create a json file called gostacking.json in the current directory

    currentStack := map[string]string{ "currentStack": stackName }

    jsonData, err := json.MarshalIndent(currentStack, "", "    ")
        if err != nil {
            fmt.Println("Error marshaling JSON:", err)
            return
        }

	// Write the JSON data to a file
	err = ioutil.WriteFile("gostacking.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("New stack created")
}

func Add(branchName string) {
     if branchName == "" {
        branchName = git.CurrentBranchName()
    }
    fmt.Println("Stack module add called", branchName)
}
