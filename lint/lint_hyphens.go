package lint

import (
	"errors"
	"iter"
	"strings"

	"github.com/goccy/go-yaml/token"
)

var hypensMaxSpacesAfter = errors.New("too many spaces after hypen")

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

	origin := ctx.nextToken.Origin
	firstNonSpace := strings.IndexFunc(origin, func(r rune) bool {
		return r != ' '
	})
	if firstNonSpace == -1 {
		firstNonSpace = len(origin)
	}
	if firstNonSpace > h.MaxSpacesAfter {
		problem := Problem{
			Line:   ctx.nextToken.Position.Line,
			Column: ctx.nextToken.Position.Column + firstNonSpace,
			Error:  newLintError(hypensMaxSpacesAfter),
		}
		if !yield(problem) {
			return
		}
	}
}
