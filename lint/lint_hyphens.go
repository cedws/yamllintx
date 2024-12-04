package lint

import (
	"errors"
	"iter"

	"github.com/goccy/go-yaml/token"
)

var ErrHypensMaxSpacesAfter = errors.New("too many spaces after hypen")

type Hyphens struct {
	MaxSpacesAfter int
}

func (h Hyphens) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		if h.MaxSpacesAfter > 0 {
			h.checkMaxSpacesAfter(ctx, yield)
		}
	}
}

func (h Hyphens) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}

func (h Hyphens) checkMaxSpacesAfter(ctx tokenContext, yield func(Problem) bool) {
	if ctx.currentToken.Type != token.SequenceEntryType || ctx.nextToken == nil {
		return
	}

	leadingSpaces := leadingSpaces(ctx.nextToken.Origin)

	if leadingSpaces > h.MaxSpacesAfter {
		problem := problem(
			ctx.nextToken.Position.Line,
			ctx.nextToken.Position.Column+leadingSpaces,
			ErrHypensMaxSpacesAfter,
		)
		if !yield(problem) {
			return
		}
	}
}
