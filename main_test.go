package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPotatofyContent(t *testing.T) {
	tc := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name: "basic case",
			input: `
			Wordle 789 6/6

			⬛⬛⬛🟩⬛
			⬛⬛⬛🟩⬛
			🟨⬛⬛🟩⬛
			🟨🟨🟨🟩⬛
			🟩🟩🟩⬛⬛
			🟩🟩🟩🟩🟩`,
			expectedOutput: `
			Wordle 789 6/6

			⬛⬛⬛🥔⬛
			⬛⬛⬛🥔⬛
			🍠⬛⬛🥔⬛
			🍠🍠🍠🥔⬛
			🥔🥔🥔⬛⬛
			🥔🥔🥔🥔🥔`,
		},
		{
			name:  "no newlines",
			input: `Wordle 790 4/6  🟨⬛⬛🟨⬛ ⬛⬛🟩🟨🟨 🟩🟨🟩⬛⬛ 🟩🟩🟩🟩🟩`,
			expectedOutput: `
			Wordle 790 4/6

			🍠⬛⬛🍠⬛
			⬛⬛🥔🍠🍠
			🥔🍠🥔⬛⬛
			🥔🥔🥔🥔🥔`,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedOutput, PotatofyWordle(tt.input))
		})
	}
}
