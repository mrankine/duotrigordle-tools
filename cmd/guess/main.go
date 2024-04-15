package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/mrankine/duotrigordle/pkg"
)

func main() {
	dictGuesses, dictSolutions := pkg.ReadWordLists()

	if len(os.Args) != 3 {
		log.Fatalf("expected guess and solution args")
	}

	guess := strings.ToUpper(os.Args[1])
	if !slices.Contains(dictGuesses, guess) {
		log.Fatalf("invalid guess")
	}

	solution := strings.ToUpper(os.Args[2])
	if !slices.Contains(dictSolutions, solution) {
		log.Fatalf("invalid solution")
	}

	response := pkg.CheckGuess(guess, solution)
	fmt.Println(pkg.FormatResponse(guess, response))
}
