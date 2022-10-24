package picket_test

import (
	"testing"

	picket "github.com/picketapi/picket-go"
)

func TestErrorRepsonseError(t *testing.T) {
	want := "test error"
	err := picket.ErrorResponse{
		Code: "code",
		Msg:  want,
	}

	if err.Error() != want {
		t.Errorf("Error() = %s, want %s", err.Error(), want)
	}
}
