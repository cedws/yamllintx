package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnchorForbidUndeclaredAliases(t *testing.T) {
	const pass1 = `
---
- &anchor
  foo: bar
- *anchor`

	const fail1 = `
---
- &anchor
	foo: bar
- *unknown`

	const fail2 = `
---
- &anchor
	foo: bar
- <<: *unknown
	extra: value`

	t.Run("Pass", func(t *testing.T) {
		for _, src := range []string{pass1} {
			lint := NewAnchors()
			lint.ForbidUndeclaredAliases = true

			err := Lint(src, lint)
			assert.NoError(t, err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		for _, src := range []string{fail1, fail2} {
			lint := NewAnchors()
			lint.ForbidUndeclaredAliases = true

			err := Lint(src, lint)
			assert.Error(t, err)
		}
	})
}

func TestAnchorForbidDuplicatedAliases(t *testing.T) {
	const pass1 = `
---
- &anchor1 Foo Bar
- &anchor2 [item 1, item 2]`

	const fail1 = `
---
- &anchor Foo Bar
- &anchor [item 1, item 2]`

	t.Run("Pass", func(t *testing.T) {
		for _, src := range []string{pass1} {
			lint := NewAnchors()
			lint.ForbidDuplicatedAnchors = true

			err := Lint(src, lint)
			assert.NoError(t, err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		for _, src := range []string{fail1} {
			lint := NewAnchors()
			lint.ForbidDuplicatedAnchors = true

			err := Lint(src, lint)
			assert.Error(t, err)
		}
	})
}

func TestAnchorForbidUnusedAnchors(t *testing.T) {
	const pass1 = `
---
- &anchor
  foo: bar
- *anchor`

	const fail1 = `
---
- &anchor
  foo: bar
- items:
  - item1
  - item2`

	t.Run("Pass", func(t *testing.T) {
		for _, src := range []string{pass1} {
			lint := NewAnchors()
			lint.ForbidUnusedAnchors = true

			err := Lint(src, lint)
			assert.NoError(t, err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		for _, src := range []string{fail1} {
			lint := NewAnchors()
			lint.ForbidUnusedAnchors = true

			err := Lint(src, lint)
			assert.Error(t, err)
		}
	})
}
