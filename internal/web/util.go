package web

import (
	"fmt"
	"net/http"
)

type ApiMessage struct {
	Code int
	Message string
}

func (m ApiMessage) Write(res http.ResponseWriter) {
	res.WriteHeader(m.Code)
	res.Write([]byte(fmt.Sprintf("{'message':'%s'}", m.Message)))
}

func (m ApiMessage) WriteWithError(res http.ResponseWriter, err error) {
	res.WriteHeader(m.Code)
	res.Write([]byte(fmt.Sprintf("{'message':'%s', 'error': '%s'}", m.Message, err.Error())))
}

var Unauthorized        = ApiMessage{ Code: http.StatusUnauthorized,        Message: "Unauthorized" }
var MethodNotAllowed    = ApiMessage{ Code: http.StatusMethodNotAllowed,    Message: "Method Not Allowed" }
var InternalServerError = ApiMessage{ Code: http.StatusInternalServerError, Message: "Internal Server Error" }
var BadRequest          = ApiMessage{ Code: http.StatusBadRequest,          Message: "Bad Request" }
