package web

import (
	"net/http"
)

func RootApiHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("API"))
}
