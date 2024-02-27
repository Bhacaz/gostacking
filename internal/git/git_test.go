package git

import (
	"errors"
	"strings"
	"testing"
)

type ExecutorStub struct {
	stubFunc func(string) (string, error)
}

func (es ExecutorStub) execCommand(command string) (string, error) {
	return es.stubFunc(command)
}

func cmdStub(f func(string) (string, error)) Commands {
	return Commands{
		executor: ExecutorStub{
			stubFunc: f,
		},
	}
}

func TestCurrentBranchName(t *testing.T) {
	t.Run("Get current branch name", func(t *testing.T) {
		gitCmd := cmdStub(func(string) (string, error) {
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
		gitCmd := cmdStub(func(string) (string, error) {
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
		gitCmd := cmdStub(func(string) (string, error) {
			return "my_feature_part1", nil
		})

		result := gitCmd.BranchExists("my_feature_part1")

		if !result {
			t.Errorf("got false, want true")
		}
	})

	t.Run("Branch does not exist", func(t *testing.T) {
		gitCmd := cmdStub(func(string) (string, error) {
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
		gitCmd := cmdStub(func(string) (string, error) {
			return "Everything up-to-date", nil
		})

		gitCmd.pushBranch("my_feature_part1")
	})

	t.Run("Error", func(t *testing.T) {
		gitCmd := cmdStub(func(string) (string, error) {
			return "", errors.New("git command error")
		})

		gitCmd.pushBranch("my_feature_part1")
	})
}

func TestCheckout(t *testing.T) {
	t.Run("Checkout branch", func(t *testing.T) {
		gitCmd := cmdStub(func(string) (string, error) {
			return "Switched to branch 'my_feature_part1'", nil
		})

		gitCmd.Checkout("my_feature_part1")
	})

	t.Run("Error", func(t *testing.T) {
		gitCmd := cmdStub(func(string) (string, error) {
			return "", errors.New("git command error")
		})

		gitCmd.Checkout("my_feature_part1")
	})
}

func TestSyncBranches(t *testing.T) {
	t.Run("Git not clean", func(t *testing.T) {
		gitCmd := cmdStub(func(cmd string) (string, error) {
			// case cmd include status
			if strings.Contains(cmd, "status") {
				return "M some_file.go", nil
			}
			t.Errorf("git command should not have been called: %s", cmd)
			return "", errors.New("git command not found")
		})

		gitCmd.SyncBranches([]string{"my_feature_part1", "my_feature_part2"}, "my_feature_part1", false)
	})

	t.Run("Fetch error", func(t *testing.T) {
		gitCmd := cmdStub(func(cmd string) (string, error) {
			// case cmd include status
			if strings.Contains(cmd, "status") {
				return "", nil
			} else if strings.Contains(cmd, "fetch") {
				return "", errors.New("git error fetching")
			}
			t.Errorf("git command should not have been called: %s", cmd)
			return "", errors.New("git command not found")
		})

		gitCmd.SyncBranches([]string{"my_feature_part1", "my_feature_part2"}, "my_feature_part1", false)
	})

	t.Run("Checkout error", func(t *testing.T) {
		gitCmd := cmdStub(func(cmd string) (string, error) {
			// case cmd include status
			if strings.Contains(cmd, "status") {
				return "", nil
			} else if strings.Contains(cmd, "fetch") {
				return "", nil
			} else if strings.Contains(cmd, "checkout") {
				return "", errors.New("git error fetching")
			}
			return "", errors.New("git command not found")
		})

		gitCmd.SyncBranches([]string{"my_feature_part1", "my_feature_part2"}, "my_feature_part1", false)
	})
}
