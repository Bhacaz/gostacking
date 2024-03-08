package stack

func stacksDataMock() *StacksData {
	return &StacksData{
		CurrentStack: "stack1",
		Stacks: []Stack{
			Stack{
				Name:     "stack1",
				Branches: []string{"branch1", "branch2"},
			},
			Stack{
				Name:     "stack2",
				Branches: []string{"branch3", "branch4"},
			},
		},
		StacksPersister: StacksPersistingStub{},
	}
}

type StacksPersistingStub struct {
	Data *StacksData
}

func (s StacksPersistingStub) LoadStacks(data *StacksData) {
	//if reflect.ValueOf(s.Data).IsZero() {
	//	s.Data = data
	//}
}

func (s StacksPersistingStub) SaveStacks(data StacksData) {
	//s.Data = &data
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
