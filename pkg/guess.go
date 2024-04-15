package pkg

const Length = 5

// Response is a 10-bit number representing the game's response to each letter
// of a guess, ordered from left to right. The low 5 bits represent Presents
// (yellow), high 5 represent Corrects (green).
type Response uint16

const AllCorrect Response = 0b1111100000
const AllPresent Response = 0b0000011111
const AllAbsent Response = 0b0000000000

// CheckGuess returns a Response for the given guess and solution. First,
// construct a pool of letters which are potentially Present (that is, all
// solution letters which are not Correct). Then for each guess letter, report
// Correct if it matches the solution letter, or Present if it's in the pool (and
// decrement the pool).
// Algorithm adapted from https://github.com/goingonit/wordle
func CheckGuess(guess, solution string) Response {
	// pool1, pool2 and mask are 26-bit masks
	var pool1, pool2, mask uint32
	for i := range Length {
		if guess[i] != solution[i] {
			mask = 1 << (solution[i] - 'A')
			if (pool1 & pool2 & mask) > 0 {
				panic("Found word with >3 repeat letters")
			}
			// logic for nth occurrence of letter:
			// 1st: pool1 is off so leave pool2 alone, flip pool1 (now on)
			// 2nd: pool1 is on so flip pool2 (now on), flip pool1 (now off)
			// 3rd: pool1 is off so leave pool2 alone, flip pool1 (now on)
			pool2 ^= pool1 & mask
			pool1 ^= mask
		}
	}

	// if there's nothing in the pool, all letters must be Correct
	if pool1 == 0 && pool2 == 0 {
		return AllCorrect
	}

	var corrects, presents, positionMask uint16
	var letterMask uint32
	for i, letter := range guess {
		positionMask = 1 << (Length - i - 1)
		letterMask = 1 << (letter - 'A')
		if guess[i] == solution[i] {
			corrects |= positionMask
		} else if ((pool2 | pool1) & letterMask) > 0 {
			presents |= positionMask
			// logic for nth occurrence of letter:
			// 3rd: pool1 is on so leave pool2 alone, flip pool1 (now off)
			// 2nd: pool1 is off so flip pool2 (now off), flip pool1 (now on)
			// 1st: pool1 is on so leave pool2 alone, flip pool1 (now off)
			pool2 ^= ^pool1 & letterMask
			pool1 ^= letterMask
		}
	}

	return Response((corrects << Length) | presents)
}
