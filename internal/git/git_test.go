package git

import (
	"errors"
	"strings"
	"testing"
)

type executorStub struct {
	stubFunc func(...string) (string, error)
}

func (es executorStub) execCommand(command ...string) (string, error) {
	return es.stubFunc(command...)
}

func cmdStub(f func(...string) (string, error)) Commands {
	return Commands{
		executor: executorStub{
			stubFunc: f,
		},
	}
}

func TestCurrentBranchName(t *testing.T) {
	t.Run("Get current branch name", func(t *testing.T) {
		gitCmd := cmdStub(func(...string) (string, error) {
			return "master", nil
		})

		result, err := gitCmd.CurrentBranchName()

		want := "master"
		if result != want {
			t.Errorf("got %s, want %s", result, want)
		}

		if err != nil {
			t.Errorf("got Error %s, want none", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		gitCmd := cmdStub(func(...string) (string, error) {
			return "", errors.New("git command error")
		})

		_, err := gitCmd.CurrentBranchName()

		if err == nil {
			t.Errorf("got none, want Error")
		}
	})
}

func TestBranchExists(t *testing.T) {
	t.Run("Branch exists", func(t *testing.T) {
		gitCmd := cmdStub(func(...string) (string, error) {
			return "my_feature_part1", nil
		})

		result := gitCmd.BranchExists("my_feature_part1")

		if !result {
			t.Errorf("got false, want true")
		}
	})

	t.Run("Branch does not exist", func(t *testing.T) {
		gitCmd := cmdStub(func(...string) (string, error) {
			return "", errors.New("fatal: Needed a single revision")
		})

		result := gitCmd.BranchExists("random_branch")

		if result {
			t.Errorf("got true, want false")
		}
	})
}

func TestPushBranch(t *testing.T) {
	t.Run("Push branch", func(t *testing.T) {
		gitCmd := cmdStub(func(...string) (string, error) {
			return "Everything up-to-date", nil
		})

		gitCmd.pushBranch("my_feature_part1")
	})

	t.Run("Error", func(t *testing.T) {
		gitCmd := cmdStub(func(...string) (string, error) {
			return "", errors.New("git command error")
		})

		gitCmd.pushBranch("my_feature_part1")
	})
}

func TestCheckout(t *testing.T) {
	t.Run("Checkout branch", func(t *testing.T) {
		gitCmd := cmdStub(func(...string) (string, error) {
			return "Switched to branch 'my_feature_part1'", nil
		})

		gitCmd.Checkout("my_feature_part1")
	})
}

func TestSyncBranches(t *testing.T) {
	t.Run("Git not clean", func(t *testing.T) {
		gitCmd := cmdStub(func(cmd ...string) (string, error) {
			if strings.Contains(strings.Join(cmd, " "), "status") {
				return "M some_file.go", nil
			}
			t.Errorf("git command should not have been called: %s", cmd)
			return "", errors.New("git command error")
		})

		gitCmd.SyncBranches([]string{"my_feature_part1", "my_feature_part2"}, "my_feature_part1", false)
	})

	t.Run("Merge error", func(t *testing.T) {
		gitCmd := cmdStub(func(cmd ...string) (string, error) {
			cmdString := strings.Join(cmd, " ")
			if strings.Contains(cmdString, "status") {
				return "", nil
			} else if strings.Contains(cmdString, "fetch") {
				return "", nil
			} else if strings.Contains(cmdString, "checkout") {
				return "", nil
			} else if strings.Contains(cmdString, "pull") {
				return "", nil
			} else if strings.Contains(cmdString, "merge") {
				if strings.Contains(cmdString, "into my_feature_part1") {
					t.Errorf("nothing should be merge in the first branch. cmd: %s", cmd)
				}
				return "", errors.New("git error merge")
			}
			t.Errorf("git command should not have been called: %s", cmd)
			return "", errors.New("git command error")
		})

		gitCmd.SyncBranches([]string{"my_feature_part1", "my_feature_part2"}, "my_feature_part1", false)
	})

	t.Run("Merge", func(t *testing.T) {
		var part1Merged = false
		var lastCheckoutMain = false

		gitCmd := cmdStub(func(cmd ...string) (string, error) {
			cmdString := strings.Join(cmd, " ")
			if strings.Contains(cmdString, "status") {
				return "", nil
			} else if strings.Contains(cmdString, "fetch") {
				return "", nil
			} else if strings.Contains(cmdString, "checkout") {
				if strings.Contains(cmdString, "main") {
					lastCheckoutMain = true
				}
				return "", nil
			} else if strings.Contains(cmdString, "pull") {
				return "", nil
			} else if strings.Contains(cmdString, "merge") {
				if strings.Contains(cmdString, "into my_feature_part1") {
					t.Errorf("nothing should be merge in the first branch. cmd: %s", cmd)
				} else if strings.Contains(cmdString, "into my_feature_part2") {
					part1Merged = true
				}
				return "", nil
			}
			t.Errorf("git command should not have been called: %s", cmdString)
			return "", errors.New("git command error")
		})

		gitCmd.SyncBranches([]string{"my_feature_part1", "my_feature_part2"}, "main", false)
		if !part1Merged {
			t.Errorf("my_feature_part1 should have been merged into my_feature_part2")
		}
		if !lastCheckoutMain {
			t.Errorf("last checkout should have been main")
		}
	})

	t.Run("Merge with push", func(t *testing.T) {
		var part1Merged = false
		var lastCheckoutMain = false

		gitCmd := cmdStub(func(cmd ...string) (string, error) {
			cmdString := strings.Join(cmd, " ")
			if strings.Contains(cmdString, "status") {
				return "", nil
			} else if strings.Contains(cmdString, "fetch") {
				return "", nil
			} else if strings.Contains(cmdString, "checkout") {
				if strings.Contains(cmdString, "main") {
					lastCheckoutMain = true
				}
				return "", nil
			} else if strings.Contains(cmdString, "pull") {
				return "", nil
			} else if strings.Contains(cmdString, "merge") {
				if strings.Contains(cmdString, "into my_feature_part1") {
					t.Errorf("nothing should be merge in the first branch. cmd: %s", cmd)
				} else if strings.Contains(cmdString, "into my_feature_part2") {
					part1Merged = true
				}
				return "", nil
			} else if strings.Contains(cmdString, "push") {
				return "", nil
			}
			t.Errorf("git command should not have been called: %s", cmd)
			return "", errors.New("git command error")
		})

		gitCmd.SyncBranches([]string{"my_feature_part1", "my_feature_part2"}, "main", true)
		if !part1Merged {
			t.Errorf("my_feature_part1 should have been merged into my_feature_part2")
		}
		if !lastCheckoutMain {
			t.Errorf("last checkout should have been main")
		}
	})
}
