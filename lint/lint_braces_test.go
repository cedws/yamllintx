package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBraces(t *testing.T) {
	tests := []struct {
		name        string
		lint        Linter
		input       string
		expectedErr error
	}{
		{
			name: "ForbidAll Pass",
			lint: Braces{
				Forbid: ForbidBracesAll,
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
			lint: Braces{
				Forbid: ForbidBracesAll,
			},
			input: `
---
object: { key1: 4, key2: 8 }
`,
			expectedErr: ErrBracesForbidden,
		},
		{
			name: "ForbidNonEmpty Pass",
			lint: Braces{
				Forbid: ForbidBracesNonEmpty,
			},
			input: `
---
object: {}
`,
			expectedErr: nil,
		},
		{
			name: "ForbidNonEmpty Fail",
			lint: Braces{
				Forbid: ForbidBracesNonEmpty,
			},
			input: `
---
object: { key1: 4, key2: 8 }
`,
			expectedErr: ErrBracesNonEmptyForbidden,
		},
		{
			name: "MinSpacesInside Pass",
			lint: Braces{
				MinSpacesInside: 1,
				MaxSpacesInside: 3,
			},
			input: `
---
object: { key1: 4, key2: 8 }
`,
			expectedErr: nil,
		},
		{
			name: "MinSpacesInside Fail",
			lint: Braces{
				MinSpacesInside: 2,
				MaxSpacesInside: 3,
			},
			input: `
---
object: { key1: 4, key2: 8 }
`,
			expectedErr: ErrBracesTooFewSpaces,
		},
		{
			name: "MaxSpacesInside Pass",
			lint: Braces{
				MinSpacesInside: 1,
				MaxSpacesInside: 3,
			},
			input: `
---
object: { key1: 4, key2: 8  }
`,
			expectedErr: nil,
		},
		{
			name: "MaxSpacesInside Fail",
			lint: Braces{
				MinSpacesInside: 1,
				MaxSpacesInside: 2,
			},
			input: `
---
object: { key1: 4, key2: 8   }
`,
			expectedErr: ErrBracesTooManySpaces,
		},
		{
			name: "MinSpacesInsideEmpty Pass",
			lint: Braces{
				MinSpacesInsideEmpty: 0,
				MaxSpacesInsideEmpty: 0,
			},
			input: `
---
object: {}
`,
			expectedErr: nil,
		},
		{
			name: "MinSpacesInsideEmpty Fail",
			lint: Braces{
				MinSpacesInsideEmpty: 1,
				MaxSpacesInsideEmpty: -1,
			},
			input: `
---
object: {}
`,
			expectedErr: ErrBracesTooFewSpacesEmpty,
		},
		{
			name: "MaxSpacesInsideEmpty Pass",
			lint: Braces{
				MinSpacesInsideEmpty: 0,
				MaxSpacesInsideEmpty: 3,
			},
			input: `
---
object: { }
`,
			expectedErr: nil,
		},
		{
			name: "MaxSpacesInsideEmpty Fail",
			lint: Braces{
				MinSpacesInsideEmpty: 0,
				MaxSpacesInsideEmpty: 1,
			},
			input: `
---
object: {   }
`,
			expectedErr: ErrBracesTooManySpacesEmpty,
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
