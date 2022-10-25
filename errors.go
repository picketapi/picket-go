// Package picketapi is a Go client for the Picket API
//
// The package is a simple wrapper around the Picket API for
// Go applications
package picketapi

type ErrorResponse struct {
	Msg  string `json:"msg"`
	Code string `json:"code"`
}

func (err ErrorResponse) Error() string {
	return err.Msg
}

func is2xxStatusCode(code int) bool {
	return code >= 200 && code < 300
}
