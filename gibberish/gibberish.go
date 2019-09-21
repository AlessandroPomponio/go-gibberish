// Package gibberish contains methods to tell whether
// the input is gibberish or not.
package gibberish

import (
	"github.com/AlessandroPomponio/go-gibberish/analysis"
	"github.com/AlessandroPomponio/go-gibberish/structs"
)

// IsGibberish returns true if the input string is likely
// to be gibberish
func IsGibberish(input string, data *structs.GibberishData) bool {
	return analysis.AverageTransitionProbability(input, data.Occurrences, data.Positions) <= data.Threshold
}
