package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const DEFAULT_WORDLIST = "2024-01-14"

func ParseResponse(input string) (Response, error) {
	if input == "G" || input == "2" {
		return AllCorrect, nil
	} else if input == "Y" || input == "1" {
		return AllPresent, nil
	} else if input == "B" || input == "0" {
		return AllAbsent, nil
	} else if len(input) == 1 {
		return 0, errors.New("invalid response: valid chars are G, Y, B, 2, 1, or 0")
	} else if len(input) != Length {
		return 0, fmt.Errorf("invalid response: length must be 1 or %d", Length)
	}

	var response Response
	for i := range Length {
		switch input[i] {
		case 'G', '2':
			response |= (1 << (Length + Length - i - 1))
		case 'Y', '1':
			response |= (1 << (Length - i - 1))
		case 'B', '0':
			continue
		default:
			return 0, errors.New("invalid response: valid chars are G, Y, B, 2, 1, or 0")
		}
	}
	return response, nil
}

func FormatResponse(word string, response Response) string {
	var r strings.Builder
	for i, char := range word {
		if response&(1<<(Length+(Length-i-1))) != 0 {
			fmt.Fprintf(&r, "\x1b[42;30m") // Correct
		} else if response&(1<<(Length-i-1)) != 0 {
			fmt.Fprintf(&r, "\x1b[43;30m") // Present
		} else {
			fmt.Fprintf(&r, "\x1b[40;37m") // Absent
		}
		fmt.Fprintf(&r, "%c\x1b[0m", char)
	}
	return r.String()
}

func FormatCodes(response Response) string {
	var c strings.Builder
	for i := range Length {
		if response&(1<<(Length+(Length-i-1))) != 0 {
			fmt.Fprintf(&c, "G") // Correct
		} else if response&(1<<(Length-i-1)) != 0 {
			fmt.Fprintf(&c, "Y") // Present
		} else {
			fmt.Fprintf(&c, "B") // Absent
		}
	}
	return c.String()
}

func FormatList(words []string, max int, separator string) string {
	count := len(words)
	suffix := ""
	if max != 0 && count > max {
		count = max
		suffix = fmt.Sprintf(" … %d more", len(words)-max)
	}
	return fmt.Sprintf("%s%s", strings.Join(words[:count], separator), suffix)
}

func PrintBoard(guesses []string, responses []Response, solution string, codes bool) {
	for i, guess := range guesses {
		fmt.Print(FormatResponse(guess, responses[i]))
		if codes {
			fmt.Printf("\t%s", FormatCodes(responses[i]))
		}
		fmt.Println()
	}
	if solution != "" {
		fmt.Printf("\x1b[1;32m%s\x1b[0m\n", solution)
	}
}

type parserFunc func(string) (string, bool)

func parser(s string) (string, bool) {
	if len(s) != Length {
		return "", false
	}
	return strings.ToUpper(s), true
}

func readStrings(path string, parser parserFunc) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	ret := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 {
			continue
		}
		s, ok := parser(t)
		if ok {
			ret = append(ret, s)
		} else {
			log.Fatalf("invalid string \"%s\"", t)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	return ret
}

func ReadWordLists() ([]string, []string) {
	wordlist, ok := os.LookupEnv("WORDLIST")
	if !ok {
		wordlist = DEFAULT_WORDLIST
	}

	exe, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}

	wordlistDir := filepath.Join(filepath.Dir(exe), "wordlists", wordlist)

	guesses := readStrings(filepath.Join(wordlistDir, "guesses.txt"), parser)
	solutions := readStrings(filepath.Join(wordlistDir, "solutions.txt"), parser)

	slices.Sort(guesses)
	slices.Sort(solutions)

	fmt.Printf("Using word list '%s' with %d guesses, %d solutions\n", wordlist, len(guesses), len(solutions))

	return guesses, solutions
}
