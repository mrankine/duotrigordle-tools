package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"slices"
	"strings"

	"github.com/mrankine/duotrigordle/pkg"
)

type SequenceInfo struct {
	guesses   []string
	responses []pkg.Response
	solution  string
}

func main() {
	dictGuesses, dictSolutions := pkg.ReadWordLists()

	codes := false
	flag.BoolVar(&codes, "c", false, "output colour codes next to board")
	flag.BoolVar(&codes, "show-codes", false, "output colour codes next to board")
	flag.Parse()

	inputArgs := flag.Args()

	if len(inputArgs) == 0 {
		log.Fatalf("no guesses provided")
	}

	var guessSequence []string
	for _, arg := range inputArgs {
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
		responseSequence := make([]pkg.Response, len(guessSequence))
		for i, guess := range guessSequence {
			if guess == solution {
				continue
			}
			responseSequence[i] = pkg.CheckGuess(guess, solution)
			remaining := []string{}
			for _, possibleSolution := range possibleSolutions {
				if pkg.CheckGuess(guess, possibleSolution) == responseSequence[i] {
					remaining = append(remaining, possibleSolution)
				}
			}
			possibleSolutions = remaining
		}

		totalPossibleSolutions += len(possibleSolutions)
		if len(possibleSolutions) == 1 {
			solvableBoards = append(solvableBoards, SequenceInfo{guesses: guessSequence, responses: responseSequence, solution: solution})
		}
	}

	slices.SortFunc(solvableBoards, func(a SequenceInfo, b SequenceInfo) int {
		for i := 0; i < len(b.responses); i++ {
			diff := int(b.responses[i]) - int(a.responses[i])
			if diff != 0 {
				return diff
			}
		}
		return 0
	})

	for _, board := range solvableBoards {
		pkg.PrintBoard(board.guesses, board.responses, board.solution, codes)
		fmt.Println()
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
