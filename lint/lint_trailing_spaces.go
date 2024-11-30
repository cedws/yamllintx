package lint

import (
	"errors"
	"iter"
	"strings"
)

type TrailingSpaces struct{}

func (t TrailingSpaces) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}

func (t TrailingSpaces) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		lastNonSpace := strings.LastIndexFunc(ctx.currentLine, func(r rune) bool {
			return r != ' '
		})

		if lastNonSpace != len(ctx.currentLine)-1 {
			problem := Problem{
				Line:   ctx.currentLineNumber,
				Column: lastNonSpace + 1,
				Error:  newLintError(errors.New("trailing spaces are forbidden")),
			}
			if !yield(problem) {
				return
			}
		}
	}
}
