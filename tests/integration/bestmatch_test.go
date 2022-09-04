package integration

import (
	"testing"
)

func Test_GetMatches_Validations(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "success"},
		{name: "error - no materials specified"},
		{name: "error - no address provided"},
		{name: "error - no phone number"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {

		})
	}
}

func Test_GetMatches_Bestmatch(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "no matches found"},
		{name: "one match in range"},
		{name: "two matches, ordered by rating"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {

		})
	}
}
