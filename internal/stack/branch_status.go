package stack
// tmp
import "github.com/Bhacaz/gostacking/internal/color"

type branchStatus struct {
	BehindRemote        bool
	AheadRemote         bool
	HasDiff             bool
	BehindDefaultBranch bool
}

func defaultBranchStatus() branchStatus {
	return branchStatus{
		BehindRemote:        false,
		AheadRemote:         false,
		HasDiff:             false,
		BehindDefaultBranch: false,
	}
}

func (bs branchStatus) Symbols() string {
	var symbolsToDisplay string
	if bs.BehindRemote {
		symbolsToDisplay += color.Teal("↓")
	}
	if bs.AheadRemote {
		symbolsToDisplay += color.Teal("↑")
	}
	if bs.HasDiff {
		symbolsToDisplay += color.Red("*")
	}
	if bs.BehindDefaultBranch {
		symbolsToDisplay += color.Magenta("*")
	}
	return symbolsToDisplay
}
