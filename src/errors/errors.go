package errors

import (
	"net/http"
)

/*
	All errors must define the handler name in c.Error()
*/

type RequestError struct {
	CODE      int
	MESSAGE   string
	META      string
	HANDLER   string
	ERRORTYPE string
}

func (err *RequestError) Error() string {
	return err.MESSAGE
}

func (err *RequestError) Code() int {
	return err.CODE
}

func (err *RequestError) Meta() map[string]string {
	return map[string]string{"META": err.META, "HANDLER": err.HANDLER, "ERRORTYPE": err.ERRORTYPE}
}

func NewError(err *RequestError, meta string, handler string) *RequestError {
	err.META = meta
	err.HANDLER = handler
	return err
}

// for returning error response in gin.JSON()
func ErrorResponse(err *RequestError) (code int, message map[string]interface{}) {
	return err.CODE, map[string]interface{}{"message": err.MESSAGE, "meta": err.META}
}

// all errors are defined here with generic message
// use NewError to overwrite message field
var (
	ErrMissingKey = &RequestError{
		ERRORTYPE: "ErrMissingKey",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "key is missing",
	}
	ErrInvalidSchema = &RequestError{
		ERRORTYPE: "ErrInvalidSchema",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "payload schema is invalid",
	}
	ErrFailedDBConn = &RequestError{
		ERRORTYPE: "ErrFailedDBConn",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "unable to establish connection with database",
	}
	ErrFailedSQLQuery = &RequestError{
		ERRORTYPE: "ErrFailedSQLQuery",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "unable to query from sql database",
	}
	ErrFailedSQLExec = &RequestError{
		ERRORTYPE: "ErrFailedSQLExec",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "unable to execute statement in sql database",
	}
	ErrAbortController = &RequestError{
		ERRORTYPE: "ErrAbortController",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "controller aborted",
	}
	ErrMissingResource = &RequestError{
		ERRORTYPE: "ErrMissingResource",
		CODE:      http.StatusNotFound,
		MESSAGE:   "resource is missing",
	}
	ErrRowScan = &RequestError{
		ERRORTYPE: "ErrRowScan",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "unable to scan into destination variables",
	}
	ErrSendMail = &RequestError{
		ERRORTYPE: "ErrSendMail",
		CODE:      http.StatusBadRequest,
		MESSAGE:   "unable to send mail via smtp",
	}
)
