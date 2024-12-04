package lint

import (
	"errors"
	"iter"
)

type TrailingSpaces struct{}

func (t TrailingSpaces) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}

func (t TrailingSpaces) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		trailingSpaces := trailingSpaces(ctx.currentLine)

		if trailingSpaces > 0 {
			problem := problem(
				ctx.currentLineNumber,
				trailingSpaces+1,
				errors.New("trailing spaces are forbidden"),
			)
			if !yield(problem) {
				return
			}
		}
	}
}
