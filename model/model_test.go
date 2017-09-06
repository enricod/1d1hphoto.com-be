package model

import (
	"testing"
)

func TestRandomStringBytes(t *testing.T) {
	l := 5
	result, _ := GenerateRandomString(l)

	if len(result) != l {
		t.Errorf("lunghezza non corretta, got %s, want %s ", len(result), 1)
	}
}
