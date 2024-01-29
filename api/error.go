package api

import (
	"fmt"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, code int, message string, err error) {
	var httpErr error
	if err == nil {
		httpErr = fmt.Errorf("%s", message)
	} else {
		if message == "" {
			httpErr = fmt.Errorf("%v", err.Error())
		}
		httpErr = fmt.Errorf("%s. Reason %v", message, err.Error())
	}

	log.Printf(httpErr.Error())
	http.Error(w, httpErr.Error(), code)
}
