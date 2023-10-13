package main

import (
	"bytes"
	"encoding/hex"
	"os/exec"
	"strings"
	"text/scanner"
	"unicode"
)

func quotedStringSplit(input string) []string {
	s := scanner.Scanner{}

	b := bytes.NewBuffer([]byte(input))

	s.Init(b)

	var output []string

	var modeSingleQuote bool
	var modeDoubleQuote bool

	var outBuf string
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '"':
			modeDoubleQuote = !modeDoubleQuote
			if modeDoubleQuote {
				outBuf = outBuf + s.TokenText()
			} else {
				output = append(output, outBuf)
			}
		case '\'':
			modeSingleQuote = !modeSingleQuote
			if modeSingleQuote {
				outBuf = outBuf + s.TokenText()
			} else {
				output = append(output, outBuf)
			}
		case ' ':
			if !modeSingleQuote && !modeDoubleQuote {
				output = append(output, outBuf)
			}
		default:
			outBuf = outBuf + s.TokenText()
		}
	}

	return output
}

func run(command string) (string, error) {

	strings.Split(command, " ")
	args := quotedStringSplit(command)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if checkError(err) {
		return "", err
	}

	var safedOutput strings.Builder
	for b := range out {
		if !unicode.IsPrint(rune(b)) {
			safedOutput.WriteString(hex.EncodeToString([]byte{byte(b)}))
		} else {
			safedOutput.WriteByte(byte(b))
		}
	}

	return string(out), err
}
