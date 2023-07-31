package files

import (
	"testing"
)

func TestIsValidFile(t *testing.T) {
	correct := ".jpg"
	wrong := ".asp"
	veryWrong := "foobar"

	if ok := isValidFiletype(correct); !ok {
		t.Errorf("expected %q to be valid, was invalid", correct)
	}

	if ok := isValidFiletype(wrong); ok {
		t.Errorf("expected %q to be invalid, was valid", wrong)
	}

	if ok := isValidFiletype(veryWrong); ok {
		t.Errorf("expected %q to be invalid, was valid", veryWrong)
	}
}
