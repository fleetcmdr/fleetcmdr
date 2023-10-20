package main

import (
	"bufio"
	"bytes"
)

func quotedStringSplit(input string) []string {

	b := bytes.NewBuffer([]byte(input))
	s := bufio.NewScanner(b)

	var output []string

	var modeSingleQuote bool
	var modeDoubleQuote bool

	var outBuf string

	s.Split(bufio.ScanWords)

	for s.Scan() {
		t := s.Text()
		//log.Printf("Token: %s", s.Text())
		outBuf = ""
		switch t {
		case `"`:
			modeDoubleQuote = !modeDoubleQuote
			if modeDoubleQuote {
				outBuf = outBuf + t
			} else {
				output = append(output, outBuf)
			}
		case `'`:
			modeSingleQuote = !modeSingleQuote
			if modeSingleQuote {
				outBuf = outBuf + t
			} else {
				output = append(output, outBuf)
			}
		case " ":
			if !modeSingleQuote && !modeDoubleQuote {
				output = append(output, outBuf)
			}
		default:
			outBuf = outBuf + t
			output = append(output, outBuf)
		}

	}

	return output
}
