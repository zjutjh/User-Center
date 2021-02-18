package util

import "net/http"

func RawHttpError(w http.ResponseWriter, err string, status int) {
	http.Error(w, err, status)
}