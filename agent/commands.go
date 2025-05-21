package main

import (
	"fmt"
	"strings"
)

type Command struct {
	UUID    string
	Name    string
	Input   string
	Output  string
	Special specialCommand
}

type specialCommand int64

const (
	specialUpgrade specialCommand = 1 << iota
)

type checkinResponse struct {
	ID             int
	Commands       []Command
	StreamActivity bool
}

// func quotedStringSplit(input string) []string {

// 	b := bytes.NewBuffer([]byte(input))
// 	s := bufio.NewScanner(b)

// 	var output []string

// 	var modeSingleQuote bool
// 	var modeDoubleQuote bool

// 	var outBuf string

// 	s.Split(bufio.ScanWords)

// 	for s.Scan() {
// 		t := s.Text()
// 		//log.Printf("Token: %s", s.Text())
// 		outBuf = ""
// 		switch t {
// 		case `"`:
// 			modeDoubleQuote = !modeDoubleQuote
// 			if modeDoubleQuote {
// 				outBuf = outBuf + t
// 			} else {
// 				output = append(output, outBuf)
// 			}
// 		case `'`:
// 			modeSingleQuote = !modeSingleQuote
// 			if modeSingleQuote {
// 				outBuf = outBuf + t
// 			} else {
// 				output = append(output, outBuf)
// 			}
// 		case " ":
// 			if !modeSingleQuote && !modeDoubleQuote {
// 				output = append(output, outBuf)
// 			}
// 		default:
// 			outBuf = outBuf + t
// 			output = append(output, outBuf)
// 		}

// 	}

// 	return output
// }

type Command struct {
	UUID    string
	Name    string
	Input   string
	Output  string
	Special specialCommand
}

type specialCommand int64

const (
	specialUpgrade specialCommand = 1 << iota
)

// ParseCommand splits a command string into tokens, preserving quoted substrings.
func quotedStringSplit(input string) ([]string, error) {
	var tokens []string
	var currentToken strings.Builder
	var inQuotes bool
	var quoteChar rune

	for i, char := range input {
		switch char {
		case '"', '\'':
			if inQuotes {
				// End quote if it matches the current quote character
				if char == quoteChar {
					inQuotes = false
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				} else {
					// Append mismatched quotes as part of the token
					currentToken.WriteRune(char)
				}
			} else {
				// Start a quoted string
				inQuotes = true
				quoteChar = char
			}
		case ' ':
			if inQuotes {
				// Treat spaces inside quotes as part of the token
				currentToken.WriteRune(char)
			} else {
				// End the current token and start a new one
				if currentToken.Len() > 0 {
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
			}
		default:
			// Append regular characters to the current token
			currentToken.WriteRune(char)
		}

		// Handle the last token if we're at the end of the input
		if i == len(input)-1 && currentToken.Len() > 0 {
			if inQuotes {
				return nil, fmt.Errorf("unterminated quote detected")
			}
			tokens = append(tokens, currentToken.String())
		}
	}

	return tokens, nil
}
