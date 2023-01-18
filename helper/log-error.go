package helper

import "log"

func LogIfError(err error) {
	if err != nil {
		log.Println(err)
	}
}
