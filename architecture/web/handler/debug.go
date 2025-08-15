package handler

import (
	"log"
	"net/http"
)

// ? debugLogHandler -
func debugLogHandler(fName string, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, fName)
}
