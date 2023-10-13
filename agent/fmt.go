package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

func BytesToHuman(n int64) string {

	switch {
	case n > PB:
		return fmt.Sprintf("%.1fPB", float64(n)/float64(PB))
	case n > TB:
		return fmt.Sprintf("%.1fTB", float64(n)/float64(TB))
	case n > GB:
		return fmt.Sprintf("%.1fGB", float64(n)/float64(GB))
	case n > MB:
		return fmt.Sprintf("%.1fMB", float64(n)/float64(MB))
	case n > KB:
		return fmt.Sprintf("%.1fKB", float64(n)/float64(KB))
	default:
		return fmt.Sprintf("%dB", n)
	}
}

func PrettyPrint(v interface{}) string {
	jsonBytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Print(err)
		return ""
	}

	return string(jsonBytes)
}
