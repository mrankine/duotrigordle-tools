package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/mrankine/duotrigordle/pkg"
)

func main() {
	dictGuesses, dictSolutions := pkg.ReadWordLists()

	// test every guess against every solution
	// build map counting number of solutions for each response
	// we are looking for guess-response pairs with one solution
	mapCountSolutions := make(map[string]map[pkg.Response]int)
	for _, i := range dictGuesses {
		mapCountSolutions[i] = make(map[pkg.Response]int)
		for _, j := range dictSolutions {
			response := pkg.CheckGuess(i, j)
			if response != pkg.AllCorrect {
				mapCountSolutions[i][response] += 1
			}
		}
	}

	bestGuesses := []string{}
	bestSolvableBoards := 0

	for _, guess := range dictGuesses {
		solvableBoards := 0
		for _, solution := range dictSolutions {
			response := pkg.CheckGuess(guess, solution)
			if mapCountSolutions[guess][response] == 1 {
				solvableBoards++
			}
		}

		if solvableBoards > bestSolvableBoards {
			bestGuesses = []string{guess}
			bestSolvableBoards = solvableBoards
		} else if solvableBoards == bestSolvableBoards {
			bestGuesses = append(bestGuesses, guess)
		}
	}

	fmt.Printf("First guess which maximises solvable boards: %s\n", strings.Join(bestGuesses, " or "))

	solvableFraction := float64(bestSolvableBoards) / float64(len(dictSolutions))
	fmt.Printf("%d of %d boards (%.1f%%) are solvable on next guess\n",
		bestSolvableBoards,
		len(dictSolutions),
		100*solvableFraction)

	const boards = 32
	fmt.Printf("%.1f%% probability of solvable board on grid of %d\n",
		100*(1-math.Pow(1-solvableFraction, float64(boards))),
		boards)
}
