package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput map[string]*votes
	}{
		{
			input:          "bmo++ increment 😊",
			expectedOutput: map[string]*votes{"bmo": &votes{1, 1}},
		},
		{
			input:          "bmo++ ykhr53++ great job!",
			expectedOutput: map[string]*votes{"bmo": &votes{1, 1}, "ykhr53": &votes{1, 1}},
		},
		{
			input:          "bmo-- decrement 😢",
			expectedOutput: map[string]*votes{"bmo": &votes{-1, 1}},
		},
		{
			input:          "bmo++ bmo-- neutral 😐",
			expectedOutput: map[string]*votes{"bmo": &votes{0, 2}},
		},
	}

	for i, tc := range tests {
		output := parse(tc.input)
		if !reflect.DeepEqual(output, tc.expectedOutput) {
			t.Errorf("Test %d: test fails", i)
		}
	}
}
