package stack

import (
    "reflect"
//     "testing"
)

func stacksDataMock() StacksData {
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

type StacksPersistingStub struct {
    data StacksData
 }

func (s *StacksPersistingStub) LoadStacks() (StacksData, error) {
    if reflect.ValueOf(s.data).IsZero() {
        s.data = stacksDataMock()
    }
    return s.data, nil
}

func (s *StacksPersistingStub) SaveStacks(data StacksData) {
    s.data = data
}

// func TestLoadStacks(t *testing.T) {
//     loader := stacksLoaderStub()
//     data, _ := loader.LoadStacks()
//
//     if data.CurrentStack != "stack1" {
//         t.Errorf("got %s, want %s", data.CurrentStack, "stack1")
//     }
//
//     if data.Stacks[0].Name != "stack1" {
//         t.Errorf("got %s, want %s", data.Stacks[0].Name, "stack1")
//     }
//
//     if data.Stacks[1].Name != "stack2" {
//         t.Errorf("got %s, want %s", data.Stacks[1].Name, "stack2")
//     }
// }
