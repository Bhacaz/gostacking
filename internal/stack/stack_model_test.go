package stack

import (
    "testing"
    )

func stacksDataStub() StacksData {
    return StacksData{
        CurrentStack: "stack1",
        Stacks: []Stack{
            Stack{
                Name: "stack1",
                Branches: []string{"branch1", "branch2"},
            },
            Stack{
                Name: "stack2",
                Branches: []string{"branch3", "branch4"},
            },
        },
    }
}

func stacksLoaderStub() StacksLoader {
    return StacksLoader{
        load: func() (StacksData, error) {
            return stacksDataStub(), nil
        },
    }
}

func TestLoadStacks(t *testing.T) {
    loader := stacksLoaderStub()
    data, _ := loader.LoadStacks()

    if data.CurrentStack != "stack1" {
        t.Errorf("got %s, want %s", data.CurrentStack, "stack1")
    }

    if data.Stacks[0].Name != "stack1" {
        t.Errorf("got %s, want %s", data.Stacks[0].Name, "stack1")
    }

    if data.Stacks[1].Name != "stack2" {
        t.Errorf("got %s, want %s", data.Stacks[1].Name, "stack2")
    }
}
