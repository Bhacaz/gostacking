package stack

import "github.com/Bhacaz/gostacking/internal/color"

func colorFunc(i int) func(...interface{}) string {
	var branchColorsSequence []func(...interface{}) string
	branchColorsSequence = append(branchColorsSequence, color.Red)
	branchColorsSequence = append(branchColorsSequence, color.Purple)
	branchColorsSequence = append(branchColorsSequence, color.Magenta)
	branchColorsSequence = append(branchColorsSequence, color.Green)
	branchColorsSequence = append(branchColorsSequence, color.Teal)

	return branchColorsSequence[i%len(branchColorsSequence)]
}

func pipesColors(index int, withDivergence bool) string {
	var result string
	for i := 0; i < index; i++ {
		colorF := colorFunc(i)
		if i == index-1 && withDivergence {
			result += colorF("|\\\n")
		} else {
			result += colorF("| ")
		}
	}
	return result
}
