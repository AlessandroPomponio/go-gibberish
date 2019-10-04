// Package training contains the functions used
// to train the model.
package training

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/AlessandroPomponio/go-gibberish/analysis"
	"github.com/AlessandroPomponio/go-gibberish/persistence"
	"github.com/AlessandroPomponio/go-gibberish/structs"
)

// TrainModel computes the probabilities of having a certain
// digraph by reading a big file.
func TrainModel(acceptedChars, trainingFileName, goodFileName, badFileName, outputFileName string) error {

	position := getRunePosition(acceptedChars)

	// Assume we have seen 10 of each character pair.  This acts as a kind of
	// prior or smoothing factor.  This way, if we see a character transition
	// live that we've never observed in the past, we won't assume the entire
	// string has 0 probability.
	occurrences := initializeOccurrencesMatrix(len(acceptedChars))

	trainingFile, err := os.Open(trainingFileName)
	if err != nil {
		return fmt.Errorf("TrainModel: unable to open training file %s", trainingFileName)
	}

	// Count the occurrences of rune pairs by reading a big file.
	tfReader := bufio.NewReader(trainingFile)
	for {

		line, _, err := tfReader.ReadLine()
		if err != nil {
			break
		}

		for _, pair := range analysis.GetDigraphs(string(line)) {

			firstPosition, firstRuneFound := position[pair.First]
			if !firstRuneFound {
				return fmt.Errorf("TrainModel: unable to find the position of the rune %s", string(pair.First))
			}

			secondPosition, secondRuneFound := position[pair.Second]
			if !secondRuneFound {
				return fmt.Errorf("TrainModel: unable to find the position of the rune %s", string(pair.First))
			}

			occurrences[firstPosition][secondPosition]++

		}

	}
	_ = trainingFile.Close()

	// Normalize the counts so that they become log probabilities.
	// We use log probabilities rather than straight probabilities to avoid
	// numeric underflow issues with long texts.
	// This contains a justification:
	// http://squarecog.wordpress.com/2009/01/10/dealing-with-underflow-in-joint-probability-calculations/
	normalizeOccurrencesMatrix(occurrences)

	// Find the probability of generating a few arbitrarily chosen good and bad phrases.
	goodProbabilities, err := averageTransitionProbabilitiesInFile(goodFileName, occurrences, position)
	if err != nil {
		return fmt.Errorf("TrainModel: error when computing good probabilities: %s", err)
	}

	badProbabilities, err := averageTransitionProbabilitiesInFile(badFileName, occurrences, position)
	if err != nil {
		return fmt.Errorf("TrainModel: error when computing bad probabilities: %s", err)
	}

	minimumGoodProbability := analysis.MinForSlice(goodProbabilities)
	maximumBadProbability := analysis.MaxForSlice(badProbabilities)

	// Make sure we are actually capable of detecting the junk.
	if minimumGoodProbability <= maximumBadProbability {
		return fmt.Errorf("TrainModel: minimumGoodProbability <= maximumBadProbability")
	}

	// Pick a threshold halfway between the worst good and best bad inputs.
	threshold := (minimumGoodProbability + maximumBadProbability) / 2

	data := structs.GibberishData{
		Occurrences: occurrences,
		Positions:   position,
		Threshold:   threshold,
	}

	err = persistence.WriteKnowledgeBase(&data, outputFileName)
	return err

}

func averageTransitionProbabilitiesInFile(fileName string, occurrences [][]float64, position map[rune]int) ([]float64, error) {

	res := make([]float64, 0, 5)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("averageTransitionProbabilitiesInFile: unable to open file %s", fileName)
	}

	fReader := bufio.NewReader(file)
	for {

		line, _, err := fReader.ReadLine()
		if err != nil {
			break
		}

		avgProb, err := analysis.AverageTransitionProbability(string(line), occurrences, position)
		if err != nil {
			break
		}

		res = append(res, avgProb)

	}

	_ = file.Close()
	return res, err

}

func getRunePosition(characters string) map[rune]int {

	position := make(map[rune]int)
	for index, currentRune := range characters {
		position[currentRune] = index
	}

	return position

}

func initializeOccurrencesMatrix(symbols int) [][]float64 {

	occurrences := make([][]float64, symbols)
	for row := range occurrences {
		occurrences[row] = make([]float64, symbols)
		for column := range occurrences[row] {
			occurrences[row][column] = 10
		}
	}

	return occurrences

}

func normalizeOccurrencesMatrix(occurrences [][]float64) {

	for _, row := range occurrences {

		sum := 0.

		for i := 0; i < len(row); i++ {
			sum += row[i]
		}

		for i := 0; i < len(row); i++ {
			row[i] = math.Log(row[i] / sum)
		}

	}

}
