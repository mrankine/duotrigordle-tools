# duotrigordle tools

This repo contains four Go programs which may be useful for players of [Duotrigordle](https://duotrigordle.com) and similar word games.

- `maximise-solvable` suggests a first guess which maximises the expected number of _solvable_ boards. A board is _solvable_ if there is only one remaining possible solution. The first guess may or may not also be a valid solution.
- `report-solvable` reports _solvable_ boards for a given sequence of one or more guesses.
- `possible-words` reports the remaining possible solutions and guesses after a sequence of guesses and responses.
- `guess` reports the game's response to a guess, given the solution.

Example usage:

- based on the latest word list, `./maximise-solvable` reports `PITON` is the guess which produces the most solvable boards (43).
- `./report-solvable PITON` prints the 43 solvable boards. For example, if the game responds to `PITON` with Green Yellow Black Green Black, the solution must be `PRIOR`.
- `./possible-words -s PITON BGBGB` reports that there are four possible solutions if the response to `PITON` is Black Green Black Green Black. They are `RIGOR`, `VIGOR`, `VISOR`, and `WIDOW`.

_Number of solvable boards_ is only one metric of guess quality. More sophisticated approaches are possible. See, for example, [There's A New Best Starter For Perfect Duotrigordle](https://www.youtube.com/watch?v=Hk5BNh1DtTU) (24 May 2022, uses old word lists).

## Build and run

Requires [Go](https://go.dev/doc/install), and a terminal with ANSI colour support.

To build:

```bash
git clone https://github.com/mrankine/duotrigordle-tools.git
cd duotrigordle-tools
go build -o . ./...
```

To run:

```bash
./maximise-solvable
```

```bash
./report-solvable guess [guess...]
```

```bash
./possible-words [-s|--show-words] guess response [guess response ...]
```

```bash
./guess guess solution
```

Where:

- `guess` is a five-letter word. Must be present in `guesses.txt`.
- `solution` is a five-letter word. Must be present in `solutions.txt`.
- `response` is the game's response to a guess. Each position of the guess is either Green (correct), Yellow (present), or Black (not in this position). Examples of valid response formats:
  - `GYGBB` corresponds to Green, Yellow, Green, Black, Black
  - `21200` same as above
  - `Y` all positions were Yellow


## Word lists

Word lists are located in `wordlists/{identifier}/{guesses,solutions}.txt`. The `wordlists` directory must be in the same directory as the executables.

The default word list is `2024-01-14`, which is consistent with Duotrigordle's 14 January 2024 update.

Specify an alternative word list using the environment variable `WORDLIST`. Examples:

```bash
export WORDLIST=2022-03-03
./maximise-solvable
```

```bash
WORDLIST=my-custom-wordlist ./maximise-solvable
```
