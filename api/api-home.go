package api

import (
	"net/http"
)

func ApiHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("parking lot api"))
	return
}
