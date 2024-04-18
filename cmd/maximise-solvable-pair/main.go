package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mrankine/duotrigordle/pkg"
)

func main() {
	dictGuesses, dictSolutions := pkg.ReadWordLists()

	var empty struct{}

	fmt.Print("Building map... ")
	// map guess -> (map response -> possible solutions)
	mapSolutions := make(map[string]map[pkg.Response][]string)
	for _, guess := range dictGuesses {
		mapSolutions[guess] = make(map[pkg.Response][]string)
		for _, solution := range dictSolutions {
			response := pkg.CheckGuess(guess, solution)
			mapSolutions[guess][response] = append(mapSolutions[guess][response], solution)
		}
	}
	fmt.Println("done")

	bestSolvableBoards := 0
	bestPair := []string{}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("Terminated! Best pair found so far is %s %d\n", bestPair, bestSolvableBoards)
		os.Exit(1)
	}()

	for i, firstGuess := range dictGuesses {
		// build a map of the letters in firstGuess. if any are repeated, skip.
		mapFirstGuessLetters := make(map[rune]struct{}, pkg.Length)
		skip := false
		for _, letter := range firstGuess {
			_, found := mapFirstGuessLetters[letter]
			if found {
				skip = true
				break
			} else {
				mapFirstGuessLetters[letter] = empty
			}
		}
		if skip {
			fmt.Printf("%s skipped\n", firstGuess)
			continue
		}

		for _, secondGuess := range dictGuesses {
			if secondGuess == firstGuess {
				continue
			}
			// skip secondGuess if it contains any of the letters in firstGuess, unlikely to be a winner
			// also skip if any repeat letters
			mapSecondGuessLetters := make(map[rune]struct{}, pkg.Length)
			skip = false
			for _, letter := range secondGuess {
				_, found := mapFirstGuessLetters[letter]
				if found {
					skip = true
					break
				} else {
					_, found := mapSecondGuessLetters[letter]
					if found {
						skip = true
						break
					} else {
						mapSecondGuessLetters[letter] = empty
					}
				}
			}
			if skip {
				continue
			}

			solvableBoards := 0

			for _, possibleSolutions := range mapSolutions[firstGuess] {
				mapResponses := make(map[pkg.Response]int)
				for _, possibleSolution := range possibleSolutions {
					response := pkg.CheckGuess(secondGuess, possibleSolution)
					mapResponses[response] += 1
				}
				for _, count := range mapResponses {
					if count == 1 {
						solvableBoards += 1
					}
				}
			}

			if solvableBoards > bestSolvableBoards {
				bestSolvableBoards = solvableBoards
				bestPair = []string{firstGuess, secondGuess}
				fmt.Printf("*** NEW BEST PAIR %s %s %d ***\n", bestPair[0], bestPair[1], bestSolvableBoards)
			}
		}
		fmt.Printf("%s done (%.2f%%)\n", firstGuess, float32(i)/float32(len(dictGuesses))*100)
	}

	fmt.Printf("Best pair found is %s %d\n", bestPair, bestSolvableBoards)
}
