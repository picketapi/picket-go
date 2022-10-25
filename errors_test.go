package picketapi_test

import (
	"testing"

	picketapi "github.com/picketapi/picket-go"
)

func TestErrorRepsonseError(t *testing.T) {
	want := "test error"
	err := picketapi.ErrorResponse{
		Code: "code",
		Msg:  want,
	}

	if err.Error() != want {
		t.Errorf("Error() = %s, want %s", err.Error(), want)
	}
}
