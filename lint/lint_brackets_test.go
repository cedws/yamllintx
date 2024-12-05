package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrackets(t *testing.T) {
	tests := []struct {
		name        string
		lint        Brackets
		input       string
		expectedErr error
	}{
		{
			name: "ForbidAll Pass",
			lint: Brackets{
				Forbid: ForbidBracketsAll,
			},
			input: `
---
object:
  key1: 4
  key2: 8
`,
			expectedErr: nil,
		},
		{
			name: "ForbidAll Fail",
			lint: Brackets{
				Forbid: ForbidBracketsAll,
			},
			input: `
---
object: [ 1, 2, abc ]
`,
			expectedErr: ErrBracketsForbidden,
		},
		{
			name: "ForbidNonEmpty Pass",
			lint: Brackets{
				Forbid: ForbidBracketsNonEmpty,
			},
			input: `
---
object: []
`,
			expectedErr: nil,
		},
		{
			name: "ForbidNonEmpty Fail",
			lint: Brackets{
				Forbid: ForbidBracketsNonEmpty,
			},
			input: `
---
object: [ 1, 2, abc ]
`,
			expectedErr: ErrBracketsNonEmptyForbidden,
		},
		{
			name: "MinSpacesInside Pass",
			lint: Brackets{
				MinSpacesInside: 1,
				MaxSpacesInside: 3,
			},
			input: `
---
object: [ 1, 2, abc ]
`,
			expectedErr: nil,
		},
		{
			name: "MinSpacesInside Fail",
			lint: Brackets{
				MinSpacesInside: 2,
				MaxSpacesInside: 3,
			},
			input: `
---
object: [ 1, 2, abc ]
`,
			expectedErr: ErrBracketsTooFewSpaces,
		},
		{
			name: "MaxSpacesInside Pass",
			lint: Brackets{
				MinSpacesInside: 1,
				MaxSpacesInside: 3,
			},
			input: `
---
object: [ 1, 2, abc   ]
`,
			expectedErr: nil,
		},
		{
			name: "MaxSpacesInside Fail",
			lint: Brackets{
				MinSpacesInside: 1,
				MaxSpacesInside: 2,
			},
			input: `
---
object: [ 1, 2, abc   ]
`,
			expectedErr: ErrBracketsTooManySpaces,
		},
		{
			name: "MinSpacesInsideEmpty Pass",
			lint: Brackets{
				MinSpacesInsideEmpty: 0,
				MaxSpacesInsideEmpty: 0,
			},
			input: `
---
object: []
`,
			expectedErr: nil,
		},
		{
			name: "MinSpacesInsideEmpty Fail",
			lint: Brackets{
				MinSpacesInsideEmpty: 1,
				MaxSpacesInsideEmpty: -1,
			},
			input: `
---
object: []
`,
			expectedErr: ErrBracketsTooFewSpacesEmpty,
		},
		{
			name: "MaxSpacesInsideEmpty Pass",
			lint: Brackets{
				MinSpacesInsideEmpty: 0,
				MaxSpacesInsideEmpty: 3,
			},
			input: `
---
object: [ ]
`,
			expectedErr: nil,
		},
		{
			name: "MaxSpacesInsideEmpty Fail",
			lint: Brackets{
				MinSpacesInsideEmpty: 0,
				MaxSpacesInsideEmpty: 1,
			},
			input: `
---
object: [   ]
`,
			expectedErr: ErrBracketsTooManySpacesEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			problem := Lint([]byte(tt.input), tt.lint)
			if problem == nil {
				assert.Nil(t, tt.expectedErr, "expected problem but got nil")
				return
			}

			assert.ErrorIs(t, problem.Error, tt.expectedErr)
		})
	}
}
