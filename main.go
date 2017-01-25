package kaa

import (
	"net/http"
)

var SNAKE_ID string

func init() {
	http.HandleFunc("/", handleInfo)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/end", handleEnd)

	SNAKE_ID = "Kaa"
}
