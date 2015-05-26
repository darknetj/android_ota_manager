package lib

import (
  "log"
)

func CheckErr(err error, msg string) {
    if err != nil {
        log.Println(msg, err)
    }
}
