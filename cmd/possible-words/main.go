package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/mrankine/duotrigordle/pkg"
)

func main() {
	dictGuesses, dictSolutions := pkg.ReadWordLists()

	var showPossible bool
	flag.BoolVar(&showPossible, "w", false, "show solutions and valid guesses")
	flag.BoolVar(&showPossible, "show-words", false, "show solutions and valid guesses")
	flag.Parse()

	inputArgs := flag.Args()

	if len(inputArgs) == 0 || len(inputArgs)%2 != 0 {
		log.Fatalln("requires pairs of guesses and responses")
	}

	guesses := []string{}
	responses := []pkg.Response{}

	for i := 0; i+1 < len(inputArgs); i += 2 {
		guess := strings.ToUpper(inputArgs[i])
		if !slices.Contains(dictGuesses, guess) {
			log.Fatalf("invalid guess")
		}
		guesses = append(guesses, guess)
		response, err := pkg.ParseResponse(strings.ToUpper(inputArgs[i+1]))
		if err != nil {
			log.Fatal(err)
		}
		responses = append(responses, response)

		fmt.Printf("%s ", pkg.FormatResponse(guess, response))
	}

	possibleGuesses := dictGuesses
	possibleSolutions := dictSolutions
	for i, guess := range guesses {
		p := []string{}
		for _, word := range possibleGuesses {
			if pkg.CheckGuess(guess, word) == responses[i] {
				p = append(p, word)
			}
		}
		possibleGuesses = p
		q := []string{}
		for _, word := range possibleSolutions {
			if pkg.CheckGuess(guess, word) == responses[i] {
				q = append(q, word)
			}
		}
		possibleSolutions = q
	}

	fmt.Printf("%d remaining possible solutions, %d remaining possible valid guesses\n",
		len(possibleSolutions),
		len(possibleGuesses))

	if showPossible && len(possibleSolutions) > 0 {
		fmt.Printf("Solutions: %s\n", pkg.FormatList(possibleSolutions, 10, " "))
	}

	if showPossible && len(possibleGuesses) > 0 {
		fmt.Printf("Valid guesses: %s\n", pkg.FormatList(possibleGuesses, 10, " "))
	}
}
