package lib

import (
	"log"
	"strings"
)

func CheckErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}
