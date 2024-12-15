package main

import (
	"context"
	"encoding/hex"
	"log"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

func run(command string) (string, error) {

	//command = fmt.Sprintf(`$(/bin/bash -c "%s")`, command)
	log.Printf("Running command \"%s\"", command)

	args := quotedStringSplit(command)
	//log.Printf("%+v", args)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)

	_ = cf

	cmd := &exec.Cmd{}
	if len(args) == 1 {
		cmd = exec.CommandContext(ctx, args[0])
	} else if len(args) > 1 {
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
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
