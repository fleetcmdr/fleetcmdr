package main

import (
	"encoding/hex"
	"log"
	"os/exec"
	"strings"
	"unicode"
)

func run(command string) (string, error) {

	//command = fmt.Sprintf(`$(/bin/bash -c "%s")`, command)
	log.Printf("Running command \"%s\"", command)

	args := quotedStringSplit(command)
	//log.Printf("%+v", args)

	var cmd *exec.Cmd
	if len(args) == 1 {
		cmd = exec.Command(args[0])
	} else if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	}

	out, err := cmd.CombinedOutput()
	if checkError(err) {
		return "", err
	}

	//log.Printf("Output: %s", string(out))

	//err = cmd.Wait()
	//if checkError(err) {
	//	return "", err
	//}

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
