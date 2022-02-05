package main

import "log"

func HandleErr(tip string, err error) {
	if err != nil {
		log.Panic(tip, err)
	}
}
