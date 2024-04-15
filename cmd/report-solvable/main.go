package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/mrankine/duotrigordle/pkg"
)

type SequenceInfo struct {
	guessSequence []string
	solution      string
	lastResponse  pkg.Response
}

func main() {
	dictGuesses, dictSolutions := pkg.ReadWordLists()

	if len(os.Args) < 2 {
		log.Fatalf("no guesses provided")
	}

	var guessSequence []string
	for _, arg := range os.Args[1:] {
		guess := strings.ToUpper(arg)
		if slices.Contains(dictGuesses, guess) {
			guessSequence = append(guessSequence, guess)
		} else {
			log.Fatalf("invalid guess '%s'", arg)
		}
	}

	fmt.Printf("Guesses: %s\n", pkg.FormatList(guessSequence, 0, " "))

	solvableBoards := []SequenceInfo{}
	totalPossibleSolutions := 0
	for _, solution := range dictSolutions {
		possibleSolutions := dictSolutions
		var response pkg.Response
		for _, guess := range guessSequence {
			if guess == solution {
				continue
			}
			response = pkg.CheckGuess(guess, solution)
			remaining := []string{}
			for _, possibleSolution := range possibleSolutions {
				if pkg.CheckGuess(guess, possibleSolution) == response {
					remaining = append(remaining, possibleSolution)
				}
			}
			possibleSolutions = remaining
		}

		totalPossibleSolutions += len(possibleSolutions)
		if len(possibleSolutions) == 1 {
			solvableBoards = append(solvableBoards, SequenceInfo{solution: solution, guessSequence: guessSequence, lastResponse: response})
		}
	}

	slices.SortFunc(solvableBoards, func(a SequenceInfo, b SequenceInfo) int {
		return int(b.lastResponse) - int(a.lastResponse)
	})

	for _, board := range solvableBoards {
		fmt.Println(pkg.FormatSequence(board.guessSequence, board.solution))
	}

	solvableFraction := float64(len(solvableBoards)) / float64(len(dictSolutions))
	fmt.Printf("%d of %d boards (%.1f%%) are solvable on next guess\n",
		len(solvableBoards),
		len(dictSolutions),
		100*solvableFraction)

	const boards = 32
	fmt.Printf("%.1f%% probability of solvable board on grid of %d\n",
		100*(1-math.Pow(1-solvableFraction, float64(boards))),
		boards)

	fmt.Printf("%.1f average possible solutions per board\n",
		float64(totalPossibleSolutions)/float64(len(dictSolutions)))
}
